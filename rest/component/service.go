package component

type Service struct {
}

// NewService 执行新的服务
func NewService(comp Component, opts []Option) *Service {
	return &Service{}
}
