package httpc

type HttpC interface {
	GetOptions() *OptionsSt
	SetOptions(opts *OptionsSt)
	Send(opts *OptionsSt) (*RespSt, error)
}
