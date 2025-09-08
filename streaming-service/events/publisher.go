package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"streaming-service/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Estructura del evento
type SongPlayedEvent struct {
	UserID uint `json:"user_id"`
	SongID uint `json:"song_id"`
}

// PublishSongPlayed envía el evento a RabbitMQ
func PublishSongPlayed(userID, songID uint) error {
	cfg := config.GetConfig()
	fmt.Printf("🚀 Intentando enviar evento a RabbitMQ: userID=%d, songID=%d\n", userID, songID)

	// Conexión
	conn, err := amqp.Dial(cfg.RabbitURL)
	if err != nil {
		fmt.Printf("❌ Error al conectar con RabbitMQ: %v\n", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("❌ Error al abrir canal: %v\n", err)
		return err
	}
	defer ch.Close()

	// Declarar exchange (fanout)
	err = ch.ExchangeDeclare(
		"song_events",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("❌ Error al declarar exchange: %v\n", err)
		return err
	}

	// Declarar cola
	queue, err := ch.QueueDeclare(
		"song_events_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("❌ Error al declarar cola: %v\n", err)
		return err
	}

	// Enlazar cola al exchange
	err = ch.QueueBind(
		queue.Name,
		"",
		"song_events",
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("❌ Error al enlazar cola: %v\n", err)
		return err
	}

	// Crear evento
	event := SongPlayedEvent{
		UserID: userID,
		SongID: songID,
	}
	body, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("❌ Error al serializar evento: %v\n", err)
		return err
	}

	// Contexto con timeout para publicación
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publicar con contexto
	err = ch.PublishWithContext(
		ctx,
		"song_events", // exchange
		"",            // routing key (fanout -> ignora)
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		fmt.Printf("❌ Error al publicar mensaje: %v\n", err)
		return err
	}

	log.Printf("📨 Evento enviado a RabbitMQ: %+v\n", event)
	return nil
}
