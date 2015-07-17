package main

type Point struct {
	Desc  string   `json:"desc"`
	Name  string   `json:"name"`
	Objs  []Object `json:"objs"`
	Props Prop2    `json:"props,"`
}

type Object struct {
	Data_type int    `json:"data_type"`
	Desc      string `json:"desc"`
	Name      string `json:"name"`
	Props     Prop   `json:"props"`
	Unit      int    `json:"unit"`
}

type Prop struct {
	Instance int    `json:"instance"`
	Type     int    `json:"type"`
	Type_str string `json:"type_str"`
}

type Prop2 struct {
	Adr       []int `json:"adr"`
	Device_id int   `json:"device_id"`
	Mac       []int `json:"mac"`
	Max_apdu  int   `json:"max_apdu"`
	Net       int   `json:"net"`
}

type Points []Point
