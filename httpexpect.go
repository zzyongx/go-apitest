package apitest

import (
	"net/http"
)

type StatusExpect struct {
	*ApiExpect
	status int
}

func (this *StatusExpect) Eq(code int) *StatusExpect {
	if code == this.status {
		return this
	}
	this.t.Fatalf("req: %s > status expect %d, got %d", this.Req(), code, this.status)
	return this
}

func (this *StatusExpect) EqAny(codes ...int) *StatusExpect {
	for _, code := range codes {
		if code == this.status {
			return this
		}
	}
	this.t.Fatalf("req: %s > status expect any %v, got %d", this.Req(), codes, this.status)
	return this
}

func (this *StatusExpect) NotEq(code int) *StatusExpect {
	if code != this.status {
		return this
	}
	this.t.Fatalf("req: %s > status %s expect not %d, got %d", this.Req(), code, this.status)
	return this
}

type HeadersExpect struct {
	*ApiExpect
	headers http.Header
}

func (this *HeadersExpect) Exist(name string) *HeadersExpect {
	if _, ok := this.headers[name]; ok {
		return this
	}
	this.t.Fatalf("req: %s > header %v not exist", this.Req(), name)
	return this
}

func (this *HeadersExpect) ExistAny(names ...string) *HeadersExpect {
	for _, name := range names {
		if _, ok := this.headers[name]; ok {
			return this
		}
	}
	this.t.Fatalf("req: %s > none of %v header exists", this.Req(), names)
	return this
}

func (this *HeadersExpect) NotExist(name string) *HeadersExpect {
	if _, ok := this.headers[name]; !ok {
		return this
	}
	this.t.Fatalf("req: %s > header %s exist", this.Req(), name)
	return this
}

func (this *HeadersExpect) Eq(name, v string) *HeadersExpect {
	if values, ok := this.headers[name]; !ok {
		this.t.Fatalf("req: %s > header %s not exist", this.Req(), name)
		return this
	} else {
		for _, value := range values {
			if v == value {
				return this
			}
		}
		this.t.Fatalf("req: %s > header %s expect %s, got %v", this.Req(), name, v, values)
		return this
	}
}

func (this *HeadersExpect) EqAny(name string, vs ...string) *HeadersExpect {
	if values, ok := this.headers[name]; !ok {
		this.t.Fatalf("req: %s > header %s not exist", this.Req(), name)
		return this
	} else {
		for _, value := range values {
			for _, v := range vs {
				if v == value {
					return this
				}
			}
		}
		this.t.Fatalf("req: %s > header %s expect any %v, got %v", this.Req(), name, vs, values)
		return this
	}
}
