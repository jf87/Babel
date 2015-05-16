package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// subscribe to backend
//http://192.168.0.101/api/newdeveloper/lights/2/
func checkForSequence(a *appContext, d Device) error {

	p := make([]float64, 1000)
	seq := make([]float64, len(d.Sequence))
	for i, k := range d.Sequence {
		seq[i] = k.Setpoint
	}

	// get current state
	for i := 0; i < 1000; i++ {
		act := strconv.Itoa(i)
		resp, err := http.Get(a.bms + "/api/actuators/" + act)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		r := int(body[0])
		p[i] = float64(r)
		fmt.Println(r)
	}

	for _, k := range seq {

		for i := 0; i < 1000; i++ {
			act := strconv.Itoa(i)
			resp, err := http.Get(a.bms + "/api/actuators/" + act)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			r := float64(int(body[0]))
			if r == k {

				p[i] = r
			}
			fmt.Println(r)

		}
	}

	return nil

}
