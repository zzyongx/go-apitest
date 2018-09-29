package apitest

import (
	"fmt"
	"strconv"
)

func interfaceToString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	} else if i, ok := v.(int); ok {
		return strconv.Itoa(i)
	} else if i, ok := v.(int64); ok {
		return strconv.FormatInt(i, 10)
	} else {
		return fmt.Sprintf("%v", v)
	}
}
