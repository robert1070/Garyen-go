package base

type (
	// 基础响应格式
	Resp struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
	}

	// 相应数据
	RespData struct {
		*Resp
		Data interface{} `json:"data"`
	}
)

func NewResp() *Resp {
	return &Resp{}
}

func (rs *Resp) Success() interface{} {
	return &Resp{
		Code: 0,
		Msg:  "success",
	}
}

func (rs *Resp) Failed(code int64, msg string) interface{} {
	return &Resp{
		Code: code,
		Msg:  msg,
	}
}

func NewRespData(i interface{}) *RespData {
	return &RespData{
		Data: i,
	}
}

func (ra *RespData) Success() interface{} {
	return &RespData{
		Resp: &Resp{
			Code: 0,
			Msg:  "success",
		},
		Data: ra.Data,
	}
}

func (ra *RespData) Failed(code int64, msg string) interface{} {
	return &Resp{
		Code: code,
		Msg:  msg,
	}
}
