package commons

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildInfo(t *testing.T) {

	buildInfo := GetBuildInfo()
	buildInfo.Name = "test"
	rr := httptest.NewRecorder()
	e := WriteJSON(http.StatusOK, buildInfo, rr)
	// Check the status code is what we expect.
	if nil != e {
		t.Error("Something went wrong with serialization")
	}

	expected := `{"name":"test"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("incorrect build format response: got %v want %v",
			rr.Body.String(), expected)
	}

}
