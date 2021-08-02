package core


func (obj *JSONObjectImpl) isObject() bool {
	return true
}

func (obj *JSONObjectImpl) returnFakeMethod(key string) Value {
	switch key {
	case "size":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return newQKValue(obj.size())
		})

	case "contain":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<1, "method object.contain must has one parameter.")
			key, ok := args[0].(string)
			assert(!ok, "method object.contain parameter must be string type")
			return newQKValue(obj.exist(key))
		})

	case "remove":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<1, "method object.remove must has one parameter.")
			key, ok := args[0].(string)
			assert(!ok, "method object.remove parameter must be string type")
			obj.remove(key)
			return nil
		})
	default:
		return NULL
	}
}