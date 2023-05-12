package krps

type SendReqSt struct {
	Records []SendReqRecordSt `json:"records"`
}

type SendReqRecordSt struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
