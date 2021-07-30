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
    return priExpr.getVar(priExpr.varname)
}

