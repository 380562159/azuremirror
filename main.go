package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
	"log"
	"flag"
	"path/filepath"
	"path"
)


var (
	contentDir *string
	ContentMap map[string][]byte // path, content
)

func InitContentMap() {
	ContentMap = make(map[string][]byte)

	mainPage, err := ReadContentByPath(filepath.Join(*contentDir, "html/main.html"))
	if err == nil {
		ContentMap["/"] = mainPage
	}

	stylePage, err := ReadContentByPath(filepath.Join(*contentDir, "css/style.css"))
	if err == nil {
		ContentMap["/style.css"] = stylePage
	}

	iconDir := filepath.Join(*contentDir, "icons")
	files, err := ioutil.ReadDir(iconDir)
	if err == nil {
		for _, file := range files {
			content, err := ioutil.ReadFile(filepath.Join(iconDir, file.Name()))
			if err == nil {
				ContentMap[path.Join("/icons", file.Name())] = content
			}
		}
	}
}


func ReadContentByPath(filePath string) ([]byte, error) {
	content, err := ioutil.ReadFile(filePath)
	return content, err
}



func rootHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimSpace(r.URL.Path)
	if urlPath == "/" {
		fp := filepath.Join(*contentDir, "html/main.html")
		c, err := ReadContentByPath(fp)
		if err != nil {
			w.WriteHeader(500)
		} else {
			w.Write(c)
			// fmt.Fprintf(w, string(c))
		}
	} else if urlPath == "/style.css" {
		fp := filepath.Join(*contentDir, "css/style.css")
		c, err := ReadContentByPath(fp)
		if err != nil {
			w.WriteHeader(500)
		} else {
			w.Header().Set("Content-Type", "text/css")
			w.Write(c)
		}
	} else {
		w.WriteHeader(404)
		fmt.Fprintf(w, "page not found")
	}
}

func GeneralHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimSpace(r.URL.Path)
	content := ContentMap[urlPath]
	if len(content) == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "page not found")
	} else {
		if strings.HasSuffix(urlPath, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(urlPath, ".png") {
			w.Header().Set("Content-Type", "image/png")
		} else {
			w.Header().Set("Content-Type", "text/html")
		}
		w.Write(content)
	}
}

func IconHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	fp := filepath.Join(*contentDir, urlPath)
	c, err := ReadContentByPath(fp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	} else {
		w.Header().Set("Content-Type", "image/png")
		// w.Header().Set("Content-Length", strconv.Itoa(len(c)))
		w.Write(c)
	}
}

func main() {
	contentDir = flag.String("contentDir", "./content", "content directory path")
	flag.Parse()

	InitContentMap()

	// fast
	http.HandleFunc("/", GeneralHandler)

	// slow
//	http.HandleFunc("/", rootHandler)
//	http.HandleFunc("/icons/", IconHandler)

	log.Fatal(http.ListenAndServe(":8010", nil))
}
