package main

//var devices Devices

type Lib struct {
	Library []Device `json:"library"`
}

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

//type Devices []Device

type Suc struct {
	Success bool `json:"success"`
}
