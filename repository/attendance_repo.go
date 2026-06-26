package repository

import(
	"context"
	"database/sql"
	"time"
)

type Attendance struct{
	ID 				int64
	UserID 			int64
	CheckInTime 	time.Time
	CheckOutTime 	sql.NullTime
	Status 			string
}

type AttendanceRepository interface {
	RecordCheckIn(ctx context.Context, a *Attendance) error
	CheckAttendanceToday(ctx context.Context, userID int64) (bool, error)
	GetAttendanceToday(ctx context.Context, userID int64) (*Attendance, error)
	RecordCheckOut(ctx context.Context, id int64, checkOutTime time.Time) error
}

type postgreAttendanceRepository struct {
	db *sql.DB
}

func NewPostgresAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &postgreAttendanceRepository{db: db}
}

func (r *postgreAttendanceRepository) RecordCheckIn(ctx context.Context, attendance *Attendance) error {
	query := `INSERT INTO attendances (user_id, check_in_time, status)
			  VALUES ($1, $2, $3)`
	
	_, err := r.db.ExecContext(ctx, query, attendance.UserID, attendance.CheckInTime, attendance.Status)
	
	if err != nil {
		return err
	}
	return nil
}

func (r *postgreAttendanceRepository) CheckAttendanceToday(ctx context.Context, userID int64) (bool, error) {
	var id int
	query := `SELECT id FROM attendances 
			  WHERE user_id = $1 AND DATE(check_in_time) = CURRENT_DATE`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id) 
	
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil 
		}
		return false, err
	}
	return true, nil
}

func (r *postgreAttendanceRepository) GetAttendanceToday(ctx context.Context, userID int64) (*Attendance, error) {
	query := `SELECT id, user_id, check_in_time, check_out_time, status
			  FROM attendances
			  WHERE user_id = $1 AND DATE(check_in_time) = CURRENT_DATE LIMIT 1`

	var a Attendance
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&a.ID, &a.UserID, &a.CheckInTime, &a.CheckOutTime, &a.Status)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *postgreAttendanceRepository) RecordCheckOut(ctx context.Context, id int64, checkOutTime time.Time) error {
	query := "UPDATE attendances SET check_out_time = $1 WHERE id = $2"

	_, err := r.db.ExecContext(ctx, query, checkOutTime, id)
	return err
}

