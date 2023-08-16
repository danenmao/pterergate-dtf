package errordef

//
// DummyError
//
type DummyError struct {
}

//
// 实现 Error 接口
//
func (e *DummyError) Error() string {
	return "Dummy"
}

//
// NotFoundError
//
type NotFoundError struct {
}

//
// 实现 Error 接口
//
func (e *NotFoundError) Error() string {
	return "NotFound"
}
