package gods

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ===================================正则表达式查找方法===================================、

// MatchResult 定义了匹配结果的结构
type MatchResult struct {
	Result     Str
	Start      int
	End        int
	GroupCount int
	group      []Str
	groupStart []int
	groupEnd   []int
	pattern    Str
	text       Str
	match      []int
}

func NewMatchResult(text, pattern Str, match []int) *MatchResult {
	groupMatch := match[2:]
	start := match[0]
	end := match[1]
	groupCount := len(groupMatch) / 2

	result := &MatchResult{
		Result:     text[start:end],
		Start:      start,
		End:        end,
		GroupCount: groupCount,
		group:      make([]Str, groupCount),
		groupStart: make([]int, groupCount),
		groupEnd:   make([]int, groupCount),
		pattern:    pattern,
		text:       text,
		match:      match,
	}

	for i := 0; i < len(groupMatch); i += 2 {
		idx, start, end := i/2, groupMatch[i], groupMatch[i+1]
		result.group[idx] = text[start:end]
		result.groupStart[idx] = start
		result.groupEnd[idx] = end
	}
	return result
}

func (o *MatchResult) String() string { return fmt.Sprintf("%#v", *o) }
func (o *MatchResult) idx(i int) int {
	preIdx := i
	if i < 0 {
		i = o.GroupCount + i
	}
	if i < 0 || i >= o.GroupCount {
		panic("索引不合法：" + strconv.Itoa(preIdx))
	}
	return i
}

// 将组名称转换成组索引
func (o *MatchResult) idxByName(name Str) int {
	return regexp.MustCompile(o.pattern.S()).SubexpIndex(name.S()) - 1
}

func (o *MatchResult) Group(i int) Str               { return o.group[o.idx(i)] }
func (o *MatchResult) GroupStart(i int) int          { return o.groupStart[o.idx(i)] }
func (o *MatchResult) GroupEnd(i int) int            { return o.groupEnd[o.idx(i)] }
func (o *MatchResult) GroupByName(name Str) Str      { return o.Group(o.idxByName(name)) }
func (o *MatchResult) GroupStartByName(name Str) int { return o.GroupStart(o.idxByName(name)) }
func (o *MatchResult) GroupEndByName(name Str) int   { return o.GroupEnd(o.idxByName(name)) }

// 根据模板字符串展开 $1,$name等引用
func (o *MatchResult) Expand(template Str) Str {
	return Str(regexp.MustCompile(o.pattern.S()).ExpandString(nil, template.S(), o.text.S(), o.match))
}

// ===================================字符串包装类===================================

// 所有操作以rune为单位，而不是字节
// 支持+ < == [:]等运算符
// Str允许接收字符串字面量，但是不允许接收string类型变量
// 可以使用Str("hello")和Str('字')，根据字符串和字符创建Str对象
type Str string

func NewStrFromInt(i int) Str       { return Str(strconv.Itoa(i)) }
func NewStrFromUint(i uint) Str     { return Str(strconv.FormatUint(uint64(i), 10)) }
func NewStrFromFloat(f float64) Str { return Str(strconv.FormatFloat(f, 'f', -1, 64)) }

// --------------------基本属性和转换--------------------
func (o Str) S() string     { return string(o) }
func (o Str) Len() int      { return utf8.RuneCountInString(o.S()) }
func (o Str) LenBytes() int { return len(o) }

// 获取第i个rune，支持负值索引
func (o Str) Get(i int) rune  { return []rune(o)[o.idx(i)] }
func (o Str) ToRunes() []rune { return []rune(o) }
func (o Str) ToBytes() []byte { return []byte(o) }

func (o Str) ParseInt() int {
	if i, err := strconv.Atoi(string(o)); err != nil {
		panic(err)
	} else {
		return i
	}
}
func (o Str) ParseFloat() float64 {
	if f, err := strconv.ParseFloat(string(o), 64); err != nil {
		panic(err)
	} else {
		return f
	}
}

// 大小写转换
func (o Str) ToLower() Str { return Str(strings.ToLower(o.S())) }
func (o Str) ToUpper() Str { return Str(strings.ToUpper(o.S())) }

// 翻转字符串
func (o Str) Reverse() Str {
	var res []rune
	for i := o.Len() - 1; i >= 0; i-- {
		res = append(res, o.Get(i))
	}
	return Str(res)
}

// 截取子串
func (o Str) Slice(start, end int) Str { return o[o.idx(start):o.idx(end)] }

// 支持负值rune索引
func (o Str) idx(i int) int {
	preIdx := i
	if i < 0 {
		i += o.Len()
	}
	if i < 0 || i >= o.Len() {
		panic("索引超出范围：" + strconv.Itoa(preIdx))
	}
	return i
}

// 将byte索引转换成rune索引，如果是负值，则原样返回
func (o Str) toRuneIdx(i int) int {
	if i < 0 {
		return i
	}
	return utf8.RuneCountInString(o[:i].S())
}

// --------------------遍历--------------------

func (o Str) ForEach(f func(rune)) Str   { return o.ForEachIdxVal(func(i int, r rune) { f(r) }) }
func (o Str) ForEachIdx(f func(int)) Str { return o.ForEachIdxVal(func(i int, r rune) { f(i) }) }
func (o Str) ForEachIdxVal(f func(int, rune)) Str {
	for i, r := range o {
		f(i, r)
	}
	return o
}

// --------------------查询--------------------

func (o Str) Has(substr Str) bool            { return o.Index(substr) >= 0 }
func (o Str) HasAny(chars Str) bool          { return strings.ContainsAny(o.S(), chars.S()) }
func (o Str) HasRune(r rune) bool            { return strings.ContainsRune(o.S(), r) }
func (o Str) HasFunc(f func(rune) bool) bool { return strings.ContainsFunc(o.S(), f) }

// 返回以rune为单位的子串索引
func (o Str) Index(substr Str) int            { return o.toRuneIdx(strings.Index(o.S(), substr.S())) }
func (o Str) IndexAny(chars Str) int          { return o.toRuneIdx(strings.IndexAny(o.S(), chars.S())) }
func (o Str) IndexRune(r rune) int            { return o.toRuneIdx(strings.IndexRune(o.S(), r)) }
func (o Str) IndexFunc(f func(rune) bool) int { return o.toRuneIdx(strings.IndexFunc(o.S(), f)) }

func (o Str) StartsWith(prefix Str) bool { return strings.HasPrefix(o.S(), prefix.S()) }
func (o Str) EndsWith(suffix Str) bool   { return strings.HasSuffix(o.S(), suffix.S()) }
func (o Str) Count(substr Str) int       { return strings.Count(o.S(), substr.S()) }

// --------------------比较--------------------

// 不区分大小写比较是否相等
func (o Str) EqualFold(other Str) bool { return strings.EqualFold(o.S(), other.S()) }

// 根据比较结果返回0、-1、1，可以直接使用不等号运算符
func (o Str) Cmp(other Str) int { return strings.Compare(o.S(), other.S()) }

// --------------------分割--------------------
func (o Str) Split(sep Str) []Str      { return o.SplitN(sep, -1) }
func (o Str) SplitAfter(sep Str) []Str { return o.SplitAfterN(sep, -1) }
func (o Str) Fields() []Str {
	var res []Str
	for _, s := range strings.Fields(o.S()) {
		res = append(res, Str(s))
	}
	return res
}

func (o Str) SplitN(sep Str, n int) []Str        { return StrArr(strings.SplitN(o.S(), sep.S(), n)) }
func (o Str) SplitAfterN(sep Str, n int) []Str   { return StrArr(strings.SplitAfterN(o.S(), sep.S(), n)) }
func (o Str) FieldsFunc(f func(rune) bool) []Str { return StrArr(strings.FieldsFunc(o.S(), f)) }

// --------------------替换--------------------

// 替换所有的匹配项
func (o Str) Replace(oldStr, newStr Str) Str {
	return Str(strings.ReplaceAll(o.S(), oldStr.S(), newStr.S()))
}

// 替换至多n个匹配项
func (o Str) ReplaceN(oldStr, newStr Str, n int) Str {
	return Str(strings.Replace(o.S(), oldStr.S(), newStr.S(), n))
}
func (o Str) Map(f func(rune) rune) Str { return Str(strings.Map(f, o.S())) }
func (o Str) MapByPairs(oldnew ...string) Str {
	return Str(strings.NewReplacer(oldnew...).Replace(o.S()))
}

// --------------------修剪--------------------

func (o Str) TrimSpace() Str                 { return Str(strings.TrimSpace(o.S())) }
func (o Str) Trim(chars Str) Str             { return Str(strings.Trim(o.S(), chars.S())) }
func (o Str) TrimLeft(chars Str) Str         { return Str(strings.TrimLeft(o.S(), chars.S())) }
func (o Str) TrimRight(chars Str) Str        { return Str(strings.TrimRight(o.S(), chars.S())) }
func (o Str) TrimFunc(f func(rune) bool) Str { return Str(strings.TrimFunc(o.S(), f)) }
func (o Str) TrimPrefix(prefix Str) Str      { return Str(strings.TrimPrefix(o.S(), prefix.S())) }
func (o Str) TrimSuffix(suffix Str) Str      { return Str(strings.TrimSuffix(o.S(), suffix.S())) }

// --------------------正则表达式--------------------

// 判断是否包含子串
func (o Str) HasReg(pattern Str) bool { return regexp.MustCompile(pattern.S()).MatchString(o.S()) }

// 分割
func (o Str) SplitReg(sep Str) []Str { return o.SplitNReg(sep, -1) }

// 分割并指定分割结果数量，n<0时，获得所有分割结果
func (o Str) SplitNReg(sep Str, n int) []Str {
	var res []Str
	for _, s := range regexp.MustCompile(sep.S()).Split(o.S(), n) {
		res = append(res, Str(s))
	}
	return res
}

// 替换

// 替换所有匹配项，识别$1 $name ${name} 的替换
func (o Str) ReplaceReg(pattern, newStr Str) Str {
	return Str(regexp.MustCompile(pattern.S()).ReplaceAllString(o.S(), newStr.S()))
}

// 不识别分组引用的替换
func (o Str) ReplaceLiteralReg(pattern, newStr Str) Str {
	return Str(regexp.MustCompile(pattern.S()).ReplaceAllLiteralString(o.S(), newStr.S()))
}

// 使用函数替换，接收一个MatchResult对象，不识别分组引用，可以使用MatchResult的Expand方法展开模板
func (o Str) ReplaceFuncReg(pattern Str, f func(result *MatchResult) Str) Str {
	var sb strings.Builder
	results := o.FindAllReg(pattern)
	i := 0
	for _, result := range results {
		sb.WriteString(o[i:result.Start].S())
		sb.WriteString(f(result).S())
		i = result.End
	}
	sb.WriteString(o[i:].S())
	return Str(sb.String())
}

// Find 函数查找第一个匹配,返回 MatchResult。如果没有匹配，返回 nil
func (o Str) FindReg(pattern Str) *MatchResult {
	match := regexp.MustCompile(pattern.S()).FindStringSubmatchIndex(o.S())
	if match == nil {
		return nil
	}

	return NewMatchResult(o, pattern, match)
}

// FindAll 函数查找所有匹配,返回 MatchResult 数组。如果没有匹配，返回 nil
func (o Str) FindAllReg(pattern Str) []*MatchResult {
	matches := regexp.MustCompile(pattern.S()).FindAllStringSubmatchIndex(o.S(), -1)
	if matches == nil {
		return nil
	}

	var results []*MatchResult
	for _, match := range matches {
		results = append(results, NewMatchResult(o, pattern, match))
	}

	return results
}

// --------------------格式化IO--------------------
func (o Str) Printf(a ...any) { fmt.Printf(o.S(), a...) }
func (o Str) Print()          { fmt.Print(o.S()) }
func (o Str) Println()        { fmt.Println(o.S()) }
func (o Str) Scanf(a ...any) {
	_, err := fmt.Scanf(o.S(), a...)
	if err != nil {
		panic(err)
	}
}

func (o Str) Format(a ...any) Str { return Str(fmt.Sprintf(o.S(), a...)) }

// --------------------其他--------------------

func (o Str) Repeat(n int) Str { return Str(strings.Repeat(o.S(), n)) }

// ===================================构造Str数组===================================

// 将string切片转换为Str切片
func StrArr(strs []string) []Str {
	var res []Str
	for _, s := range strs {
		res = append(res, Str(s))
	}
	return res
}
