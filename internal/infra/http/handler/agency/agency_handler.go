package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"log"

	domainAgency "go-ddd/internal/domain/agency"
	"go-ddd/internal/dto"
	usecaseAgency "go-ddd/internal/usecase/agency"
)

type AgencyHandler struct {
	createAgencyUseCase *usecaseAgency.CreateAgencyUseCase
	getAgencyUseCase    *usecaseAgency.GetAgencyUseCase
	listAgenciesUseCase *usecaseAgency.ListAgenciesUseCase
	updatedAgencyUseCase *usecaseAgency.UpdateAgencyUseCase
	deleteAgencyUseCase *usecaseAgency.deleteAgencyUseCase
}

func NewAgencyHandler(
	createAgencyUseCase *usecaseAgency.CreateAgencyUseCase, 
	getAgencyUseCase *usecaseAgency.GetAgencyUseCase, 
	listAgenciesUseCase *usecaseAgency.ListAgenciesUseCase,
	updatedAgencyUseCase *usecaseAgency.UpdateAgencyUseCase,
	deleteAgencyUseCase *usecaseAgency.DeleteAgencyUseCase) *AgencyHandler {
	return &AgencyHandler{
		createAgencyUseCase: 	createAgencyUseCase,
		getAgencyUseCase: 		getAgencyUseCase,
		listAgenciesUseCase:	listAgenciesUseCase,
		updatedAgencyUseCase: 	updatedAgencyUseCase,
		deleteAgencyUseCase:	deleteAgencyUseCase
	}
}

// GetByID lida com GET /agencies/{id}
func (h *AgencyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Captura o parâmetro {id} diretamente pelo net/http nativo do Go (1.22+)
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id parameter is required"})
		return
	}

	output, err := h.getAgencyUseCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainAgency.ErrAgencyNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "agency not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

// ListAll lida com GET /agencies
func (h *AgencyHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	output, err := h.listAgenciesUseCase.Execute(r.Context())
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

// Create lida com a requisição HTTP POST /agencies
func (h *AgencyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAgencyDTO

	// 1. Decodifica o JSON do body da requisição no DTO de entrada
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// 2. Executa o Use Case passando o DTO
	output, err := h.createAgencyUseCase.Execute(r.Context(), input)
	if err != nil {
		// Trata erros conhecidos do domínio com códigos de status apropriados
		if errors.Is(err, domainAgency.ErrInvalidAgency) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Erro interno não esperado (ex: falha no banco de dados)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// 3. Retorna a resposta HTTP 201 Created com o DTO de saída em formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

// Update lida com PUT /agencies/{id}
func (h *AgencyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id parameter is required"})
		return
	}

	var input dto.UpdateAgencyDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	output, err := h.updateAgencyUseCase.Execute(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domainAgency.ErrAgencyNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "agency not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

// Delete lida com DELETE /agencies/{id}
func (h *AgencyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id parameter is required"})
		return
	}

	err := h.deleteAgencyUseCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainAgency.ErrAgencyNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "agency not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// 204 No Content para exclusões bem-sucedidas
	w.WriteHeader(http.StatusNoContent)
}