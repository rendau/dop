package smss

import (
	"github.com/rendau/dop/adapters/client/httpc"
)

type St struct {
	httpc httpc.HttpC
}

func New(httpc httpc.HttpC) *St {
	return &St{
		httpc: httpc,
	}
}

func (s *St) Send(phone string, msg string) bool {
	return s.send(&SendReqSt{
		To:   phone,
		Text: msg,
		Sync: true,
	})
}

func (s *St) SendAsync(phone string, msg string) bool {
	return s.send(&SendReqSt{
		To:   phone,
		Text: msg,
		Sync: false,
	})
}

func (s *St) send(req *SendReqSt) bool {
	_, err := s.httpc.Send(&httpc.OptionsSt{
		Method: "POST",
		Uri:    "send",

		ReqObj: req,
	})

	return err == nil
}
