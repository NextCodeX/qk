package core

import (
	"testing"
	"fmt"
)

func TestParse4ComplexExpr(t *testing.T)  {
	var arr []Token
	arr = append(arr, varToken("a"))
	arr = append(arr, symbolToken("["))
	arr = append(arr, newToken("1", Int))
	arr = append(arr, symbolToken("]"))
	arr = parse4ComplexTokens(arr)
	fmt.Println(tokensString(arr), arr[0].TokenTypeName())
}

func Test_tailMerge2(t *testing.T)  {
	var arr []Token
	arr = append(arr, varToken("a"))
	arr = append(arr, symbolToken("="))
	arr = tailMerge2(arr, AddSelf)
	fmt.Println(tokensString(arr))
	fmt.Println(arr[len(arr)-1].TokenTypeName())
}

func Test_tailMerge(t *testing.T)  {
	var arr []Token
	arr = append(arr, varToken("a"))
	arr = append(arr, symbolToken("="))
	arr = tailMerge(arr, symbolToken("="))
	fmt.Println(tokensString(arr))
}
