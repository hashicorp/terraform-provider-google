package google

import (
	"fmt"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func getResourceName(resourceId *cloudresourcemanager.ResourceId) string {
	if resourceId == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s", resourceId.Type, resourceId.Id)
}
