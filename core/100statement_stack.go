package core

// stack 在Quick中是存放变量的一个地方
// 对于函数, 它自身就是stack, 它的stack与parent不一致, 它的parent是父函数(上一层stack)
// 对于非函数的statement, 它们的parent就是stack, parent与stack是一致的, 皆是父函数

// ValueStack 为statement, expression提供变量操作的接口
type ValueStack struct {
    stack Function
}

func (vs *ValueStack) getVar(name string) Value {
    var level = vs.stack
    var res Value
    for level != nil {
        varMap := level.getLocalVars()
        if varMap == nil {
            break
        }
        res = varMap.get(name)
        if res == nil {
            level = level.getParent()
        } else {
            break
        }
    }
    return res
}

func (vs *ValueStack) setVar(name string, value Value) {
    var level = vs.stack
    var res Value
    for level != nil {
        varMap := level.getLocalVars()
        if varMap == nil {
            break
        }
        res = varMap.get(name)
        if res == nil {
            level = level.getParent()
        } else {
            varMap.add(name, value)
            break
        }
    }
    if res == nil {
        vs.stack.getLocalVars().add(name, value)
    }
}


