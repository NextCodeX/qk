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
	list      []BinaryExpression
	ExpressionAdapter
}

func newMultiExpressionImpl() MultiExpression {
	return &MultiExpressionImpl{}
}

func (mulExpr *MultiExpressionImpl) setParent(p Function) {
	mulExpr.ExpressionAdapter.setParent(p)

	for _, subExpr := range mulExpr.list {
		subExpr.setParent(p)
	}

	mulExpr.finalExpr.setParent(p)
}

func (mulExpr *MultiExpressionImpl) execute() Value {
	// 每次执行多元表达式之前，初始化临时变量池
	mulExpr.setLocalVar(tmpVarsKey, emptyJsonObject())
	//defer func() {mulExpr.setVar(tmpVarsKey, nil)}()

	res := mulExpr.recursiveEvalMultiExpression(mulExpr.finalExpr, mulExpr.list)
	return res
}

func (mulExpr *MultiExpressionImpl) getfinalExpr() BinaryExpression {
	return mulExpr.finalExpr
}
func (mulExpr *MultiExpressionImpl) getExprs() []BinaryExpression {
	return mulExpr.list
}

func (mulExpr *MultiExpressionImpl) recursiveEvalMultiExpression(curExpr BinaryExpression, exprList []BinaryExpression) Value {
	left := curExpr.leftExpr()
	right := curExpr.rightExpr()
	if left.isConst() && right.isConst() {
		return curExpr.execute()
	}

	mulExpr.calculateIfNotExist(left, exprList)
	if curExpr.isOr() {
		if leftVal := left.execute(); toBoolean(leftVal) {
			curExpr.resultCache(leftVal)
			return leftVal
		}
	}
	if curExpr.isAnd() {
		if leftVal := left.execute(); !toBoolean(leftVal) {
			curExpr.resultCache(leftVal)
			return leftVal
		}
	}

	mulExpr.calculateIfNotExist(right, exprList)

	return curExpr.execute()
}

// 检查相应的变量是否已计算，若未计算则进行计算
func (mulExpr *MultiExpressionImpl) calculateIfNotExist(primaryExpr PrimaryExpression, exprList []BinaryExpression) {
	varExpr, ok := primaryExpr.(*VarPrimaryExpression)
	if !ok || !isTmpVar(varExpr.name()) || mulExpr.getVar(varExpr.name()) != NULL {
		return
	}
	nextExpr := mulExpr.getNextExprForMultiExpression(varExpr.name(), exprList)
	if nextExpr != nil {
		mulExpr.recursiveEvalMultiExpression(nextExpr, exprList)
	}
}

func (mulExpr *MultiExpressionImpl) getNextExprForMultiExpression(varname string, exprList []BinaryExpression) BinaryExpression {
	for _, subExpr := range exprList {
		receiver := subExpr.getReceiver()
		if receiver != nil {
			if vp, ok := receiver.(*VarPrimaryExpression); ok && vp.name() == varname {
				return subExpr
			}
		}

		leftExpr := subExpr.leftExpr()
		if subExpr.isAssign() && leftExpr.isVar() && leftExpr.(*VarPrimaryExpression).name() == varname {
			subExpr.execute()
			return nil
		}
	}
	runtimeExcption("executeMultiExpression Exception: no expression for ", varname)
	return nil
}
