package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var su = make(chan int)

type Events struct {
	t []time.Time
	v []float64
}

// subscribe to backend FIXME JF this needs to be seriously cleaned up
func checkForSequence(a *appContext, d Device) error {
	fmt.Println("checkForSequence")
	found := false
	p := make([]int, 1000)
	e := make([]Events, 1000)
	var matches []int
	t_total := 0
	//matches := make([]int, 1000)
	seq := make([]float64, len(d.Sequence))
	for i, k := range d.Sequence {
		seq[i] = k.Setpoint
		t_total = t_total + k.Time
	}
	t_total = t_total * 2
	tt := time.Duration(t_total) * time.Second
	t0 := time.Now()
	for !found && (time.Since(t0) < tt) {
		// get current state
		for i := 0; i < 1000; i++ {
			act := strconv.Itoa(i)
			resp, err := http.Get(a.bms + "/api/actuators/" + act)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			r := string(body[0])
			ri, err := strconv.Atoi(r)
			if float64(ri) == seq[p[i]] {
				fmt.Println("FOUND SETPOINT ", ri)
				fmt.Println("SEQUENCE", p[i])
				if e[i].t == nil {
					e[i].t = make([]time.Time, 5)
					e[i].v = make([]float64, 5)
				}
				e[i].t[p[i]] = time.Now()
				e[i].v[p[i]] = float64(ri)
				p[i] += 1
				if p[i] == len(seq) {
					found = true
					matches = append(matches, i)
					fmt.Println("len seq", len(seq))
					fmt.Println("BMS Setpoint is: ", i)
					fmt.Println("FUCKING A")
				}
			}
		}
		time.Sleep(1000 * time.Millisecond)

	}
	if !found {
		fmt.Println("not found :-()")
		if time.Since(t0) > tt {
			fmt.Println("we ran out of time before a match was found")
		}
		su <- 2
	} else {
		su <- 1
	}

	for _, v := range matches {
		fmt.Println(v)
	}
	for _, v := range matches {
		fmt.Printf("Match for Actuator %v:\n", v)
		for i, _ := range seq {
			fmt.Printf("with setpoint: %v at time %v\n", e[v].v[i], e[v].t[i])
		}
		fmt.Printf("\n")

	}

	return nil

}
