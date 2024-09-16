package main

import (
	"fmt"
	"io"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "No route specified, please specify a valid endpoint.")
}

func main() {
	fmt.Println("Hello World!")
}
