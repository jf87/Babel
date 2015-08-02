package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var tries int

type Point struct {
	Name         string
	ID           int
	Count        int
	Done         bool
	TimeReadings map[int]float64
	Readings     []Reading
	Twins        []string
}

type Reading struct {
	T int
	V float64
}

var id int
var files []string

func addFiles(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, "txt") {
		files = append(files, path)
	}
	return nil
}

func main() {
	tries = 60

	flag.Parse()
	root := flag.Arg(0)

	err := filepath.Walk(root, addFiles)
	if err != nil {
		fmt.Println("ERROR")
		log.Fatal(err)
	}
	var csvs []*csv.Reader
	for _, v := range files {
		file, err := os.Open(v)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		c := csv.NewReader(file)
		csvs = append(csvs, c)
	}
	var points map[string]Point
	points = make(map[string]Point)
	for i := 0; i <= 1008; i++ {
		id = 0
		readValues(csvs, points)
		analyze(points)
		result(points, i)
	}

}

func readValues(csvs []*csv.Reader, points map[string]Point) {

	for i := 0; i < tries; i++ {
		newValue(csvs, points)
	}
}

func analyze(points map[string]Point) {
	for key, point := range points {
		//fmt.Println(point.Name)
		//fmt.Println("new point")
		for _, r := range point.Readings {
			var newTwins []string
			for _, twin := range point.Twins {
				stillTwin := false
				val, ok := points[twin]
				if val.Done {
					continue
				}
				if ok {
					vall, okk := val.TimeReadings[r.T]
					if okk {
						if vall == r.V {
							//same value in same time
							stillTwin = true
						} else {
							stillTwin = false
							// also not a twin anymore
						}
					} else {
						stillTwin = false
						// twin is not a twin anymore
					}
				} else {
					stillTwin = false
				}
				if stillTwin {
					newTwins = append(newTwins, twin)
				}
			}
			point.Twins = newTwins
			if len(newTwins) < 2 {
				//	log.Fatal("DONE")
				point.Done = true
			} else {
				point.Count++
			}
			points[key] = point
		}
	}

}

func result(points map[string]Point, i int) {
	//fmt.Printf("ID, UUID, Time (s), Siblings Left\n")
	no_points := 0
	not_matched := 0
	for _, point := range points {
		no_points++
		l := len(point.Twins)
		//	t := point.Count * 10
		//	n := point.Name[len(point.Name)-40 : len(point.Name)-4]
		if l != 0 {
			not_matched++
			//		t = tries * 10
		}
		//fmt.Printf("%v, %v, %v, %v\n", point.ID, n, t, l)
	}
	//fmt.Printf("Total, Matched, Not Matched\n")
	var perc float64
	perc = (float64(no_points-not_matched) / float64(no_points))

	fmt.Printf("%v, %v, %v, %v, %.2f\n", tries*10*(i+1), no_points, (no_points - not_matched), not_matched, perc)

}

func newValue(csvs []*csv.Reader, points map[string]Point) {

	for k, c := range csvs {
		rec, err := c.Read()
		if err != nil {
			//fmt.Println("end of line")
			//log.Fatal(err)
			continue
		}
		tmp, err := strconv.ParseFloat(strings.TrimSpace(rec[0]), 64)
		if err != nil {
			fmt.Println(err)
		}
		t := int(tmp)
		t = t / 10000

		v, err := strconv.ParseFloat(strings.TrimSpace(rec[1]), 64)
		if err != nil {
			fmt.Println(err)
			//log.Fatal(err)
		}

		point, ok := points[files[k]]
		if !ok {
			var p Point
			var r Reading
			p.Count = 1
			p.ID = id
			id++
			p.Name = files[k]
			r.T = t
			r.V = v
			p.TimeReadings = make(map[int]float64)
			p.TimeReadings[t] = v
			p.Twins = files
			p.Readings = append(p.Readings, r)
			points[p.Name] = p
		} else {

			var r Reading
			r.T = t
			r.V = v
			point.TimeReadings[t] = v
			point.Readings = append(point.Readings, r)
			points[files[k]] = point
		}
	}

	/*
		for _, s := range scanners {

			s.Scan()
			fmt.Println(s.Text())
			if err := s.Err(); err != nil {
				log.Fatal(err)
			}
		}
	*/
}
