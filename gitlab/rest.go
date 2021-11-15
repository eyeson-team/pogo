package gitlab

import (
	"encoding/json"
	"io"
	"net/http"
)

// Job defines the JSON formatted data structure received from the GitLab API.
type Job struct {
	Tags []string `json:"tag_list"`
}

// ApiClient provides a REST API client for a GitLab instance.
type ApiClient struct {
	Url      string
	JobToken string
}

// GetJob fetches job details from the gitlab API.
//
// ref. https://docs.gitlab.com/ee/api/jobs.html#get-job-tokens-job
func (r *ApiClient) GetJob() (*Job, error) {
	resp, err := http.Get(r.Url + "/job?job_token=" + r.JobToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var job Job
	err = json.Unmarshal(body, &job)
	return &job, err
}

// NewApiClient provides a new API client using the current environment of the
// running job.
func NewApiClient() *ApiClient {
	return &ApiClient{Url: GetApiUrl(), JobToken: GetJobToken()}
}
