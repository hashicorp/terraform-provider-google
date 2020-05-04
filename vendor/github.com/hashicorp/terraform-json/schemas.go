package tfjson

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

// ProviderSchemasFormatVersion is the version of the JSON provider
// schema format that is supported by this package.
const ProviderSchemasFormatVersion = "0.1"

// ProviderSchemas represents the schemas of all providers and
// resources in use by the configuration.
type ProviderSchemas struct {
	// The version of the plan format. This should always match the
	// ProviderSchemasFormatVersion constant in this package, or else
	// an unmarshal will be unstable.
	FormatVersion string `json:"format_version,omitempty"`

	// The schemas for the providers in this configuration, indexed by
	// provider type. Aliases are not included, and multiple instances
	// of a provider in configuration will be represented by a single
	// provider here.
	Schemas map[string]*ProviderSchema `json:"provider_schemas,omitempty"`
}

// Validate checks to ensure that ProviderSchemas is present, and the
// version matches the version supported by this library.
func (p *ProviderSchemas) Validate() error {
	if p == nil {
		return errors.New("provider schema data is nil")
	}

	if p.FormatVersion == "" {
		return errors.New("unexpected provider schema data, format version is missing")
	}

	if ProviderSchemasFormatVersion != p.FormatVersion {
		return fmt.Errorf("unsupported provider schema data format version: expected %q, got %q", PlanFormatVersion, p.FormatVersion)
	}

	return nil
}

func (p *ProviderSchemas) UnmarshalJSON(b []byte) error {
	type rawSchemas ProviderSchemas
	var schemas rawSchemas

	err := json.Unmarshal(b, &schemas)
	if err != nil {
		return err
	}

	*p = *(*ProviderSchemas)(&schemas)

	return p.Validate()
}

// ProviderSchema is the JSON representation of the schema of an
// entire provider, including the provider configuration and any
// resources and data sources included with the provider.
type ProviderSchema struct {
	// The schema for the provider's configuration.
	ConfigSchema *Schema `json:"provider,omitempty"`

	// The schemas for any resources in this provider.
	ResourceSchemas map[string]*Schema `json:"resource_schemas,omitempty"`

	// The schemas for any data sources in this provider.
	DataSourceSchemas map[string]*Schema `json:"data_source_schemas,omitempty"`
}

// Schema is the JSON representation of a particular schema
// (provider configuration, resources, data sources).
type Schema struct {
	// The version of the particular resource schema.
	Version uint64 `json:"version"`

	// The root-level block of configuration values.
	Block *SchemaBlock `json:"block,omitempty"`
}

// SchemaBlock represents a nested block within a particular schema.
type SchemaBlock struct {
	// The attributes defined at the particular level of this block.
	Attributes map[string]*SchemaAttribute `json:"attributes,omitempty"`

	// Any nested blocks within this particular block.
	NestedBlocks map[string]*SchemaBlockType `json:"block_types,omitempty"`
}

// SchemaNestingMode is the nesting mode for a particular nested
// schema block.
type SchemaNestingMode string

const (
	// SchemaNestingModeSingle denotes single block nesting mode, which
	// allows a single block of this specific type only in
	// configuration. This is generally the same as list or set types
	// with a single-element constraint.
	SchemaNestingModeSingle SchemaNestingMode = "single"

	// SchemaNestingModeList denotes list block nesting mode, which
	// allows an ordered list of blocks where duplicates are allowed.
	SchemaNestingModeList SchemaNestingMode = "list"

	// SchemaNestingModeSet denotes set block nesting mode, which
	// allows an unordered list of blocks where duplicates are
	// generally not allowed. What is considered a duplicate is up to
	// the rules of the set itself, which may or may not cover all
	// fields in the block.
	SchemaNestingModeSet SchemaNestingMode = "set"

	// SchemaNestingModeMap denotes map block nesting mode. This
	// creates a map of all declared blocks of the block type within
	// the parent, keying them on the label supplied in the block
	// declaration. This allows for blocks to be declared in the same
	// style as resources.
	SchemaNestingModeMap SchemaNestingMode = "map"
)

// SchemaBlockType describes a nested block within a schema.
type SchemaBlockType struct {
	// The nesting mode for this block.
	NestingMode SchemaNestingMode `json:"nesting_mode,omitempty"`

	// The block data for this block type, including attributes and
	// subsequent nested blocks.
	Block *SchemaBlock `json:"block,omitempty"`

	// The lower limit on items that can be declared of this block
	// type.
	MinItems uint64 `json:"min_items,omitempty"`

	// The upper limit on items that can be declared of this block
	// type.
	MaxItems uint64 `json:"max_items,omitempty"`
}

// SchemaAttribute describes an attribute within a schema block.
type SchemaAttribute struct {
	// The attribute type.
	AttributeType cty.Type `json:"type,omitempty"`

	// The description field for this attribute.
	Description string `json:"description,omitempty"`

	// If true, this attribute is required - it has to be entered in
	// configuration.
	Required bool `json:"required,omitempty"`

	// If true, this attribute is optional - it does not need to be
	// entered in configuration.
	Optional bool `json:"optional,omitempty"`

	// If true, this attribute is computed - it can be set by the
	// provider. It may also be set by configuration if Optional is
	// true.
	Computed bool `json:"computed,omitempty"`

	// If true, this attribute is sensitive and will not be displayed
	// in logs. Future versions of Terraform may encrypt or otherwise
	// treat these values with greater care than non-sensitive fields.
	Sensitive bool `json:"sensitive,omitempty"`
}
