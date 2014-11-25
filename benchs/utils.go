package benchs

import (
	"database/sql"
	"fmt"
	"os"
)

type Model struct {
	Id      int `qbs:"pk" sql:"pk"`
	Name    string
	Title   string
	Fax     string
	Web     string
	Age     int
	Aight   bool
	Counter int64
}

func NewModel() *Model {
	m := new(Model)
	m.Name = "Orm Benchmark"
	m.Title = "Just a Benchmark for fun"
	m.Fax = "99909990"
	m.Web = "http://beego.me"
	m.Age = 100
	m.Aight = true
	m.Counter = 1000

	return m
}

var (
	ORM_MULTI    int
	ORM_MAX_IDLE int
	ORM_MAX_CONN int
	ORM_SOURCE   string
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func wrapExecute(b *B, cbk func()) {
	b.StopTimer()
	defer b.StartTimer()
	cbk()
}

func initDB() {
	sqls := []string{
		"DROP TABLE IF EXISTS model",
		"CREATE TABLE model (" +
			"id SERIAL PRIMARY KEY," +
			"name text," +
			"title text," +
			"fax text," +
			"web text," +
			"age int," +
			"aight  boolean," +
			"counter bigint" +
			") ",
	}

	DB, err := sql.Open("postgres", ORM_SOURCE)
	checkErr(err)
	defer DB.Close()

	err = DB.Ping()
	checkErr(err)

	for _, sql := range sqls {
		//fmt.Printf(sql)
		_, err = DB.Exec(sql)
		checkErr(err)
	}
}
