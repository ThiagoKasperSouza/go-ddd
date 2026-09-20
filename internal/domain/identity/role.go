package identity


// Role representa o tipo customizado para os papéis de usuário
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
)

// String para implementar a interface Stringer (opcional)
func (r Role) String() string {
	return string(r)
}

// IsValid garante que a role recebida de fora é válida
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}