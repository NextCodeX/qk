package core

type MultiExpression interface {
    execute() Value
    recursiveEvalMultiExpression(nextExpr BinaryExpression, exprList []BinaryExpression) Value
    calculateIfNotExist(primaryExpr PrimaryExpression, exprList []BinaryExpression)
    getNextExprForMultiExpression(varname string, exprList []BinaryExpression) BinaryExpression


    getfinalExpr() BinaryExpression
    getExprs() []BinaryExpression

    Expression
}

type MultiExpressionImpl struct {
    finalExpr BinaryExpression
    list []BinaryExpression
    ExpressionAdapter
}

func newMultiExpressionImpl() MultiExpression {
    return &MultiExpressionImpl{}
}

func (mulExpr *MultiExpressionImpl) setStack(stack Function) {
    mulExpr.stack = stack

    for _, subExpr := range mulExpr.list {
        subExpr.setStack(stack)
    }

    mulExpr.finalExpr.setStack(stack)
}

func (mulExpr *MultiExpressionImpl) execute() Value {
    // 每次执行多元表达式之前，初始化临时变量池
    mulExpr.setVar(tmpVarsKey, emptyJsonObject())
    res := mulExpr.recursiveEvalMultiExpression(mulExpr.finalExpr, mulExpr.list)
    return res
}

func (mulExpr *MultiExpressionImpl) getfinalExpr() BinaryExpression {
    return mulExpr.finalExpr
}
func (mulExpr *MultiExpressionImpl) getExprs() []BinaryExpression {
    return mulExpr.list
}


func (mulExpr *MultiExpressionImpl) recursiveEvalMultiExpression(nextExpr BinaryExpression, exprList []BinaryExpression) Value {
    left := nextExpr.leftExpr()
    right := nextExpr.rightExpr()
    if left.isConst() && right.isConst() {
        return nextExpr.execute()
    }

    mulExpr.calculateIfNotExist(left, exprList)
    mulExpr.calculateIfNotExist(right, exprList)

    return nextExpr.execute()
}

// 检查相应的变量是否已计算，若未计算则进行计算
func (mulExpr *MultiExpressionImpl) calculateIfNotExist(primaryExpr PrimaryExpression, exprList []BinaryExpression) {
    if !primaryExpr.isVar() {
        return
    }

    varExpr := primaryExpr.(*VarPrimaryExpression)
    if !varExpr.execute().isNULL() {
        return
    }
    nextExpr := mulExpr.getNextExprForMultiExpression(varExpr.getName(), exprList)
    if nextExpr != nil {
        mulExpr.recursiveEvalMultiExpression(nextExpr, exprList)
    }
}

func (mulExpr *MultiExpressionImpl) getNextExprForMultiExpression(varname string, exprList []BinaryExpression) BinaryExpression {
    for _, subExpr := range exprList {
        receiver := subExpr.getReceiver()
        if receiver != nil && receiver.raw()[0].raw() == varname {
            return subExpr
        }

        leftExpr := subExpr.leftExpr()
        if subExpr.isAssign() && leftExpr.isVar() && leftExpr.(*VarPrimaryExpression).getName() == varname {
            subExpr.execute()
            return nil
        }
    }
    runtimeExcption("executeMultiExpression Exception: no expression for ", varname)
    return nil
}