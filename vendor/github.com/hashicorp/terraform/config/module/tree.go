package module

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/terraform/tfdiags"

	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform/config"
)

// RootName is the name of the root tree.
const RootName = "root"

// Tree represents the module import tree of configurations.
//
// This Tree structure can be used to get (download) new modules, load
// all the modules without getting, flatten the tree into something
// Terraform can use, etc.
type Tree struct {
	name     string
	config   *config.Config
	children map[string]*Tree
	path     []string
	lock     sync.RWMutex

	// version is the final version of the config loaded for the Tree's module
	version string
	// source is the "source" string used to load this module. It's possible
	// for a module source to change, but the path remains the same, preventing
	// it from being reloaded.
	source string
	// parent allows us to walk back up the tree and determine if there are any
	// versioned ancestor modules which may effect the stored location of
	// submodules
	parent *Tree
}

// NewTree returns a new Tree for the given config structure.
func NewTree(name string, c *config.Config) *Tree {
	return &Tree{config: c, name: name}
}

// NewEmptyTree returns a new tree that is empty (contains no configuration).
func NewEmptyTree() *Tree {
	t := &Tree{config: &config.Config{}}

	// We do this dummy load so that the tree is marked as "loaded". It
	// should never fail because this is just about a no-op. If it does fail
	// we panic so we can know its a bug.
	if err := t.Load(&Storage{Mode: GetModeGet}); err != nil {
		panic(err)
	}

	return t
}

// NewTreeModule is like NewTree except it parses the configuration in
// the directory and gives it a specific name. Use a blank name "" to specify
// the root module.
func NewTreeModule(name, dir string) (*Tree, error) {
	c, err := config.LoadDir(dir)
	if err != nil {
		return nil, err
	}

	return NewTree(name, c), nil
}

// Config returns the configuration for this module.
func (t *Tree) Config() *config.Config {
	return t.config
}

// Child returns the child with the given path (by name).
func (t *Tree) Child(path []string) *Tree {
	if t == nil {
		return nil
	}

	if len(path) == 0 {
		return t
	}

	c := t.Children()[path[0]]
	if c == nil {
		return nil
	}

	return c.Child(path[1:])
}

// Children returns the children of this tree (the modules that are
// imported by this root).
//
// This will only return a non-nil value after Load is called.
func (t *Tree) Children() map[string]*Tree {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.children
}

// DeepEach calls the provided callback for the receiver and then all of
// its descendents in the tree, allowing an operation to be performed on
// all modules in the tree.
//
// Parents will be visited before their children but otherwise the order is
// not defined.
func (t *Tree) DeepEach(cb func(*Tree)) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	t.deepEach(cb)
}

func (t *Tree) deepEach(cb func(*Tree)) {
	cb(t)
	for _, c := range t.children {
		c.deepEach(cb)
	}
}

// Loaded says whether or not this tree has been loaded or not yet.
func (t *Tree) Loaded() bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.children != nil
}

// Modules returns the list of modules that this tree imports.
//
// This is only the imports of _this_ level of the tree. To retrieve the
// full nested imports, you'll have to traverse the tree.
func (t *Tree) Modules() []*Module {
	result := make([]*Module, len(t.config.Modules))
	for i, m := range t.config.Modules {
		result[i] = &Module{
			Name:      m.Name,
			Version:   m.Version,
			Source:    m.Source,
			Providers: m.Providers,
		}
	}

	return result
}

// Name returns the name of the tree. This will be "<root>" for the root
// tree and then the module name given for any children.
func (t *Tree) Name() string {
	if t.name == "" {
		return RootName
	}

	return t.name
}

// Load loads the configuration of the entire tree.
//
// The parameters are used to tell the tree where to find modules and
// whether it can download/update modules along the way.
//
// Calling this multiple times will reload the tree.
//
// Various semantic-like checks are made along the way of loading since
// module trees inherently require the configuration to be in a reasonably
// sane state: no circular dependencies, proper module sources, etc. A full
// suite of validations can be done by running Validate (after loading).
func (t *Tree) Load(s *Storage) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	children, err := t.getChildren(s)
	if err != nil {
		return err
	}

	// Go through all the children and load them.
	for _, c := range children {
		if err := c.Load(s); err != nil {
			return err
		}
	}

	// Set our tree up
	t.children = children

	return nil
}

func (t *Tree) getChildren(s *Storage) (map[string]*Tree, error) {
	children := make(map[string]*Tree)

	// Go through all the modules and get the directory for them.
	for _, m := range t.Modules() {
		if _, ok := children[m.Name]; ok {
			return nil, fmt.Errorf(
				"module %s: duplicated. module names must be unique", m.Name)
		}

		// Determine the path to this child
		modPath := make([]string, len(t.path), len(t.path)+1)
		copy(modPath, t.path)
		modPath = append(modPath, m.Name)

		log.Printf("[TRACE] module source: %q", m.Source)

		// add the module path to help indicate where modules with relative
		// paths are being loaded from
		s.output(fmt.Sprintf("- module.%s", strings.Join(modPath, ".")))

		// Lookup the local location of the module.
		// dir is the local directory where the module is stored
		mod, err := s.findRegistryModule(m.Source, m.Version)
		if err != nil {
			return nil, err
		}

		// The key is the string that will be used to uniquely id the Source in
		// the local storage.  The prefix digit can be incremented to
		// invalidate the local module storage.
		key := "1." + t.versionedPathKey(m)
		if mod.Version != "" {
			key += "." + mod.Version
		}

		// Check for the exact key if it's not a registry module
		if !mod.registry {
			mod.Dir, err = s.findModule(key)
			if err != nil {
				return nil, err
			}
		}

		if mod.Dir != "" && s.Mode != GetModeUpdate {
			// We found it locally, but in order to load the Tree we need to
			// find out if there was another subDir stored from detection.
			subDir, err := s.getModuleRoot(mod.Dir)
			if err != nil {
				// If there's a problem with the subdir record, we'll let the
				// recordSubdir method fix it up.  Any other filesystem errors
				// will turn up again below.
				log.Println("[WARN] error reading subdir record:", err)
			}

			fullDir := filepath.Join(mod.Dir, subDir)

			child, err := NewTreeModule(m.Name, fullDir)
			if err != nil {
				return nil, fmt.Errorf("module %s: %s", m.Name, err)
			}
			child.path = modPath
			child.parent = t
			child.version = mod.Version
			child.source = m.Source
			children[m.Name] = child
			continue
		}

		// Split out the subdir if we have one.
		// Terraform keeps the entire requested tree, so that modules can
		// reference sibling modules from the same archive or repo.
		rawSource, subDir := getter.SourceDirSubdir(m.Source)

		// we haven't found a source, so fallback to the go-getter detectors
		source := mod.url
		if source == "" {
			source, err = getter.Detect(rawSource, t.config.Dir, getter.Detectors)
			if err != nil {
				return nil, fmt.Errorf("module %s: %s", m.Name, err)
			}
		}

		log.Printf("[TRACE] detected module source %q", source)

		// Check if the detector introduced something new.
		// For example, the registry always adds a subdir of `//*`,
		// indicating that we need to strip off the first component from the
		// tar archive, though we may not yet know what it is called.
		source, detectedSubDir := getter.SourceDirSubdir(source)
		if detectedSubDir != "" {
			subDir = filepath.Join(detectedSubDir, subDir)
		}

		output := ""
		switch s.Mode {
		case GetModeUpdate:
			output = fmt.Sprintf("  Updating source %q", m.Source)
		default:
			output = fmt.Sprintf("  Getting source %q", m.Source)
		}
		s.output(output)

		dir, ok, err := s.getStorage(key, source)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("module %s: not found, may need to run 'terraform init'", m.Name)
		}

		log.Printf("[TRACE] %q stored in %q", source, dir)

		// expand and record the subDir for later
		fullDir := dir
		if subDir != "" {
			fullDir, err = getter.SubdirGlob(dir, subDir)
			if err != nil {
				return nil, err
			}

			// +1 to account for the pathsep
			if len(dir)+1 > len(fullDir) {
				return nil, fmt.Errorf("invalid module storage path %q", fullDir)
			}
			subDir = fullDir[len(dir)+1:]
		}

		// add new info to the module record
		mod.Key = key
		mod.Dir = dir
		mod.Root = subDir

		// record the module in our manifest
		if err := s.recordModule(mod); err != nil {
			return nil, err
		}

		child, err := NewTreeModule(m.Name, fullDir)
		if err != nil {
			return nil, fmt.Errorf("module %s: %s", m.Name, err)
		}
		child.path = modPath
		child.parent = t
		child.version = mod.Version
		child.source = m.Source
		children[m.Name] = child
	}

	return children, nil
}

// Path is the full path to this tree.
func (t *Tree) Path() []string {
	return t.path
}

// String gives a nice output to describe the tree.
func (t *Tree) String() string {
	var result bytes.Buffer
	path := strings.Join(t.path, ", ")
	if path != "" {
		path = fmt.Sprintf(" (path: %s)", path)
	}
	result.WriteString(t.Name() + path + "\n")

	cs := t.Children()
	if cs == nil {
		result.WriteString("  not loaded")
	} else {
		// Go through each child and get its string value, then indent it
		// by two.
		for _, c := range cs {
			r := strings.NewReader(c.String())
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				result.WriteString("  ")
				result.WriteString(scanner.Text())
				result.WriteString("\n")
			}
		}
	}

	return result.String()
}

// Validate does semantic checks on the entire tree of configurations.
//
// This will call the respective config.Config.Validate() functions as well
// as verifying things such as parameters/outputs between the various modules.
//
// Load must be called prior to calling Validate or an error will be returned.
func (t *Tree) Validate() tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	if !t.Loaded() {
		diags = diags.Append(fmt.Errorf(
			"tree must be loaded before calling Validate",
		))
		return diags
	}

	// Terraform core does not handle root module children named "root".
	// We plan to fix this in the future but this bug was brought up in
	// the middle of a release and we don't want to introduce wide-sweeping
	// changes at that time.
	if len(t.path) == 1 && t.name == "root" {
		diags = diags.Append(fmt.Errorf(
			"root module cannot contain module named 'root'",
		))
		return diags
	}

	// Validate our configuration first.
	diags = diags.Append(t.config.Validate())

	// If we're the root, we do extra validation. This validation usually
	// requires the entire tree (since children don't have parent pointers).
	if len(t.path) == 0 {
		if err := t.validateProviderAlias(); err != nil {
			diags = diags.Append(err)
		}
	}

	// Get the child trees
	children := t.Children()

	// Validate all our children
	for _, c := range children {
		childDiags := c.Validate()
		diags = diags.Append(childDiags)
		if diags.HasErrors() {
			continue
		}
	}

	// Go over all the modules and verify that any parameters are valid
	// variables into the module in question.
	for _, m := range t.config.Modules {
		tree, ok := children[m.Name]
		if !ok {
			// This should never happen because Load watches us
			panic("module not found in children: " + m.Name)
		}

		// Build the variables that the module defines
		requiredMap := make(map[string]struct{})
		varMap := make(map[string]struct{})
		for _, v := range tree.config.Variables {
			varMap[v.Name] = struct{}{}

			if v.Required() {
				requiredMap[v.Name] = struct{}{}
			}
		}

		// Compare to the keys in our raw config for the module
		for k, _ := range m.RawConfig.Raw {
			if _, ok := varMap[k]; !ok {
				diags = diags.Append(fmt.Errorf(
					"module %q: %q is not a valid argument",
					m.Name, k,
				))
			}

			// Remove the required
			delete(requiredMap, k)
		}

		// If we have any required left over, they aren't set.
		for k, _ := range requiredMap {
			diags = diags.Append(fmt.Errorf(
				"module %q: missing required argument %q",
				m.Name, k,
			))
		}
	}

	// Go over all the variables used and make sure that any module
	// variables represent outputs properly.
	for source, vs := range t.config.InterpolatedVariables() {
		for _, v := range vs {
			mv, ok := v.(*config.ModuleVariable)
			if !ok {
				continue
			}

			tree, ok := children[mv.Name]
			if !ok {
				diags = diags.Append(fmt.Errorf(
					"%s: reference to undefined module %q",
					source, mv.Name,
				))
				continue
			}

			found := false
			for _, o := range tree.config.Outputs {
				if o.Name == mv.Field {
					found = true
					break
				}
			}
			if !found {
				diags = diags.Append(fmt.Errorf(
					"%s: %q is not a valid output for module %q",
					source, mv.Field, mv.Name,
				))
			}
		}
	}

	return diags
}

// versionedPathKey returns a path string with every levels full name, version
// and source encoded. This is to provide a unique key for our module storage,
// since submodules need to know which versions of their ancestor modules they
// are loaded from.
// For example, if module A has a subdirectory B, if module A's source or
// version is updated B's storage key must reflect this change in order for the
// correct version of B's source to be loaded.
func (t *Tree) versionedPathKey(m *Module) string {
	path := make([]string, len(t.path)+1)
	path[len(path)-1] = m.Name + ";" + m.Source
	// We're going to load these in order for easier reading and debugging, but
	// in practice they only need to be unique and consistent.

	p := t
	i := len(path) - 2
	for ; i >= 0; i-- {
		if p == nil {
			break
		}
		// we may have been loaded under a blank Tree, so always check for a name
		// too.
		if p.name == "" {
			break
		}
		seg := p.name
		if p.version != "" {
			seg += "#" + p.version
		}

		if p.source != "" {
			seg += ";" + p.source
		}

		path[i] = seg
		p = p.parent
	}

	key := strings.Join(path, "|")
	return key
}

// treeError is an error use by Tree.Validate to accumulates all
// validation errors.
type treeError struct {
	Name     []string
	Errs     []error
	Children []*treeError
}

func (e *treeError) Add(err error) {
	e.Errs = append(e.Errs, err)
}

func (e *treeError) AddChild(err *treeError) {
	e.Children = append(e.Children, err)
}

func (e *treeError) ErrOrNil() error {
	if len(e.Errs) > 0 || len(e.Children) > 0 {
		return e
	}
	return nil
}

func (e *treeError) Error() string {
	name := strings.Join(e.Name, ".")
	var out bytes.Buffer
	fmt.Fprintf(&out, "module %s: ", name)

	if len(e.Errs) == 1 {
		// single like error
		out.WriteString(e.Errs[0].Error())
	} else {
		// multi-line error
		for _, err := range e.Errs {
			fmt.Fprintf(&out, "\n    %s", err)
		}
	}

	if len(e.Children) > 0 {
		// start the next error on a new line
		out.WriteString("\n  ")
	}
	for _, child := range e.Children {
		out.WriteString(child.Error())
	}

	return out.String()
}
