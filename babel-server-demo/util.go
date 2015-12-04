package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/twinj/uuid"
)

func CreateUUID() string {
	u := uuid.NewV4()
	return u.String()
}

func ValidateUUID(aUUID string) (uuid.UUID, error) {
	return uuid.ParseUUID(aUUID)
}

func epochToTime(epoch float64) (time.Time, error) {
	epochi := int64(epoch)
	return time.Unix(0, epochi*int64(time.Millisecond)), nil
}

/*
func fakeActuation(a *appContext, d Device) {
	fmt.Println("fakeActuation")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	url := a.bms + "/api/actuators/" + strconv.Itoa(r.Intn(1000))

	for _, s := range d.Sequence {
		fmt.Println("URL:>", url)
		payload := `{"Value":` + strconv.Itoa(int(s.Setpoint)) + `}`
		fmt.Println("PAYLOAD:>", payload)

		var jsonStr = []byte(payload)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		//time.Sleep(time.Duration(s.Time * 1000))
		time.Sleep(5000 * time.Millisecond)
	}

}

*/
func WriteCSV(data [][]string, create bool) {
	var err error
	var file *os.File
	if create {
		file, err = os.Create("result.csv")
	} else {
		file, err = os.OpenFile("result.csv", os.O_RDWR|os.O_APPEND, 0)
	}

	if err != nil {
		log.Fatal("Cannot create file ", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Fatal("Cannot write to file ", err)
		}
	}

	writer.Flush()

}
