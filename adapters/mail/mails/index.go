package mails

import (
	"github.com/rendau/dop/adapters/client/httpc"
	"github.com/rendau/dop/adapters/mail"
)

type St struct {
	httpc httpc.HttpC
}

func New(httpc httpc.HttpC) *St {
	return &St{
		httpc: httpc,
	}
}

func (m *St) Send(data *mail.SendReqSt) bool {
	_, err := m.httpc.Send(&httpc.OptionsSt{
		Method: "POST",
		Uri:    "send",

		ReqObj: data,
	})

	return err == nil
}
