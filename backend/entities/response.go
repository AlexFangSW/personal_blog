package entities

type Ret struct {
	Error  any `json:"error"`
	Status int `json:"status"`
	Msg    any `json:"msg"`
}

func NewRet(err any, status int, msg any) *Ret {
	return &Ret{
		Error:  err,
		Status: status,
		Msg:    msg,
	}
}
