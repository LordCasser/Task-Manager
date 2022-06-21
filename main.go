package main

import (
	"embed"
	"github.com/LordCasser/onefile"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static
var static embed.FS

func main() {
	err := os.MkdirAll("./resources", 0777)
	if err != nil {
		log.Println(err)
		return
	}
	overwrite := &onefile.Overwrite{
		Fsys: nil,
		Pair: map[string]string{},
	}
	fsys, _ := fs.Sub(static, "static")
	handle := onefile.New(fsys, overwrite, "index.html")
	http.Handle("/", handle)
	http.HandleFunc("/store", storeHandle)
	http.HandleFunc("/load", loadHandle)
	_ = http.ListenAndServe(":8080", nil)
}

func storeHandle(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("./resources/data.json", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(file, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func loadHandle(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("./resources/data.json", os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
