package api

import (
	"embed"
	"net/http"
	"time"

	"github.com/ugent-library/people-service/modifiedfs"
)

//go:embed openapi.yaml
var openapiFS embed.FS

func OpenapiFileServer() http.Handler {
	return http.FileServer(modifiedfs.FSWithStatModified(openapiFS, time.Now()))
}
