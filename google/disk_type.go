package google

import (
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

// readDiskType finds the disk type with the given name.
func readDiskType(c *Config, zone *compute.Zone, project, name string) (*compute.DiskType, error) {
	diskType, err := c.clientCompute.DiskTypes.Get(project, zone.Name, name).Do()
	if err == nil && diskType != nil && diskType.SelfLink != "" {
		return diskType, nil
	} else {
		return nil, err
	}
}

// readRegionDiskType finds the disk type with the given name.
func readRegionDiskType(c *Config, region *compute.Region, project, name string) (*computeBeta.DiskType, error) {
	diskType, err := c.clientComputeBeta.RegionDiskTypes.Get(project, region.Name, name).Do()
	if err == nil && diskType != nil && diskType.SelfLink != "" {
		return diskType, nil
	} else {
		return nil, err
	}
}
