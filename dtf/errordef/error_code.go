package errordef

// 错误码定义
const (
	Error_Success          = 0
	Error_InvalidParameter = 0x100001
)

// 错误串定义
const (
	Error_Msg_Success          = "Success"
	Error_Msg_InvalidParameter = "InvalidParameter"
	Error_Msg_ParsingParam     = "ParsingParamError"
	Error_Msg_OperationFailed  = "OperationFailed"
)

// 错误码映射表
var ErrorMsgMap = map[uint32]string{
	Error_Success:          Error_Msg_InvalidParameter,
	Error_InvalidParameter: Error_Msg_InvalidParameter,
}
