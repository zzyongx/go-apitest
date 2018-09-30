package t

import (
	"encoding/json"
	"fmt"
	apitest "github.com/zzyongx/go-apitest"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var api *apitest.Api

const TOKEN string = "90#@xw"

type bookSt struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var BOOKSELF map[string][]*bookSt = make(map[string][]*bookSt, 0)

func getCookie(r *http.Request, name string) string {
	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func init() {
	api = apitest.NewHttpTest("http://127.0.0.1:8282")

	// simulation test server
	http.HandleFunc("/static/config.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`
{"user": "zzyongx", "color": ["green", "gray"]}
`))
	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		user := r.PostForm.Get("user")
		pass := r.PostForm.Get("password")

		if user == "zzyongx" && pass == "123465" {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			userCookie := http.Cookie{Name: "user", Value: "zzyongx", Expires: expiration,
				Domain: "example.com", Path: "/"}
			tokenCookie := http.Cookie{Name: "token", Value: TOKEN, Expires: expiration,
				Domain: "example.com", Path: "/"}
			http.SetCookie(w, &userCookie)
			http.SetCookie(w, &tokenCookie)
			w.Write([]byte(`{"code": 0, "message": "OK"}`))
		} else {
			w.Write([]byte(`{"code": 403, "message": "authentication error"}`))
		}
	})

	http.HandleFunc("/api/bookself/book", func(w http.ResponseWriter, r *http.Request) {
		token := getCookie(r, "token")
		if token != TOKEN {
			w.Write([]byte(`{"code":401, "message": "unauthorized"}`))
			return
		}
		user := getCookie(r, "user")

		r.ParseForm()
		book := r.PostForm.Get("name")

		id := len(BOOKSELF[user])
		BOOKSELF[user] = append(BOOKSELF[user], &bookSt{Id: id, Name: book})

		w.Write([]byte(fmt.Sprintf(`{"code": 0, "data": {"id": %d}}`, id)))
	})

	http.HandleFunc("/api/bookself/book/0", func(w http.ResponseWriter, r *http.Request) {
		user := getCookie(r, "user")

		if r.Method == "PUT" {
			r.ParseForm()
			book := r.PostForm.Get("name")
			for i, bookSt := range BOOKSELF[user] {
				if bookSt.Id == 0 {
					bookSt.Name = book
					BOOKSELF[user][i] = bookSt
				}
			}
		} else if r.Method == "DELETE" {
			books := BOOKSELF[user]
			BOOKSELF[user] = nil
			for _, bookSt := range books {
				if bookSt.Id != 0 {
					BOOKSELF[user] = append(BOOKSELF[user], bookSt)
				}
			}
		}
		w.Write([]byte(fmt.Sprintf(`{"code": 0, "message": "OK"}`)))
	})

	http.HandleFunc("/api/bookself/books", func(w http.ResponseWriter, r *http.Request) {
		user := getCookie(r, "user")
		books := BOOKSELF[user]

		limitStr := r.URL.Query().Get("limit")
		limit := int64(-1)
		if len(limitStr) > 0 {
			limit, _ = strconv.ParseInt(limitStr, 10, 64)
		}

		if order := r.URL.Query().Get("order"); order != "" {
			var asc bool
			if order == "asc" {
				asc = true
			} else if order == "desc" {
				asc = false
			} else {
				w.Write([]byte(fmt.Sprintf(`{"code": 400, "message": "order must be asc or desc"}`)))
				return
			}

			sort.Slice(books, func(i, j int) bool {
				x, y := books[i], books[j]
				if asc {
					return x.Name < y.Name
				} else {
					return !(x.Name < y.Name)
				}
			})
		}

		if limit > 0 && limit < int64(len(books)) {
			books = books[0:limit]
		}

		booksJson, _ := json.Marshal(books)
		w.Write([]byte(fmt.Sprintf(`{"code": 0, "data": %s}`, string(booksJson))))
	})

	go func() {
		if err := http.ListenAndServe(":8282", nil); err != nil {
			panic("start mock server failed: " + err.Error())
		}
	}()
	time.Sleep(time.Second) // wait listen ok
}
