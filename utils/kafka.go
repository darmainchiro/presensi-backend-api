package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type AttendanceEvent struct {
	UserID		int64  `json:"user_id"`
	Timestamp	string `json:"timestamp"`
	Type		string `json:"type"`
}

func PublishAttendanceEvent(event AttendanceEvent) error {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	w := &kafka.Writer{
		Addr: 		kafka.TCP(broker),
		Topic:		"attendance_events",
		Balancer:	&kafka.LeastBytes{},
	}
	defer w.Close()

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key: []byte("presensi"),
			Value: eventBytes,
		},
	)

	if err != nil {
		log.Printf("Gagal mengirim pesan ke Kafka: %v", err)
		return err
	}

	log.Println("Pesan berhasil diterbitkan ke Kafka!")
	return nil
}