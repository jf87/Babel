package main

import (
	"flag"
	"fmt"
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
				if i%3 == 0 {
					a.actuators[i] = r.Intn(100)
					//a.actuators[i] = 77
				}
			}
		}
		time.Sleep(1000 * time.Millisecond)
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, _ := range a {
		a[i] = r.Intn(100)
	}

	context := &appContext{actuators: a, useractuators: u}

	fmt.Println("randomize")
	go randomize(context)
	fmt.Println("starting router")
	router := NewRouter(context)
	log.Fatal(http.ListenAndServe(":"+*port, router))

}
