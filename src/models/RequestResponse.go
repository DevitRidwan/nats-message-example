package models

// Response
type (
	Response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
)

// Request
type (
	RequestQueue struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	RequestProduce struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Message  string `json:"message"`
	}
)
