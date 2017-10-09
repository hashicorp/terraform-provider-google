package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
)

func dataSourceGoogleComputeSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeSnapshotRead,

		Schema: map[string]*schema.Schema{
			//"filter": dataSourceFiltersSchema(),
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_disk_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"snapshot_encryption_key_sha256": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_disk_encryption_key_sha256": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_size_gb": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"storage_bytes": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"storage_bytes_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"licenses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleComputeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	labels := expandLabels(d)
	log.Printf("[DEBUG] Labels %s", labels)

	if len(labels) > 0 {
		filter := ""
		log.Printf("[DEBUG] Labels length : %d", len(labels))
		for k, v := range labels {
			log.Printf("[DEBUG] Label key : '%s', value : '%s'", k, v)
			filter = fmt.Sprintf("%s(labels.%s eq %s)", filter, k, v)
		}
		log.Printf("[DEBUG] Labels filter : %s", filter)
		snapshotList, err := config.clientCompute.Snapshots.List(project).Filter(filter).Do()
		if err != nil {
			return fmt.Errorf("error while listing snapshots with filter %s: %s", filter, err)
		}
		log.Printf("[DEBUG] SnapshotList length : %d", len(snapshotList.Items))

		if len(snapshotList.Items) > 1 {
			return fmt.Errorf("too many snapshots found with these labels")
		} else if len(snapshotList.Items) == 0 {
			return fmt.Errorf("no snapshot found with these labels")
		} else if v, ok := d.GetOk("name"); ok {
			foundSnapName := snapshotList.Items[0].Name
			if v != foundSnapName {
				return fmt.Errorf("error, name different from snapshot name found with these labels : %s != %s", v, foundSnapName)
			}

		}
		d.Set("name", snapshotList.Items[0].Name)
	}

	snapshot, err := config.clientCompute.Snapshots.Get(
		project, d.Get("name").(string)).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't exist anymore

			return fmt.Errorf("Snapshot Not Found : %s", d.Get("name"))
		}

		return fmt.Errorf("Error reading snapshot: %s", err)
	}
	d.Set("self_link", snapshot.SelfLink)
	d.Set("description", snapshot.Description)
	d.Set("source_disk_link", snapshot.SourceDisk)

	if snapshot.SnapshotEncryptionKey != nil && snapshot.SnapshotEncryptionKey.Sha256 != "" {
		d.Set("snapshot_encryption_key_sha256", snapshot.SnapshotEncryptionKey.Sha256)
	}

	if snapshot.SourceDiskEncryptionKey != nil && snapshot.SourceDiskEncryptionKey.Sha256 != "" {
		d.Set("source_disk_encryption_key_sha256", snapshot.SourceDiskEncryptionKey.Sha256)
	}

	d.Set("source_disk_id", snapshot.SourceDiskId)
	d.Set("status", snapshot.Status)
	d.Set("storage_size", snapshot.StorageBytes)
	d.Set("storage_size_status", snapshot.StorageBytesStatus)
	d.Set("disk_size", snapshot.DiskSizeGb)
	d.Set("labels", snapshot.Labels)

	d.SetId(snapshot.Name)
	return nil
}
