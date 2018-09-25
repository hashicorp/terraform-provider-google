package random

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInteger() *schema.Resource {
	return &schema.Resource{
		Create: CreateInteger,
		Read:   RepopulateInteger,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: ImportInteger,
		},

		Schema: map[string]*schema.Schema{
			"keepers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"min": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"max": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"seed": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"result": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func CreateInteger(d *schema.ResourceData, meta interface{}) error {
	min := d.Get("min").(int)
	max := d.Get("max").(int)
	seed := d.Get("seed").(string)

	if max <= min {
		return fmt.Errorf("Minimum value needs to be smaller than maximum value")
	}
	rand := NewRand(seed)
	number := rand.Intn((max+1)-min) + min

	d.Set("result", number)
	d.SetId(strconv.Itoa(number))

	return nil
}

func RepopulateInteger(d *schema.ResourceData, _ interface{}) error {
	return nil
}

func ImportInteger(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 3 && len(parts) != 4 {
		return nil, fmt.Errorf("Invalid import usage: expecting {result},{min},{max} or {result},{min},{max},{seed}")
	}

	result, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"result\": {{err}}", err)
	}
	d.Set("result", result)

	min, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"min\": {{err}}", err)
	}
	d.Set("min", min)

	max, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"max\": {{err}}", err)
	}
	d.Set("max", max)

	if len(parts) == 4 {
		d.Set("seed", parts[3])
	}
	d.SetId(parts[0])

	return []*schema.ResourceData{d}, nil
}
