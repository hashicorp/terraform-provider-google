// affectedtests determines, for a given GitHub PR, which acceptance tests it affects.
//
// Example usage: git diff HEAD~ > tmp.diff && go run affectedtests.go -diff tmp.diff
//
// It is also possible to get the diff from a PR: go run affectedtests.go -pr 2771
// However, this mode only reads the changed files from the PR and does not (currently)
// take into account new resources/tests that might have been added in this PR.
//
// This script currently only works for changes to resources.
// It is a TODO to make it work for changes to tests, data sources, and common utilities.
// It also currently does not pick up tests that use configs from other files.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

func main() {
	diff := flag.String("diff", "", "file containing git diff to use when determining changed files")
	pr := flag.Uint("pr", 0, "PR # to use to determine changed files")
	flag.Parse()
	if (*pr == 0 && *diff == "") || (*pr != 0 && *diff != "") {
		fmt.Println("Exactly one of -pr and -diff must be set")
		flag.Usage()
		os.Exit(1)
	}

	_, scriptPath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get current working directory")
	}
	tpgDir := scriptPath
	for !strings.HasPrefix(filepath.Base(tpgDir), "terraform-provider-") && tpgDir != "/" {
		tpgDir = filepath.Clean(tpgDir + "/..")
	}
	if tpgDir == "/" {
		log.Fatal("Script was run outside of google provider directory")
	}
	repo := strings.TrimPrefix(filepath.Base(tpgDir), "terraform-provider-")
	googleDir := tpgDir + "/" + repo

	providerFiles, err := readProviderFiles(googleDir)
	if err != nil {
		log.Fatal(err)
	}

	var diffVal string
	if *diff == "" {
		diffVal, err = getDiffFromPR(*pr, repo)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		d, err := ioutil.ReadFile(*diff)
		if err != nil {
			log.Fatal(err)
		}
		diffVal = string(d)
	}

	tests := map[string]struct{}{}
	for _, r := range getChangedResourcesFromDiff(diffVal, repo) {
		rn, err := getResourceName(r, googleDir, providerFiles)
		if err != nil {
			log.Fatal(err)
		}
		if rn == "" {
			log.Fatalf("Could not find resource represented by %s", r)
		}
		log.Printf("File %s matches resource %s", r, rn)
		ts, err := getTestsAffectedBy(rn, googleDir)
		if err != nil {
			log.Fatal(err)
		}
		for _, t := range ts {
			tests[t] = struct{}{}
		}
	}
	testnames := []string{}
	for tn := range tests {
		testnames = append(testnames, tn)
	}
	sort.Strings(testnames)
	for _, tn := range testnames {
		fmt.Println(tn)
	}
}

func readProviderFiles(googleDir string) ([]string, error) {
	pfs := []string{}
	dir, err := ioutil.ReadDir(googleDir)
	if err != nil {
		return nil, err
	}
	for _, f := range dir {
		if strings.HasPrefix(f.Name(), "provider") {
			p, err := ioutil.ReadFile(googleDir + "/" + f.Name())
			if err != nil {
				return nil, err
			}
			pfs = append(pfs, string(p))
		}
	}
	return pfs, nil
}

func getDiffFromPR(pr uint, repo string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://github.com/terraform-providers/terraform-provider-%s/pull/%d.diff", repo, pr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getChangedResourcesFromDiff(diff, repo string) []string {
	results := []string{}
	for _, l := range strings.Split(diff, "\n") {
		if strings.HasPrefix(l, "+++ b/") {
			log.Println("Found addition: " + l)
			fName := strings.TrimPrefix(l, "+++ b/"+repo+"/")
			if strings.HasPrefix(fName, "resource_") && !strings.HasSuffix(fName, "_test.go") {
				results = append(results, fName)
			}
		}
	}
	log.Printf("PR contains resource files %v", results)
	return results
}

func getResourceName(fName, googleDir string, providerFiles []string) (string, error) {
	resourceFile, err := parser.ParseFile(token.NewFileSet(), googleDir+"/"+fName, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}
	// Loop through all the top-level objects in the resource file.
	// One of them is the resource definition: something like resourceComputeInstance()
	for k := range resourceFile.Scope.Objects {
		// Matches the line in the provider file where the resource is defined,
		// e.g. "google_compute_instance":     resourceComputeInstance()
		re := regexp.MustCompile(`"(.*)":\s*` + k + `\(\)`)

		// Check all the provider files to see if they have a line that matches
		// that regexp. If so, return the resource name.
		for _, pf := range providerFiles {
			sm := re.FindStringSubmatch(pf)
			if len(sm) > 1 {
				log.Println("Full match is " + sm[0])
				return sm[1], nil
			}
		}
	}

	return "", nil
}

func getTestsAffectedBy(rn, googleDir string) ([]string, error) {
	lines, err := getLinesContainingResourceName(rn, googleDir)
	if err != nil {
		return nil, err
	}

	results := []string{}
	for _, line := range lines {
		fset := token.NewFileSet()
		p, err := parser.ParseFile(fset, line.file, nil, parser.AllErrors)
		if err != nil {
			return nil, err
		}

		// Find the top-level func containing this offset
		def := findFuncContainingOffset(line.offset, fset, p)
		if def == "" {
			// We couldn't find the place in the file that contains this offset, just skip and move on
			continue
		}

		// Go back through and find the test that calls the definition we just found
		results = append(results, findTestsCallingFunc(p, def)...)
	}
	return results, nil
}

func findFuncContainingOffset(offset int, fset *token.FileSet, p *ast.File) string {
	for k, sc := range p.Scope.Objects {
		d := sc.Decl.(ast.Node)
		if fset.Position(d.Pos()).Offset < offset && offset < fset.Position(d.End()).Offset {
			return k
		}
	}
	return ""
}

func findTestsCallingFunc(p *ast.File, funcName string) []string {
	results := []string{}
	for objName, sc := range p.Scope.Objects {
		if !strings.HasPrefix(objName, "Test") {
			continue
		}
		d, ok := sc.Decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// Starting at each Test, see if there's a path to the func we just found.
		ast.Inspect(d, func(n ast.Node) bool {
			if n, ok := n.(*ast.Ident); ok {
				if n.Name == funcName {
					results = append(results, objName)
				}
			}
			return true
		})
	}
	return results
}

type location struct {
	file   string
	offset int
}

func getLinesContainingResourceName(rn, googleDir string) ([]location, error) {
	results := []location{}
	resDef := regexp.MustCompile(fmt.Sprintf(`resource "%s"`, rn))
	dir, err := ioutil.ReadDir(googleDir)
	if err != nil {
		return nil, err
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		fPath := googleDir + "/" + f.Name()
		contents, err := ioutil.ReadFile(fPath)
		if err != nil {
			return nil, err
		}
		matches := resDef.FindAllIndex(contents, -1)
		for _, loc := range matches {
			// the full match is at contents[loc[0]:loc[1]], but we only need one value
			results = append(results, location{fPath, loc[0]})
		}
	}
	return results, nil
}
