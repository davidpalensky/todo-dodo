package logic

// This error just wraps a string
type LogicError struct {
	Kind string
	// Developer friendly message
	Msg string
}

func (err *LogicError) Error() string {
	return "kind: " + err.Kind + ", msg: " + err.Msg
}
