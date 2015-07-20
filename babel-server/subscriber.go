package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var su = make(chan int)
var sync = make(chan int)

var active bool

type Events struct {
	t []time.Time
	v []float64
}

func monitorBMS(a *appContext, d Device) error {
	log.Printf(
		"%s\t%s\t%s\t%v",
		"START",
		"monitorBMS",
		"started function",
		time.Since(start).Nanoseconds(),
	)
	fmt.Println("monitorBMS")
	sync <- 1
	active = true
	fmt.Println("now active")
	time.Sleep(15000 * time.Millisecond)
	tt := time.Duration(15) * time.Second
	t0 := time.Now()
	var br BabelReadings
	br = make(map[string]BabelReading)
	log.Printf(
		"%s\t%s\t%s\t%v",
		"START",
		"monitorBMS",
		"started to query BMS",
		time.Since(start).Nanoseconds(),
	)
	i := 0
	for time.Since(t0) < tt { //for now just loop until time is over
		resp, err := http.Get(a.bms)
		fmt.Println(resp)
		fmt.Println(err)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		br, err = decodeSmapReadings(body, br)
		i++
		time.Sleep(500 * time.Millisecond)
	}
	active = false
	fmt.Println("now !active")
	log.Printf(
		"%s\t%s\t%s\t%v",
		"STOP",
		"monitorBMS",
		"stopped to query BMS",
		time.Since(start).Nanoseconds(),
	)
	log.Printf(
		"%s\t%s\t%v\t%v",
		"STATE",
		"GETRequests send to smap",
		i,
		time.Since(start).Nanoseconds(),
	)
	matches, err := checkForValue(a, d, br)
	fmt.Printf("matches %v", matches)
	j, _ := reducePoints(a, matches) //changes reduced points in a

	if matches == nil || len(matches) == 0 {
		log.Printf(
			"%s\t%s\t%s\t%v",
			"STOP",
			"monitorBMS",
			"stopped monitorBMS - 0 matches",
			time.Since(start).Nanoseconds(),
		)
		su <- 0
	} else if len(matches) > 1 {
		log.Printf(
			"%s\t%s\t%s\t%v",
			"STOP",
			"monitorBMS",
			"stopped monitorBMS, matches: "+strconv.Itoa(j),
			time.Since(start).Nanoseconds(),
		)
		su <- j
	} else {
		log.Printf(
			"%s\t%s\t%s\t%v",
			"STOP",
			"monitorBMS",
			"stopped monitorBMS - 1 matches",
			time.Since(start).Nanoseconds(),
		)
		su <- 1
	}
	return err
}

//check for a value that the user reads
func checkForValue(a *appContext, d Device, br BabelReadings) (BabelReadings, error) {
	log.Printf(
		"%s\t%s\t%s\t%v",
		"START",
		"checkForValue",
		"started checkForValue",
		time.Since(start).Nanoseconds(),
	)
	fmt.Println("checkForValue")
	match := false
	var br_new BabelReadings
	br_new = make(map[string]BabelReading)
	for _, v := range br {
		for _, va := range v.Readings {
			fl, err := strconv.ParseFloat(d.Value, 64)
			if err != nil {
				fmt.Println("ERRORRRR")
				return br_new, err
			}
			if va[1] == fl {
				match = true
			}
		}
		if match {
			br_new[v.PointName] = v
			match = false
		}
	}
	log.Printf(
		"%s\t%s\t%s\t%v",
		"STOP",
		"checkForValue",
		"stopped function",
		time.Since(start).Nanoseconds(),
	)
	return br_new, nil
}

// reduce points that need to be queried, based on prior readings
func reducePoints(a *appContext, br BabelReadings) (int, error) {
	log.Printf(
		"%s\t%s\t%s\t%v",
		"START",
		"reducePoints",
		"started function",
		time.Since(start).Nanoseconds(),
	)

	fmt.Println("reducePoints")
	//create index
	var prr Points
	i := 0
	j := 0
	for _, v := range a.points_reduced {
		var o Objects
		for _, va := range v.Objs {
			i++
			_, ok := br[va.Name]
			if ok {
				j++
				o = append(o, va)
			}
		}
		if len(o) > 0 {
			v.Objs = o
			prr = append(prr, v)
		}
	}
	a.points_reduced = prr
	log.Printf(
		"%s\t%s\t%v\t%v",
		"STATE",
		"bms_points",
		i,
		time.Since(start).Nanoseconds(),
	)
	fmt.Printf("reduced to %v\n", prr)
	log.Printf(
		"%s\t%s\t%v\t%v",
		"STATE",
		"bms_points",
		j,
		time.Since(start).Nanoseconds(),
	)
	log.Printf(
		"%s\t%s\t%s\t%v",
		"STOP",
		"reducePoints",
		"stopped function",
		time.Since(start).Nanoseconds(),
	)
	return j, nil
}

func decodeSmapReadings(jsonRaw []byte, br BabelReadings) (BabelReadings, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(jsonRaw, &obj); err != nil {
		return br, err
	}
	for _, v := range obj {
		var smapReading SmapReading
		//fmt.Println("value %v\n", value)
		if err := json.Unmarshal(v, &smapReading); err != nil {
			return br, err
		}
		if smapReading.UUID != "" && smapReading.Metadata.PointName != "" {
			if val, ok := br[smapReading.Metadata.PointName]; ok {
				val.Readings = append(val.Readings, smapReading.Readings[0])
				br[smapReading.Metadata.PointName] = val
			} else {
				fmt.Printf("%v, %v, %v", smapReading.UUID, smapReading.Readings, smapReading.Metadata.PointName)
				var b BabelReading
				b.UUID = smapReading.UUID
				b.Readings = smapReading.Readings
				b.PointName = smapReading.Metadata.PointName
				br[b.PointName] = b
			}
		}
	}
	return br, nil
}

/*
func checkForSequence(a *appContext, d Device) error {
	fmt.Println("checkForSequence")
	//provide initial point list turn back on for point reduction
	sync <- 1
	active = true
	fmt.Print("now active")
	time.Sleep(2000 * time.Millisecond)
	//keep track of time dependent on provided sequence
	t_total := 0
	for _, k := range d.Sequence {
		t_total = t_total + k.Time
	}
	t_total = t_total * 2 //twice the time seems like a good timeout
	tt := time.Duration(t_total) * time.Second
	t0 := time.Now()

	var readings []SmapReading //all readings
	var points map[string]int  //points that we need to consider
	points = make(map[string]int)
	for time.Since(t0) < tt { //for now just loop until time is over
		resp, err := http.Get(a.bms)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		readings, points, err = DecodeSmapJson(body, readings, points)
		if err != nil {
			fmt.Printf("ERROR %v\n", err)
		}
		fmt.Println(readings)
		//reducePoints(a, d, readings)
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

//use this for now for thermostat
func checkForIncrease(a *appContext, d Device) error {
	fmt.Println("checkForIncrease")
	sync <- 1
	active = true
	fmt.Print("now active")
	time.Sleep(2000 * time.Millisecond)
	t_total := d.Sequence[0].Time
	t_total = t_total * 2 //twice the time seems like a good timeout
	tt := time.Duration(t_total) * time.Second
	t0 := time.Now()

	var readings []SmapReading
	//var points map[string]int
	//points = make(map[string]int)
	for time.Since(t0) < tt { //for now just loop until time is over
		resp, err := http.Get(a.bms)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		//var matches map[string]int

		readings, err = DecodeSmapJson2(body, readings)
		if err != nil {
			fmt.Printf("ERROR %v\n", err)
		}
		fmt.Println(readings)
		//sync <- 1
		time.Sleep(500 * time.Millisecond)
	}
	active = false

	//check received readings for user-reading and reduce points
	matches, err := checkForValue(a, d, readings)
	//reducePoints(a, matches) FIXME take back

	sync <- 1
	active = true
	fmt.Print("now active")
	time.Sleep(2000 * time.Millisecond)
	t_total = 0
	for _, k := range d.Sequence[1:] {
		t_total = t_total + k.Time
	}
	t_total = t_total * 2 //twice the time seems like a good timeout
	tt = time.Duration(t_total) * time.Second
	t0 = time.Now()

	var readings2 []SmapReading
	var points2 map[string]int
	points2 = make(map[string]int)

	for time.Since(t0) < tt { //for now just loop until time is over
		resp, err := http.Get(a.bms)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		//var matches map[string]int

		readings2, points2, err = DecodeSmapJson(body, readings2, points2)
		if err != nil {
			fmt.Printf("ERROR %v\n", err)
		}
		fmt.Println(readings)
		//sync <- 1
		time.Sleep(500 * time.Millisecond)
	}
	active = false

	matches, err = findIncrease(readings, points2, d)
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

func DecodeSmapJson2(jsonRaw []byte, readings []SmapReading) ([]SmapReading, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(jsonRaw, &obj); err != nil {
		return readings, err
	}
	for key, value := range obj {
		var smapReading SmapReading

		fmt.Printf("key %v\n", key)
		//fmt.Println("value %v\n", value)
		if err := json.Unmarshal(value, &smapReading); err != nil {
			return readings, err
		}
		if smapReading.UUID != "" {
			if val, ok := points[smapReading.UUID]; ok {
				fmt.Println("just new reading")
				readings[val].Readings = append(readings[val].Readings, smapReading.Readings[0])

			}

		}
		if val, ok := points[smapReading.UUID]; ok {
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
	return readings, nil
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

func findIncrease(readings []SmapReading, points map[string]int, d Device) ([]SmapReading, error) {
	fmt.Println("findIncrease")
	var matches []SmapReading
	match := false

	for k, v := range readings {
		for _, va := range v.Readings {
			if va[1] > (d.Value + 0.1) {
				match = true
			}
		}
		if match {
			matches = append(matches, readings[k])
		}
		match = false
	}

	return matches, nil
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
*/
