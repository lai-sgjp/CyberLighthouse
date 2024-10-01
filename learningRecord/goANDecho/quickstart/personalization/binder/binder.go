package main

import (
	"fmt"
	"log"
	"net/http"
)

func main () {
	http.HandleFunc("/",func(w http.ResponseWriter, req *http.Request) {
		req.ParseMultipartForm(32 << 20)

		data := map[string]interface{}{
			"form": req.form,
			"post_form": req.PostForm,
		}

		fmt.Fprintln(w,data)
	})
	log.Fatal(http.ListenAndServe(":2020",nil))
}