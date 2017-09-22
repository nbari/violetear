package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"

	"github.com/nbari/violetear"
)

func hello(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 1000000; i++ {
		math.Pow(36, 89)
	}
	fmt.Fprint(w, "Hello!")
}

func main() {
	router := violetear.New()
	router.HandleFunc("/", hello)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Fatal(http.ListenAndServe(":8080", router))
}
