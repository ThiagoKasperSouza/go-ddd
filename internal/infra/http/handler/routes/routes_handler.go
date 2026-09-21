package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"log"

	domainRoute "go-ddd/internal/domain/routes"
	usecaseRoutes "go-ddd/internal/usecase/routes"
	
)

type RouteHandler struct {
	getRouteUseCase    *usecaseRoutes.GetRouteUseCase
	listRoutesUseCase *usecaseRoutes.ListRoutesUseCase
}

func NewRouteHandler(
	getRouteUseCase *usecaseRoutes.GetRouteUseCase, 
	listRoutesUseCase *usecaseRoutes.ListRoutesUseCase) *RouteHandler {
	return &RouteHandler{
		getRouteUseCase: 		getRouteUseCase,
		listRoutesUseCase:	listRoutesUseCase,
	}
}

// GetByID lida com GET /routes/{id}
func (h *RouteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Captura o parâmetro {id} diretamente pelo net/http nativo do Go (1.22+)
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id parameter is required"})
		return
	}

	output, err := h.getRouteUseCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainRoute.ErrRouteNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Route not found"})
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

// ListAll lida com GET /routes
func (h *RouteHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	output, err := h.listRoutesUseCase.Execute(r.Context())
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
