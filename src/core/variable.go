package core


type Variable struct{
	name string
	val *Value
}

func newVar(name string, rawVal interface{}) *Variable {
	res := &Variable{
		name: name,
		val:  newVal(rawVal),
	}

	return res
}

func toVar(name string, rawVal *Value) *Variable {
	res := &Variable{
		name: name,
		val:  rawVal,
	}
	return res
}
