package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"io/ioutil"
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
		defer f.Close()

		_, err = os.Stat(UPLOAD_DIR)

		if err != nil && os.IsNotExist(err) {
			errDir := os.Mkdir(UPLOAD_DIR, 0777)
			if errDir != nil {
				http.Error(w, errDir.Error(), http.StatusInternalServerError)
			}
		}

		t, err := os.Create(UPLOAD_DIR + "/" + filename)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError+3)
			return
		}
		http.Redirect(w, r, "/view?id="+filename, http.StatusFound)
	}

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-type", "image")
	http.ServeFile(w, r, imagePath)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir(UPLOAD_DIR)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var listHtml string
	for _, fileInfo := range fileInfoArr {
		imgid := fileInfo.Name()
		listHtml += "<li><a target=\"_blank\" href=\"/view?id="+ imgid +"\">" + imgid + "</a></li>"
	}
	w.Header().Set("Content-type", "text/html")
	io.WriteString(w, "<ol>" + listHtml + "</ol>")

}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}



func main() {

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/", listHandler)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
}
