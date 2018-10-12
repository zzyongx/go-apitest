package apitest

import (
	"regexp"
)

type PlainExpect struct {
	*ApiExpect
	data []byte
}

func (this *PlainExpect) Contains(restr string) *PlainExpect {
	re, err := regexp.Compile(restr)
	if err != nil {
		this.Fatalf("compile regex %s error: %s", restr, err)
	}
	if !re.Match(this.data) {
		this.Fatalf("expect %s contains %s", string(this.data), restr)
	}
	return this
}

func (this *PlainExpect) NotContains(restr string) *PlainExpect {
	re, err := regexp.Compile(restr)
	if err != nil {
		this.Fatalf("compile regex %s error: %s", restr, err)
	}
	if re.Match(this.data) {
		this.Fatalf("expect %s doesn't contain %s", string(this.data), restr)
	}
	return this
}
