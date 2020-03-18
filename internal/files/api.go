package files

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: only images
func (a *app) upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Println(err)
		return
	}

	file, header, err := r.FormFile("f")
	if err != nil {
		log.Println(err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: use different service for image storing
	f, err := os.OpenFile(filepath.Join("./images", header.Filename), os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	// TODO: fill rewrite file
	if _, err := f.Write(data); err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
