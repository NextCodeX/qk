package core

type PrimaryExpressionType int
const (
    VarPrimaryExpressionType PrimaryExpressionType = 1 << iota
    ConstPrimaryExpressionType
    DynamicStrPrimaryExpressionType
    ArrayPrimaryExpressionType
    ObjectPrimaryExpressionType
    ChainCallPrimaryExpressionType
    ElementPrimaryExpressionType
    AttibutePrimaryExpressionType
    FunctionCallPrimaryExpressionType
    MethodCallPrimaryExpressionType
    NotPrimaryExpressionType
    ExprPrimaryExpressionType
    OtherPrimaryExpressionType PrimaryExpressionType = 0
)

type PrimaryExpression interface {
    getName() string
    notFlag() bool
    setNotFlag(flag bool)
    addType(t PrimaryExpressionType)

    execute() Value

    isVar() bool
    isConst() bool
    isDynamicStr() bool
    isArray() bool
    isObject() bool
    isChainCall() bool
    isElement() bool
    isAttibute() bool
    isFunctionCall() bool
    isMethodCall() bool
    isNot() bool
    isExpr() bool
    isOther() bool

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

func (priExpr *PrimaryExpressionImpl) isElement() bool {
    return (priExpr.t & ElementPrimaryExpressionType) == ElementPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isAttibute() bool {
    return (priExpr.t & AttibutePrimaryExpressionType) == AttibutePrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isFunctionCall() bool {
    return (priExpr.t & FunctionCallPrimaryExpressionType) == FunctionCallPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isMethodCall() bool {
    return (priExpr.t & MethodCallPrimaryExpressionType) == MethodCallPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isNot() bool {
    return (priExpr.t & NotPrimaryExpressionType) == NotPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isExpr() bool {
    return (priExpr.t & ExprPrimaryExpressionType) == ExprPrimaryExpressionType
}

func (priExpr *PrimaryExpressionImpl) isOther() bool {
    return (priExpr.t & OtherPrimaryExpressionType) == OtherPrimaryExpressionType
}
