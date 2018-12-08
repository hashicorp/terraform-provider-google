package openstack

import (
	"fmt"
	"log"
	"sort"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/snapshots"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceBlockStorageSnapshotV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockStorageSnapshotV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"volume_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"most_recent": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Computed values
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceBlockStorageSnapshotV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.blockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	listOpts := snapshots.ListOpts{
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		VolumeID: d.Get("volume_id").(string),
	}

	allPages, err := snapshots.List(client, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query snapshots: %s", err)
	}

	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve snapshots: %s", err)
	}

	if len(allSnapshots) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var snapshot snapshots.Snapshot
	if len(allSnapshots) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] Multiple results found and `most_recent` is set to: %t", recent)

		if recent {
			snapshot = dataSourceBlockStorageV2MostRecentSnapshot(allSnapshots)
		} else {
			log.Printf("[DEBUG] Multiple results found: %#v", allSnapshots)
			return fmt.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		snapshot = allSnapshots[0]
	}

	return dataSourceBlockStorageSnapshotV2Attributes(d, snapshot)
}

func dataSourceBlockStorageSnapshotV2Attributes(d *schema.ResourceData, snapshot snapshots.Snapshot) error {

	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("description", snapshot.Description)
	d.Set("size", snapshot.Size)
	d.Set("status", snapshot.Status)
	d.Set("volume_id", snapshot.VolumeID)

	if err := d.Set("metadata", snapshot.Metadata); err != nil {
		log.Printf("[DEBUG] Unable to set metadata for snapshot %s: %s", snapshot.ID, err)
	}

	return nil
}

func dataSourceBlockStorageV2MostRecentSnapshot(snapshots []snapshots.Snapshot) snapshots.Snapshot {
	sortedSnapshots := snapshots
	sort.Sort(blockStorageV2SnapshotSort(sortedSnapshots))
	return sortedSnapshots[len(sortedSnapshots)-1]
}
