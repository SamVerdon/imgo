package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func genFileName() string {
	hasher := sha1.New()
	nowString := fmt.Sprint(time.Now())
	hasher.Write([]byte(nowString))
	hashedName := hex.EncodeToString(hasher.Sum(nil))
	return hashedName[:10] + ".png"
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "400, Bad request", http.StatusBadRequest)
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		fmt.Fprintf(w, "ParseMultipartForm() err: %v", err)
	}
	image, _, err := r.FormFile("image")
	if err != nil {
		fmt.Fprintf(w, "FormFile() err: %v", err)
	}
	defer image.Close()
	fileName := genFileName()
	f, err := os.OpenFile("./assets/images/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(w, "OpenFile() err: %v", err)
	}
	defer f.Close()
	fmt.Fprint(w, r.Host+"/images/"+fileName)
	io.Copy(f, image)
}

func greet(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "<h1>Hello World! %s</h1>", time.Now())
	fmt.Fprintf(w, "<title>Test</title>")
}

func main() {
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/", greet)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./assets/images"))))
	http.ListenAndServe(":80", nil)
}
