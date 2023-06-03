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

func (r *Response) DefaultNotAcceptable() {
	r.Status = false
	r.Code = 406
	r.Message = "input not valid"
	r.Data = nil
}

func (r *Response) DefaultInternalError() {
	r.Status = false
	r.Code = 500
	r.Message = "request failed, server error"
	r.Data = nil
}
