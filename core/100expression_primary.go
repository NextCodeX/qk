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
    NestedPrimaryExpressionType
    NotPrimaryExpressionType
    SelfIncrPrimaryExpressionType
    SelfDecrPrimaryExpressionType
    TernaryOperatorPrimaryExpressionType
)

type PrimaryExpression interface {
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

    Expression
}

type PrimaryExpressionImpl struct {
    t PrimaryExpressionType
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

    if res == nil {
        res = NULL
    }
    return res
}

func (priExpr *PrimaryExpressionImpl) isVar() bool {
    return priExpr.t & VarPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isConst() bool {
    return priExpr.t & ConstPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isDynamicStr() bool {
    return priExpr.t & DynamicStrPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isArray() bool {
    return priExpr.t & ArrayPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isObject() bool {
    return priExpr.t & ObjectPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isChainCall() bool {
    return priExpr.t & ChainCallPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isElemFunctionCall() bool {
    return priExpr.t & ElemFunctionCallPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isFunction() bool {
    return priExpr.t & FunctionPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isSubList() bool {
    return priExpr.t & SubListPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isElement() bool {
    return priExpr.t & ElementPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isFunctionCall() bool {
    return priExpr.t & FunctionCallPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isNot() bool {
    return priExpr.t & NotPrimaryExpressionType > 0
}

func (priExpr *PrimaryExpressionImpl) isTernaryOperator() bool {
    return priExpr.t & TernaryOperatorPrimaryExpressionType > 0
}
