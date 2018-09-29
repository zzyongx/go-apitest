package apitest

import (
	"github.com/oliveagle/jsonpath"
)

type JsonExpect struct {
	*ApiExpect
	data string
	obj  interface{}
}

func (this *JsonExpect) MustJsonPathLookup(jpath string) interface{} {
	if pat, err := jsonpath.Compile(jpath); err != nil {
		this.t.Fatalf("invalid jsonpath %s error: %s", jpath, err)
	} else {
		if r, err := pat.Lookup(this.obj); err != nil {
			this.t.Fatalf("lookup %s in %s error: %s", jpath, this.data)
		} else {
			return r
		}
	}
	return nil
}

func (this *JsonExpect) Eq(jpath string, value interface{}) *JsonExpect {
	this.MustJsonPathLookup(jpath)
	return this
}
