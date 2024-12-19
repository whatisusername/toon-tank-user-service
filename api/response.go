package api

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func successResponse(data interface{}) response {
	return response{
		Success: true,
		Message: "Success",
		Data:    data,
	}
}

func errorResponse(err error) response {
	return response{
		Success: false,
		Message: err.Error(),
	}
}
