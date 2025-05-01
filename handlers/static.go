package handlers

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func NoDirListingHandler(fileServerHandler http.Handler, fileSystem fs.FS, stripPrefix string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trimmedPath := strings.TrimPrefix(r.URL.Path, stripPrefix)
		trimmedPath = path.Clean("/" + trimmedPath)
		trimmedPath = strings.TrimPrefix(trimmedPath, "/")

		if !fs.ValidPath(trimmedPath) {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}

		if trimmedPath == "" || strings.Contains(trimmedPath, "..") {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}

		f, err := fileSystem.Open(trimmedPath)
		if err == nil {
			defer f.Close()
			stat, err := f.Stat()
			if err == nil && stat.IsDir() {
				http.Error(w, "403 Forbidden", http.StatusForbidden)
				return
			}
		}

		fileServerHandler.ServeHTTP(w, r)
	})
}
