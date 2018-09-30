package t

import (
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	api.Get("/static/config.json").Expect(t).
		Status().Eq(200).
		Headers().Eq("Content-Type", "application/json; charset=utf-8").
		EqAny("Content-Type", "application/json", "application/json; charset=utf-8").
		Exist("Content-Length").ExistAny("NotFound", "Content-Length").
		Json().Eq("$.user", "zzyongx")

	api.Post("/api/login").Form("user", "zzyongx").Form("password", "123456").Expect(t).
		Json().Eq("$.code", 403)

	var token string
	oneYearAfter := time.Now().Add(365 * 24 * time.Hour)
	api.Post("/api/login").Form("user", "zzyongx").Form("password", "123465").Expect(t).
		Cookies("user").Value("zzyongx").Domain("example.com").Expires(time.Now(), oneYearAfter).
		Cookies("token").Value("90#@xw").Domain("example.com").Expires(time.Now(), oneYearAfter).StoreValue(&token).
		Json().Eq("$.code", 0)

	if token != "90#@xw" {
		t.Fatalf("cookie storevalue bug, expect 90#@xw, got %s", token)
	}

	book1Name, book2Name := "The Art of Computer Program", "Structure and Interpretation of Computer Programs"

	api.Post("/api/bookself/book").Form("name", book1Name).Expect(t).
		Json().Eq("$.code", 401)

	var book1Id, book2Id int64
	api.Post("/api/bookself/book").Form("name", book1Name).Cookie("user", "zzyongx").Cookie("token", token).Expect(t).
		Json().StoreInt("$.data.id", &book1Id)

	api.Cookies["user"] = "zzyongx"
	api.Cookies["token"] = token

	api.Post("/api/bookself/book").Form("name", book2Name).Expect(t).
		Json().StoreInt("$.data.id", &book2Id)

	api.Get("/api/bookself/books?order=desc").Param("limit", 1).Expect(t).
		Json().Eq("$.data[0].name", book1Name)

	book1Name = "The Art of Computer Programing"
	api.Put("/api/bookself/book/{id}", book1Id).Form("name", book1Name).Expect(t).
		Status().Eq(200)

	api.Get("/api/bookself/books?order=desc").Param("limit", 1).Expect(t).
		Json().Eq("$.data[0].name", book1Name)

	api.Get("/api/bookself/books").Param("order", "asc").Expect(t).
		Json().Eq("$.data[-1].name", book1Name)

	api.Delete("/api/bookself/book/{id}", book1Id).Expect(t).
		Json().Eq("$.code", 0)

	api.Get("/api/bookself/books").Expect(t).
		Json().Eq("$.data.length()", 1)
}

func TestMockServer(t *testing.T) {
	//	time.Sleep(time.Second * 300)
}
