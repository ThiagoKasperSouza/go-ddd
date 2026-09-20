package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	usecaseIdentity "go-ddd/internal/usecase/identity"
	domainIdentity "go-ddd/internal/domain/identity"
	"go-ddd/internal/dto"
)

type AuthHandler struct {
	loginUseCase      *usecaseIdentity.LoginUseCase
	createUserUseCase *usecaseIdentity.CreateUserUseCase
}

func NewAuthHandler(
	loginUseCase *usecaseIdentity.LoginUseCase,
	createUserUseCase *usecaseIdentity.CreateUserUseCase,
) *AuthHandler {
	return &AuthHandler{
		loginUseCase:      loginUseCase,
		createUserUseCase: createUserUseCase,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	output, err := h.loginUseCase.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, usecaseIdentity.ErrInvalidCredentials) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
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

// Register lida com POST /register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := h.createUserUseCase.Execute(r.Context(), input)
	if err != nil {
		// Trata usuário já existente (Conflito - 409)
		if errors.Is(err, usecaseIdentity.ErrUserAlreadyExists) {
			respondWithError(w, http.StatusConflict, "email already registered")
			return
		}

		// Trata erro de validação de senha fraca/curta (Bad Request - 400)
		if errors.Is(err, domainIdentity.ErrPasswordTooShort) || errors.Is(err, domainIdentity.ErrInvalidEmail) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Retorna 201 Created com os dados do usuário criado (sem a senha/hash)
	respondWithJSON(w, http.StatusCreated, output)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}