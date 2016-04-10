package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	UPLOAD_DIR = "./uploads"
	VIEWS_DIR  = "./views"
)

var templates map[string]*template.Template

func init() {
	fileInfoArr, err := ioutil.ReadDir(VIEWS_DIR)
	if err != nil {
		panic(err)
		return
	}

	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()

		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}

		templatePath = VIEWS_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templatePath] = t

	}

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		w.Header().Set("content-type", "text/html;charset=utf-8")
		if err := renderHtml(w, VIEWS_DIR+"/upload", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

		var saveFile string
		var now time.Time
		now = time.Now()
		saveFile = strconv.Itoa(now.Day()) + "-" + strconv.Itoa(int(now.Unix())) + path.Ext(filename)

		t, err := os.Create(UPLOAD_DIR + "/" + saveFile)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError+3)
			return
		}
		http.Redirect(w, r, "/view?id="+saveFile, http.StatusFound)
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

	locals := make(map[string]interface{})
	images := []string{}
	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
	}

	w.Header().Set("Content-type", "text/html")
	locals["images"] = images
	if err := renderHtml(w, VIEWS_DIR+"/list", locals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func renderHtml(w http.ResponseWriter, viewPath string, locals map[string]interface{}) (err error) {
	err = templates[viewPath+".html"].Execute(w, locals)
	return
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
