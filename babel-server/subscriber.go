package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var su = make(chan int)
var sync = make(chan int)

var active bool

type Events struct {
	t []time.Time
	v []float64
}

// subscribe to backend FIXME JF this needs to be seriously cleaned up
func checkForSequence(a *appContext, d Device) error {
	fmt.Println("checkForSequence")
	//provide initial point list turn back on for point reduction
	sync <- 1
	active = true
	fmt.Print("now active")
	time.Sleep(2000 * time.Millisecond)
	//keep track of time dependent on provided sequence
	t_total := 0
	seq := make([]float64, len(d.Sequence))
	for i, k := range d.Sequence {
		seq[i] = k.Setpoint
		t_total = t_total + k.Time
	}
	t_total = t_total * 2 //twice the time seems like a good timeout
	tt := time.Duration(t_total) * time.Second
	t0 := time.Now()

	var readings []SmapReading
	var points map[string]int
	points = make(map[string]int)
    fmt.Println("before for loop")
	for time.Since(t0) < tt { //for now just loop until time is over
		resp, err := http.Get(a.bms)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
        fmt.Println("now got BMS values")
		//var matches map[string]int

		readings, points, err = DecodeSmapJson(body, readings, points)
		if err != nil {
			fmt.Printf("ERROR %v\n", err)
		}
		fmt.Println(readings)
		reducePoints(a, d, readings)
		//sync <- 1
		time.Sleep(500 * time.Millisecond)
	}
	active = false

	matches, err := findMatch(readings, points, d)
	if err != nil {
		fmt.Printf("ERROR %v\n", err)
	}

	fmt.Printf("matches %v\n", matches)

	if matches == nil {
		su <- 2
	} else {
		su <- 1
	}

	for _, v := range matches {
		fmt.Println(v)
	}
	for _, v := range matches {
		fmt.Printf("Match for Actuator %v:\n", v)
	}
	return nil
}

func DecodeSmapJson(jsonRaw []byte, readings []SmapReading, points map[string]int) ([]SmapReading, map[string]int, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(jsonRaw, &obj); err != nil {
		return readings, points, err
	}
	for key, value := range obj {
		var smapReading SmapReading

		fmt.Printf("key %v\n", key)
		//fmt.Println("value %v\n", value)
		if err := json.Unmarshal(value, &smapReading); err != nil {
			return readings, points, err
		}
		if val, ok := points[key]; ok {
			fmt.Println("just new value")
			readings[val].Readings = append(readings[val].Readings, smapReading.Readings[0])
		} else {
			fmt.Println("new point")
			if smapReading.UUID != "" {
				smapReading.Resource = key
				fmt.Printf("smapReading %v\n", smapReading)
				readings = append(readings, smapReading)
				points[key] = len(readings) - 1
			}
		}
	}
	return readings, points, nil
}

// reduce points that need to be queried, based on prior readings
func reducePoints(a *appContext, d Device, readings []SmapReading) {
	fmt.Println("NOT IMPLEMENTED")
}

func findMatch(readings []SmapReading, points map[string]int, d Device) ([]SmapReading, error) {
	var matches []SmapReading
	p := make([]int, len(readings))
	seq := make([]float64, len(d.Sequence))
	for i, k := range d.Sequence {
		seq[i] = k.Setpoint
	}
	for k, v := range readings {
		for _, va := range v.Readings {
			//time := v.Readings[0][0]
			set := va[1]
			fmt.Printf("setpoint %v\n", set)
			if !(p[k] >= len(seq)) {
				if set == seq[p[k]] {
					fmt.Println("found one stepoint matching")
					p[k] += 1
					if p[k] == len(seq) {
						fmt.Println("MATCH")
						fmt.Println(readings[k])

						matches = append(matches, readings[k])
					}

				} else {
					if p[k] != 0 {
						if set != seq[p[k]-1] {
							fmt.Println("pattern does not match, resetting counter...")
							p[k] = 0
						} else {
							fmt.Println("same value read, good for now")
						}
					}
				}
			}
		}
	}
	return matches, nil
}
