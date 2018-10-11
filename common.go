package apitest

import (
	"database/sql"
	"fmt"
	mysql "github.com/go-sql-driver/mysql"
	"math"
	"reflect"
	"strconv"
)

func interfaceToString(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	} else if i, ok := v.(int); ok {
		return strconv.Itoa(i), nil
	} else if i, ok := v.(int64); ok {
		return strconv.FormatInt(i, 10), nil
	} else {
		return "", fmt.Errorf("%v is not string", v)
	}
}

func MustInterfaceToString(v interface{}) string {
	if s, err := interfaceToString(v); err != nil {
		return fmt.Sprintf("%v", v)
	} else {
		return s
	}
}

func interfaceToStrings(v interface{}) ([]string, error) {
	if vs, ok := v.([]interface{}); ok {
		ss := make([]string, 0)
		for _, x := range vs {
			if s, err := interfaceToString(x); err != nil {
				return nil, fmt.Errorf("%v is not a string array")
			} else {
				ss = append(ss, s)
			}
		}
		return ss, nil
	} else {
		return nil, fmt.Errorf("%v is not array", v)
	}
}

func interfaceToInt(v interface{}) (int64, error) {
	if i32, ok := v.(int); ok {
		return int64(i32), nil
	} else if i64, ok := v.(int64); ok {
		return i64, nil
	} else if f32, ok := v.(float32); ok {
		return int64(f32), nil
	} else if f64, ok := v.(float64); ok {
		return int64(f64), nil
	} else {
		return 0, fmt.Errorf("%v is not int64", v)
	}
}

func interfaceToInts(v interface{}) ([]int64, error) {
	if vs, ok := v.([]interface{}); ok {
		ii := make([]int64, 0)
		for _, x := range vs {
			if i, err := interfaceToInt(x); err != nil {
				return nil, fmt.Errorf("%v is not a int array")
			} else {
				ii = append(ii, i)
			}
		}
		return ii, nil
	} else {
		return nil, fmt.Errorf("%v is not array", v)
	}
}

func interfaceToFloat(v interface{}) (float64, error) {
	if i32, ok := v.(int); ok {
		return float64(i32), nil
	} else if i64, ok := v.(int64); ok {
		return float64(i64), nil
	} else if f32, ok := v.(float32); ok {
		return float64(f32), nil
	} else if f64, ok := v.(float64); ok {
		return f64, nil
	} else {
		return 0, fmt.Errorf("%v is not float64", v)
	}
}

func interfaceToFloats(v interface{}) ([]float64, error) {
	if vs, ok := v.([]interface{}); ok {
		ff := make([]float64, 0)
		for _, x := range vs {
			if f, err := interfaceToFloat(x); err != nil {
				return nil, fmt.Errorf("%v is not a float array")
			} else {
				ff = append(ff, f)
			}
		}
		return ff, nil
	} else {
		return nil, fmt.Errorf("%v is not array", v)
	}
}

func stringContains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func stringOnlyContains(values []string, want string) bool {
	for _, value := range values {
		if value != want {
			return false
		}
	}
	return len(values) > 0
}

func intContains(values []int64, want int64) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func intOnlyContains(values []int64, want int64) bool {
	for _, value := range values {
		if value != want {
			return false
		}
	}
	return len(values) > 0
}

func floatContains(values []float64, want float64) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func floatOnlyContains(values []float64, want float64) bool {
	for _, value := range values {
		if value != want {
			return false
		}
	}
	return len(values) > 0
}

func dbExec(c *sql.DB, sql string, values ...interface{}) (int64, int64, error) {
	if result, err := c.Exec(sql, values...); err != nil {
		return 0, 0, err
	} else {
		if rowsAffected, err := result.RowsAffected(); err != nil {
			return 0, 0, err
		} else {
			if lastInsertId, err := result.LastInsertId(); err != nil {
				return 0, 0, err
			} else {
				return rowsAffected, lastInsertId, err
			}
		}
	}
}

type DbRow struct {
	row map[string]interface{}
}

func NewDbRow(fields []string, values []interface{}) *DbRow {
	row := make(map[string]interface{}, 0)
	for i, f := range fields {
		row[f] = values[i]
	}
	return &DbRow{row: row}
}

func (this *DbRow) GetString(f string) (string, error) {
	if len(this.row) == 0 {
		return "", fmt.Errorf("unknow field %s in empty row", f)
	}

	if v, ok := this.row[f]; ok {
		if ns, ok := v.(*sql.NullString); ok {
			if ns.Valid {
				return ns.String, nil
			} else {
				return "", nil
			}
		} else {
			return "", fmt.Errorf("field %s type error: %s not string", f, reflect.TypeOf(v))
		}
	} else {
		return "", fmt.Errorf("unknow field %s", f)
	}
}

func (this *DbRow) MustGetString(f string) string {
	if s, err := this.GetString(f); err == nil {
		return s
	} else {
		panic(err.Error())
	}
}

func (this *DbRow) GetInt(f string) (int64, error) {
	if len(this.row) == 0 {
		return math.MinInt64, fmt.Errorf("unknow field %s in empty row", f)
	}

	if v, ok := this.row[f]; ok {
		if ni, ok := v.(*sql.NullInt64); ok {
			if ni.Valid {
				return ni.Int64, nil
			} else {
				return 0, nil
			}
		} else {
			return math.MinInt64, fmt.Errorf("type error: not int")
		}
	} else {
		return math.MinInt64, fmt.Errorf("unknow field %s", f)
	}
}

func (this *DbRow) MustGetInt(f string) int64 {
	if i, err := this.GetInt(f); err == nil {
		return i
	} else {
		panic(err.Error)
	}
}

func (this *DbRow) GetFloat(f string) (float64, error) {
	if len(this.row) == 0 {
		return math.MaxFloat64, fmt.Errorf("unknow field %s in empty row", f)
	}

	if v, ok := this.row[f]; ok {
		if nf, ok := v.(*sql.NullFloat64); ok {
			if nf.Valid {
				return nf.Float64, nil
			} else {
				return 0.0, nil
			}
		} else {
			return math.MaxFloat64, fmt.Errorf("field %s type error: %s not float", f, reflect.TypeOf(v))
		}
	} else {
		return math.MaxFloat64, fmt.Errorf("unknow field %s", f)
	}
}

func (this *DbRow) MustGetFloat(f string) float64 {
	if nf, err := this.GetFloat(f); err == nil {
		return nf
	} else {
		panic(err.Error())
	}
}

func buildRowValue(colTypes []*sql.ColumnType) ([]interface{}, error) {
	row := make([]interface{}, len(colTypes))
	for i, typ := range colTypes {
		dbType := typ.ScanType()
		if dbType == reflect.TypeOf(sql.NullString{}) || dbType == reflect.TypeOf(sql.RawBytes{}) {
			row[i] = &sql.NullString{}
		} else if dbType == reflect.TypeOf(sql.NullInt64{}) ||
			dbType.Kind() == reflect.Int32 || dbType.Kind() == reflect.Int64 ||
			dbType.Kind() == reflect.Uint32 || dbType.Kind() == reflect.Uint64 ||
			dbType.Kind() == reflect.Int16 || dbType.Kind() == reflect.Int8 ||
			dbType.Kind() == reflect.Uint16 || dbType.Kind() == reflect.Uint8 {

			row[i] = &sql.NullInt64{}
		} else if dbType == reflect.TypeOf(sql.NullFloat64{}) {
			row[i] = &sql.NullFloat64{}
		} else if dbType == reflect.TypeOf(mysql.NullTime{}) {
			row[i] = &mysql.NullTime{}
		} else {
			return nil, fmt.Errorf("unknow type %v", dbType)
		}
	}
	return row, nil
}

func dbQuery(db *sql.DB, sql string, params ...interface{}) ([]*DbRow, error) {
	rows, err := db.Query(sql, params...)
	if err != nil {
		return nil, fmt.Errorf("query %s error: %s", sql, err)
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("query %s error: %s", sql, err)
	}

	tbl := make([]*DbRow, 0)

	fields := make([]string, len(colTypes))
	for i, typ := range colTypes {
		fields[i] = typ.Name()
	}

	if _, err := buildRowValue(colTypes); err != nil {
		return nil, fmt.Errorf("query %s error: %s", sql, err)
	}

	for rows.Next() {
		row, _ := buildRowValue(colTypes)

		if err := rows.Scan(row...); err != nil {
			return nil, fmt.Errorf("query %s error: %s", sql, err)
		}

		m := NewDbRow(fields, row)
		tbl = append(tbl, m)
	}
	return tbl, nil
}
