package main

//var devices Devices

type Lib struct {
	Library []Device `json:"library"`
}

type Device struct {
	Kind     string `json:"kind"`
	Location string `json:"location"`
	Model    string `json:"model,omitempty"`
	//Sequence     []Seq   `json:"sequence"`
	Bacnet_types []int  `json:"bacnet_types,omitempty"`
	Value        string `json:"value,omitempty"`
	UUID         string `json:"uuid,omitempty"`
}

type Seq struct {
	Instruction string  `json:"instruction"`
	Setpoint    float64 `json:"setpoint"`
	Time        int     `json:"time"`
}

//type Devices []Device

type Result struct {
	Result string `json:"result"`
}
