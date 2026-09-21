package middleware

import (
	"context"
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	Timestamp  time.Time `json:"ts"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Status     int       `json:"status"`
	DurationMs float64   `json:"duration_ms"`
	UserID     string    `json:"user_id"`
	IP         string    `json:"ip"`
	BytesOut   int       `json:"bytes_out"`
}

// Publisher desacopla o request do Kafka: Publish nunca bloqueia.
// Se o buffer encher (Kafka lento/fora), o evento é descartado e contado.
type Publisher struct {
	ch      chan Event
	w       *kafka.Writer
	done    chan struct{}
	dropped atomic.Uint64
}

func NewPublisher(brokers []string, topic string) *Publisher {
	p := &Publisher{
		ch: make(chan Event, 10_000),
		w: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireOne,
		},
		done: make(chan struct{}),
	}
	go p.loop()
	return p
}

func (p *Publisher) Publish(e Event) {
	select {
	case p.ch <- e:
	default:
		if n := p.dropped.Add(1); n%1000 == 1 {
			log.Printf("monitor: buffer cheio, %d eventos descartados até agora", n)
		}
	}
}

func (p *Publisher) loop() {
	defer close(p.done)

	const maxBatch = 200
	batch := make([]kafka.Message, 0, maxBatch)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := p.w.WriteMessages(ctx, batch...); err != nil {
			log.Printf("monitor: erro ao enviar %d eventos: %v", len(batch), err)
		}
		batch = batch[:0]
	}

	for {
		select {
		case e, ok := <-p.ch:
			if !ok {
				flush()
				return
			}
			b, err := json.Marshal(e)
			if err != nil {
				continue
			}
			batch = append(batch, kafka.Message{Key: []byte(e.UserID), Value: b})
			if len(batch) >= maxBatch {
				flush()
			}
		case <-tick.C:
			flush()
		}
	}
}

// Close deve ser chamado DEPOIS do shutdown do servidor HTTP
// (Publish após Close causa panic). Faz flush do que sobrou.
func (p *Publisher) Close() error {
	close(p.ch)
	<-p.done
	return p.w.Close()
}