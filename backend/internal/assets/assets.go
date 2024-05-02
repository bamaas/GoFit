package assets

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

// Embed the build directory from the frontend.
//go:embed all:static
var BuildFs embed.FS

var BuildPath string = "static"

// Get the subtree of the embedded files with `static` directory as a root.
func BuildHTTPFS() http.FileSystem {
    build, err := fs.Sub(BuildFs, BuildPath)
    if err != nil {
        log.Fatal(err)		// TODO: fix error handling
    }
    return http.FS(build)
}