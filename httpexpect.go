package apitest

import (
	"net/http"
	"testing"
	"time"
)

type StatusExpect struct {
	*ApiExpect
	status int
}

func (this *StatusExpect) Eq(code int) *StatusExpect {
	if code == this.status {
		return this
	}
	this.Fatalf("status expect %d, got %d", code, this.status)
	return this
}

func (this *StatusExpect) EqAny(codes ...int) *StatusExpect {
	for _, code := range codes {
		if code == this.status {
			return this
		}
	}
	this.Fatalf("status expect any %v, got %d", codes, this.status)
	return this
}

func (this *StatusExpect) NotEq(code int) *StatusExpect {
	if code != this.status {
		return this
	}
	this.Fatalf("status %s expect not %d, got %d", code, this.status)
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
	this.Fatalf("header %v not exist", name)
	return this
}

func (this *HeadersExpect) ExistAny(names ...string) *HeadersExpect {
	for _, name := range names {
		if _, ok := this.headers[name]; ok {
			return this
		}
	}
	this.Fatalf("none of %v header exists", names)
	return this
}

func (this *HeadersExpect) NotExist(name string) *HeadersExpect {
	if _, ok := this.headers[name]; !ok {
		return this
	}
	this.Fatalf("header %s exist", name)
	return this
}

func (this *HeadersExpect) Eq(name, v string) *HeadersExpect {
	if values, ok := this.headers[name]; !ok {
		this.Fatalf("header %s not exist", name)
		return this
	} else {
		for _, value := range values {
			if v == value {
				return this
			}
		}
		this.Fatalf("header %s expect %s, got %v", name, v, values)
		return this
	}
}

func (this *HeadersExpect) EqAny(name string, vs ...string) *HeadersExpect {
	if values, ok := this.headers[name]; !ok {
		this.Fatalf("header %s not exist", name)
		return this
	} else {
		for _, value := range values {
			for _, v := range vs {
				if v == value {
					return this
				}
			}
		}
		this.Fatalf("header %s expect any %v, got %v", name, vs, values)
		return this
	}
}

type CookiesExpect struct {
	*ApiExpect
	name    string
	cookies []*http.Cookie
}

func (this *CookiesExpect) getCookie() []*http.Cookie {
	cookies := make([]*http.Cookie, 0)
	for _, cookie := range this.cookies {
		if cookie.Name == this.name {
			cookies = append(cookies, cookie)
		}
	}
	return cookies
}

func (this *CookiesExpect) StoreValue(value *string) *CookiesExpect {
	for _, cookie := range this.getCookie() {
		*value = cookie.Value
		return this
	}
	this.Fatalf("cookie %s's value was not found, got %v", this.name, this.cookies)
	return this
}

func (this *CookiesExpect) Value(value string) *CookiesExpect {
	for _, cookie := range this.getCookie() {
		if cookie.Value == value {
			return this
		}
	}
	this.Fatalf("cookie %s's value expect %s, got %v", this.name, value, this.cookies)
	return this
}

func (this *CookiesExpect) Domain(domain string) *CookiesExpect {
	for _, cookie := range this.getCookie() {
		if cookie.Domain == domain {
			return this
		}
	}
	this.Fatalf("cookie %s's domain expect %s, got %v", this.name, domain, this.cookies)
	return this
}

func (this *CookiesExpect) Expires(startAt, endAt time.Time) *CookiesExpect {
	for _, cookie := range this.getCookie() {
		if cookie.Expires.After(startAt) && cookie.Expires.Before(endAt) {
			return this
		}
	}
	this.Fatalf("cookie %s's expires expect between %s, got %v", this.name, startAt, endAt, this.cookies)
	return this
}

func (this *CookiesExpect) Test(test func(t *testing.T, cookies []*http.Cookie)) *CookiesExpect {
	test(this.t, this.getCookie())
	return this
}
