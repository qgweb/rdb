// ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝
//          常用函数
//　＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝
package common

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"reflect"
	"time"
	"unicode"
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// 打印完整sql
func PrintSql(sql string, values ...interface{}) {
	var formattedValues []string
	sqlRegexp := regexp.MustCompile(`(\$\d+)|\?`)
	for _, value := range values {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() {
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok {
				formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
			} else if b, ok := value.([]byte); ok {
				if str := string(b); isPrintable(str) {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
				} else {
					formattedValues = append(formattedValues, "'<binary>'")
				}
			} else if r, ok := value.(driver.Valuer); ok {
				if value, err := r.Value(); err == nil && value != nil {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			} else {
				formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
			}
		} else {
			formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
		}
	}

	var formattedValuesLength = len(formattedValues)
	var nsql = ""
	for index, value := range sqlRegexp.Split(sql, -1) {
		nsql += value
		if index < formattedValuesLength {
			nsql += formattedValues[index]
		}
	}

	fmt.Println("\033[0;31m[", nsql, "]\033[0m")
}
