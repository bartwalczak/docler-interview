package server

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestHealthHandler(t *testing.T) {
	e := httpexpect.New(t, serverURL)
	e.GET("/check").Expect().Status(http.StatusOK)
}
