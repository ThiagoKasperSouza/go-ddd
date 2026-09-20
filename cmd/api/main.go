package main

import (
	"database/sql"
	"fmt"

	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"

	"go-ddd/internal/infra/db/postgres"
	usecaseAgency "go-ddd/internal/usecase/agency"
	handlerAgency "go-ddd/internal/infra/http/handler/agency"
)

// Função auxiliar para ler variáveis de ambiente com fallback padrão
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func configureDb() (*sql.DB, error) {
	// 1. Configuração e Conexão com o Banco de Dados
	dbURL := getEnv("DATABASE_URL", "")
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close() // Fecha o pool se o Ping falhar antes de retornar o erro
		return nil, err
	}

	log.Println("Conexão com PostgreSQL estabelecida com sucesso!")
	return db, nil
}

func injectAgency(db *sql.DB) *handlerAgency.AgencyHandler {
	// 2. Injeção de Dependências (De baixo para cima)

	// Repositório
	agencyRepo := postgres.NewAgencyRepository(db)

	// Casos de Uso (Use Cases)
	createAgencyUseCase := usecaseAgency.NewCreateAgencyUseCase(agencyRepo)
	getAgencyUseCase := usecaseAgency.NewGetAgencyUseCase(agencyRepo)
	listAgenciesUseCase := usecaseAgency.NewListAgenciesUseCase(agencyRepo)
	updateAgencyUseCase := usecaseAgency.NewUpdateAgencyUseCase(agencyRepo)
	deleteAgencyUseCase := usecaseAgency.NewDeleteAgencyUseCase(agencyRepo)

	// Handler Web
	return handlerAgency.NewAgencyHandler(
		createAgencyUseCase,
		getAgencyUseCase,
		listAgenciesUseCase,
		updateAgencyUseCase,
		deleteAgencyUseCase
	)

}

func main() {
	// 1. Carrega o arquivo .env
	// Se o arquivo .env não existir (ex: no ambiente de produção/Docker), ele ignora o erro
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	db,err := configureDb()
	if err != nil {
		fmt.Println("Erro ao configurar banco: ",err)
	}
	agencyHandler := injectAgency(db)
	// 3. Configuração de Rotas com o ServeMux Nativo (Go 1.22+)
	mux := http.NewServeMux()

	// O Go 1.22+ aceita métodos HTTP e parâmetros entre chaves {id} nativamente
	mux.HandleFunc("POST /agencies", agencyHandler.Create)
	mux.HandleFunc("GET /agencies", agencyHandler.ListAll)
	mux.HandleFunc("GET /agencies/{id}", agencyHandler.GetByID)
	mux.HandleFunc("PUT /agencies/{id}", agencyHandler.Update)
	mux.HandleFunc("DELETE /agencies/{id}", agencyHandler.Delete)

	// 4. Inicialização do Servidor HTTP
	port := getEnv("PORT", "8080")
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("Servidor rodando na porta %s...", port)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}

}