# Go DDD API - GTFS
API desenvolvida em **Go** seguindo os princípios de **Domain-Driven Design (DDD)** e arquitetura em camadas para o gerenciamento de tabelas do padrão GTFS.

## 🛠️ Tecnologias Utilizadas

- **Go** (Golang)
- **PostgreSQL** (Banco de dados relacional)
- **SQL** para scripts de migração e criação de tabelas
- **Kafka, Prometheus e Grafana** Observabilidade

## 📁 Estrutura do Projeto:
<img width="270" height="421" alt="image"  src="https://github.com/user-attachments/assets/326630d3-0ee6-4497-88f4-6bea28fe693c" />
<br>

## Diagrama Agency:
<br>
<img width="1000"  alt="agency" src="https://github.com/user-attachments/assets/b06a2ac4-fa4c-4a15-9cd1-916c3f5792dc" />

## Diagrama Routes:
<br>
<img width="1000" alt="routes" src="https://github.com/user-attachments/assets/eb5921cb-1e05-4b6d-9785-870712163127" />

## Diagrama Trips:
<img width="1000" alt="trips" src="https://github.com/user-attachments/assets/45f9f43e-d0d4-469f-8d28-54bb754fbd8c" />

## Diagrama Shapes:
<img width="1000" alt="shapes" src="https://github.com/user-attachments/assets/eb275a15-3ba2-46cb-91c6-c04425731c92" />


🏗️ Camadas da Arquitetura
Domain (internal/domain): Contém o coração da aplicação. As regras de negócio fundamentais e as interfaces de repositórios estão concentradas aqui, sem dependências externas.

Use Cases (internal/usecase): Implementa o fluxo das funcionalidades do sistema (CRUD de agências). Interage diretamente com as entidades do domínio.

DTO (internal/dto): Mapeia as estruturas de dados necessárias para a entrada e saída da API.

Infrastructure (internal/infra): Contém os detalhes técnicos e integrações com o mundo externo (PostgreSQL, Handlers HTTP, Roteamento).

Cmd (cmd/api): Ponto de entrada onde as dependências são injetadas e o servidor é inicializado.

🚀 Como Executar o Projeto
Pré-requisitos
Go (v1.18+ recomendado)

PostgreSQL rodando localmente ou via Docker

1. Configurar o Banco de Dados
Execute os scripts SQL presentes no diretório sql/ para criar o banco de dados e as tabelas necessárias:

create_db.sql — Criação do banco de dados.

create_tables.sql — Criação das tabelas relativas à aplicação.

create_index.sql — Criação dos índices.

copy_files.sql — Inserção/carga de dados iniciais (se aplicável).

2. Instalar Dependências
No diretório raiz do projeto, execute:

```Bash
go mod download
```
3. Rodar a Aplicação
Para iniciar o servidor HTTP:

```Bash
go run cmd/api/main.go
```

## Collection de endpoints Insomnia:
Importe o arquivo **Go-ddd collection** no Insomnia

## Refs:
<br>
- https://medium.com/@tomascdmota/the-difference-between-a-dto-an-entity-and-a-model-6f7265b395c9

