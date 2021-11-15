package gitlab

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetServiceImagesCanHandleEmptyValue(t *testing.T) {
	images, err := GetServiceImages()

	if err != nil {
		t.Errorf("Failed to fetch service images with %v.", err)
	}
	if len(images) != 0 {
		t.Error("No images should be found.")
	}
}

func TestGetServiceImagesProvideNames(t *testing.T) {
	sample := "[{\"name\":\"redis:latest\",\"alias\":\"\",\"entrypoint\":null,\"command\":null},{\"name\":\"my-postgres:9.4\",\"alias\":\"pg\",\"entrypoint\":[\"path\",\"to\",\"entrypoint\"],\"command\":[\"path\",\"to\",\"cmd\"]}]"
	os.Setenv("CUSTOM_ENV_CI_JOB_SERVICES", sample)
	images, err := GetServiceImages()

	if err != nil {
		t.Errorf("Failed to fetch service images with %v", err)
	}

	if images[0].Name != "redis:latest" {
		t.Errorf("Expects redis:latest received %v", images[0].Name)
	}
}

func testServeJobApi() error {
	fixtureFile := "../fixtures/rest/job.json"
	data, err := os.ReadFile(fixtureFile)
	if err != nil {
		return err
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(data))
	}))
	os.Setenv("CUSTOM_ENV_CI_API_V4_URL", srv.URL)
	return nil
}

func TestJobTagsMatchMatch(t *testing.T) {
	if err := testServeJobApi(); err != nil {
		t.Error("Could not start fake API")
	}
	res, err := JobTagsMatch(&[]string{"test", "pogo runner"})
	if err != nil {
		t.Errorf("Tags match should not raise error, got %v", err)
	}
	if res == false {
		t.Error("Failed to match given tags")
	}
}

func TestJobTagsMatchNoMatch(t *testing.T) {
	if err := testServeJobApi(); err != nil {
		t.Error("Could not start fake API")
	}
	res, err := JobTagsMatch(&[]string{"test", "docker runner"})
	if err != nil {
		t.Errorf("Tags match should not raise error, got %v", err)
	}
	if res == true {
		t.Error("Failed to not match given tags")
	}
}
