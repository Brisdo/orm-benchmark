package benchs

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var db gorm.DB

func init() {
	st := NewSuite("gorm")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, GormInsert)
		st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, GormInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, GormUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, GormRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, GormReadSlice)
		db, _ = gorm.Open("postgres", ORM_SOURCE)
		db.DB().SetMaxIdleConns(ORM_MAX_IDLE)
		db.DB().SetMaxOpenConns(ORM_MAX_CONN)

		db.SingularTable(true)
		db.AutoMigrate(&Model{})

	}
}

func GormInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()

		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		db.Create(m)
	}
}

func GormInsertMulti(b *B) {
	fmt.Println("Not support multi insert")
}

func GormUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		db.Create(m)
	})

	for i := 0; i < b.N; i++ {
		db.Save(m)
	}
}

func GormRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		db.Create(m)
	})

	for i := 0; i < b.N; i++ {
		db.First(m, m.Id)
	}
}

func GormReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < 100; i++ {
			m.Id = 0
			db.Create(m)
		}
	})
	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := db.Where("id > ?", 0).Limit(100).Find(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
