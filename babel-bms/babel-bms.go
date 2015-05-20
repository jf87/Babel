package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type appContext struct {
	actuators     []int // NOTE we don't need a pointer here because we will pass the struct itself as pointer
	useractuators map[int]int
}

type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func randomize(a *appContext) {

	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i, _ := range a.actuators {
			_, ok := a.useractuators[i]
			if !ok {
				a.actuators[i] = r.Intn(100)
			}
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

func main() {
	var (
		port   = flag.String("port", "8889", "Port to listen on (optional).")
		points = flag.Int("points", 1000, "Number of actuation points that should be created (optional).")
	)

	flag.Parse()
	var a []int
	a = make([]int, *points)
	var u map[int]int
	u = make(map[int]int)
	//u = make([]int, 10) NOTE no need to make as we append() later and append takes care of this

	context := &appContext{actuators: a, useractuators: u}

	go randomize(context)
	router := NewRouter(context)
	log.Fatal(http.ListenAndServe(":"+*port, router))

}
