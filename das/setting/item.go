package setting

type Listener struct {
	Port uint16 `json:"port"`
}

type RegisterCentre struct {
	Ip   string `json:"ip"`
	Port uint16 `json:"port"`
}