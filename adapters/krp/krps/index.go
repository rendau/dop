package krps

import (
	"github.com/rendau/dop/adapters/client/httpc"
	"github.com/rendau/dop/adapters/logger"
	"net/http"
)

type St struct {
	lg    logger.Lite
	httpc httpc.HttpC
}

func New(lg logger.Lite, httpc httpc.HttpC) *St {
	return &St{
		lg:    lg,
		httpc: httpc,
	}
}

func (s *St) SendJson(topic, key string, value any) error {
	return s.SendManyJson(topic, key, []any{value})
}

func (s *St) SendManyJson(topic, key string, value []any) error {
	reqObj := &SendReqSt{
		Records: make([]SendReqRecordSt, len(value)),
	}

	for i, v := range value {
		reqObj.Records[i] = SendReqRecordSt{
			Key:   key,
			Value: v,
		}
	}

	_, err := s.httpc.Send(&httpc.OptionsSt{
		Uri:    "topics/" + topic,
		Method: "POST",
		Headers: http.Header{
			"Content-Type": []string{"application/vnd.kafka.json.v2+json"},
		},
		LogPrefix: "topics/" + topic + "(json)",
		ReqObj:    reqObj,
	})

	return err
}
