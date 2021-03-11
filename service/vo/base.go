package vo

type (
	// 基础响应格式
	BaseResp struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
	}

	// 相应数据
	BaseRespData struct {
		*BaseResp
		Data *List `json:"data"`
	}

	// data 数据结构
	List struct {
		Total int64       `json:"total"`
		List  interface{} `json:"list"`
	}
)

const (
	BaseRespType = iota + 1
	BaseRespDataType
)

func HandlerRespVo(resp int8) BaseRespVo {
	var baseRespVo BaseRespVo
	switch resp {
	case BaseRespType:
		baseRespVo = new(BaseResp)
	case BaseRespDataType:
		baseRespVo = new(BaseRespData)
	}
	return baseRespVo
}

type BaseRespVo interface {
	Success() interface{}
	Failed(code int64, msg string) interface{}
}

func (b *BaseResp) Success() interface{} {
	return &BaseResp{
		Code: 0,
		Msg:  "success",
	}
}

func (b *BaseResp) Failed(code int64, msg string) interface{} {
	return &BaseResp{
		Code: code,
		Msg:  msg,
	}
}

func (b *BaseRespData) Success() interface{} {
	return &BaseRespData{
		BaseResp: &BaseResp{
			Code: 0,
			Msg:  "success",
		},
		Data: b.Data,
	}
}

func (b *BaseRespData) Failed(code int64, msg string) interface{} {
	return &BaseResp{
		Code: code,
		Msg:  msg,
	}
}
