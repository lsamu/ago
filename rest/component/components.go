package component

// Components 多个组件
type Components struct {
	comps []ComponentOption
}

// Register 注册
func (cc *Components) Register(c Component, oo ...Option) {
	cc.comps = append(cc.comps, ComponentOption{cc: c, op: oo})
}

// List 列表
func (cc *Components) List() []ComponentOption {
	return cc.comps
}
