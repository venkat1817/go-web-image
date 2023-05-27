package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	// Create the "gallery" directory if it doesn't exist
	err := os.MkdirAll("gallery", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", index)
	http.Handle("/gallery/", http.StripPrefix("/gallery", http.FileServer(http.Dir("./gallery"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		mf, fh, err := req.FormFile("nf")
		if err != nil {
			fmt.Println(err)
		}
		defer mf.Close()
		// create sha for file name
		ext := strings.Split(fh.Filename, ".")[1]
		h := sha1.New()
		io.Copy(h, mf)
		fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
		// create new file
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		path := filepath.Join(wd, "gallery", fname)
		nf, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer nf.Close()
		mf.Seek(0, 0)
		io.Copy(nf, mf)
	}

	file, err := os.Open("gallery")
	if err != nil {
		log.Println("Failed opening directory:", err)
	}

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	file.Close()

	tpl.ExecuteTemplate(w, "index.html", list)
}
