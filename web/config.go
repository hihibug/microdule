package web

type Config struct {
	Mode       string `json:"mode" yaml:"mode"`
	LogColType bool   `json:"LogColType" yaml:"LogColType"`
	LogPath    string `json:"LogPath" yaml:"LogPath"`
	UseHtml    bool   `json:"UseHtml" yaml:"UseHtml"`
	DelimsStr  string `json:"DelimsStr" yaml:"DelimsStr"`
	DelimsEnd  string `json:"DelimsEnd" yaml:"DelimsEnd"`
	StaticPath string `json:"StaticPath" yaml:"StaticPath"`
	TmpPath    string `json:"TmpPath" yaml:"TmpPath"`
	Addr       string `json:"addr" yaml:"addr"`
}
