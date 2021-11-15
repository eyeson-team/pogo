package gitlab

import (
	"encoding/json"
	"os"
	"path"
	"regexp"
	"strings"
)

// ServiceImage is used to decode service information that is provided by
// GitLab in JSON encoded form through the CUSTOM_ENV_CI_JOB_SERVICES variable.
type ServiceImage struct {
	Name       string   `json:name`
	Alias      string   `json:alias`
	Entrypoint []string `json:entrypoint`
	Command    []string `json:command`
}

// GetContainerName constructs a unique readable name by the given environment
// variables that are set by GitLab. It consists of the runner identifier, the
// project identifier, the unique concurrent project identifier and the current
// job identifier.
func GetContainerName() string {
	parts := []string{
		"runner",
		os.Getenv("CUSTOM_ENV_CI_RUNNER_ID"),
		"project",
		os.Getenv("CUSTOM_ENV_CI_PROJECT_ID"),
		"concurrent",
		os.Getenv("CUSTOM_ENV_CI_CONCURRENT_PROJECT_ID"),
		"job",
		os.Getenv("CUSTOM_ENV_CI_JOB_ID"),
	}
	return strings.Join(parts, "-")
}

// GetContainerImage returns the container image required by the current CI
// job.
func GetContainerImage() string {
	return os.Getenv("CUSTOM_ENV_CI_JOB_IMAGE")
}

// GetJobToken returns the job token.
func GetJobToken() string {
	return os.Getenv("CUSTOM_ENV_CI_JOB_TOKEN")
}

// GetApiUrl returns the api url.
func GetApiUrl() string {
	return os.Getenv("CUSTOM_ENV_CI_API_V4_URL")
}

// GetServiceContainerName does provide a service container name that can be
// used for a given service. Note that it reduces all characters to
// alphanumeric in order to avoid troubles for the container name.
func GetServiceContainerName(serviceName string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	parts := []string{GetContainerName(), "service", reg.ReplaceAllString(serviceName, "")}
	return strings.Join(parts, "-")
}

// GetServiceImages provides all services contained in the CI variable
// CUSTOM_ENV_CI_JOB_SERVICES. This function can fail if the variable does
// not contain valid JSON encoded data, but does allow an empty string.
func GetServiceImages() ([]ServiceImage, error) {
	res := []ServiceImage{}
	serviceJson := os.Getenv("CUSTOM_ENV_CI_JOB_SERVICES")
	if serviceJson == "" {
		return res, nil
	}
	err := json.Unmarshal([]byte(serviceJson), &res)
	return res, err
}

// JobTagsMatch checks if the current tags of a job intersect with a given
// list of tags.
func JobTagsMatch(tags *[]string) (bool, error) {
	job, err := NewApiClient().GetJob()
	if err != nil {
		return false, err
	}
	for _, tag := range *tags {
		for _, jobTag := range job.Tags {
			if tag == jobTag {
				return true, nil
			}
		}
	}
	return false, nil
}

// GetPathSlug provides a unique path by project using the CI environment
// variable.
func GetPathSlug() string {
	return os.Getenv("CUSTOM_ENV_CI_PROJECT_PATH_SLUG")
}

// GetBuildsDir provides the project build directory.
func GetBuildsDir() string {
	return path.Join("/builds", GetPathSlug())
}

// GetCacheDir provides the project cache directory.
func GetCacheDir() string {
	return path.Join("/cache", GetPathSlug())
}
