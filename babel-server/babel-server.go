package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type appContext struct {
	bms            string
	library        *Lib
	points         *Points
	points_reduced Points
}

type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func main() {
	var (
		port = flag.String("port", "8888", "Port to listen on (optional)")
		lib  = flag.String("lib", "", "Library with definition of types and patterns.")
		bms  = flag.String("bms", "", "Url to BMS Points")
		db   = flag.String("db", "", "Bacnet DB file")
	)

	flag.Parse()
	if *lib == "" || *bms == "" {
		fmt.Fprintln(os.Stderr, "Missing library (-lib), Bacnet DB file (-db) and/or bms flag (-bms)")
		fmt.Fprintln(os.Stderr, `Usage:
      babel-server [flags]
Flags:`)
		flag.PrintDefaults()

		os.Exit(1)
	}

	//var devices Devices
	devices := getLibrary(*lib)
	points := getPoints(*db)
	//testSmap("data/smap.json")

	context := &appContext{bms: *bms, library: devices, points: points}
	router := NewRouter(context)
	fmt.Println("Babel-Server has started...")
	log.Fatal(http.ListenAndServe(":"+*port, router))
}

func getLibrary(filename string) *Lib {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var l Lib
	if err = json.Unmarshal(file, &l); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		os.Exit(1)
	}
	return &l
}

func getPoints(filename string) *Points {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var p Points
	if err = json.Unmarshal(file, &p); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		os.Exit(1)
	}
	return &p
}

/*
func testSmap(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var readings []SmapReading
	var points map[string]int
	points = make(map[string]int)

	//var matches map[string]int

	readings, points, err = DecodeSmapJson(file, readings, points)
	if err != nil {
		fmt.Printf("err %v\n", err)
	}
	fmt.Println(readings)
}
*/
