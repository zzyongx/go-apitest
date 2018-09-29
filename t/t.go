package t

import (
	apitest "github.com/zzyongx/go-apitest"
	"net/http"
	"time"
)

var api *apitest.Api

func init() {
	api = apitest.NewHttpTest("http://127.0.0.1:8282")

	// simulation test server
	http.HandleFunc("/static/config.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`
{"user": "zzyongx", "color": ["green", "gray"]}
`))
	})

	go func() {
		if err := http.ListenAndServe(":8282", nil); err != nil {
			panic("start mock server failed: " + err.Error())
		}
	}()
	time.Sleep(time.Second) // wait listen ok
}
