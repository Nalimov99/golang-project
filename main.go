package main

import (
	"fmt"
	"net/http"
)

func main() {
	h := http.HandlerFunc(Echo)

	if err := http.ListenAndServe("localhost:3020", h); err != nil {
		fmt.Print("Error")
	}
}

func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.Method, " ", r.URL)
}
