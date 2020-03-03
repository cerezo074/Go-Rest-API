package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	homePath string = "/"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != homePath {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Asset not found \n"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Running API V1\n"))
}

func main() {
	http.HandleFunc(homePath, rootHandler)
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
