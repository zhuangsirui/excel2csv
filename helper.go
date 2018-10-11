package main

import (
	"regexp"
	"strconv"
)

var floatRegex = regexp.MustCompile(`^\d+\.\d+$`)

// roundFloat 会依照 cell 值判断是否为疑似浮点数。如果是，则将其转换为浮点类型以
// 后，再功过 strconv.FormatFloat 转换回字符串以处理在 excel 中数值表示不精确的
// 问题
//
// excel 原始文件中出问题的数据举例：
//     `<c r="N65"><v>1.1000000000000001</v></c>`
func roundFloat(value string) string {
	if !floatRegex.MatchString(value) {
		return value
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
