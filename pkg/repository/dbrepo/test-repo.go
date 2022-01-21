package dbrepo

import (
	"errors"
	"time"

	"github.com/jinyanomura/bookings/pkg/models"
)

// InsertReservation inserts a reservation into the database.
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("invalid room_id")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database.
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.StartDate.Format("2006-01-02") == "2222-02-02" {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if the room is available on given dates.
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns slice of available rooms if any on given date range. 
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets a room by ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

// GetUserByID returns user information specified by given id.
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	return models.User{}, nil
}

// UpdateUser updates user in the database.
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// UpdateReservation updates reseervation details in the database.
func (m *testDBRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

// Authenticate authenticates users
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0, "", nil
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

// AllNewReservations returns a slice of all reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

// GetReservationByID returns one reservation specified by ID.
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	return models.Reservation{}, nil
}

// DeleteReservation deletes reservations specified by id.
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation updates processed property for given reservation.
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}