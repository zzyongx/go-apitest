package apitest

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type Db struct {
	dbs     map[string]*sql.DB
	request *DbRequest
}

type DbRequest struct {
	op     string
	db     *sql.DB
	sql    string
	params []interface{}
}

func NewDbReqest(op string, db *sql.DB, sql string, params []interface{}) *DbRequest {
	return &DbRequest{
		op:     op,
		db:     db,
		sql:    sql,
		params: params,
	}
}

func NewDbTest(cnf *IniCnf, keys ...string) (*Db, error) {
	ctx := &Db{
		dbs: make(map[string]*sql.DB, 0),
	}
	for _, key := range keys {
		db, err := sql.Open("mysql", cnf.GetString(key))
		if err != nil {
			return nil, err
		}
		ctx.dbs[key] = db
	}
	return ctx, nil
}

func MustNewDbTest(cnf *IniCnf, keys ...string) *Db {
	db, err := NewDbTest(cnf, keys...)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func (this *Db) Exec(name string, sql string, params ...interface{}) *Db {
	if db, ok := this.dbs[name]; ok {
		this.request = NewDbReqest("w", db, sql, params)
	} else {
		panic(fmt.Sprintf("db %s was not found", name))
	}
	return this
}

func (this *Db) Query(name string, sql string, params ...interface{}) *Db {
	if db, ok := this.dbs[name]; ok {
		this.request = NewDbReqest("r", db, sql, params)
	} else {
		panic(fmt.Sprintf("db %s was not found", name))
	}
	return this
}

type DbExpect struct {
	t *testing.T

	rowsAffected int64
	lastInsertId int64

	rows []*DbRow
}

func (this *Db) Expect(t *testing.T) *DbExpect {
	expect := &DbExpect{t: t}
	if this.request.op == "w" {
		if rowsAffected, lastInsertId, err := dbExec(this.request.db, this.request.sql, this.request.params...); err != nil {
			t.Fatalf("exec %s error: %s", this.request.sql, err)
		} else {
			expect.rowsAffected = rowsAffected
			expect.lastInsertId = lastInsertId
		}
	} else {
		if rows, err := dbQuery(this.request.db, this.request.sql, this.request.params...); err != nil {
			t.Fatalf("query %s error: %s", this.request.sql, err)
		} else {
			expect.rows = rows
		}
	}
	return expect
}

func (this *DbExpect) RowNumEq(num int) *DbExpect {
	if num != len(this.rows) {
		this.t.Fatalf("db rows expect %d, got %d", num, len(this.rows))
	}
	return this
}

func (this *DbExpect) Eq(field string, value interface{}) *DbExpect {
	if v, ok := value.(string); ok {
		if v != this.rows[0].MustGetString(field) {
			this.t.Fatalf("field %s expect %s, got %s", field, v, this.rows[0].MustGetString(field))
		}
	} else if v, err := interfaceToInt(value); err == nil {
		if v != this.rows[0].MustGetInt(field) {
			this.t.Fatalf("field %s expect %d, got %d", field, v, this.rows[0].MustGetInt(field))
		}
	} else if v, err := interfaceToFloat(value); err == nil {
		if v != this.rows[0].MustGetFloat(field) {
			this.t.Fatalf("field %s expect %f, got %f", field, v, this.rows[0].MustGetFloat(field))
		}
	} else {
		this.t.Fatalf("unsupport type %#v", value)
	}
	return this
}

func (this *DbExpect) Test(test func(t *testing.T, rows []*DbRow)) *DbExpect {
	test(this.t, this.rows)
	return this
}
