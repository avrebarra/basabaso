package server

type RespCode string

type Resp struct {
	Code    RespCode    `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (e *Resp) AddMsg(msg string) {
	e.Message += ": " + msg
}

func (e Resp) Normalize() Resp {
	// send empty object instead of nil
	if e.Data == nil {
		e.Data = map[string]interface{}{}
	}
	return e
}

// ***

type RespKind int

const (
	RCSuccess RespKind = iota
	RCPending
	RCUnexpected
)

var RespPresets map[RespKind]Resp = map[RespKind]Resp{
	RCUnexpected: {Code: "0", Message: "unexpected failure"},
	RCSuccess:    {Code: "1", Message: "operation success"},
	RCPending:    {Code: "2", Message: "operation pending"},
}
