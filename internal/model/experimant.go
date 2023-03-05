package model

type Experiment struct {
	Status  string `json:"status"`
	Name    string `json:"name"`
	Teacher string `json:"teacher"`
	Address string `json:"address"`
	Date    int64  `json:"date"`
	Section string `json:"section"`
	Week    string `json:"week"`
	Time    string `json:"time"`
}
