package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccRedisInstance_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_update(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstance_update2(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstance_update2(name, false),
			},
		},
	})
}

// Validate that read replica is enabled on the instance without having to recreate
func TestAccRedisInstance_updateReadReplicasMode(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceReadReplicasUnspecified(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstanceReadReplicasEnabled(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstanceReadReplicasUnspecified(name, false),
			},
		},
	})
}

/* Validate that read replica is enabled on the instance without recreate
 * and secondaryIp is auto provisioned when passed as 'auto' */
func TestAccRedisInstance_updateReadReplicasModeWithAutoSecondaryIp(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceReadReplicasUnspecified(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstanceReadReplicasEnabledWithAutoSecondaryIP(name, true),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstanceReadReplicasUnspecified(name, false),
			},
		},
	})
}

func testAccRedisInstanceReadReplicasUnspecified(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 5
	tier = "STANDARD_HA"
  region         = "us-central1"
	%s
  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
}
`, name, lifecycleBlock)
}

func testAccRedisInstanceReadReplicasEnabled(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 5
  tier = "STANDARD_HA"
  region         = "us-central1"
	%s
  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  read_replicas_mode = "READ_REPLICAS_ENABLED"
  secondary_ip_range = "10.79.0.0/28"
	}
`, name, lifecycleBlock)
}

func testAccRedisInstanceReadReplicasEnabledWithAutoSecondaryIP(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 5
  tier = "STANDARD_HA"
  region         = "us-central1"
	%s
  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  read_replicas_mode = "READ_REPLICAS_ENABLED"
  secondary_ip_range = "auto"
}
`, name, lifecycleBlock)
}

func TestAccRedisInstance_regionFromLocation(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", RandInt(t))

	// Pick a zone that isn't in the provider-specified region so we know we
	// didn't fall back to that one.
	region := "us-west1"
	zone := "us-west1-a"
	if acctest.GetTestRegionFromEnv() == "us-west1" {
		region = "us-central1"
		zone = "us-central1-a"
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_regionFromLocation(name, zone),
				Check:  resource.TestCheckResourceAttr("google_redis_instance.test", "region", region),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedisInstance_redisInstanceAuthEnabled(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_redisInstanceAuthEnabled(context),
			},
			{
				ResourceName:            "google_redis_instance.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccRedisInstance_redisInstanceAuthDisabled(context),
			},
			{
				ResourceName:            "google_redis_instance.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func TestSecondaryIpDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"empty strings": {
			Old:                "",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"auto range": {
			Old:                "",
			New:                "auto",
			ExpectDiffSuppress: false,
		},
		"auto on already applied range": {
			Old:                "10.0.0.0/28",
			New:                "auto",
			ExpectDiffSuppress: true,
		},
		"same ranges": {
			Old:                "10.0.0.0/28",
			New:                "10.0.0.0/28",
			ExpectDiffSuppress: true,
		},
		"different ranges": {
			Old:                "10.0.0.0/28",
			New:                "10.1.2.3/28",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if secondaryIpDiffSuppress("whatever", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestAccRedisInstance_downgradeRedisVersion(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstance_redis5(name),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisInstance_redis4(name),
			},
			{
				ResourceName:      "google_redis_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitRedisInstance_redisVersionIsDecreasing(t *testing.T) {
	t.Parallel()
	type testcase struct {
		name       string
		old        interface{}
		new        interface{}
		decreasing bool
	}
	tcs := []testcase{
		{
			name:       "stays the same",
			old:        "REDIS_4_0",
			new:        "REDIS_4_0",
			decreasing: false,
		},
		{
			name:       "increases",
			old:        "REDIS_4_0",
			new:        "REDIS_5_0",
			decreasing: false,
		},
		{
			name:       "nil vals",
			old:        nil,
			new:        "REDIS_4_0",
			decreasing: false,
		},
		{
			name:       "corrupted",
			old:        "REDIS_4_0",
			new:        "REDIS_banana",
			decreasing: false,
		},
		{
			name:       "decreases",
			old:        "REDIS_6_0",
			new:        "REDIS_4_0",
			decreasing: true,
		},
	}

	for _, tc := range tcs {
		decreasing := isRedisVersionDecreasingFunc(tc.old, tc.new)
		if decreasing != tc.decreasing {
			t.Errorf("%s: expected decreasing to be %v, but was %v", tc.name, tc.decreasing, decreasing)
		}
	}
}

func testAccRedisInstance_update(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "pre-update"
  memory_size_gb = 1
  region         = "us-central1"
	%s

  labels = {
    my_key    = "my_val"
    other_key = "other_val"
  }

  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  redis_version = "REDIS_4_0"
}
`, name, lifecycleBlock)
}

func testAccRedisInstance_update2(name string, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "post-update"
  memory_size_gb = 1
	%s

  labels = {
    my_key    = "my_val"
    other_key = "new_val"
  }

  redis_configs = {
    maxmemory-policy       = "noeviction"
    notify-keyspace-events = ""
  }
  redis_version = "REDIS_5_0"
}
`, name, lifecycleBlock)
}

func testAccRedisInstance_regionFromLocation(name, zone string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  memory_size_gb = 1
  location_id    = "%s"
}
`, name, zone)
}

func testAccRedisInstance_redisInstanceAuthEnabled(context map[string]interface{}) string {
	return Nprintf(`
resource "google_redis_instance" "cache" {
  name           = "tf-test-memory-cache%{random_suffix}"
  memory_size_gb = 1
  auth_enabled = true
}
`, context)
}

func testAccRedisInstance_redisInstanceAuthDisabled(context map[string]interface{}) string {
	return Nprintf(`
resource "google_redis_instance" "cache" {
  name           = "tf-test-memory-cache%{random_suffix}"
  memory_size_gb = 1
  auth_enabled = false
}
`, context)
}

func testAccRedisInstance_redis5(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 1
  region         = "us-central1"

  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  redis_version = "REDIS_5_0"
}
`, name)
}

func testAccRedisInstance_redis4(name string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "test" {
  name           = "%s"
  display_name   = "redissss"
  memory_size_gb = 1
  region         = "us-central1"

  redis_configs = {
    maxmemory-policy       = "allkeys-lru"
    notify-keyspace-events = "KEA"
  }
  redis_version = "REDIS_4_0"
}
`, name)
}
