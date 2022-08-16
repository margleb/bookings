package dbrepo

import (
	"context"
	"database/sql"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/models"
	"github.com/margleb/booking/internal/repository"
	"time"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRep создает новое postgres подключение в репозитории
func NewPostgresRep(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		App: a,
		DB:  conn,
	}
}

// InsertReservation добавляет бронирование в базу данных
func (p *postgresDbRepo) InsertReservation(res models.Reservation) error {

	// используем подключение 3 cекунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// затем если ничего не происходит, отключаем
	defer cancel()

	// добавляем в базу данных значение
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8. $9)`
	_, err := p.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now())

	if err != nil {
		return err
	}

	return nil
}
