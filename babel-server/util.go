package main

import (
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

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsPoint(s []Point, name string) bool {
	for _, a := range s {
		if a.Name == name {
			return true
		}
	}
	return false
}
