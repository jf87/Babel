package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type appContext struct {
	db  *sql.DB
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
	d2 := Device{Kind: "thermostat", Location: "", Sequence: ss2}
	devices = append(devices, d2)

}

func main() {
	var (
		//user = flag.String("user", "", "The database user.")
		//pw   = flag.String("pw", "", "The password to the db.")
		//db   = flag.String("db", "", "The database name.")
		port = flag.String("port", "8888", "Port to listen on (optional)")
	)

	flag.Parse()

	/*
				if len(os.Args) == 1 || *user == "" || *db == "" || *pw == "" {
					fmt.Fprintln(os.Stderr, "Babel-server - Connecting Humans and Buildings")
					fmt.Fprintln(os.Stderr, "Too few arguments.")
					fmt.Fprintln(os.Stderr, `Usage:
			  go-adder [flags]
			Flags:`)
					flag.PrintDefaults()

					os.Exit(1)
				}

		mydb, err := ConnectDB(*db, *user, *pw)
		if err != nil {
			panic(err)
		}
		context := &appContext{
			db: mydb,
		}

		context.InitDB()
	*/
	context := &appContext{bms: "http://localhost:8889"}
	router := NewRouter(context)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
