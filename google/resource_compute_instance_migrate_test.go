package google

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestComputeInstanceMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"v0.4.2 and earlier": {
			StateVersion: 0,
			Attributes: map[string]string{
				"metadata.#":           "2",
				"metadata.0.foo":       "bar",
				"metadata.1.baz":       "qux",
				"metadata.2.with.dots": "should.work",
			},
			Expected: map[string]string{
				"metadata.foo":       "bar",
				"metadata.baz":       "qux",
				"metadata.with.dots": "should.work",
			},
		},
		"change scope from list to set": {
			StateVersion: 1,
			Attributes: map[string]string{
				"service_account.#":          "1",
				"service_account.0.email":    "xxxxxx-compute@developer.gserviceaccount.com",
				"service_account.0.scopes.#": "4",
				"service_account.0.scopes.0": "https://www.googleapis.com/auth/compute",
				"service_account.0.scopes.1": "https://www.googleapis.com/auth/datastore",
				"service_account.0.scopes.2": "https://www.googleapis.com/auth/devstorage.full_control",
				"service_account.0.scopes.3": "https://www.googleapis.com/auth/logging.write",
			},
			Expected: map[string]string{
				"service_account.#":                   "1",
				"service_account.0.email":             "xxxxxx-compute@developer.gserviceaccount.com",
				"service_account.0.scopes.#":          "4",
				"service_account.0.scopes.1693978638": "https://www.googleapis.com/auth/devstorage.full_control",
				"service_account.0.scopes.172152165":  "https://www.googleapis.com/auth/logging.write",
				"service_account.0.scopes.299962681":  "https://www.googleapis.com/auth/compute",
				"service_account.0.scopes.3435931483": "https://www.googleapis.com/auth/datastore",
			},
		},
		"add new create_timeout attribute": {
			StateVersion: 2,
			Attributes:   map[string]string{},
			Expected: map[string]string{
				"create_timeout": "4",
			},
		},
		// "replace disk with boot disk": {
		// 	StateVersion: 3,
		// 	Attributes: map[string]string{
		// 		"disk.#":                            "1",
		// 		"disk.0.disk":                       "disk-1",
		// 		"disk.0.type":                       "pd-ssd",
		// 		"disk.0.auto_delete":                "false",
		// 		"disk.0.size":                       "12",
		// 		"disk.0.device_name":                "device-name",
		// 		"disk.0.disk_encryption_key_raw":    "encrypt-key",
		// 		"disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		// 	},
		// 	Expected: map[string]string{
		// 		"boot_disk.#":                            "1",
		// 		"boot_disk.0.auto_delete":                "false",
		// 		"boot_disk.0.device_name":                "device-name",
		// 		"boot_disk.0.disk_encryption_key_raw":    "encrypt-key",
		// 		"boot_disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		// 		"boot_disk.0.source":                     "disk-1",
		// 	},
		// },
		// "replace disk with attached disk": {
		// 	StateVersion: 3,
		// 	Attributes: map[string]string{
		// 		"boot_disk.#":                       "1",
		// 		"disk.#":                            "1",
		// 		"disk.0.disk":                       "path/to/disk",
		// 		"disk.0.device_name":                "device-name",
		// 		"disk.0.disk_encryption_key_raw":    "encrypt-key",
		// 		"disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		// 	},
		// 	Expected: map[string]string{
		// 		"boot_disk.#":                                "1",
		// 		"attached_disk.#":                            "1",
		// 		"attached_disk.0.source":                     "path/to/disk",
		// 		"attached_disk.0.device_name":                "device-name",
		// 		"attached_disk.0.disk_encryption_key_raw":    "encrypt-key",
		// 		"attached_disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		// 	},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "i-abc123",
			Attributes: tc.Attributes,
		}
		is, err := resourceComputeInstanceMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, tc.Expected[k], k, is.Attributes[k], is.Attributes)
			}
		}

		for k, v := range is.Attributes {
			if tc.Expected[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, tc.Expected[k], k, is.Attributes[k], is.Attributes)
			}
		}
	}
}

func TestComputeInstanceMigrateState_empty(t *testing.T) {
	var is *terraform.InstanceState
	var meta interface{}

	// should handle nil
	is, err := resourceComputeInstanceMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if is != nil {
		t.Fatalf("expected nil instancestate, got: %#v", is)
	}

	// should handle non-nil but empty
	is = &terraform.InstanceState{}
	is, err = resourceComputeInstanceMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
}
