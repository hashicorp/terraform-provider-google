package random

import (
	"crypto/rand"
	"math/big"
	"sort"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceString() *schema.Resource {
	return &schema.Resource{
		Create:        CreateString,
		Read:          ReadString,
		Delete:        schema.RemoveFromState,
		MigrateState:  resourceRandomStringMigrateState,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"keepers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"length": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"special": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"upper": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"lower": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"number": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"min_numeric": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},

			"min_upper": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},

			"min_lower": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},

			"min_special": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},

			"override_special": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func CreateString(d *schema.ResourceData, meta interface{}) error {
	const numChars = "0123456789"
	const lowerChars = "abcdefghijklmnopqrstuvwxyz"
	const upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var specialChars = "!@#$%&*()-_=+[]{}<>:?"

	length := d.Get("length").(int)
	upper := d.Get("upper").(bool)
	minUpper := d.Get("min_upper").(int)
	lower := d.Get("lower").(bool)
	minLower := d.Get("min_lower").(int)
	number := d.Get("number").(bool)
	minNumeric := d.Get("min_numeric").(int)
	special := d.Get("special").(bool)
	minSpecial := d.Get("min_special").(int)
	overrideSpecial := d.Get("override_special").(string)

	if overrideSpecial != "" {
		specialChars = overrideSpecial
	}

	var chars = string("")
	if upper {
		chars += upperChars
	}
	if lower {
		chars += lowerChars
	}
	if number {
		chars += numChars
	}
	if special {
		chars += specialChars
	}

	minMapping := map[string]int{
		numChars:     minNumeric,
		lowerChars:   minLower,
		upperChars:   minUpper,
		specialChars: minSpecial,
	}
	var result = make([]byte, 0, length)
	for k, v := range minMapping {
		s, err := generateRandomBytes(&k, v)
		if err != nil {
			return errwrap.Wrapf("error generating random bytes: {{err}}", err)
		}
		result = append(result, s...)
	}
	s, err := generateRandomBytes(&chars, length-len(result))
	if err != nil {
		return errwrap.Wrapf("error generating random bytes: {{err}}", err)
	}
	result = append(result, s...)
	order := make([]byte, len(result))
	if _, err := rand.Read(order); err != nil {
		return errwrap.Wrapf("error generating random bytes: {{err}}", err)
	}
	sort.Slice(result, func(i, j int) bool {
		return order[i] < order[j]
	})

	d.Set("result", string(result))
	d.SetId("none")
	return nil
}

func generateRandomBytes(charSet *string, length int) ([]byte, error) {
	bytes := make([]byte, length)
	setLen := big.NewInt(int64(len(*charSet)))
	for i := range bytes {
		idx, err := rand.Int(rand.Reader, setLen)
		if err != nil {
			return nil, err
		}
		bytes[i] = (*charSet)[idx.Int64()]
	}
	return bytes, nil
}

func ReadString(d *schema.ResourceData, meta interface{}) error {
	return nil
}
