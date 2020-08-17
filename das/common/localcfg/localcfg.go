package localcfg

//
// log config
//
type Log struct {
	Level      int    `json:"level"`
	MaxAge     int    `json:"max_age"`
	MaxSize    int    `json:"max_size"`
	FileName   string `json:"file_name"`
	MaxBackups int    `json:"max_backups"`
}

//
// config center
//
type ApolloParam struct {
	Addr       string   `json:"addr"`
	AppId      string   `json:"app_id"`
	CfgServer  string   `json:"cfg_server"`
	Cluster    string   `json:"cluster"`
}

//
// local basic config
//
type LocalCfg struct {
	Log    *Log         `json:"log"`
	Apollo *ApolloParam `json:"apollo"`
}