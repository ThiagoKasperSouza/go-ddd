package trips

import (
	"encoding/json"
	"errors"
	"net/http"
	"log"

	domainTrip "go-ddd/internal/domain/trips"
	usecaseTrip "go-ddd/internal/usecase/trips"
	
)

type TripHandler struct {
	getTripUseCase    *usecaseTrip.GetTripUseCase
	listTripsUseCase *usecaseTrip.ListTripsUseCase
}

func NewTripHandler(
	getTripUseCase *usecaseTrip.GetTripUseCase, 
	listTripsUseCase *usecaseTrip.ListTripsUseCase) *TripHandler {
	return &TripHandler{
		getTripUseCase: 		getTripUseCase,
		listTripsUseCase:	listTripsUseCase,
	}
}

// GetByID lida com GET /Trip/{id}
func (h *TripHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Captura o parâmetro {id} diretamente pelo net/http nativo do Go (1.22+)
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id parameter is required"})
		return
	}

	output, err := h.getTripUseCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainTrip.ErrTripNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Trip not found"})
			return
		}
		log.Printf("Erro no GetById UseCase: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

// ListAll lida com GET /Trip
func (h *TripHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	output, err := h.listTripsUseCase.Execute(r.Context())
	if err != nil {
		log.Printf("Erro no ListAll UseCase: %v\n", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
