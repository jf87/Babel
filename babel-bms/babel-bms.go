package main

import (
	"database/sql"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type appContext struct {
	db *sql.DB
}

type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

var actuators [1000]int

func init() {
}

func randomize() {

	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for i, _ := range actuators {
			actuators[i] = r.Intn(10)
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

func main() {
	var (
		port = flag.String("port", "8889", "Port to listen on (optional)")
	)

	flag.Parse()

	go randomize()
	context := &appContext{}
	router := NewRouter(context)
	log.Fatal(http.ListenAndServe(":"+*port, router))

}
