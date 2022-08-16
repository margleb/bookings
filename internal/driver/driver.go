package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB соединение с базой данных
type DB struct {
	SQL *sql.DB
}

// Переменная соединения
var dbConn = &DB{}

// Максимальное кол-во открытых соединений
const maxOpenDBConn = 10

// Максимальное кол-во простаивающих соединений
const maxIdleDBConn = 5

// Максимальное время существования соединения
const maxConnLifetime = 5 * time.Minute

// ConnectSQL - устанавливает соединение сс базой данных
func ConnectSQL(dns string) (*DB, error) {

	// создаем соединение
	d, err := NewDatabase(dns)
	if err != nil {
		// panic(err)
	}

	// устанавливаем настройки соединения
	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxConnLifetime)

	// присваиваем в объект соединение
	dbConn.SQL = d

	// тестируем соединение
	err = TestConn(d)
	if err != nil {
		return nil, err
	}

	return dbConn, err
}

// NewDatabase - создает соединение
func NewDatabase(dns string) (*sql.DB, error) {

	// создаем соединение
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	// проверяем существование соединения
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// TestConn - тестируем соединение
func TestConn(d *sql.DB) error {

	// проверяем существование соединения
	err := d.Ping()

	if err != nil {
		return err
	}

	return nil
}
