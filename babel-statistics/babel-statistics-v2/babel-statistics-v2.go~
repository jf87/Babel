package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var files_processed []string
var csvs []*csv.Reader
var INTERVALL int
var START int
var END int
var files []string
var pos [1024][2]int

func addFiles(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, "txt") {
		files = append(files, path)
	}
	return nil
}

func init() {
	INTERVALL = 3600
	START = 133663933
	END = 1391895250
}

func main() {

	flag.Parse()
	root := flag.Arg(0)

	err := filepath.Walk(root, addFiles)
	if err != nil {
		log.Fatal(err)
	}

	preProcess()

}

func preProcess() {
	var tooFar bool
	//	var startdone, enddone bool
	for _, v := range files {
		file, err := os.Open(v)
		defer file.Close()
		if err != nil {
			fmt.Println("too many open files")
			log.Fatal(err)
		}
		c := csv.NewReader(file)
		csvs = append(csvs, c)
	}
	//	startdone = false
	//enddone = false
	for {
		for f_no, c := range csvs {
			tooFar = false
			recs, err := c.ReadAll()
			if err != nil {
				log.Fatal(err)
			}
			for k, v := range recs {
				//tooFar = false
				tf, err := strconv.ParseFloat(strings.TrimSpace(v[0]), 64)
				if err != nil {
					log.Fatal(err)
				}
				ti := int(tf / 10000)
				if ti == START {
					pos[f_no][0] = k
					fmt.Printf("time match: %v\n", ti)
					break
				} else if ti > START {
					START = ti
					pos[f_no][0] = k
					tooFar = true
					fmt.Printf("too far: t=%v, p=%v\n", ti, k)
					break
				} else {
					pos[f_no][0] = -1
				}
			}
			if tooFar {
				fmt.Println("breaking out loop")
				//tooFar = false
				break
			}
		}
		fmt.Printf("loop finished: toofar:%v, current t: %v\n", tooFar, START)
		if !tooFar {
			fmt.Println("finished")
			break
		}
	}
	csvs = nil
	for k, v := range files {
		file, err := os.Open(v)
		defer file.Close()
		if err != nil {
			fmt.Println("too many open files")
			log.Fatal(err)
		}
		sc := bufio.NewScanner(file)
		line_no := pos[k][0]
		i := 0
		write := false
		if line_no != -1 {
			f, err := os.Create(v + "-new.txt")
			//f, err := ioutil.WriteFile(v+"-new.txt", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			for sc.Scan() {
				if i == line_no {
					write = true
				}
				if write {
					if _, err = f.WriteString(sc.Text() + "\n"); err != nil {
						log.Fatal(err)
					}
				}
				i++
			}
		}
	}
}

/*
		c := csv.NewReader(file)
		csvs = append(csvs, c)
	}
	for k, c := range csvs {

		file, err := os.Create(files[k] + "-new.txt")
		if err != nil {
			log.Fatal("Cannot create file ", err)
		}
		defer file.Close()

		recs, err := c.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		var recs_f [][]float64
		for t, r := range recs {

		}
		fmt.Println(pos)
		fmt.Println(pos[k][0])
		recs_cut := append(recs[:pos[k][0]], recs[(pos[k][0])+1:]...)
		w := csv.NewWriter(file)
		w.WriteAll(recs_cut) // calls Flush internally
*/
