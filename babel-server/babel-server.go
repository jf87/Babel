package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type appContext struct {
	bms string
}

type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func init() {
	s1 := Seq{Instruction: "Turn the light on.", Setpoint: 1, Time: 5}
	s2 := Seq{Instruction: "Turn the light off.", Setpoint: 0, Time: 5}
	s3 := Seq{Instruction: "Turn the light on.", Setpoint: 1, Time: 5}
	var ss []Seq
	ss = append(ss, s1)
	ss = append(ss, s2)
	ss = append(ss, s3)
	d := Device{Kind: "light", Location: "", Sequence: ss}
	devices = append(devices, d)

	s4 := Seq{Instruction: "Set the setpoint to 42.", Setpoint: 42, Time: 30}
	s5 := Seq{Instruction: "Set the setpoint to 23.", Setpoint: 23, Time: 30}
	s6 := Seq{Instruction: "Set the setpoint to 16.", Setpoint: 16, Time: 30}
	var ss2 []Seq
	ss2 = append(ss2, s4)
	ss2 = append(ss2, s5)
	ss2 = append(ss2, s6)
	d2 := Device{Kind: "thermostat", Location: "", Temperature: " ", Sequence: ss2}
	devices = append(devices, d2)

}

func main() {
	var (
		port = flag.String("port", "8888", "Port to listen on (optional)")
	)

	flag.Parse()

	context := &appContext{bms: "http://localhost:8889"}
	router := NewRouter(context)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
