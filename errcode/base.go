package jarviserrcode

const (
	FILEREADSIZEINVALID = 1
)

// GetErrCodeString -
func GetErrCodeString(errcode int) string {
	switch errcode {
	case FILEREADSIZEINVALID:
		return "Invalid FileSize & BytesRead."
	}

	return ""
}
