package main

var devices Devices

type Device struct {
	Kind        string `json:"kind"`
	Location    string `json:"location"`
	Temperature string `json:"temperature,omitempty"`
	Sequence    []Seq  `json:"sequence"`
}

type Seq struct {
	Instruction string  `json:"instruction"`
	Setpoint    float64 `json:"setpoint"`
	Time        int     `json:"time"`
}

type Devices []Device
