package errutil

var (
	ErrHTTPRequest     = NewInternalError("http request error")
	ErrJSONDecode      = NewInternalError("json decode error")
	ErrTimeParse       = NewInternalError("time parse error")
	ErrGetProgramNotOK = NewInternalError("http get program status code not ok")
	// 分類できない系
	ErrInternal = NewInternalError("internal something error")
)
