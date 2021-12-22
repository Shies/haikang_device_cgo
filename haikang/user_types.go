package haikang

type PingParam struct {
	Welcome interface{} `json:"welcome"` // hello world
}

type RealPlayInfo struct {
	SGetIP   string
	LUserID  int
	LChannel int
}
