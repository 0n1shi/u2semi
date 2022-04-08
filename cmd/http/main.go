package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)

	for k, v := range r.Header {
		fmt.Printf("%s: %s\n", k, strings.Join(v, " "))
	}

	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("\n%s\n", string(body))
}
