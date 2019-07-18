package google

// readDiskType finds the disk type with the given name.
func readDiskType(c *Config, d TerraformResourceData, name string) (*ZonalFieldValue, error) {
	return parseZonalFieldValue("diskTypes", name, "project", "zone", d, c, false)
}

// readRegionDiskType finds the disk type with the given name.
func readRegionDiskType(c *Config, d TerraformResourceData, name string) (*RegionalFieldValue, error) {
	return parseRegionalFieldValue("diskTypes", name, "project", "region", "zone", d, c, false)
}
