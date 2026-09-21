package shape

import (
	"encoding/json"
	"errors"
	"net/http"
	"log"

	domainShape "go-ddd/internal/domain/shapes"
	usecaseShape "go-ddd/internal/usecase/shapes"
	
)

type ShapeHandler struct {
	getShapeUseCase    *usecaseShape.GetShapeUseCase
	listShapesUseCase *usecaseShape.ListShapesUseCase
}

func NewShapeHandler(
	getShapeUseCase *usecaseShape.GetShapeUseCase, 
	listShapesUseCase *usecaseShape.ListShapesUseCase) *ShapeHandler {
	return &ShapeHandler{
		getShapeUseCase: 		getShapeUseCase,
		listShapesUseCase:	listShapesUseCase,
	}
}

// GetByID lida com GET /shape/{id}
func (h *ShapeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Captura o parâmetro {id} diretamente pelo net/http nativo do Go (1.22+)
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id parameter is required"})
		return
	}

	output, err := h.getShapeUseCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainShape.ErrShapeNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Shape not found"})
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

// ListAll lida com GET /shape
func (h *ShapeHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	output, err := h.listShapesUseCase.Execute(r.Context())
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
