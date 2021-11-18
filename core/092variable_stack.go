package core

// cur 在Quick中是存放变量的一个地方
// 对于函数, 它自身就是stack, 它的stack与parent不一致, 它的parent是父函数(上一层stack)
// 对于非函数的statement, 它们的parent就是stack, parent与stack是一致的, 皆是父函数

// ValueStack 为statement, expression提供变量操作的接口
type ValueStack struct {
	cur Frame
}

func (vs *ValueStack) getVar(name string) Value {
	var level = vs.cur

	for level != nil {
		if varMap := level.varList(); varMap != nil {
			if val := varMap.get(name); val != nil {
				return val
			}
		}

		level = level.parentFrame()
	}
	return NULL
}

func (vs *ValueStack) setVar(name string, value Value) {
	var level = vs.cur
	var flag bool
	for level != nil {
		if varMap := level.varList(); varMap != nil {
			if val := varMap.get(name); val != nil {
				varMap.add(name, value)
				flag = true
				break
			}
		}

		level = level.parentFrame()
	}
	if !flag {
		vs.cur.varList().add(name, value)
	}
}
