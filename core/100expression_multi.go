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
    return mulExpr.recursiveEvalMultiExpression(mulExpr.finalExpr, mulExpr.list)
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
    varname := primaryExpr.getName()
    variable := mulExpr.getVar(varname)
    if variable != nil {
        return
    }
    nextExpr := mulExpr.getNextExprForMultiExpression(varname, exprList)
    if nextExpr == nil {
        runtimeExcption("executeMultiExpression Exception")
    }
    mulExpr.recursiveEvalMultiExpression(nextExpr, exprList)
}

func (mulExpr *MultiExpressionImpl) getNextExprForMultiExpression(varname string, exprList []BinaryExpression) BinaryExpression {
    for _, subExpr := range exprList {
        if subExpr.getReceiver() == varname {
            return subExpr
        }
    }
    return nil
}