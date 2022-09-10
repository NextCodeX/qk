package core

type ValueAdapter struct{}

func (this *ValueAdapter) isNULL() bool {
	return false
}
func (this *ValueAdapter) isByteArray() bool {
	return false
}
func (this *ValueAdapter) isInt() bool {
	return false
}
func (this *ValueAdapter) isFloat() bool {
	return false
}
func (this *ValueAdapter) isBoolean() bool {
	return false
}
func (this *ValueAdapter) isString() bool {
	return false
}
func (this *ValueAdapter) isAny() bool {
	return false
}
func (this *ValueAdapter) isJsonArray() bool {
	return false
}
func (this *ValueAdapter) isJsonObject() bool {
	return false
}
func (this *ValueAdapter) isFunction() bool {
	return false
}
func (this *ValueAdapter) isObject() bool {
	return false
}
