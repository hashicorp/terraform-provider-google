package google

import (
	"log"
	"strings"

	"github.com/hashicorp/errwrap"
	"google.golang.org/api/googleapi"
)

func transformSQLDatabaseReadError(err error) error {
	if gErr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error); ok {
		if gErr.Code == 400 && strings.Contains(gErr.Message, "Invalid request since instance is not running") {
			// This error occurs when attempting a GET after deleting the sql database and sql instance. It leads to to
			// inconsistent behavior as handleNotFoundError(...) expects an error code of 404 when a resource does not
			// exist. To get the desired behavior from handleNotFoundError, modify the return code to 404 so that
			// handleNotFoundError(...) will treat this as a NotFound error
			gErr.Code = 404
		}

		log.Printf("[DEBUG] Transformed SQLDatabase error")
		return gErr
	}

	return err
}
