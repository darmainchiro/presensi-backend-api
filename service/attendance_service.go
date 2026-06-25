package service

import(
	"context"
	"errors"
	"time"
	"backend-api/repository"
	"backend-api/utils"
)

type AttendanceService interface {
	ProcessCheckIn(ctx context.Context, userID int64) error
	ProcessCheckOut(ctx context.Context, userID int64) error

}

type attendanceService struct {
	repo repository.AttendanceRepository
}

func NewAttendanceService(repo repository.AttendanceRepository) AttendanceService {
	return&attendanceService{repo: repo}
}

func (s *attendanceService) ProcessCheckIn(ctx context.Context, userID int64) error{
	if userID <= 0 {
		return errors.New("user id tidak valid")
	}

	hasCheckIn, err := s.repo.CheckAttendanceToday(ctx, userID)
	if err != nil {
		return errors.New("terjadi kesalahan sistem saat memvalidasi data absensi")
	}

	if hasCheckIn {
		return errors.New("anda sudah melakukan check in hari ini")
	}

	currentTime := time.Now()
	status := "Hadir"

	deadline := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 9, 0, 0, 0, currentTime.Location())
	if currentTime.After(deadline){
		status = "Terlambat"
	}

	attendanceRecord := &repository.Attendance{
		UserID: userID,
		CheckInTime: currentTime,
		Status: status,	
	}

	err = s.repo.RecordCheckIn(ctx, attendanceRecord)
	if err != nil {
		return err
	}

	go func() {
		event := utils.AttendanceEvent{
			UserID: 	userID,
			Timestamp:	currentTime.Format("2006-01-02 15:04:05"),
			Type:		"CHECK_IN",
		}
		utils.PublishAttendanceEvent(event)
	}()

	return nil
}

func (s *attendanceService) ProcessCheckOut(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errors.New("user id tidak valid")
	}

	attendance, err := s.repo.GetAttendanceToday(ctx, userID)
	if attendance == nil {
		return errors.New("anda belum check in")
	}

	if attendance.CheckOutTime.Valid {
		return errors.New("anda sudah check out")
	}

	currentTime := time.Now()
	err = s.repo.RecordCheckOut(ctx, attendance.ID, currentTime)
	if err != nil{
		return errors.New("terjadi kesalahan sistem saat menyimpan data check out")
	}
	return nil

}