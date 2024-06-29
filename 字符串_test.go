package main

import (
	"fmt"
	"testing"
)

func Test字符串包装类(t *testing.T) {
	s := Str("hello world")
	fmt.Println(s[1:4])
	fmt.Println(s < "zzz", s < "abc") // 支持string的操作符

	s = "你好世界"
	fmt.Println("rune子串的位置", s.Index("世界"))

	fmt.Println("\n================正则表达式搜索===================")
	var s1 Str = "aab aaab"
	fmt.Println(s1.FindReg("(?P<key>a*)(?P<name>b)").GroupByName("name"))
	fmt.Println(s1.ReplaceFuncReg("(?P<key>a*)(?P<name>b)", func(result *MatchResult) Str {
		return result.Expand("key:${key},name:${name}")
	}))

	fmt.Println("\n================遍历===================")
	s.ForEach(func(r rune) {
		fmt.Printf("%c\n", r)
	})
}
