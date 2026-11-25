package bigquery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestBigQueryDatasetAccessPatchRetriesWithFreshEtag(t *testing.T) {
	var patchCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, fmt.Sprintf("request method = %q, want PATCH", r.Method), http.StatusInternalServerError)
			return
		}

		patchCount++
		expectedEtag := fmt.Sprintf("etag-%d", patchCount)
		if got := r.Header.Get("If-Match"); got != expectedEtag {
			t.Errorf("If-Match = %q, want %q", got, expectedEtag)
		}
		if patchCount == 1 {
			w.WriteHeader(http.StatusPreconditionFailed)
			_, _ = w.Write([]byte(`{"error":{"code":412,"message":"Precondition Failed"}}`))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	etagReads := 0
	encoder := func(_ *schema.ResourceData, _ interface{}, _ map[string]interface{}) (map[string]interface{}, string, error) {
		etagReads++
		return map[string]interface{}{"access": []interface{}{}}, fmt.Sprintf("etag-%d", etagReads), nil
	}

	_, err := resourceBigQueryDatasetAccessPatchWithRetry(
		nil,
		&transport_tpg.Config{Client: server.Client(), SleepUnit: time.Millisecond},
		nil,
		"",
		"",
		server.URL,
		time.Second,
		encoder,
	)
	if err != nil {
		t.Fatalf("resourceBigQueryDatasetAccessPatchWithRetry() error = %v", err)
	}
	if got, want := patchCount, 2; got != want {
		t.Errorf("PATCH request count = %d, want %d", got, want)
	}
	if got, want := etagReads, 2; got != want {
		t.Errorf("ETag read count = %d, want %d", got, want)
	}
}

func TestSetBigQueryDatasetEtagHeader(t *testing.T) {
	testCases := []struct {
		name string
		etag string
		want string
	}{
		{name: "present", etag: "etag-value", want: "etag-value"},
		{name: "missing", etag: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			headers := make(http.Header)
			setBigQueryDatasetEtagHeader(headers, testCase.etag)
			if got := headers.Get("If-Match"); got != testCase.want {
				t.Errorf("If-Match = %q, want %q", got, testCase.want)
			}
		})
	}
}
