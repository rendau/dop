package httpc

type HttpC interface {
	GetOptions() OptionsSt
	Send(opts OptionsSt) (*RespSt, error)
	SendJson(opts OptionsSt) (*RespSt, error)
	SendRecvJson(opts OptionsSt) (*RespSt, error)
	SendJsonRecvJson(opts OptionsSt) (*RespSt, error)
}
