package apitest

import (
	"github.com/oliveagle/jsonpath"
	"reflect"
	"strings"
	"testing"
)

type JsonExpect struct {
	*ApiExpect
	data string
	obj  interface{}
}

type StringMatchFunc func(string, string) bool
type InterfaceMatchFunc func(interface{}, interface{}) bool

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

func (this *JsonExpect) Match(jpath string, matchFunc InterfaceMatchFunc, value interface{}) *JsonExpect {
	inter := this.MustJsonPathLookup(jpath)
	if inter == nil {
		if matchFunc(inter, value) {
			return this
		} else {
			this.Fatalf("lookup %s in %s error: null", jpath, this.data)
		}
	} else if iv, err := interfaceToString(inter); err == nil {
		if v, err := interfaceToString(value); err == nil {
			if matchFunc(iv, v) {
				return this
			} else {
				this.Fatalf("lookup %s in %s error: expect %s, got %s", jpath, this.data, v, iv)
			}
		} else {
			this.Fatalf("lookup %s in %s error: expect %v, got string %s", jpath, this.data, value, iv)
		}
	} else if iv, err := interfaceToInt(inter); err == nil {
		if v, err := interfaceToInt(value); err == nil {
			if matchFunc(iv, v) {
				return this
			} else {
				this.Fatalf("lookup %s in %s error: expect %d, got %d", jpath, this.data, v, iv)
			}
		} else {
			this.Fatalf("lookup %s in %s error: expect %v, got int64 %d", jpath, this.data, value, iv)
		}
	} else if iv, err := interfaceToFloat(inter); err == nil {
		if v, err := interfaceToFloat(value); err == nil {
			if matchFunc(iv, v) {
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

func (this *JsonExpect) Eq(jpath string, value interface{}) *JsonExpect {
	return this.Match(jpath, func(x interface{}, y interface{}) bool {
		if xx, ok := x.(string); ok {
			yy, _ := y.(string)
			return xx == yy
		} else if xx, ok := x.(int64); ok {
			yy, _ := y.(int64)
			return xx == yy
		} else if xx, ok := x.(float64); ok {
			yy, _ := y.(float64)
			return xx == yy
		} else {
			return x == y
		}
	}, value)
}

func (this *JsonExpect) MatchAnyString(jpath string, matchFunc StringMatchFunc, values ...string) *JsonExpect {
	inter := this.MustJsonPathLookup(jpath)
	v := MustInterfaceToString(inter)
	for _, value := range values {
		if matchFunc(v, value) {
			return this
		}
	}
	this.Fatalf("lookup %s in %s error: expect any %v, got %s", jpath, this.data, values, v)
	return this
}

func (this *JsonExpect) EqAnyString(jpath string, values ...string) *JsonExpect {
	return this.MatchAnyString(jpath, func(got, want string) bool { return got == want }, values...)
}

func (this *JsonExpect) EqAnyInt(jpath string, values ...int64) *JsonExpect {
	inter := this.MustJsonPathLookup(jpath)
	v, err := interfaceToInt(inter)
	if err != nil {
		this.Fatalf("lookup %s in %s error: %s", jpath, this.data, err)
	}

	for _, value := range values {
		if value == v {
			return this
		}
	}
	this.Fatalf("lookup %s in %s error: expect any %v, got %s", jpath, this.data, values, v)
	return this
}

func (this *JsonExpect) EqAnyFloat(jpath string, values ...float64) *JsonExpect {
	inter := this.MustJsonPathLookup(jpath)
	v, err := interfaceToFloat(inter)
	if err != nil {
		this.Fatalf("lookup %s in %s error: %s", jpath, this.data, err)
	}

	for _, value := range values {
		if value == v {
			return this
		}
	}
	this.Fatalf("lookup %s in %s error: expect any %v, got %s", jpath, this.data, values, v)
	return this
}

func (this *JsonExpect) ContainsString(jpath string, value string) *JsonExpect {
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

func (this *JsonExpect) IsNull(jpath string) *JsonExpect {
	v := this.MustJsonPathLookup(jpath)
	if v != nil {
		this.Fatalf("lookup %s in %s is null", jpath, this.data)
	}
	return this
}

func (this *JsonExpect) NotNull(jpath string) *JsonExpect {
	v := this.MustJsonPathLookup(jpath)
	if v == nil {
		this.Fatalf("lookup %s in %s is null", jpath, this.data)
	}
	return this
}

func (this *JsonExpect) MustJsonPahAsInt(jpath string) int64 {
	i, err := interfaceToInt(this.MustJsonPathLookup(jpath))
	if err != nil {
		this.Fatalf("lookup %s in %s is not int", jpath, this.data)
	}
	return i
}

func (this *JsonExpect) Gt(jpath string, value int64) *JsonExpect {
	if i := this.MustJsonPahAsInt(jpath); !(i > value) {
		this.Fatalf("%d was not Gt %d", i, value)
	}
	return this
}

func (this *JsonExpect) Ge(jpath string, value int64) *JsonExpect {
	if i := this.MustJsonPahAsInt(jpath); !(i >= value) {
		this.Fatalf("%d was not Gt %d", i, value)
	}
	return this
}

func (this *JsonExpect) Lt(jpath string, value int64) *JsonExpect {
	if i := this.MustJsonPahAsInt(jpath); !(i < value) {
		this.Fatalf("%d was not Lt %d", i, value)
	}
	return this
}

func (this *JsonExpect) Le(jpath string, value int64) *JsonExpect {
	if i := this.MustJsonPahAsInt(jpath); !(i <= value) {
		this.Fatalf("%d was not Le %d", i, value)
	}
	return this
}

func (this *JsonExpect) MustJsonPahAsStringArray(jpath string) []string {
	inter := this.MustJsonPathLookup(jpath)
	values, err := interfaceToStrings(inter)
	if err != nil {
		this.Fatalf("lookup %s in %s error: %s", jpath, this.data, err)
	}
	return values
}

func (this *JsonExpect) MustJsonPahAsIntArray(jpath string) []int64 {
	inter := this.MustJsonPathLookup(jpath)
	values, err := interfaceToInts(inter)
	if err != nil {
		this.Fatalf("lookup %s in %s error: %s", jpath, this.data, err)
	}
	return values
}

func (this *JsonExpect) MustJsonPahAsFloatArray(jpath string) []float64 {
	inter := this.MustJsonPathLookup(jpath)
	values, err := interfaceToFloats(inter)
	if err != nil {
		this.Fatalf("lookup %s in %s error: %s", jpath, this.data, err)
	}
	return values
}

func (this *JsonExpect) StringContains(jpath string, value string) *JsonExpect {
	values := this.MustJsonPahAsStringArray(jpath)
	if !stringContains(values, value) {
		this.Fatalf("looup %s in %s error: %v !contains %s", jpath, this.data, values, value)
	}
	return this
}

func (this *JsonExpect) IntContains(jpath string, value int64) *JsonExpect {
	values := this.MustJsonPahAsIntArray(jpath)
	if !intContains(values, value) {
		this.Fatalf("looup %s in %s error: %v !contains %s", jpath, this.data, values, value)
	}
	return this
}

func (this *JsonExpect) FloatContains(jpath string, value float64) *JsonExpect {
	values := this.MustJsonPahAsFloatArray(jpath)
	if !floatContains(values, value) {
		this.Fatalf("looup %s in %s error: %v !contains %s", jpath, this.data, values, value)
	}
	return this
}

func (this *JsonExpect) Contains(jpath string, value interface{}) *JsonExpect {
	if v, ok := value.(string); ok {
		return this.StringContains(jpath, v)
	} else if v, ok := value.(int64); ok {
		return this.IntContains(jpath, v)
	} else if v, ok := value.(int); ok {
		return this.IntContains(jpath, int64(v))
	} else if v, ok := value.(float64); ok {
		return this.FloatContains(jpath, v)
	} else if v, ok := value.(float32); ok {
		return this.FloatContains(jpath, float64(v))
	} else {
		this.Fatalf("unsupport value type %s", reflect.TypeOf(value))
	}
	return this
}

func (this *JsonExpect) StringOnlyContains(jpath string, value string) *JsonExpect {
	values := this.MustJsonPahAsStringArray(jpath)
	if !stringOnlyContains(values, value) {
		this.Fatalf("lookup %s in %s error: %v contains other than %s", jpath, this.data, values, value)
	}
	return this
}

func (this *JsonExpect) IntOnlyContains(jpath string, value int64) *JsonExpect {
	values := this.MustJsonPahAsIntArray(jpath)
	if !intOnlyContains(values, value) {
		this.Fatalf("lookup %s in %s error: %v contains other than %s", jpath, this.data, values, value)
	}
	return this
}

func (this *JsonExpect) FloatOnlyContains(jpath string, value float64) *JsonExpect {
	values := this.MustJsonPahAsFloatArray(jpath)
	if !floatOnlyContains(values, value) {
		this.Fatalf("lookup %s in %s error: %v contains other than %s", jpath, this.data, values, value)
	}
	return this
}

func (this *JsonExpect) OnlyContains(jpath string, value interface{}) *JsonExpect {
	if v, ok := value.(string); ok {
		return this.StringOnlyContains(jpath, v)
	} else if v, ok := value.(int64); ok {
		return this.IntOnlyContains(jpath, v)
	} else if v, ok := value.(int); ok {
		return this.IntOnlyContains(jpath, int64(v))
	} else if v, ok := value.(float64); ok {
		return this.FloatOnlyContains(jpath, v)
	} else if v, ok := value.(float32); ok {
		return this.FloatOnlyContains(jpath, float64(v))
	} else {
		this.Fatalf("unsupport value type %s", reflect.TypeOf(value))
	}
	return this
}

func (this *JsonExpect) Test(jpath string, test func(t *testing.T, v interface{})) *JsonExpect {
	test(this.t, this.MustJsonPathLookup(jpath))
	return this
}
