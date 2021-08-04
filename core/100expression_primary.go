package core

type PrimaryExpressionType int
const (
    VarPrimaryExpressionType PrimaryExpressionType = 1 << iota
    ConstPrimaryExpressionType
    DynamicStrPrimaryExpressionType
    ArrayPrimaryExpressionType
    ObjectPrimaryExpressionType
    ChainCallPrimaryExpressionType
    ElemFunctionCallPrimaryExpressionType
    FunctionPrimaryExpressionType
    SubListPrimaryExpressionType
    ElementPrimaryExpressionType
    FunctionCallPrimaryExpressionType
    NotPrimaryExpressionType
    ExprPrimaryExpressionType
)

type PrimaryExpression interface {
    getName() string
    notFlag() bool
    setNotFlag(flag bool)
    addType(t PrimaryExpressionType)
    typeName() string

    execute() Value

    isVar() bool
    isConst() bool
    isDynamicStr() bool
    isArray() bool
    isObject() bool
    isChainCall() bool
    isElemFunctionCall() bool
    isFunction() bool
    isSubList() bool
    isElement() bool
    isFunctionCall() bool
    isNot() bool
    isExpr() bool

    Expression
}

type PrimaryExpressionImpl struct {
    t PrimaryExpressionType
    not   bool             // 是否进行非处理
    doExec func() Value
    ExpressionAdapter
}

func newPrimaryExpression() PrimaryExpression {
    return &PrimaryExpressionImpl{}
}

func (priExpr *PrimaryExpressionImpl) execute() Value {
    if priExpr == nil {
        return NULL
    }
    var res Value

    if priExpr.doExec != nil {
        res = priExpr.doExec()
    } else {
        runtimeExcption("ExpressionExecutor#evalPrimaryExpr: unknown primary expression type")
    }

    if priExpr.isNot() { // 是否为非类型表达式
        flag := toBoolean(res)
        if priExpr.not { // 是否需要对表达式进行非逻辑运算
            return newQKValue(!flag)
        } else {
            return newQKValue(flag)
        }
    } else {
        if res == nil {
            res = NULL
        }
        return res
    }
}

func (priExpr *PrimaryExpressionImpl) typeName() string {
    var typeNames string
    if priExpr.isVar() {
        typeNames += "variable, "
    }
    if priExpr.isConst() {
        typeNames += "constant, "
    }
    if priExpr.isDynamicStr() {
        typeNames += "dynamic string, "
    }
    if priExpr.isArray() {
        typeNames += "json array, "
    }
    if priExpr.isObject() {
        typeNames += "json object, "
    }
    if priExpr.isChainCall() {
        typeNames += "chain call, "
    }
    if priExpr.isElemFunctionCall() {
        typeNames += "element and function call mixture, "
    }
    if priExpr.isFunction() {
        typeNames += "function, "
    }
    if priExpr.isSubList() {
        typeNames += "subList, "
    }
    if priExpr.isElement() {
        typeNames += "element, "
    }
    if priExpr.isFunctionCall() {
        typeNames += "function call, "
    }
    if priExpr.isNot() {
        typeNames += "not operator, "
    }
    if priExpr.isExpr() {
        typeNames += "expression, "
    }

    if typeNames == "" {
        return "unknown expression type"
    } else {
        return typeNames
    }
}

func (priExpr *PrimaryExpressionImpl) getName() string {
    return ""
}
func (priExpr *PrimaryExpressionImpl) notFlag() bool {
    return priExpr.not
}
func (priExpr *PrimaryExpressionImpl) setNotFlag(flag bool) {
    priExpr.not = flag
}

func (priExpr *PrimaryExpressionImpl) addType(t PrimaryExpressionType) {
    priExpr.t = priExpr.t | t
}

func (priExpr *PrimaryExpressionImpl) isVar() bool {
    return (priExpr.t & VarPrimaryExpressionType) == VarPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isConst() bool {
    return (priExpr.t & ConstPrimaryExpressionType) == ConstPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isDynamicStr() bool {
    return (priExpr.t & DynamicStrPrimaryExpressionType) == DynamicStrPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isArray() bool {
    return (priExpr.t & ArrayPrimaryExpressionType) == ArrayPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isObject() bool {
    return (priExpr.t & ObjectPrimaryExpressionType) == ObjectPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isChainCall() bool {
    return (priExpr.t & ChainCallPrimaryExpressionType) == ChainCallPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isElemFunctionCall() bool {
    return (priExpr.t & ElemFunctionCallPrimaryExpressionType) == ElemFunctionCallPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isFunction() bool {
    return (priExpr.t & FunctionPrimaryExpressionType) == FunctionPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isSubList() bool {
    return (priExpr.t & SubListPrimaryExpressionType) == SubListPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isElement() bool {
    return (priExpr.t & ElementPrimaryExpressionType) == ElementPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isFunctionCall() bool {
    return (priExpr.t & FunctionCallPrimaryExpressionType) == FunctionCallPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isNot() bool {
    return (priExpr.t & NotPrimaryExpressionType) == NotPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isExpr() bool {
    return (priExpr.t & ExprPrimaryExpressionType) == ExprPrimaryExpressionType
}
