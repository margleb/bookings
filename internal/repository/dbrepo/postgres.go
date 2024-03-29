package dbrepo

import (
	"context"
	"errors"
	"github.com/margleb/booking/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func (m *postgresDbRepo) AllUsers() bool {
	return true
}

// InsertReservation добавляет бронирование в базу данных
func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {

	// используем подключение 3 cекунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// затем если ничего не происходит, отключаем
	defer cancel()

	var newId int

	// добавляем в базу данных значение
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction добавляет ограничения для комнаты
func (m *postgresDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	// используем подключение 3 cекунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// затем если ничего не происходит, отключаем
	defer cancel()

	// добавляем в базу данных значение
	stmt := `insert into room_restrictions(start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_Id) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, stmt, r.StartDate, r.EndDate, r.RoomID, r.ReservationID, time.Now(), time.Now(), r.RestrictionID)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID - проверяем, свободна ли комната на данный диапазон дат
func (m *postgresDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {

	// время выполнения запроса не более трех секунд
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// запрос к базе данных
	query := `select count(*) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date`

	// кол-во полученных строк
	var numRows int

	// выполняем запрос к базе
	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	// если комната свободная
	if numRows == 0 {
		return true, nil
	} else {
		return false, nil
	}

}

// SearchAvailabilityForAllRooms - возвращает slice доступных комнат для бронирования
func (m *postgresDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	// время выполнения запроса не более трех секунд
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	// запрос
	query := `select r.id, r.room_name from rooms r where
	r.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`

	// выполняем запрос
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	// проходимся циклом и добавляем комнаты в slice
	for rows.Next() {

		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}

		if err = rows.Err(); err != nil {
			log.Fatal("Error scanning rows", err)
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

// GetRoomByID - gets a room by ID
func (m *postgresDbRepo) GetRoomByID(id int) (models.Room, error) {

	// время выполнения запроса не более трех секунд
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}

// GetUserByID - получает пользователя по ID
func (m *postgresDbRepo) GetUserByID(id int) (models.User, error) {
	// используем подключение 3 cекунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// затем если ничего не происходит, отключаем
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
    from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user in the database
func (m *postgresDbRepo) UpdateUser(u models.User) error {
	// используем подключение 3 cекунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// затем если ничего не происходит, отключаем
	defer cancel()

	query := `
    	update users set first_name = $1,
		last_name = $2,
		email = $3,
		access_level = $4,
		updated_at = $5
    `

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate - authenticate a user
func (m *postgresDbRepo) Authenticate(email, testPassword string) (int, string, error) {
	// используем подключение 3 cекунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// затем если ничего не происходит, отключаем
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
