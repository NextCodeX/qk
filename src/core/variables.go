package core


type Variables map[string]*Variable

func newVariables() Variables {
	return make(map[string]*Variable)
}

func (vs *Variables) isEmpty() bool {
	return vs == nil || len(*vs) < 1
}

func (vs *Variables) add(v *Variable) {
	(*vs)[v.name] = v
}

func (vs *Variables) get(name string) *Variable {
	if vs.isEmpty() {
		return nil
	}
	res, ok := (*vs)[name]
	if ok {
		return res
	}
	return nil
}
