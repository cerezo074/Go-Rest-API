package main

import (
	"fmt"
	"github/cerezo074/GoAPI/handlers"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc(handlers.UsersPathSlashed, handlers.UsersRotuer)
	http.HandleFunc(handlers.UsersPath, handlers.UsersRotuer)
	http.HandleFunc(handlers.HomePath, handlers.RootHandler)
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
