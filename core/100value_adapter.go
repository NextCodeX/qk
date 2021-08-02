package core

type ValueAdapter struct {}

func (valAdapter *ValueAdapter) isNULL() bool {
    return false
}
func (valAdapter *ValueAdapter) isInt() bool {
    return false
}
func (valAdapter *ValueAdapter) isFloat() bool {
    return false
}
func (valAdapter *ValueAdapter) isBoolean() bool {
    return false
}
func (valAdapter *ValueAdapter) isString() bool {
    return false
}
func (valAdapter *ValueAdapter) isAny() bool {
    return false
}
func (valAdapter *ValueAdapter) isClass() bool {
    return false
}
func (valAdapter *ValueAdapter) isJsonArray() bool {
    return false
}
func (valAdapter *ValueAdapter) isJsonObject() bool {
    return false
}
func (valAdapter *ValueAdapter) isFunction() bool {
    return false
}
func (valAdapter *ValueAdapter) isObject() bool {
    return false
}