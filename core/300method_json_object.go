package core

func evalJSONObjectMethod(obj JSONObject, method string, args []interface{}) (res *Value) {
	if method == "size" {
		return newQKValue(obj.size())
	}
	if method == "contain" {
		assert(len(args)<1, "method object.contain must has one parameter.")
		key, ok := args[0].(string)
		assert(!ok, "method object.contain parameter must be string type")
		return newQKValue(obj.exist(key))
	}
	if method == "remove" {
		assert(len(args)<1, "method object.remove must has one parameter.")
		key, ok := args[0].(string)
		assert(!ok, "method object.remove parameter must be string type")
		obj.remove(key)
		return
	}

	return nil
}