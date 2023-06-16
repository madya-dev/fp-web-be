package response

type Response struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Response) DefaultOK() {
	r.Status = true
	r.Code = 200
	r.Message = ""
	r.Data = nil
}

func (r *Response) DefaultCreated() {
	r.Status = true
	r.Code = 201
	r.Message = "data created"
	r.Data = nil
}

func (r *Response) DefaultBadRequest() {
	r.Status = false
	r.Code = 400
	r.Message = "input not valid"
	r.Data = nil
}
func (r *Response) DefaultUnauthorized() {
	r.Status = false
	r.Code = 401
	r.Message = "authorization failed"
	r.Data = nil
}

func (r *Response) DefaultNotFound() {
	r.Status = false
	r.Code = 404
	r.Message = "data not found"
	r.Data = nil
}

func (r *Response) DefaultConflict() {
	r.Status = false
	r.Code = 409
	r.Message = "input data conflict"
	r.Data = nil
}

func (r *Response) DefaultInternalError() {
	r.Status = false
	r.Code = 500
	r.Message = "request failed, server error"
	r.Data = nil
}
