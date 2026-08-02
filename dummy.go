package dummy

import (
	"context"
	"net/http"
)

func BadContext() {
	ctx := context.Background()
	_ = ctx
}

func BadLoop() {
	for i := 0; i < 10; i++ {
		go func() {
			// leak
		}()
	}
}

func BadServer() {
	http.ListenAndServe(":8080", nil)
}

type DB struct{}

func (db *DB) Query(q string) {}

func BadDB() {
	db := &DB{}
	db.Query("SELECT * FROM users") // SELECT *

	for i := 0; i < 10; i++ {
		db.Query("SELECT id FROM users WHERE id = 1") // Query inside loop
	}
}

func BadAPI(w http.ResponseWriter, r *http.Request) {
	// massive handler > 20 lines
	a := 1
	b := 2
	c := 3
	d := 4
	e := 5
	f := 6
	g := 7
	h := 8
	i := 9
	j := 10
	k := 11
	l := 12
	m := 13
	n := 14
	o := 15
	p := 16
	q := 17
	s := 18
	t := 19
	u := 20
	v := 21
	_ = a + b + c + d + e + f + g + h + i + j + k + l + m + n + o + p + q + s + t + u + v
}
