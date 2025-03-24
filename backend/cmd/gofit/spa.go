package main

import (
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/bamaas/gofit/internal/assets"
)

func (app *application) spaHandler(w http.ResponseWriter, r *http.Request) {
	serveIndex := func (){
        index, err := assets.BuildFs.ReadFile(filepath.Join(assets.BuildPath, "index.html"))
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusAccepted)
        _, err = w.Write(index)
        if err != nil {
            http.Error(w, "Error writing response", http.StatusInternalServerError)
            return
        }
	}

    f, err := assets.BuildFs.Open(filepath.Join(assets.BuildPath, r.URL.Path))
	// if file is not found, serve index.html
    if errors.Is(err, fs.ErrNotExist) {
		serveIndex()
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	// if file is a directory, serve index.html
	stat, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if stat.IsDir() {
		serveIndex()
        return
	}

    defer f.Close()
    http.FileServer(assets.BuildHTTPFS()).ServeHTTP(w, r)
}