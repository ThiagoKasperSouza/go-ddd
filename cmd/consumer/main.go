// Consumer: lê eventos do tópico Kafka e os expõe como métricas Prometheus em /metrics.
// Fica em cmd/consumer/main.go, no mesmo módulo da API (reaproveita monitor.Event).
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"

	"go-ddd/internal/infra/http/middleware" // ajuste para o module path do seu go.mod
)

var (
	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_requests_total",
		Help: "Requisições autenticadas por método, rota e status.",
	}, []string{"method", "route", "status"})

	duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "api_request_duration_seconds",
		Help:    "Latência das requisições por método e rota.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
	}, []string{"method", "route"})

	respBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_response_bytes_total",
		Help: "Bytes enviados nas respostas por método e rota.",
	}, []string{"method", "route"})

	lag = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "api_events_consumer_lag",
		Help: "Mensagens pendentes de consumo no tópico (saúde do pipeline).",
	})

	invalid = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "api_events_invalid_total",
		Help: "Mensagens do Kafka que não puderam ser decodificadas.",
	})
)

// routeLimiter protege o Prometheus de cardinalidade infinita: se por algum motivo
// o path vier cru (ex.: /users/1, /users/2...), rotas além do limite viram "other".
type routeLimiter struct {
	seen map[string]struct{}
	max  int
}

func (l *routeLimiter) normalize(route string) string {
	if _, ok := l.seen[route]; ok {
		return route
	}
	if len(l.seen) >= l.max {
		return "other"
	}
	l.seen[route] = struct{}{}
	return route
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	brokers := strings.Split(env("KAFKA_BROKERS", "localhost:9094"), ",")
	topic := env("KAFKA_TOPIC", "api.events")
	group := env("KAFKA_GROUP", "metrics-consumer")
	addr := env("METRICS_ADDR", ":2112")

	prometheus.MustRegister(requests, duration, respBytes, lag, invalid)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        group,
		MinBytes:       1,
		MaxBytes:       1 << 20,
		MaxWait:        500 * time.Millisecond,
		CommitInterval: time.Second,
		// Só vale na primeira execução do grupo (sem offset salvo): não reprocessa
		// histórico, que inflaria os contadores com tráfego antigo.
		StartOffset: kafka.LastOffset,
	})
	defer r.Close()

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				lag.Set(float64(r.Stats().Lag))
			}
		}
	}()

	lim := &routeLimiter{seen: make(map[string]struct{}), max: 500}
	log.Printf("consumer: lendo %q (grupo %q), métricas em %s/metrics", topic, group, addr)

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Printf("consumer: %v", err)
			time.Sleep(time.Second)
			continue
		}

		var e middleware.Event
		if err := json.Unmarshal(msg.Value, &e); err != nil {
			invalid.Inc()
			continue
		}

		route := lim.normalize(e.Path)
		requests.WithLabelValues(e.Method, route, strconv.Itoa(e.Status)).Inc()
		duration.WithLabelValues(e.Method, route).Observe(e.DurationMs / 1000)
		respBytes.WithLabelValues(e.Method, route).Add(float64(e.BytesOut))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}