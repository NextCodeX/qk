package core

type VarPrimaryExpression struct {
	varname string
	PrimaryExpressionImpl
}

func newVarPrimaryExpression(varname string) PrimaryExpression {
	expr := &VarPrimaryExpression{}
	expr.t = VarPrimaryExpressionType
	expr.varname = varname
	expr.doExec = expr.doExecute
	return expr
}

// 获取变量名
func (priExpr *VarPrimaryExpression) name() string {
	return priExpr.varname
}

func (priExpr *VarPrimaryExpression) doExecute() Value {
	return priExpr.getVar(priExpr.varname)
}

func (priExpr *VarPrimaryExpression) getField(obj Object) Value {
	return obj.get(priExpr.varname)
}

func (priExpr *VarPrimaryExpression) beAssigned(res Value) {
	priExpr.setVar(priExpr.varname, res)
}

func (priExpr *VarPrimaryExpression) assign(object JSONObject, res Value) {
	object.put(priExpr.varname, res)
}

func (priExpr *VarPrimaryExpression) nameIs(s string) bool {
	return priExpr.varname == s
}
