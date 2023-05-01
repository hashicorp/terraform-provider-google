package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigQueryRoutine_bigQueryRoutine_Update(t *testing.T) {
	t.Parallel()

	dataset := fmt.Sprintf("tfmanualdataset%s", RandString(t, 10))
	routine := fmt.Sprintf("tfmanualroutine%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryRoutineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryRoutine_bigQueryRoutine(dataset, routine),
			},
			{
				ResourceName:      "google_bigquery_routine.sproc",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryRoutine_bigQueryRoutine_Update(dataset, routine),
			},
			{
				ResourceName:      "google_bigquery_routine.sproc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigQueryRoutine_bigQueryRoutine(dataset, routine string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_routine" "sproc" {
  dataset_id = google_bigquery_dataset.test.dataset_id
  routine_id     = "%s"
  routine_type = "SCALAR_FUNCTION"
  language = "SQL"
  definition_body = "1"
}
`, dataset, routine)
}

func testAccBigQueryRoutine_bigQueryRoutine_Update(dataset, routine string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_routine" "sproc" {
  dataset_id = google_bigquery_dataset.test.dataset_id
  routine_id     = "%s"
  routine_type = "SCALAR_FUNCTION"
  language = "JAVASCRIPT"
  definition_body = "CREATE FUNCTION multiplyInputs return x*y;"
  arguments {
    name = "x"
    data_type = "{\"typeKind\" :  \"FLOAT64\"}"
  }
  arguments {
    name = "y"
    data_type = "{\"typeKind\" :  \"FLOAT64\"}"
  }

  return_type = "{\"typeKind\" :  \"FLOAT64\"}"
}
`, dataset, routine)
}
