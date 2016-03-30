package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	UPLOAD_DIR = "./uploads"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "text/html;charset=utf-8")
		io.WriteString(w, `
            <form method="post" action="/upload" enctype="multipart/form-data">
                Choose an image to upload: <input name="image" type="file"/>
                <input type="submit" value="Upload"/>
            </form>
        `)
		return
	}

	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError+1)
			return
		}
		filename := h.Filename
		t, err := os.Create(UPLOAD_DIR + "/" + filename)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError+2)
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError+3)
			return
		}
		http.Redirect(w, r, "/view?id="+filename, http.StatusFound)
	}

}
func main() {

	http.HandleFunc("/upload", uploadHandler)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
}
