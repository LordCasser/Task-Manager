package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/LordCasser/onefile"
	"github.com/pkg/browser"
)

//go:embed static
var static embed.FS

const Port = "8088"
const StorePath = "./resources/data.json"

func main() {
	err := os.MkdirAll("./resources", 0777)
	//file, err := os.Open(StorePath)
	//defer func() { file.Close() }()
	//if err != nil && os.IsNotExist(err) {
	//	file, err = os.Create(StorePath)
	//	if err != nil {
	//		log.Panicln("create data file error", err)
	//	}
	//}
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

	url := "http://localhost:" + Port
	go func() {
		err = browser.OpenURL(url)
		if err != nil {
			log.Println("browser error:", err)
		}
	}()
	http.Handle("/", handle)
	http.HandleFunc("/store", storeHandle)
	http.HandleFunc("/load", loadHandle)
	log.Println("[+]", "server has start", url)
	err = http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Println(err)
	}

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
