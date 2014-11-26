package benchs

import (
	"database/sql"
	"fmt"
	"strings"
)

var raw *sql.DB

const (
	rawInsertBaseSQL   = "INSERT INTO model (name, title, fax, web, age, aight, counter) VALUES "
	rawInsertValuesSQL = "($1, $2, $3, $4, $5, $6, $7)"
	rawInsertReturnId  = " RETURNING id"
	rawInsertSQL       = rawInsertBaseSQL + rawInsertValuesSQL
	rawUpdateSQL       = "UPDATE model SET name=$1, title=$2, fax=$3, web=$4, age=$5, aight=$6, counter=$7 WHERE id=$8"
	rawSelectSQL       = "SELECT id, name, title, fax, web, age, aight, counter FROM model WHERE id=$1"
	rawSelectMultiSQL  = "SELECT id, name, title, fax, web, age, aight, counter FROM model WHERE id>0 LIMIT 100"
)

func init() {
	st := NewSuite("raw")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, RawInsert)
		st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, RawInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, RawUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, RawRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, RawReadSlice)
		raw, _ = sql.Open("postgres", ORM_SOURCE)
	}
}

func RawInsert(b *B) {
	var m *Model
	var stmt *sql.Stmt
	sqlInsertNew := rawInsertSQL + rawInsertReturnId
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		stmt, err = raw.Prepare(sqlInsertNew)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	defer stmt.Close()

	for i := 0; i < b.N; i++ {
		err := stmt.QueryRow(m.Name, m.Title, m.Fax, m.Web, m.Age, m.Aight, m.Counter).Scan(&m.Id)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func rawInsert(m *Model) error {
	sqlInsertNew := rawInsertSQL + rawInsertReturnId
	err := raw.QueryRow(sqlInsertNew, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Aight, m.Counter).Scan(&m.Id)
	if err != nil {
		return err
	}
	return err
}

func RawInsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()

		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
	})

	for i := 0; i < b.N; i++ {
		nFields := 7
		query := rawInsertBaseSQL + strings.Repeat(rawInsertValuesSQL+",", len(ms)-1) + rawInsertValuesSQL
		args := make([]interface{}, len(ms)*nFields)
		for j := range ms {
			offset := j * nFields
			args[offset+0] = ms[j].Name
			args[offset+1] = ms[j].Title
			args[offset+2] = ms[j].Fax
			args[offset+3] = ms[j].Web
			args[offset+4] = ms[j].Age
			args[offset+5] = ms[j].Aight
			args[offset+6] = ms[j].Counter
		}
		m := ms[0]
		_, err := raw.Exec(query, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Aight, m.Counter)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}

	}
}

func RawUpdate(b *B) {
	var m *Model
	var stmt *sql.Stmt
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		rawInsert(m)
		fmt.Println("id is ", m.Id)
		stmt, err = raw.Prepare(rawUpdateSQL)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	defer stmt.Close()

	for i := 0; i < b.N; i++ {
		_, err := stmt.Exec(m.Name, m.Title, m.Fax, m.Web, m.Age, m.Aight, m.Counter, m.Id)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func RawRead(b *B) {
	var m *Model
	var stmt *sql.Stmt
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		rawInsert(m)
		stmt, err = raw.Prepare(rawSelectSQL)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	defer stmt.Close()

	for i := 0; i < b.N; i++ {
		var mout Model
		err := stmt.QueryRow(m.Id).Scan(
			&mout.Id,
			&mout.Name,
			&mout.Title,
			&mout.Fax,
			&mout.Web,
			&mout.Age,
			&mout.Aight,
			&mout.Counter,
		)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func RawReadSlice(b *B) {
	var m *Model
	var stmt *sql.Stmt
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		for i := 0; i < 100; i++ {
			err = rawInsert(m)
			if err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
		stmt, err = raw.Prepare(rawSelectMultiSQL)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	defer stmt.Close()

	for i := 0; i < b.N; i++ {
		var j int
		models := make([]Model, 100)
		rows, err := stmt.Query()
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
		for j = 0; rows.Next() && j < len(models); j++ {
			err = rows.Scan(
				&models[j].Id,
				&models[j].Name,
				&models[j].Title,
				&models[j].Fax,
				&models[j].Web,
				&models[j].Age,
				&models[j].Aight,
				&models[j].Counter,
			)
			if err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
		models = models[:j]
		if err = rows.Err(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
		if err = rows.Close(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
