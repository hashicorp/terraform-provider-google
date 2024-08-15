// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package artifactregistry

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages#DockerImage
type DockerImage struct {
	name           string
	uri            string
	tags           []string
	imageSizeBytes string
	mediaType      string
	uploadTime     string
	buildTime      string
	updateTime     string
}

func DataSourceArtifactRegistryDockerImage() *schema.Resource {

	return &schema.Resource{
		Read: DataSourceArtifactRegistryDockerImageRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The region of the artifact registry repository. For example, "us-west1".`,
			},
			"repository_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The last part of the repository name to fetch from.`,
			},
			"image_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The image name to fetch.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fully qualified name of the fetched image.`,
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI to access the image.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `All tags associated with the image.`,
			},
			"image_size_bytes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Calculated size of the image in bytes.`,
			},
			"media_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Media type of this image.`,
			},
			"upload_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time, as a RFC 3339 string, the image was uploaded.`,
			},
			"build_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time, as a RFC 3339 string, this image was built.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time, as a RFC 3339 string, this image was updated.`,
			},
		},
	}
}

func DataSourceArtifactRegistryDockerImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	var res DockerImage

	imageName, tag, digest := parseImage(d.Get("image_name").(string))

	if digest != "" {
		// fetch image by digest
		// https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages/get
		imageUrlSafe := url.QueryEscape(imageName)
		urlRequest, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{ArtifactRegistryBasePath}}projects/{{project}}/locations/{{location}}/repositories/{{repository_id}}/dockerImages/%s@%s", imageUrlSafe, digest))
		if err != nil {
			return fmt.Errorf("Error setting api endpoint")
		}

		resGet, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    urlRequest,
			UserAgent: userAgent,
		})
		if err != nil {
			return err
		}

		res = convertResponseToStruct(resGet)
	} else {
		// fetch the list of images, ordered by update time
		// https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages/list
		urlRequest, err := tpgresource.ReplaceVars(d, config, "{{ArtifactRegistryBasePath}}projects/{{project}}/locations/{{location}}/repositories/{{repository_id}}/dockerImages")
		if err != nil {
			return fmt.Errorf("Error setting api endpoint")
		}

		urlRequest, err = transport_tpg.AddQueryParams(urlRequest, map[string]string{"orderBy": "update_time desc"})
		if err != nil {
			return err
		}

		res, err = retrieveAndFilterImages(d, config, urlRequest, userAgent, imageName, tag)
		if err != nil {
			return err
		}
	}

	// set the schema data using the response
	if err := d.Set("name", res.name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("self_link", res.uri); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}

	if err := d.Set("tags", res.tags); err != nil {
		return fmt.Errorf("Error setting tags: %s", err)
	}

	if err := d.Set("image_size_bytes", res.imageSizeBytes); err != nil {
		return fmt.Errorf("Error setting image_size_bytes: %s", err)
	}

	if err := d.Set("media_type", res.mediaType); err != nil {
		return fmt.Errorf("Error setting media_type: %s", err)
	}

	if err := d.Set("upload_time", res.uploadTime); err != nil {
		return fmt.Errorf("Error setting upload_time: %s", err)
	}

	if err := d.Set("build_time", res.buildTime); err != nil {
		return fmt.Errorf("Error setting build_time: %s", err)
	}

	if err := d.Set("update_time", res.updateTime); err != nil {
		return fmt.Errorf("Error setting update_time: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "{{ArtifactRegistryBasePath}}projects/{{project}}/locations/{{location}}/repositories/{{repository_id}}/dockerImages/{{image_name}}")
	if err != nil {
		return fmt.Errorf("Error constructing the data source id: %s", err)
	}

	d.SetId(id)

	return nil
}

func parseImage(image string) (imageName string, tag string, digest string) {
	splitByAt := strings.Split(image, "@")
	splitByColon := strings.Split(image, ":")

	if len(splitByAt) == 2 {
		imageName = splitByAt[0]
		digest = splitByAt[1]
	} else if len(splitByColon) == 2 {
		imageName = splitByColon[0]
		tag = splitByColon[1]
	} else {
		imageName = image
	}

	return imageName, tag, digest
}

func retrieveAndFilterImages(d *schema.ResourceData, config *transport_tpg.Config, urlRequest string, userAgent string, imageName string, tag string) (DockerImage, error) {
	// Paging through the list method until either:
	// if a tag was provided, the matching image name and tag pair
	// otherwise, return the first matching image name

	for {
		resListImages, token, err := retrieveListOfDockerImages(config, urlRequest, userAgent)
		if err != nil {
			return DockerImage{}, err
		}

		var resFiltered []DockerImage
		for _, image := range resListImages {
			if strings.Contains(image.name, "/"+url.QueryEscape(imageName)+"@") {
				resFiltered = append(resFiltered, image)
			}
		}

		if tag != "" {
			for _, image := range resFiltered {
				for _, iterTag := range image.tags {
					if iterTag == tag {
						return image, nil
					}
				}
			}
		} else if len(resFiltered) > 0 {
			return resFiltered[0], nil
		}

		if token == "" {
			return DockerImage{}, fmt.Errorf("Requested image was not found.")
		}

		urlRequest, err = transport_tpg.AddQueryParams(urlRequest, map[string]string{"pageToken": token})
		if err != nil {
			return DockerImage{}, err
		}
	}
}

func retrieveListOfDockerImages(config *transport_tpg.Config, urlRequest string, userAgent string) ([]DockerImage, string, error) {
	resList, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    urlRequest,
		UserAgent: userAgent,
	})
	if err != nil {
		return make([]DockerImage, 0), "", err
	}

	if nextPageToken, ok := resList["nextPageToken"].(string); ok {
		return flattenDataSourceListResponse(resList), nextPageToken, nil
	} else {
		return flattenDataSourceListResponse(resList), "", nil
	}
}

func flattenDataSourceListResponse(res map[string]interface{}) []DockerImage {
	var dockerImages []DockerImage

	resDockerImages, _ := res["dockerImages"].([]interface{})

	for _, resImage := range resDockerImages {
		image, _ := resImage.(map[string]interface{})
		dockerImages = append(dockerImages, convertResponseToStruct(image))
	}

	return dockerImages
}

func convertResponseToStruct(res map[string]interface{}) DockerImage {
	var dockerImage DockerImage

	if name, ok := res["name"].(string); ok {
		dockerImage.name = name
	}

	if uri, ok := res["uri"].(string); ok {
		dockerImage.uri = uri
	}

	if tags, ok := res["tags"].([]interface{}); ok {
		var stringTags []string

		for _, tag := range tags {
			strTag := tag.(string)
			stringTags = append(stringTags, strTag)
		}
		dockerImage.tags = stringTags
	}

	if imageSizeBytes, ok := res["imageSizeBytes"].(string); ok {
		dockerImage.imageSizeBytes = imageSizeBytes
	}

	if mediaType, ok := res["mediaType"].(string); ok {
		dockerImage.mediaType = mediaType
	}

	if uploadTime, ok := res["uploadTime"].(string); ok {
		dockerImage.uploadTime = uploadTime
	}

	if buildTime, ok := res["buildTime"].(string); ok {
		dockerImage.buildTime = buildTime
	}

	if updateTime, ok := res["updateTime"].(string); ok {
		dockerImage.updateTime = updateTime
	}

	return dockerImage
}
