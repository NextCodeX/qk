package core

import (
    "fmt"
    "strconv"
)

var (
    funcList = make(map[string]Function)
    mainFunc = newFunc()
)


func Compile(ts []Token) {
    mainFunc.raw = ts
    for i:=0; i<len(ts); {
        t := ts[i]
        var endIndex int
        var stmt *Statement

        if t.t != Identifier {
            goto next_loop
        }
        switch t.str {
        case "func":
        case "if":
           stmt, endIndex = parseIfStatement(i, ts)
        case "for":
           stmt, endIndex = parseForStatement(i, ts)
        case "switch":
        default:
            endIndex = nextBoundary(i, ts)
            if endIndex>0 {
                stmt = newStatement(ExpressionStatement, ts[i:endIndex])
            }

        }
        if endIndex>0 {
            parseStatement(stmt)
            mainFunc.addStatement(stmt)
            i = endIndex
        }


        next_loop:
        i++
    }
}

func parseIfStatement(currentIndex int, ts []Token) (*Statement, int) {
    stmt := &Statement{t:IfStatement}

    index := nextSymbol(currentIndex, ts, "{")
    stmt.condition = &Expression{
        t:     BinaryExpression,
        raw:   ts[currentIndex+1:index],
    }

    scopeOpenCount := 1
    var endIndex int
    for i:=index+1; i<=len(ts); i++ {
        t := ts[i]
        if t.str == "{" {
            scopeOpenCount++
        }
        if t.str == "}" {
            scopeOpenCount--
            endIndex = i
            if scopeOpenCount == 0 {
                break
            }
        }
        stmt.raw = append(stmt.raw, t)
    }
    return stmt, endIndex
}

func parseForStatement(currentIndex int, ts []Token) (*Statement, int) {
    stmt := &Statement{t:IfStatement}
    index := nextSymbol(currentIndex, ts, "{")
    exprs := splitExpression(ts[currentIndex+1:index])
    stmt.setHeaderInfo(exprs)
    scopeOpenCount := 1
    var endIndex int
    for i:=index+1; i<=len(ts); i++ {
        t := ts[i]
        if t.str == "{" {
            scopeOpenCount++
        }
        if t.str == "}" {
            scopeOpenCount--
            endIndex = i
            if scopeOpenCount == 0 {
                break
            }
        }
        stmt.raw = append(stmt.raw, t)
    }
    return stmt, endIndex
}

func splitExpression(ts []Token) []*Expression {
    res := make([]*Expression, 3)
    if !hasSymbol(ts, ";") {
        res[1] = &Expression{
            t:     BinaryExpression,
            raw:   ts,
        }
        return res
    }
    index := nextSymbol(0, ts, ";")
    if index > 2 {
        res[0] = &Expression{
            t:     BinaryExpression,
            raw:   ts[:index],
        }
    }
    preIndex := index+1
    index = nextSymbol(preIndex, ts, ";")
    if index - preIndex > 2 {
        res[1] = &Expression{
            t:     BinaryExpression,
            raw:   ts[preIndex:index],
        }
    }
    preIndex = index+1
    index = nextSymbol(preIndex, ts, ";")
    if index - preIndex > 2 {
        res[2] = &Expression{
            t:     BinaryExpression,
            raw:   ts[preIndex:index],
        }
    }
    return res
}


func nextBoundary(currentIndex int, ts []Token) int {
    return nextSymbol(currentIndex, ts, ";")
}

func nextSymbol(currentIndex int, ts []Token, s string) int {
    for i:=currentIndex; i<len(ts); i++ {
        t := ts[i]
        if t.t == Symbol && t.str == s {
            return i
        }
    }
    return -1
}

func hasSymbol(ts []Token, s string) bool {
    for i:=0; i<len(ts); i++ {
        t := ts[i]
        if t.t == Symbol && t.str == s {
            return true
        }
    }
    return false
}


func parseStatement(stmt *Statement) {
    ts := stmt.raw

    switch {
    case stmt.isExpressionStatement():
        parseExpressionStatement(ts, stmt)

    case stmt.isIfStatement():
    case stmt.isForStatement():
    case stmt.isSwitchStatement():
    case stmt.isReturnStatement():
    }

}

func parseExpressionStatement(ts []Token, stmt *Statement) {
    var expr *Expression
    tlen := len(ts)
    if tlen < 2 {
        return
    }
    if tlen == 2 {
        if ts[0].isIdentifier() && (ts[1].str=="++" || ts[1].str=="--") {

        }
        return
    }

    // 去括号
    if ts[0].str == "(" && ts[tlen-1].str == ")" {
        ts = ts[1:tlen-1]
    }
    tlen = len(ts)

    // 处理基本表达式
    if tlen == 3 {
        expr = parseTerm(ts)
        if expr != nil {
            stmt.addExpression(expr)
        }
        return
    }

    // 处理多维表达式
    expr, next := parseListExpression(ts, stmt)
    if expr != nil {
        stmt.addExpression(expr)
    }
    if next != nil {
        parseExpressionStatement(next, stmt)
    }
}

func parseListExpression(ts []Token, stmt *Statement) (*Expression, []Token) {
    first := ts[0]
    mid := ts[1]
    third := ts[2]
    fourth := ts[3]
    var expr *Expression
    if first.isIdentifier() && mid.str == "=" {
        if first.isTmp() && len(ts) == 5 {
            fifth := ts[4]
            exprType := getExpressionType(fourth.str)
            expr = &Expression{
                t:          exprType,
                left:       parsePrimaryExpression(&third),
                right:      parsePrimaryExpression(&fifth),
                tmpname:    first.str,
            }
            return expr, nil
        }

        expr = newListExpression()
        tmp := &Expression{
            t:     AssignExpression | TmpExpression,
            left:  PrimaryExpr{t: Varname, name: first.str},
            right: PrimaryExpr{t: Fill},
        }
        expr.list = append(expr.list, tmp)

        return expr, ts[3:]
    }
    b11 := isMulDiv(mid.str)
    b12 := isAddSub(mid.str) && isAddSub(fourth.str)
    if b11 || b12 {
        exprType := getExpressionType(mid.str)
        tmpname := getTmpname(stmt)
        expr = &Expression{
            t:          TmpExpression | exprType,
            left:       parsePrimaryExpression(&first),
            right:      parsePrimaryExpression(&third),
            tmpname:    tmpname,
        }
        head := tmpToken(tmpname)
        next := insert(head, ts[3:])
        return expr, next
    }
    if isAddSub(mid.str) && isMulDiv(fourth.str) {
        exprType := getExpressionType(mid.str)
        tmpname := getTmpname(stmt)
        tail := tmpToken(tmpname)
        expr = &Expression{
            t:          TmpExpression | exprType,
            left:       parsePrimaryExpression(&first),
            right:      parsePrimaryExpression(&tail),
        }
        assignToken := symbolToken("=")
        next := insert2(tail, assignToken, ts[2:])
        return expr, next
    }
    return nil, nil
}

func isAddSub(op string) bool {
    return op == "+" || op == "-"
}

func isMulDiv(op string) bool {
    return op == "*" || op == "/"
}

func getTmpname(stmt *Statement) string {
    name := fmt.Sprintf("tmp.%v", stmt.tmpcount)
    stmt.tmpcount++
    return name
}

func parseTerm(ts []Token) *Expression{
    first := ts[0]
    mid := ts[1]
    third := ts[2]
    left := parsePrimaryExpression(&first)
    right := parsePrimaryExpression(&third)
    var exprType ExpressionType
    switch {
    case  first.isIdentifier() && mid.str == "=":
        exprType = AssignExpression

    case mid.str == "+":
        exprType = AddExpression

    case mid.str == "-":
        exprType = SubExpression

    case mid.str == "*":
        exprType = MulExpression

    case mid.str == "/":
        exprType = DivExpression

    default:
        return nil
    }

    expr := &Expression{
        t:     exprType,
        left:  left,
        right: right,
    }
    return expr
}

func getExpressionType(op string) ExpressionType {
    switch op {
    case "+": return AddExpression
    case "-": return SubExpression
    case "*": return MulExpression
    case "/": return DivExpression

    default:
        return -1
    }
}

func parsePrimaryExpression(t *Token) PrimaryExpr {
    v := tokenToValue(t)
    var res PrimaryExpr
    if v == nil {
        res = PrimaryExpr{t: Varname, name: t.str}
    } else {
        res = PrimaryExpr{t: Const, result: v}
    }
    return res
}

func tokenToValue(t *Token) *Value {
    var v Value
    if t.isFloat() {
        f, err := strconv.ParseFloat(t.str, 64)
        exitOnError(err)
        v = newVal(f)
        return &v
    }
    if t.isInt() {
        i, err := strconv.Atoi(t.str)
        exitOnError(err)
        v = newVal(i)
        return &v
    }
    if t.isStr() {
        v = newVal(fmt.Sprintf("%v", t.str))
        return &v
    }
    if t.str == "true" || t.str == "false" {
        b, err := strconv.ParseBool(t.str)
        exitOnError(err)
        v = newVal(b)
        return &v
    }
    return nil
}
