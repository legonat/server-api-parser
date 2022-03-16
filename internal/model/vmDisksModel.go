package model

type DiskDb struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	Vm   string `json:"vm"`
}

type Disk struct {
	Id int `json:"id"`
	DiskDb
}

type Vm struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

type VmDiscs struct {
	Vm
	Discs []Disk
}

type DiskResults struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Disk `json:"results"`
}

type VmResults struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Vm   `json:"results"`
}
