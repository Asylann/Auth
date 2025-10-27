package response

type Response struct {
	Data interface{} `json:"data,omitempty"`
	Err  string      `json:"err,omitempty"`
}
