package fbhttp

import (
	"net/http"

	"github.com/filebrowser/filebrowser/v2/files"
)

var metadataHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file, err := files.NewFileInfo(&files.FileOptions{
		Fs:      d.user.Fs,
		Path:    r.URL.Path,
		Checker: d,
	})
	if err != nil {
		return errToStatus(err), err
	}
	if file.IsDir {
		return http.StatusBadRequest, nil
	}

	fd, err := file.Fs.Open(file.Path)
	if err != nil {
		return errToStatus(err), err
	}
	defer fd.Close()

	metadata, err := files.DetectMetadata(file.Extension, fd)
	if err != nil {
		return errToStatus(err), err
	}

	return renderJSON(w, r, metadata)
})
