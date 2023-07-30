package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

func load(filepath string) *image.NRGBA {
	imageFile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		log.Fatal(err)
	}
	return img.(*image.NRGBA)
}

func genFileName(incomingName string) (string, error) {
	r, _ := regexp.Compile("png$|jpeg$")
	extension := r.FindString(incomingName)
	if extension == "" {
		return "", errors.New("no valid file extension found")
	}
	hasher := sha1.New()
	nowString := fmt.Sprint(time.Now())
	hasher.Write([]byte(nowString))
	hashedName := hex.EncodeToString(hasher.Sum(nil))
	return hashedName[:10] + "." + extension, nil
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "400, Bad request", http.StatusBadRequest)
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		fmt.Fprintf(w, "ParseMultipartForm() err: %v", err)
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Fprintf(w, "FormFile() err: %v", err)
	}
	defer image.Close()
	fileName, err := genFileName(handler.Filename)
	if err != nil {
		fmt.Fprint(w, err)
	}
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
	http.ListenAndServe(":8080", nil)
}
