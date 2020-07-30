package core

import "module"

var moduleFuncs map[string]*module.FunctionExecutor

func init()  {
	moduleFuncs = module.Load()
}