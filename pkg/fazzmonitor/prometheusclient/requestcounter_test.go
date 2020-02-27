package prometheusclient_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/payfazz/go-apt/pkg/fazzmonitor/prometheusclient"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Test_InstrumentHandlerCounter(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest("", "/", nil)
	prometheusclient.InstrumentHandlerCounter("/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
	promhttp.Handler().ServeHTTP(w, r)
	substr := `http_requests_total{code="200",method="GET",path="/users"} 1`
	if !strings.Contains(w.Body.String(), substr) {
		t.Fatalf("expect to contain = %s", substr)
	}
}
