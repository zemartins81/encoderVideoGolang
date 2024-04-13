package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
	"github.com/zemartins81/encoderVideoGolang/domain"
	"log"
)

type Database struct {
	DB            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "Test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debug = true

	connection, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("Test db errir: %v", err)
	}
	return connection
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error

	if d.Env != "Test" {
		d.DB, err = gorm.Open(d.DbType, d.Dsn)
	} else {
		d.DB, err = gorm.Open(d.DbTypeTest, d.DsnTest)
	}
	if err != nil {
		return nil, err
	}
	if d.Debug {
		d.DB.LogMode(true)
	}

	if d.AutoMigrateDb {
		d.DB.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.DB.Model(domain.Job{}).AddForeignKey("viceo_id", "videos (id)", "CASCADE", "CASCADE")
	}

	return d.DB, nil
}
