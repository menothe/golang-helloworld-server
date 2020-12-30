package main

import (
    "fmt"
    "net/http"
    "os"
)

func helloWorld(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Hello World")
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func main() {
    http.HandleFunc("/", helloWorld)
    http.ListenAndServe(getPort(), nil)
}
