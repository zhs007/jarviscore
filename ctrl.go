package jarviscore

// Ctrl -
type Ctrl interface {
	Run(command []byte) (result []byte, err error)
}
