package gitlab

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetJob(t *testing.T) {
	fixtureFile := "../fixtures/rest/job.json"
	data, err := os.ReadFile(fixtureFile)
	if err != nil {
		t.Errorf("Could not read fixture file %v", fixtureFile)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if r.URL.Path != "/job" {
		// 	t.Errorf("Request path should be /job, got %v", r.URL.Path)
		// }
		// if r.URL.RawQuery != "?job_token=<job-token>" {
		// 	t.Errorf("Request path should be /job, got %v", r.URL.RawQuery)
		// }
		fmt.Fprintf(w, string(data))
	}))
	client := ApiClient{Url: srv.URL, JobToken: "<job-token>"}
	job, err := client.GetJob()
	if err != nil {
		t.Errorf("Could not fetch job details %v", err)
	}
	if len(job.Tags) != 2 {
		t.Errorf("Expected job to have 2 tags, got %v", len(job.Tags))
	}
}
