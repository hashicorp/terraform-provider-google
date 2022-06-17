// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	osconfig "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/osconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccOsConfigOsPolicyAssignment_FixedOsPolicyAssignment(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  getTestProjectFromEnv(),
		"zone":          getTestZoneFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOsConfigOsPolicyAssignmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOsConfigOsPolicyAssignment_FixedOsPolicyAssignment(context),
			},
			{
				ResourceName:            "google_os_config_os_policy_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rollout.0.min_wait_duration"},
			},
			{
				Config: testAccOsConfigOsPolicyAssignment_FixedOsPolicyAssignmentUpdate0(context),
			},
			{
				ResourceName:            "google_os_config_os_policy_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rollout.0.min_wait_duration"},
			},
			{
				Config: testAccOsConfigOsPolicyAssignment_FixedOsPolicyAssignmentUpdate1(context),
			},
			{
				ResourceName:            "google_os_config_os_policy_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rollout.0.min_wait_duration"},
			},
		},
	})
}
func TestAccOsConfigOsPolicyAssignment_PercentOsPolicyAssignment(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  getTestProjectFromEnv(),
		"zone":          getTestZoneFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOsConfigOsPolicyAssignmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOsConfigOsPolicyAssignment_PercentOsPolicyAssignment(context),
			},
			{
				ResourceName:            "google_os_config_os_policy_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rollout.0.min_wait_duration"},
			},
			{
				Config: testAccOsConfigOsPolicyAssignment_PercentOsPolicyAssignmentUpdate0(context),
			},
			{
				ResourceName:            "google_os_config_os_policy_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rollout.0.min_wait_duration"},
			},
		},
	})
}

func testAccOsConfigOsPolicyAssignment_FixedOsPolicyAssignment(context map[string]interface{}) string {
	return Nprintf(`
resource "google_os_config_os_policy_assignment" "primary" {
  instance_filter {
    all = false

    exclusion_labels {
      labels = {
        label-two = "value-two"
      }
    }

    inclusion_labels {
      labels = {
        label-one = "value-one"
      }
    }

    inventories {
      os_short_name = "centos"
      os_version    = "8.*"
    }
  }

  location = "%{zone}"
  name     = "tf-test-assignment%{random_suffix}"

  os_policies {
    id   = "policy"
    mode = "VALIDATION"

    resource_groups {
      resources {
        id = "apt"

        pkg {
          desired_state = "INSTALLED"

          apt {
            name = "bazel"
          }
        }
      }

      resources {
        id = "deb1"

        pkg {
          desired_state = "INSTALLED"

          deb {
            source {
              local_path = "$HOME/package.deb"
            }
          }
        }
      }

      resources {
        id = "deb2"

        pkg {
          desired_state = "INSTALLED"

          deb {
            source {
              allow_insecure = true

              remote {
                uri             = "ftp.us.debian.org/debian/package.deb"
                sha256_checksum = "3bbfd1043cd7afdb78cf9afec36c0c5370d2fea98166537b4e67f3816f256025"
              }
            }

            pull_deps = true
          }
        }
      }

      resources {
        id = "deb3"

        pkg {
          desired_state = "INSTALLED"

          deb {
            source {
              gcs {
                bucket     = "test-bucket"
                object     = "test-object"
                generation = 1
              }
            }

            pull_deps = true
          }
        }
      }

      resources {
        id = "yum"

        pkg {
          desired_state = "INSTALLED"

          yum {
            name = "gstreamer-plugins-base-devel.x86_64"
          }
        }
      }

      resources {
        id = "zypper"

        pkg {
          desired_state = "INSTALLED"

          zypper {
            name = "gcc"
          }
        }
      }

      resources {
        id = "rpm1"

        pkg {
          desired_state = "INSTALLED"

          rpm {
            source {
              local_path = "$HOME/package.rpm"
            }

            pull_deps = true
          }
        }
      }

      resources {
        id = "rpm2"

        pkg {
          desired_state = "INSTALLED"

          rpm {
            source {
              allow_insecure = true

              remote {
                uri             = "https://mirror.jaleco.com/centos/8.3.2011/BaseOS/x86_64/os/Packages/efi-filesystem-3-2.el8.noarch.rpm"
                sha256_checksum = "3bbfd1043cd7afdb78cf9afec36c0c5370d2fea98166537b4e67f3816f256025"
              }
            }
          }
        }
      }

      resources {
        id = "rpm3"

        pkg {
          desired_state = "INSTALLED"

          rpm {
            source {
              gcs {
                bucket     = "test-bucket"
                object     = "test-object"
                generation = 1
              }
            }
          }
        }
      }

      inventory_filters {
        os_short_name = "centos"
        os_version    = "8.*"
      }
    }

    resource_groups {
      resources {
        id = "apt-to-deb"

        pkg {
          desired_state = "INSTALLED"

          apt {
            name = "bazel"
          }
        }
      }

      resources {
        id = "deb-local-path-to-gcs"

        pkg {
          desired_state = "INSTALLED"

          deb {
            source {
              local_path = "$HOME/package.deb"
            }
          }
        }
      }

      resources {
        id = "googet"

        pkg {
          desired_state = "INSTALLED"

          googet {
            name = "gcc"
          }
        }
      }

      resources {
        id = "msi1"

        pkg {
          desired_state = "INSTALLED"

          msi {
            source {
              local_path = "$HOME/package.msi"
            }

            properties = ["REBOOT=ReallySuppress"]
          }
        }
      }

      resources {
        id = "msi2"

        pkg {
          desired_state = "INSTALLED"

          msi {
            source {
              allow_insecure = true

              remote {
                uri             = "https://remote.uri.com/package.msi"
                sha256_checksum = "3bbfd1043cd7afdb78cf9afec36c0c5370d2fea98166537b4e67f3816f256025"
              }
            }
          }
        }
      }

      resources {
        id = "msi3"

        pkg {
          desired_state = "INSTALLED"

          msi {
            source {
              gcs {
                bucket     = "test-bucket"
                object     = "test-object"
                generation = 1
              }
            }
          }
        }
      }
    }

    allow_no_resource_group_match = false
    description                   = "A test os policy"
  }

  rollout {
    disruption_budget {
      fixed = 1
    }

    min_wait_duration = "3.5s"
  }

  description = "A test os policy assignment"
  project     = "%{project_name}"
}


`, context)
}

func testAccOsConfigOsPolicyAssignment_FixedOsPolicyAssignmentUpdate0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_os_config_os_policy_assignment" "primary" {
  instance_filter {
    all = false

    inventories {
      os_short_name = ""
      os_version    = "9.*"
    }
  }

  location = "%{zone}"
  name     = "tf-test-assignment%{random_suffix}"

  os_policies {
    id   = "policy"
    mode = "ENFORCEMENT"

    resource_groups {
      resources {
        id = "apt"

        pkg {
          desired_state = "INSTALLED"

          apt {
            name = "firefox"
          }
        }
      }

      resources {
        id = "new-deb1"

        pkg {
          desired_state = "REMOVED"

          deb {
            source {
              local_path = "$HOME/new-package.deb"
            }
          }
        }
      }

      resources {
        id = "new-deb2"

        pkg {
          desired_state = "REMOVED"

          deb {
            source {
              allow_insecure = false

              remote {
                uri             = "ftp.us.debian.org/debian/new-package.deb"
                sha256_checksum = "9f8e5818ccb47024d01000db713c0a333679b64678ff5fe2d9bea0a23014dd54"
              }
            }

            pull_deps = false
          }
        }
      }

      resources {
        id = "new-yum"

        pkg {
          desired_state = "REMOVED"

          yum {
            name = "vlc.x86_64"
          }
        }
      }

      resources {
        id = "new-zypper"

        pkg {
          desired_state = "REMOVED"

          zypper {
            name = "ModemManager"
          }
        }
      }

      resources {
        id = "new-rpm1"

        pkg {
          desired_state = "REMOVED"

          rpm {
            source {
              local_path = "$HOME/new-package.rpm"
            }

            pull_deps = false
          }
        }
      }

      resources {
        id = "new-rpm2"

        pkg {
          desired_state = "REMOVED"

          rpm {
            source {
              allow_insecure = false

              remote {
                uri             = "https://mirror.jaleco.com/centos/8.3.2011/BaseOS/x86_64/os/Packages/NetworkManager-adsl-1.26.0-12.el8_3.x86_64.rpm"
                sha256_checksum = "9f8e5818ccb47024d01000db713c0a333679b64678ff5fe2d9bea0a23014dd54"
              }
            }
          }
        }
      }

      resources {
        id = "new-rpm3"

        pkg {
          desired_state = "REMOVED"

          rpm {
            source {
              gcs {
                bucket     = "new-test-bucket"
                object     = "new-test-object"
                generation = 2
              }
            }
          }
        }
      }

      inventory_filters {
        os_short_name = ""
        os_version    = "9.*"
      }
    }

    resource_groups {
      resources {
        id = "apt-to-deb"

        pkg {
          desired_state = "REMOVED"

          deb {
            source {
              local_path = "$HOME/new-package.deb"
            }
          }
        }
      }

      resources {
        id = "deb-local-path-to-gcs"

        pkg {
          desired_state = "REMOVED"

          deb {
            source {
              gcs {
                bucket     = "new-test-bucket"
                object     = "new-test-object"
                generation = 2
              }
            }
          }
        }
      }

      resources {
        id = "new-googet"

        pkg {
          desired_state = "REMOVED"

          googet {
            name = "julia"
          }
        }
      }

      resources {
        id = "new-msi1"

        pkg {
          desired_state = "REMOVED"

          msi {
            source {
              local_path = "$HOME/new-package.msi"
            }

            properties = ["ACTION=INSTALL"]
          }
        }
      }

      resources {
        id = "new-msi2"

        pkg {
          desired_state = "REMOVED"

          msi {
            source {
              allow_insecure = false

              remote {
                uri             = "https://remote.uri.com/new-package.msi"
                sha256_checksum = "9f8e5818ccb47024d01000db713c0a333679b64678ff5fe2d9bea0a23014dd54"
              }
            }
          }
        }
      }

      resources {
        id = "new-msi3"

        pkg {
          desired_state = "REMOVED"

          msi {
            source {
              gcs {
                bucket     = "new-test-bucket"
                object     = "new-test-object"
                generation = 2
              }
            }
          }
        }
      }
    }

    allow_no_resource_group_match = true
    description                   = "An updated test os policy"
  }

  rollout {
    disruption_budget {
      fixed = 2
    }

    min_wait_duration = "7.5s"
  }

  description = "An updated test os policy assignment"
  project     = "%{project_name}"
}


`, context)
}

func testAccOsConfigOsPolicyAssignment_FixedOsPolicyAssignmentUpdate1(context map[string]interface{}) string {
	return Nprintf(`
resource "google_os_config_os_policy_assignment" "primary" {
  instance_filter {
    all = true
  }

  location = "%{zone}"
  name     = "tf-test-assignment%{random_suffix}"

  os_policies {
    id   = "policy"
    mode = "VALIDATION"

    resource_groups {
      resources {
        id = "apt-to-yum"

        repository {
          apt {
            archive_type = "DEB"
            components   = ["doc"]
            distribution = "debian"
            uri          = "https://atl.mirrors.clouvider.net/debian"
            gpg_key      = ".gnupg/pubring.kbx"
          }
        }
      }

      resources {
        id = "yum"

        repository {
          yum {
            base_url     = "http://centos.s.uw.edu/centos/"
            id           = "yum"
            display_name = "yum"
            gpg_keys     = ["RPM-GPG-KEY-CentOS-7"]
          }
        }
      }

      resources {
        id = "zypper"

        repository {
          zypper {
            base_url     = "http://mirror.dal10.us.leaseweb.net/opensuse"
            id           = "zypper"
            display_name = "zypper"
            gpg_keys     = ["sample-key-uri"]
          }
        }
      }

      resources {
        id = "goo"

        repository {
          goo {
            name = "goo"
            url  = "https://foo.com/googet/bar"
          }
        }
      }

      resources {
        id = "exec1"

        exec {
          validate {
            interpreter = "SHELL"
            args        = ["arg1"]

            file {
              local_path = "$HOME/script.sh"
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter = "SHELL"
            args        = ["arg1"]

            file {
              allow_insecure = true

              remote {
                uri             = "https://www.example.com/script.sh"
                sha256_checksum = "c7938fed83afdccbb0e86a2a2e4cad7d5035012ca3214b4a61268393635c3063"
              }
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "exec2"

        exec {
          validate {
            interpreter = "SHELL"
            args        = ["arg1"]

            file {
              allow_insecure = true

              remote {
                uri             = "https://www.example.com/script.sh"
                sha256_checksum = "c7938fed83afdccbb0e86a2a2e4cad7d5035012ca3214b4a61268393635c3063"
              }
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter = "SHELL"
            args        = ["arg1"]

            file {
              local_path = "$HOME/script.sh"
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "exec3"

        exec {
          validate {
            interpreter = "SHELL"

            file {
              allow_insecure = true

              gcs {
                bucket     = "test-bucket"
                object     = "test-object"
                generation = 1
              }
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter      = "SHELL"
            output_file_path = "$HOME/out"
            script           = "pwd"
          }
        }
      }

      resources {
        id = "exec4"

        exec {
          validate {
            interpreter      = "SHELL"
            output_file_path = "$HOME/out"
            script           = "pwd"
          }

          enforce {
            interpreter = "SHELL"

            file {
              allow_insecure = true

              gcs {
                bucket     = "test-bucket"
                object     = "test-object"
                generation = 1
              }
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "file1"

        file {
          path  = "$HOME/file"
          state = "PRESENT"

          file {
            local_path = "$HOME/file"
          }
        }
      }
    }

    resource_groups {
      resources {
        id = "file2"

        file {
          path  = "$HOME/file"
          state = "PRESENT"

          file {
            allow_insecure = true

            remote {
              uri             = "https://www.example.com/file"
              sha256_checksum = "c7938fed83afdccbb0e86a2a2e4cad7d5035012ca3214b4a61268393635c3063"
            }
          }
        }
      }

      resources {
        id = "file3"

        file {
          path  = "$HOME/file"
          state = "PRESENT"

          file {
            gcs {
              bucket     = "test-bucket"
              object     = "test-object"
              generation = 1
            }
          }
        }
      }

      resources {
        id = "file4"

        file {
          path    = "$HOME/file"
          state   = "PRESENT"
          content = "sample-content"
        }
      }
    }
  }

  rollout {
    disruption_budget {
      percent = 1
    }

    min_wait_duration = "3.5s"
  }

  description = "A test os policy assignment"
  project     = "%{project_name}"
}


`, context)
}

func testAccOsConfigOsPolicyAssignment_PercentOsPolicyAssignment(context map[string]interface{}) string {
	return Nprintf(`
resource "google_os_config_os_policy_assignment" "primary" {
  instance_filter {
    all = true
  }

  location = "%{zone}"
  name     = "tf-test-assignment%{random_suffix}"

  os_policies {
    id   = "policy"
    mode = "VALIDATION"

    resource_groups {
      resources {
        id = "apt-to-yum"

        repository {
          apt {
            archive_type = "DEB"
            components   = ["doc"]
            distribution = "debian"
            uri          = "https://atl.mirrors.clouvider.net/debian"
            gpg_key      = ".gnupg/pubring.kbx"
          }
        }
      }

      resources {
        id = "yum"

        repository {
          yum {
            base_url     = "http://centos.s.uw.edu/centos/"
            id           = "yum"
            display_name = "yum"
            gpg_keys     = ["RPM-GPG-KEY-CentOS-7"]
          }
        }
      }

      resources {
        id = "zypper"

        repository {
          zypper {
            base_url     = "http://mirror.dal10.us.leaseweb.net/opensuse"
            id           = "zypper"
            display_name = "zypper"
            gpg_keys     = ["sample-key-uri"]
          }
        }
      }

      resources {
        id = "goo"

        repository {
          goo {
            name = "goo"
            url  = "https://foo.com/googet/bar"
          }
        }
      }

      resources {
        id = "exec1"

        exec {
          validate {
            interpreter = "SHELL"
            args        = ["arg1"]

            file {
              local_path = "$HOME/script.sh"
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter = "SHELL"
            args        = ["arg1"]

            file {
              allow_insecure = true

              remote {
                uri             = "https://www.example.com/script.sh"
                sha256_checksum = "c7938fed83afdccbb0e86a2a2e4cad7d5035012ca3214b4a61268393635c3063"
              }
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "exec2"

        exec {
          validate {
            interpreter = "SHELL"
            args        = ["arg1"]

            file {
              allow_insecure = true

              remote {
                uri             = "https://www.example.com/script.sh"
                sha256_checksum = "c7938fed83afdccbb0e86a2a2e4cad7d5035012ca3214b4a61268393635c3063"
              }
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter = "SHELL"
            args        = ["arg1"]

            file {
              local_path = "$HOME/script.sh"
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "exec3"

        exec {
          validate {
            interpreter = "SHELL"

            file {
              allow_insecure = true

              gcs {
                bucket     = "test-bucket"
                object     = "test-object"
                generation = 1
              }
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter      = "SHELL"
            output_file_path = "$HOME/out"
            script           = "pwd"
          }
        }
      }

      resources {
        id = "exec4"

        exec {
          validate {
            interpreter      = "SHELL"
            output_file_path = "$HOME/out"
            script           = "pwd"
          }

          enforce {
            interpreter = "SHELL"

            file {
              allow_insecure = true

              gcs {
                bucket     = "test-bucket"
                object     = "test-object"
                generation = 1
              }
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "file1"

        file {
          path  = "$HOME/file"
          state = "PRESENT"

          file {
            local_path = "$HOME/file"
          }
        }
      }
    }

    resource_groups {
      resources {
        id = "file2"

        file {
          path  = "$HOME/file"
          state = "PRESENT"

          file {
            allow_insecure = true

            remote {
              uri             = "https://www.example.com/file"
              sha256_checksum = "c7938fed83afdccbb0e86a2a2e4cad7d5035012ca3214b4a61268393635c3063"
            }
          }
        }
      }

      resources {
        id = "file3"

        file {
          path  = "$HOME/file"
          state = "PRESENT"

          file {
            gcs {
              bucket     = "test-bucket"
              object     = "test-object"
              generation = 1
            }
          }
        }
      }

      resources {
        id = "file4"

        file {
          path    = "$HOME/file"
          state   = "PRESENT"
          content = "sample-content"
        }
      }
    }
  }

  rollout {
    disruption_budget {
      percent = 1
    }

    min_wait_duration = "3.5s"
  }

  description = "A test os policy assignment"
  project     = "%{project_name}"
}


`, context)
}

func testAccOsConfigOsPolicyAssignment_PercentOsPolicyAssignmentUpdate0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_os_config_os_policy_assignment" "primary" {
  instance_filter {
    all = false

    exclusion_labels {
      labels = {
        label-two = "value-two"
      }
    }

    inclusion_labels {
      labels = {
        label-one = "value-one"
      }
    }
  }

  location = "%{zone}"
  name     = "tf-test-assignment%{random_suffix}"

  os_policies {
    id   = "new-policy"
    mode = "VALIDATION"

    resource_groups {
      resources {
        id = "apt-to-yum"

        repository {
          yum {
            base_url     = "http://mirrors.rcs.alaska.edu/centos/"
            id           = "new-yum"
            display_name = "new-yum"
            gpg_keys     = ["RPM-GPG-KEY-CentOS-Debug-7"]
          }
        }
      }

      resources {
        id = "new-yum"

        repository {
          yum {
            base_url     = "http://mirrors.rcs.alaska.edu/centos/"
            id           = "new-yum"
            display_name = "new-yum"
            gpg_keys     = ["RPM-GPG-KEY-CentOS-Debug-7"]
          }
        }
      }

      resources {
        id = "new-zypper"

        repository {
          zypper {
            base_url     = "http://mirror.vtti.vt.edu/opensuse/"
            id           = "new-zypper"
            display_name = "new-zypper"
            gpg_keys     = ["new-sample-key-uri"]
          }
        }
      }

      resources {
        id = "new-goo"

        repository {
          goo {
            name = "new-goo"
            url  = "https://foo.com/googet/baz"
          }
        }
      }

      resources {
        id = "new-exec1"

        exec {
          validate {
            interpreter = "POWERSHELL"
            args        = ["arg2"]

            file {
              local_path = "$HOME/script.bat"
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter = "POWERSHELL"
            args        = ["arg2"]

            file {
              allow_insecure = false

              remote {
                uri             = "https://www.example.com/script.bat"
                sha256_checksum = "9f8e5818ccb47024d01000db713c0a333679b64678ff5fe2d9bea0a23014dd54"
              }
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "new-exec2"

        exec {
          validate {
            interpreter = "POWERSHELL"
            args        = ["arg2"]

            file {
              allow_insecure = false

              remote {
                uri             = "https://www.example.com/script.bat"
                sha256_checksum = "9f8e5818ccb47024d01000db713c0a333679b64678ff5fe2d9bea0a23014dd54"
              }
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter = "POWERSHELL"
            args        = ["arg2"]

            file {
              local_path = "$HOME/script.bat"
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "new-exec3"

        exec {
          validate {
            interpreter = "POWERSHELL"

            file {
              allow_insecure = false

              gcs {
                bucket     = "new-test-bucket"
                object     = "new-test-object"
                generation = 2
              }
            }

            output_file_path = "$HOME/out"
          }

          enforce {
            interpreter      = "POWERSHELL"
            output_file_path = "$HOME/out"
            script           = "dir"
          }
        }
      }

      resources {
        id = "new-exec4"

        exec {
          validate {
            interpreter      = "POWERSHELL"
            output_file_path = "$HOME/out"
            script           = "dir"
          }

          enforce {
            interpreter = "POWERSHELL"

            file {
              allow_insecure = false

              gcs {
                bucket     = "new-test-bucket"
                object     = "new-test-object"
                generation = 2
              }
            }

            output_file_path = "$HOME/out"
          }
        }
      }

      resources {
        id = "new-file1"

        file {
          path  = "$HOME/new-file"
          state = "PRESENT"

          file {
            local_path = "$HOME/new-file"
          }
        }
      }
    }

    resource_groups {
      resources {
        id = "new-file2"

        file {
          path  = "$HOME/new-file"
          state = "CONTENTS_MATCH"

          file {
            allow_insecure = false

            remote {
              uri             = "https://www.example.com/new-file"
              sha256_checksum = "9f8e5818ccb47024d01000db713c0a333679b64678ff5fe2d9bea0a23014dd54"
            }
          }
        }
      }

      resources {
        id = "new-file3"

        file {
          path  = "$HOME/new-file"
          state = "CONTENTS_MATCH"

          file {
            gcs {
              bucket     = "new-test-bucket"
              object     = "new-test-object"
              generation = 2
            }
          }
        }
      }

      resources {
        id = "new-file4"

        file {
          path    = "$HOME/new-file"
          state   = "CONTENTS_MATCH"
          content = "new-sample-content"
        }
      }
    }
  }

  rollout {
    disruption_budget {
      percent = 2
    }

    min_wait_duration = "3.5s"
  }

  description = "A test os policy assignment"
  project     = "%{project_name}"
}


`, context)
}

func testAccCheckOsConfigOsPolicyAssignmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_os_config_os_policy_assignment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &osconfig.OSPolicyAssignment{
				Location:           dcl.String(rs.Primary.Attributes["location"]),
				Name:               dcl.String(rs.Primary.Attributes["name"]),
				Description:        dcl.String(rs.Primary.Attributes["description"]),
				Project:            dcl.StringOrNil(rs.Primary.Attributes["project"]),
				Baseline:           dcl.Bool(rs.Primary.Attributes["baseline"] == "true"),
				Deleted:            dcl.Bool(rs.Primary.Attributes["deleted"] == "true"),
				Etag:               dcl.StringOrNil(rs.Primary.Attributes["etag"]),
				Reconciling:        dcl.Bool(rs.Primary.Attributes["reconciling"] == "true"),
				RevisionCreateTime: dcl.StringOrNil(rs.Primary.Attributes["revision_create_time"]),
				RevisionId:         dcl.StringOrNil(rs.Primary.Attributes["revision_id"]),
				RolloutState:       osconfig.OSPolicyAssignmentRolloutStateEnumRef(rs.Primary.Attributes["rollout_state"]),
				Uid:                dcl.StringOrNil(rs.Primary.Attributes["uid"]),
			}

			client := NewDCLOsConfigClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetOSPolicyAssignment(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_os_config_os_policy_assignment still exists %v", obj)
			}
		}
		return nil
	}
}
