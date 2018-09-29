package t

import (
	"testing"
	_ "time"
)

func TestLogin(t *testing.T) {
	api.Get("/static/config.json").Expect(t).
		Status().Eq(200).
		Headers().Eq("Content-Type", "application/json; charset=utf-8").
		EqAny("Content-Type", "application/json", "application/json; charset=utf-8").
		Exist("Content-Length").ExistAny("NotFound", "Content-Length").
		Json().Eq("$.user", "zzyongx")

	api.Post("/api/login")
}

func TestMockServer(t *testing.T) {
	//	time.Sleep(time.Second * 300)
}
