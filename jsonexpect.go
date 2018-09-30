package apitest

import (
	"github.com/oliveagle/jsonpath"
	"reflect"
	"strings"
)

type JsonExpect struct {
	*ApiExpect
	data string
	obj  interface{}
}

func (this *JsonExpect) MustJsonPathLookup(jpath string) interface{} {
	defer func() {
		if r := recover(); r != nil {
			//			fmt.Println(r)
		}
	}()

	length := false
	if strings.HasSuffix(jpath, ".length()") {
		jpath = jpath[0 : len(jpath)-len(".length()")]
		length = true
	}

	if pat, err := jsonpath.Compile(jpath); err != nil {
		this.Fatalf("invalid jsonpath %s error: %s", jpath, err)
	} else {
		if r, err := pat.Lookup(this.obj); err != nil {
			this.Fatalf("lookup %s in %s error: %s", jpath, this.data, err)
		} else {
			maySlice := reflect.ValueOf(r)
			if maySlice.Kind() == reflect.Slice && length {
				return maySlice.Len()
			} else {
				return r
			}
		}
	}
	return nil
}

func (this *JsonExpect) Eq(jpath string, value interface{}) *JsonExpect {
	inter := this.MustJsonPathLookup(jpath)
	if inter == nil {
		if value == nil {
			return this
		} else {
			this.Fatalf("lookup %s in %s error: null", jpath, this.data)
		}
	} else if iv, err := interfaceToString(inter); err == nil {
		if v, err := interfaceToString(value); err == nil {
			if iv == v {
				return this
			} else {
				this.Fatalf("lookup %s in %s error: expect %s, got %s", jpath, this.data, v, iv)
			}
		} else {
			this.Fatalf("lookup %s in %s error: expect %v, got string %s", jpath, this.data, value, iv)
		}
	} else if iv, err := interfaceToInt(inter); err == nil {
		if v, err := interfaceToInt(value); err == nil {
			if int64(iv) == v {
				return this
			} else {
				this.Fatalf("lookup %s in %s error: expect %d, got %d", jpath, this.data, v, iv)
			}
		} else {
			this.Fatalf("lookup %s in %s error: expect %v, got int64 %d", jpath, this.data, value, iv)
		}
	} else if iv, err := interfaceToFloat(inter); err == nil {
		if v, err := interfaceToFloat(value); err == nil {
			if float64(iv) == v {
				return this
			} else {
				this.Fatalf("lookup %s in %s error: expect %f, got %f", jpath, this.data, v, iv)
			}
		} else {
			this.Fatalf("lookup %s in %s error: expect %v, got float64 %f", jpath, this.data, value, iv)
		}
	} else {
		this.Fatalf("lookup %s in %s error: expect %v, got %v", jpath, this.data, value, inter)
	}
	return this
}

func intStore(i32 int, i *int64, f *float64) {
	if i != nil {
		*i = int64(i32)
	} else if f != nil {
		*f = float64(i32)
	}
}

func int64Store(i64 int64, i *int64, f *float64) {
	if i != nil {
		*i = int64(i64)
	} else if f != nil {
		*f = float64(i64)
	}
}

func floatStore(f32 float32, i *int64, f *float64) {
	if i != nil {
		*i = int64(f32)
	} else if f != nil {
		*f = float64(f32)
	}
}

func float64Store(f64 float64, i *int64, f *float64) {
	if i != nil {
		*i = int64(f64)
	} else if f != nil {
		*f = float64(f64)
	}
}

func (this *JsonExpect) storeIntFloat(jpath string, i *int64, f *float64) *JsonExpect {
	v := this.MustJsonPathLookup(jpath)

	if i32, ok := v.(int); ok {
		intStore(i32, i, f)
	} else if i64, ok := v.(int64); ok {
		int64Store(i64, i, f)
	} else if f32, ok := v.(float32); ok {
		floatStore(f32, i, f)
	} else if f64, ok := v.(float64); ok {
		float64Store(f64, i, f)
	} else {
		this.Fatalf("lookup %s in %s error: %#v's type is %s", jpath, this.data, v, reflect.ValueOf(v).String())
	}
	return this
}

func (this *JsonExpect) StoreInt(jpath string, ret *int64) *JsonExpect {
	return this.storeIntFloat(jpath, ret, nil)
}

func (this *JsonExpect) StoreFloat(jpath string, ret *float64) *JsonExpect {
	return this.storeIntFloat(jpath, nil, ret)
}

func (this *JsonExpect) StoreString(jpath string, ret *string) *JsonExpect {
	v := this.MustJsonPathLookup(jpath)
	if s, ok := v.(string); ok {
		*ret = s
	} else {
		this.Fatalf("lookup %s in %s error: %#v's type is %s", jpath, this.data, v, reflect.ValueOf(v).String())
	}
	return this
}
