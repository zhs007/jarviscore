package jarviscore

// Ctrl -
type Ctrl interface {
	Run(command string) (result string, err error)
}
