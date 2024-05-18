package assets

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

// Embed database migrations files
//go:embed all:migrations
var MigrationsFs embed.FS
var MigrationsPath string = "migrations"

// Embed static frontend files
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

// Get the subtree of the embedded files with `migrations` directory as a root.
// func BuildMigrationsFs() source.Driver {
//     d, err := iofs.New(MigrationsFs, MigrationsPath)
//     if err != nil {
//         log.Fatal(err)		// TODO: fix error handling
//     }
//     return d
// }