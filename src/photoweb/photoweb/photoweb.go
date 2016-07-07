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
	"runtime/debug"
)

const (
	UPLOAD_DIR = "./uploads"
	VIEWS_DIR  = "./views"
	checkFlag = 0x0001
)

var templates = make(map[string]*template.Template)

func init() {
	fileInfoArr, err := ioutil.ReadDir(VIEWS_DIR)

	check(err)

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

		renderHtml(w, VIEWS_DIR+"/upload", nil)

	}

	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		check(err)
		filename := h.Filename
		defer f.Close()

		_, err = os.Stat(UPLOAD_DIR)

		if err != nil && os.IsNotExist(err) {
			errDir := os.Mkdir(UPLOAD_DIR, 0777)
			check(errDir)
		}

		var saveFile string
		var now time.Time
		now = time.Now()
		saveFile = strconv.Itoa(now.Day()) + "-" + strconv.Itoa(int(now.Unix())) + path.Ext(filename)

		t, err := os.Create(UPLOAD_DIR + "/" + saveFile)

		check(err)
		defer t.Close()

		_, err = io.Copy(t, f)
		check(err)

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

	check(err)

	locals := make(map[string]interface{})
	images := []string{}
	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
	}

	w.Header().Set("Content-type", "text/html")
	locals["images"] = images
	renderHtml(w, VIEWS_DIR+"/list", locals)

}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func renderHtml(w http.ResponseWriter, viewPath string, locals map[string]interface{}) {
	err := templates[viewPath+".html"].Execute(w, locals)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func wrapHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + " " + r.Proto + " " + r.URL.Scheme + r.Host + r.URL.Path)
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				log.Println("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}

}

func staticDirHandle(mux *http.ServeMux, prefix string, staticDir string, flag int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix) - 1:]
		log.Println(r.Method + " " + r.Proto + " " + r.URL.Scheme + r.Host + file)
		if (checkFlag & flag) == 0 {
			if exists := isExists(file); !exists {
				http.NotFound(w, r)
				return
			}
		}
		http.ServeFile(w, r, file)

	})
}


func main() {
	mux := http.NewServeMux()

	staticDirHandle(mux, "/assets/", "./public", 0)
	mux.HandleFunc("/upload", wrapHandler(uploadHandler))
	mux.HandleFunc("/view", wrapHandler(viewHandler))
	mux.HandleFunc("/", wrapHandler(listHandler))
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
}
