// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestComputeInstanceTemplate_reorderDisks(t *testing.T) {
	t.Parallel()

	cBoot := map[string]interface{}{
		"source": "boot-source",
	}
	cFallThrough := map[string]interface{}{
		"auto_delete": true,
	}
	cDeviceName := map[string]interface{}{
		"device_name": "disk-1",
	}
	cScratch := map[string]interface{}{
		"type": "SCRATCH",
	}
	cSource := map[string]interface{}{
		"source": "disk-source",
	}
	cScratchNvme := map[string]interface{}{
		"type":      "SCRATCH",
		"interface": "NVME",
	}

	aBoot := map[string]interface{}{
		"source": "boot-source",
		"boot":   true,
	}
	aScratchNvme := map[string]interface{}{
		"device_name": "scratch-1",
		"type":        "SCRATCH",
		"interface":   "NVME",
	}
	aSource := map[string]interface{}{
		"device_name": "disk-2",
		"source":      "disk-source",
	}
	aScratchScsi := map[string]interface{}{
		"device_name": "scratch-2",
		"type":        "SCRATCH",
		"interface":   "SCSI",
	}
	aFallThrough := map[string]interface{}{
		"device_name": "disk-3",
		"auto_delete": true,
		"source":      "fake-source",
	}
	aFallThrough2 := map[string]interface{}{
		"device_name": "disk-4",
		"auto_delete": true,
		"source":      "fake-source",
	}
	aDeviceName := map[string]interface{}{
		"device_name": "disk-1",
		"auto_delete": true,
		"source":      "fake-source-2",
	}
	aNoMatch := map[string]interface{}{
		"device_name": "disk-2",
		"source":      "disk-source-doesn't-match",
	}

	cases := map[string]struct {
		ConfigDisks    []interface{}
		ApiDisks       []map[string]interface{}
		ExpectedResult []map[string]interface{}
	}{
		"all disks represented": {
			ApiDisks: []map[string]interface{}{
				aBoot, aScratchNvme, aSource, aScratchScsi, aFallThrough, aDeviceName,
			},
			ConfigDisks: []interface{}{
				cBoot, cFallThrough, cDeviceName, cScratch, cSource, cScratchNvme,
			},
			ExpectedResult: []map[string]interface{}{
				aBoot, aFallThrough, aDeviceName, aScratchScsi, aSource, aScratchNvme,
			},
		},
		"one non-match": {
			ApiDisks: []map[string]interface{}{
				aBoot, aNoMatch, aScratchNvme, aScratchScsi, aFallThrough, aDeviceName,
			},
			ConfigDisks: []interface{}{
				cBoot, cFallThrough, cDeviceName, cScratch, cSource, cScratchNvme,
			},
			ExpectedResult: []map[string]interface{}{
				aBoot, aFallThrough, aDeviceName, aScratchScsi, aScratchNvme, aNoMatch,
			},
		},
		"two fallthroughs": {
			ApiDisks: []map[string]interface{}{
				aBoot, aScratchNvme, aFallThrough, aSource, aScratchScsi, aFallThrough2, aDeviceName,
			},
			ConfigDisks: []interface{}{
				cBoot, cFallThrough, cDeviceName, cScratch, cFallThrough, cSource, cScratchNvme,
			},
			ExpectedResult: []map[string]interface{}{
				aBoot, aFallThrough, aDeviceName, aScratchScsi, aFallThrough2, aSource, aScratchNvme,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Disks read using d.Get will always have values for all keys, so set those values
			for _, disk := range tc.ConfigDisks {
				d := disk.(map[string]interface{})
				for _, k := range []string{"auto_delete", "boot"} {
					if _, ok := d[k]; !ok {
						d[k] = false
					}
				}
				for _, k := range []string{"device_name", "disk_name", "interface", "mode", "source", "type"} {
					if _, ok := d[k]; !ok {
						d[k] = ""
					}
				}
			}

			// flattened disks always set auto_delete, boot, device_name, interface, mode, source, and type
			for _, d := range tc.ApiDisks {
				for _, k := range []string{"auto_delete", "boot"} {
					if _, ok := d[k]; !ok {
						d[k] = false
					}
				}

				for _, k := range []string{"device_name", "interface", "mode", "source"} {
					if _, ok := d[k]; !ok {
						d[k] = ""
					}
				}
				if _, ok := d["type"]; !ok {
					d["type"] = "PERSISTENT"
				}
			}

			result := reorderDisks(tc.ConfigDisks, tc.ApiDisks)
			if !reflect.DeepEqual(tc.ExpectedResult, result) {
				t.Errorf("reordering did not match\nExpected: %+v\nActual: %+v", tc.ExpectedResult, result)
			}
		})
	}
}

func TestComputeInstanceTemplate_scratchDiskSizeCustomizeDiff(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Typee       string // misspelled on purpose, type is a special symbol
		DiskType    string
		DiskSize    int
		Interfacee  string
		ExpectError bool
	}{
		"scratch disk correct size 1": {
			Typee:       "SCRATCH",
			DiskType:    "local-ssd",
			DiskSize:    375,
			Interfacee:  "NVME",
			ExpectError: false,
		},
		"scratch disk correct size 2": {
			Typee:       "SCRATCH",
			DiskType:    "local-ssd",
			DiskSize:    3000,
			Interfacee:  "NVME",
			ExpectError: false,
		},
		"scratch disk incorrect size": {
			Typee:       "SCRATCH",
			DiskType:    "local-ssd",
			DiskSize:    300,
			Interfacee:  "NVME",
			ExpectError: true,
		},
		"scratch disk incorrect interface": {
			Typee:       "SCRATCH",
			DiskType:    "local-ssd",
			DiskSize:    3000,
			Interfacee:  "SCSI",
			ExpectError: true,
		},
		"non-scratch disk": {
			Typee:       "PERSISTENT",
			DiskType:    "",
			DiskSize:    300,
			Interfacee:  "NVME",
			ExpectError: false,
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			After: map[string]interface{}{
				"disk.#":              1,
				"disk.0.type":         tc.Typee,
				"disk.0.disk_type":    tc.DiskType,
				"disk.0.disk_size_gb": tc.DiskSize,
				"disk.0.interface":    tc.Interfacee,
			},
		}
		err := resourceComputeInstanceTemplateScratchDiskCustomizeDiffFunc(d)
		if tc.ExpectError && err == nil {
			t.Errorf("%s failed, expected error but was none", tn)
		}
		if !tc.ExpectError && err != nil {
			t.Errorf("%s failed, found unexpected error: %s", tn, err)
		}
	}
}
