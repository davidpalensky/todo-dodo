package action

// Helper for telling the go compiler to calm down a little
func DoNothing(a any) {

}

// This error just wraps a string
type ActionError struct {
	Kind string
	// Developer friendly message
	Msg string
}

func (err *ActionError) Error() string {
	return "kind: " + err.Kind + ", msg: " + err.Msg
}
