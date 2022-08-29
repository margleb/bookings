package dbrepo

import (
	"errors"
	"github.com/margleb/booking/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation добавляет бронирование в базу данных
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

// InsertRoomRestriction добавляет ограничения для комнаты
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

// SearchAvailabilityByDatesByRoomID - проверяем, свободна ли комната на данный диапазон дат
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms - возвращает slice доступных комнат для бронирования
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID - gets a room by ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("Some error")
	}
	return room, nil
}
