package repository

import "github.com/margleb/booking/internal/models"

// DatabaseRepo интерфейс, содержащий функции для работы с базой данных
type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
}
