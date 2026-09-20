package identity

import (
	"context"
	"errors"

	domainIdentity "go-ddd/internal/domain/identity"
	"go-ddd/internal/dto"
)

var ErrUserAlreadyExists = errors.New("user with this email already exists")

type CreateUserUseCase struct {
	userRepo       domainIdentity.UserRepository
	passwordHasher domainIdentity.PasswordHasher
}

func NewCreateUserUseCase(
	userRepo domainIdentity.UserRepository,
	passwordHasher domainIdentity.PasswordHasher,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.CreateUserDTO) (*dto.UserOutputDTO, error) {
	// 1. Validação básica de tamanho de senha em texto puro
	if len(input.Password) < 6 {
		return nil, domainIdentity.ErrPasswordTooShort
	}

	// 2. Verifica se o e-mail já está cadastrado
	existingUser, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 3. Gera o HASH BCRYPT da senha
	hashedPassword, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	// 4. Converte as roles DTO para o Enum de Domínio
	var domainRoles []domainIdentity.Role
	for _, r := range input.Roles {
		role := domainIdentity.Role(r)
		if role.IsValid() {
			domainRoles = append(domainRoles, role)
		}
	}

	// 5. Cria a Entidade de Usuário contendo o HASH
	user, err := domainIdentity.NewUser(
		input.ID,
		input.Name,
		input.Email,
		hashedPassword,
		domainRoles,
	)
	if err != nil {
		return nil, err
	}

	// 6. Salva no banco de dados
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// 7. Retorna o DTO de resposta (NUNCA retornando a senha ou o hash)
	var outputRoles []string
	for _, r := range user.Roles() {
		outputRoles = append(outputRoles, string(r))
	}

	return &dto.UserOutputDTO{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email(),
		Roles:     outputRoles,
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt(),
	}, nil
}