package public

import (
	"embed"
	"io/fs"
	"net/http"
	"time"

	"github.com/ugent-library/people-service/old/modifiedfs"
)

//go:embed swagger-ui-5.1.0/*.js swagger-ui-5.1.0/*.css swagger-ui-5.1.0/*.html swagger-ui-5.1.0/*.png
var assets embed.FS

func swaggerUI() fs.FS {
	fs, err := fs.Sub(assets, "swagger-ui-5.1.0")
	if err != nil {
		panic(err)
	}
	return fs
}

func SwaggerFileServer() http.Handler {
	return http.FileServer(modifiedfs.FSWithStatModified(swaggerUI(), time.Now()))
}
