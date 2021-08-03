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

func (priExpr *VarPrimaryExpression) getName() string {
    return priExpr.varname
}

func (priExpr *VarPrimaryExpression) doExecute() Value {
    //fmt.Println("exec varExpr: ", priExpr.varname, priExpr.getVar(priExpr.varname), tokensString(priExpr.ts), priExpr.getStack())
    return priExpr.getVar(priExpr.varname)
}

func (priExpr *VarPrimaryExpression) getAttribute(obj Object) Value {
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

