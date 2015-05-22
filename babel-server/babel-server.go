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
	bms     string
	library *Devices
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
	)

	flag.Parse()
	if *lib == "" || *bms == "" {
		fmt.Fprintln(os.Stderr, "Missing library (-lib) and/or bms flag (-bms)")
		fmt.Fprintln(os.Stderr, `Usage:
  go-adder [flags]
Flags:`)
		flag.PrintDefaults()

		os.Exit(1)
	}

	//var devices Devices
	devices := getLibrary(*lib)

	context := &appContext{bms: *bms, library: devices}
	router := NewRouter(context)
	fmt.Println("Babel-Server has started...")
	log.Fatal(http.ListenAndServe(":"+*port, router))
}

func getLibrary(filename string) *Devices {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var d Devices
	if err = json.Unmarshal(file, &d); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		os.Exit(1)
	}
	return &d
}
