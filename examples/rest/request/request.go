package request

type UserRequest struct {
    Type string `json:"type" form:"type"`
}

func (r *UserRequest) Check() error{
    return nil
}