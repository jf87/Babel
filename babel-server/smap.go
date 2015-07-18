package main

import "encoding/json"

type SmapBacnet struct {
	Value json.RawMessage
}

type SmapReading struct {
	Resource string      `json:resource`
	Readings [][]float64 `json:readings`
	UUID     string      `json:"uuid"`
}

type SmapMeta struct {
	UUID       string    `json:"uuid"`
	Properties Property  `json:",omitempty"`
	Path       string    `json:",omitempty"` //use path as resource path
	Metadata   Metadatum `json:",omitempty"`
	Readings   []string  `json:"Readings"`
	//Actuator   Actuator  `json:",omitempty"`
}

type Property struct {
	Timezone      string `json:",omitempty"`
	UnitofMeasure string `json:",omitempty"`
	ReadingType   string `json:",omitempty"`
}

type Metadatum struct {
	SourceName string      `json:",omitempty"`
	Instrument Instrument  `json:",omitempty"`
	Location   Location    `json:",omitempty"`
	Extra      interface{} `json:",omitempty"` //NOTE using interface because extra has no pre-known structure
}

type Instrument struct {
	Manufacturer   string `json:",omitempty"`
	Model          string `json:",omitempty"`
	SamplingPeriod string `json:",omitempty"`
}

type Location struct {
	Building    string `json:",omitempty"`
	Campus      string `json:",omitempty"`
	Floor       string `json:",omitempty"`
	Section     string `json:",omitempty"`
	Room        string `json:",omitempty"`
	Coordinates string `json:", omitempty"`
}

type Actuator struct {
	States   string `json:",omitempty"`
	Values   string `json:",omitempty"`
	Model    string `json:",omitempty"`
	MinValue string `json:",omitempty"`
	MaxValue string `json:",omitempty"`
}

type SmapMetas []SmapMeta
