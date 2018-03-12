package configschema

import (
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var mapLabelNames = []string{"key"}

// DecoderSpec returns a hcldec.Spec that can be used to decode a HCL Body
// using the facilities in the hcldec package.
//
// The returned specification is guaranteed to return a value of the same type
// returned by method ImpliedType, but it may contain null or unknown values if
// any of the block attributes are defined as optional and/or computed
// respectively.
func (b *Block) DecoderSpec() hcldec.Spec {
	ret := hcldec.ObjectSpec{}
	if b == nil {
		return ret
	}

	for name, attrS := range b.Attributes {
		switch {
		case attrS.Computed && attrS.Optional:
			// In this special case we use an unknown value as a default
			// to get the intended behavior that the result is computed
			// unless it has been explicitly set in config.
			ret[name] = &hcldec.DefaultSpec{
				Primary: &hcldec.AttrSpec{
					Name: name,
					Type: attrS.Type,
				},
				Default: &hcldec.LiteralSpec{
					Value: cty.UnknownVal(attrS.Type),
				},
			}
		case attrS.Computed:
			ret[name] = &hcldec.LiteralSpec{
				Value: cty.UnknownVal(attrS.Type),
			}
		default:
			ret[name] = &hcldec.AttrSpec{
				Name:     name,
				Type:     attrS.Type,
				Required: attrS.Required,
			}
		}
	}

	for name, blockS := range b.BlockTypes {
		if _, exists := ret[name]; exists {
			// This indicates an invalid schema, since it's not valid to
			// define both an attribute and a block type of the same name.
			// However, we don't raise this here since it's checked by
			// InternalValidate.
			continue
		}

		childSpec := blockS.Block.DecoderSpec()

		switch blockS.Nesting {
		case NestingSingle:
			ret[name] = &hcldec.BlockSpec{
				TypeName: name,
				Nested:   childSpec,
				Required: blockS.MinItems == 1 && blockS.MaxItems >= 1,
			}
		case NestingList:
			ret[name] = &hcldec.BlockListSpec{
				TypeName: name,
				Nested:   childSpec,
				MinItems: blockS.MinItems,
				MaxItems: blockS.MaxItems,
			}
		case NestingSet:
			ret[name] = &hcldec.BlockSetSpec{
				TypeName: name,
				Nested:   childSpec,
				MinItems: blockS.MinItems,
				MaxItems: blockS.MaxItems,
			}
		case NestingMap:
			ret[name] = &hcldec.BlockMapSpec{
				TypeName:   name,
				Nested:     childSpec,
				LabelNames: mapLabelNames,
			}
		default:
			// Invalid nesting type is just ignored. It's checked by
			// InternalValidate.
			continue
		}
	}

	return ret
}
