package repository

import "github.com/jinyanomura/bookings/pkg/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}