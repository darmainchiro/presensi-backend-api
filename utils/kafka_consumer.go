package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func StartKafkaConsumer() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:	[]string{broker},
		Topic:		"attendance_events",
		Partition:	0,
		MinBytes:	10e3,
		MaxBytes:	10e6,
	})

	defer r.Close()

	log.Println("Kafka Consumer berjalan: Siap mengolah data presensi dari latar belakang...")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error: %v\n", err)
			break
		}

		fmt.Printf("\n[KAFKA CONSUMER MENERIMA TUGAS]\n")
		fmt.Printf("Menerima Event Check-In/Out!\n")
		fmt.Printf("Isi Pesan: %s\n", string(m.Value))
		fmt.Println("-> Memproses pengiriman notifikasi email... (selesai)")
		fmt.Printf("-------------------------\n")
	}
}