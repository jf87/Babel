package main

import ()

var devices Devices

type Device struct {
	Kind     string `json:"kind"`
	Location string `json:"location"`
	Sequence []Seq  `json:"actuators"`
}

type Seq struct {
	Instruction string  `json:"instruction"`
	Setpoint    float64 `json: "setpoint"`
	Time        int     `json:"time"`
}

type Devices []Device
