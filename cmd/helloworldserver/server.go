package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Hello world!</h>", r.FormValue("name"))
	})

	http.ListenAndServe(":8888", nil)
}
