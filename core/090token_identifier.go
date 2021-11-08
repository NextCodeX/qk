package core


// 变量名，函数名
type NameToken struct {
	TokenAdapter
}
func newNameToken(raw string, row int) Token {
	t := &NameToken{}
	t.lineIndex = row
	t.val = raw
	t.typName = "Identifier"
	return t
}
func (t *NameToken) toExpr() PrimaryExpression {
	return newVarPrimaryExpression(t.name())
}
func (t *NameToken) name() string {
	return t.val.(string)
}

// 关键字(if, for...)
type KeyToken struct {
	TokenAdapter
}
func newKeyToken(raw string, row int) Token {
	t := &KeyToken{}
	t.lineIndex = row
	t.val = raw
	t.typName = "Key Word"
	return t
}

func (t *KeyToken) assertKey(k string) bool {
	return t.val.(string) == k
}
func (t *KeyToken) assertKeys(ks ...string) bool {
	origin := t.val.(string)
	for _, k := range ks {
		if origin == k {
			return true
		}
	}
	return false
}
