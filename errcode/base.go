package jarviserrcode

const (
	FILEREADSIZEINVALID = 1
	PEERADDREMPTY       = 2
)

// GetErrCodeString -
func GetErrCodeString(errcode int) string {
	switch errcode {
	case FILEREADSIZEINVALID:
		return "invalid filesize & bytesread."
	case PEERADDREMPTY:
		return "peeraddr is empty."
	}

	return ""
}
