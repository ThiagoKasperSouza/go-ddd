package identity

import (
	"context"
	"errors"
	"os"

	domainIdentity "go-ddd/internal/domain/identity"
	"go-ddd/internal/dto"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type LoginUseCase struct {
	userRepo       domainIdentity.UserRepository
	passwordHasher domainIdentity.PasswordHasher
}

func NewLoginUseCase(
	userRepo domainIdentity.UserRepository,
	passwordHasher domainIdentity.PasswordHasher,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input dto.LoginInputDTO) (*dto.LoginOutputDTO, error) {
	// 1. Busca o usuário no banco de dados
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Compara a senha informada com o hash gravado no banco
	if !uc.passwordHasher.Compare(input.Password, user.PasswordHash()) {
		return nil, ErrInvalidCredentials
	}

	// 3. Pega a secret key do ambiente
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev"
	}

	// 4. ONDE O GenerateJWT É USADO: Gerando o token assinado para o usuário autenticado
	token, err := domainIdentity.GenerateJWT(user.ID(), user.Roles(), jwtSecret)
	if err != nil {
		return nil, err
	}

	return &dto.LoginOutputDTO{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}