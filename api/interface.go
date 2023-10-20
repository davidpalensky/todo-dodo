package api

// Every response from Todo-dodo's api should include the meta data defined here.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	// Should be empty if success = 1
	Err_msg string `json:"err_msg"`
}
