package logic

// This error just wraps a string
type ActionError struct {
	Kind string
	// Developer friendly message
	Msg string
}

func (err *ActionError) Error() string {
	return "kind: " + err.Kind + ", msg: " + err.Msg
}
