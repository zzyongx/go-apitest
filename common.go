package apitest

import (
	"fmt"
	"strconv"
)

func interfaceToString(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	} else if i, ok := v.(int); ok {
		return strconv.Itoa(i), nil
	} else if i, ok := v.(int64); ok {
		return strconv.FormatInt(i, 10), nil
	} else {
		return "", fmt.Errorf("%v is not string", v)
	}
}

func MustInterfaceToString(v interface{}) string {
	if s, err := interfaceToString(v); err != nil {
		return fmt.Sprintf("%v", v)
	} else {
		return s
	}
}

func interfaceToInt(v interface{}) (int64, error) {
	if i32, ok := v.(int); ok {
		return int64(i32), nil
	} else if i64, ok := v.(int64); ok {
		return i64, nil
	} else if f32, ok := v.(float32); ok {
		return int64(f32), nil
	} else if f64, ok := v.(float64); ok {
		return int64(f64), nil
	} else {
		return 0, fmt.Errorf("%v is not int64", v)
	}
}

func interfaceToFloat(v interface{}) (float64, error) {
	if i32, ok := v.(int); ok {
		return float64(i32), nil
	} else if i64, ok := v.(int64); ok {
		return float64(i64), nil
	} else if f32, ok := v.(float32); ok {
		return float64(f32), nil
	} else if f64, ok := v.(float64); ok {
		return f64, nil
	} else {
		return 0, fmt.Errorf("%v is not float64", v)
	}
}
