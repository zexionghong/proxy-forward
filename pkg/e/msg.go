package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "error",
	INVALID_PARAMS: "request params error",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
