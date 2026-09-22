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
	"go-ddd/internal/infra/http/middleware"
	domainIdentity "go-ddd/internal/domain/identity"
	handlerAgency "go-ddd/internal/infra/http/handler/agency"
	handlerAuth "go-ddd/internal/infra/http/handler/auth"
	handlerRoutes "go-ddd/internal/infra/http/handler/routes"
	handlerShapes "go-ddd/internal/infra/http/handler/shapes"
	handlerTrips "go-ddd/internal/infra/http/handler/trips"
	usecaseAgency "go-ddd/internal/usecase/agency"
	usecaseIdentity "go-ddd/internal/usecase/identity"
	usecaseRoutes "go-ddd/internal/usecase/routes"
	usecaseShapes "go-ddd/internal/usecase/shapes"
	usecaseTrips "go-ddd/internal/usecase/trips"
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
		deleteAgencyUseCase,
	)

}

func injectLogin(db *sql.DB) *handlerAuth.AuthHandler{
	userRepo := postgres.NewUserRepository(db)

	// 2. Inicializa Serviços de Domínio
	passwordHasher := domainIdentity.NewBcryptHasher()
	createUserUseCase := usecaseIdentity.NewCreateUserUseCase(userRepo, passwordHasher)
	loginUseCase := usecaseIdentity.NewLoginUseCase(userRepo, passwordHasher)
	return handlerAuth.NewAuthHandler(loginUseCase, createUserUseCase)
}

func injectRoutes(db *sql.DB) *handlerRoutes.RouteHandler{
	routesRepo := postgres.NewRouteRepository(db)

	// 2. Inicializa Serviços de Domínio
	getRouteUseCase := usecaseRoutes.NewGetRouteUseCase(routesRepo)
	listRoutesUseCase := usecaseRoutes.NewListRoutesUseCase(routesRepo)
	return handlerRoutes.NewRouteHandler(getRouteUseCase, listRoutesUseCase)
}

func injectShapes(db *sql.DB) *handlerShapes.ShapeHandler{
	shapesRepo := postgres.NewShapeRepository(db)

	// 2. Inicializa Serviços de Domínio
	getShapeUseCase := usecaseShapes.NewGetShapeUseCase(shapesRepo)
	listShapesUseCase := usecaseShapes.NewListShapesUseCase(shapesRepo)
	return handlerShapes.NewShapeHandler(getShapeUseCase, listShapesUseCase)
}

func injectTrips(db *sql.DB) *handlerTrips.TripHandler{
	tripsRepo := postgres.NewTripRepository(db)

	// 2. Inicializa Serviços de Domínio
	getTripUseCase := usecaseTrips.NewGetTripUseCase(tripsRepo)
	listTripsUseCase := usecaseTrips.NewListTripsUseCase(tripsRepo)
	return handlerTrips.NewTripHandler(getTripUseCase, listTripsUseCase)
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
	authHandler := injectLogin(db)
	routesHandler := injectRoutes(db)
	shapesHandler := injectShapes(db)
	tripsHandler := injectTrips(db)
	// 3. Configuração de Rotas com o ServeMux Nativo (Go 1.22+)
	mux := http.NewServeMux()

	// O Go 1.22+ aceita métodos HTTP e parâmetros entre chaves {id} nativamente
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /users", authHandler.Register)
	mux.Handle("GET /trips", middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleUser)(
				http.HandlerFunc(tripsHandler.ListAll),
			),
		),
	)
	mux.Handle("GET /trips/{id}", middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleUser)(
				http.HandlerFunc(tripsHandler.GetByID),
			),
		),
	)
	mux.Handle("GET /shapes", middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleUser)(
				http.HandlerFunc(shapesHandler.ListAll),
			),
		),
	)
	mux.Handle("GET /shapes/{id}", middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleUser)(
				http.HandlerFunc(shapesHandler.GetByID),
			),
		),
	)
	mux.Handle("GET /routes", middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleUser)(
				http.HandlerFunc(routesHandler.ListAll),
			),
		),
	)
	mux.Handle("GET /routes/{id}", middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleUser)(
				http.HandlerFunc(routesHandler.GetByID),
			),
		),
	)
	mux.Handle("POST /agencies", 
		middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleAdmin)(
				http.HandlerFunc(agencyHandler.Create),
			),
		),
	)
	mux.Handle("GET /agencies", middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleUser)(
				http.HandlerFunc(agencyHandler.ListAll),
			),
		),
	)
	mux.Handle("GET /agencies/{id}", middleware.EnsureAuthenticated(
			middleware.RequireRole(domainIdentity.RoleUser)(
				http.HandlerFunc(agencyHandler.GetByID),
			),
		),
	)
	mux.Handle("PUT /agencies/{id}", 
	middleware.EnsureAuthenticated(
		middleware.RequireRole(domainIdentity.RoleAdmin)(
			http.HandlerFunc(agencyHandler.Update),
			),
		),
	)
	mux.Handle("DELETE /agencies/{id}", 
	middleware.EnsureAuthenticated(
		middleware.RequireRole(domainIdentity.RoleAdmin)(
			http.HandlerFunc(agencyHandler.Delete),
		),
	),
)

	// 4. Inicialização do Servidor HTTP
	port := getEnv("PORT", "8080")
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("Servidor rodando na porta %s...", port)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}

}