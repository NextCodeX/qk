package core

import (
	"fmt"
	"os"
	"path/filepath"
)

type Interpreter struct {
	main          *CustomFunction
	scriptFile    string
	scriptContent []byte
}

func newInterpreter(scriptPath string) *Interpreter {
	bs, err := os.ReadFile(scriptPath)
	if err != nil {
		errorf("failed to read %v;\n error info: %v", scriptPath, err)
	}
	return &Interpreter{
		main:          newMainFunction(),
		scriptFile:    scriptPath,
		scriptContent: bs,
	}
}

func (this *Interpreter) tryIgnoreFirstLine() []byte {
	// 用于linux文本可执行文件
	// 忽略以‘#’字符开头的首行
	if len(this.scriptContent) < 1 {
		return this.scriptContent
	}
	if this.scriptContent[0] == '#' {
		for i, ch := range this.scriptContent {
			if ch == '\n' {
				return this.scriptContent[i+1:]
			}
		}
	}
	return this.scriptContent
}

// 脚本解析执行
func (this *Interpreter) run() {
	// 忽略以‘#’字符开头的首行
	bs := this.tryIgnoreFirstLine()

	// 词法分析
	ts := ParseTokens(bs)
	if DEBUG {
		printTokensByLine(ts)
	}

	// 语法分析
	this.main.setTokenList(ts)
	Compile(this.main)

	// 初始化内置变量及函数
	this.initInternalVars()
	// 程序执行
	this.main.execute()

	// 等待所有协程执行完，再结束程序
	goroutineWaiter.Wait()
}

// 用于加载内置变量与内置函数
func (this *Interpreter) initInternalVars() {
	initVars := make(map[string]Value)
	// 添加命令行参数
	rawCmdArgs := os.Args
	if len(rawCmdArgs) > 2 {
		var arr []Value
		for _, rawCmdArg := range rawCmdArgs[2:] {
			arr = append(arr, newQKValue(rawCmdArg))
		}
		initVars["cmdArgs"] = array(arr)
	} else {
		initVars["cmdArgs"] = emptyArray()
	}

	// 提供qk执行文件所在的目录
	if executable, err := os.Executable(); err == nil {
		rootDir := filepath.Dir(executable)
		initVars["qkDir"] = newQKValue(rootDir)
	} else {
		fmt.Println(err)
	}

	// 当前命令行所在的路径，与`pwd`等同(工作路径)
	if cwd, err := os.Getwd(); err == nil {
		initVars["pwd"] = newQKValue(cwd)
	} else {
		fmt.Println(err)
	}

	// 提供当前脚本文件所在的目录
	if dir, err := filepath.Abs(filepath.Dir(this.scriptFile)); err == nil {
		initVars["root"] = newQKValue(dir)
	} else {
		fmt.Println(err)
	}

	// 添加POST请求常用的Content-Type
	mimes := make(map[string]Value)
	mimes["txt"] = newQKValue("text/plain")
	mimes["json"] = newQKValue("application/json")
	mimes["form"] = newQKValue("application/x-www-form-urlencoded")
	mimes["data"] = newQKValue("multipart/form-data")
	initVars["mime"] = jsonObject(mimes)

	// 将内部函数注册到主函数(main)的内部栈中
	fnSet := newInternalFunctionSet(this)
	for fname, ifunc := range fnSet.internalFuntions {
		initVars[fname] = ifunc
	}

	this.main.setInternalVars(initVars)
}
