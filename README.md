# 🌎 EcoPonto API (Backend)

API de backend para o projeto EcoPonto — uma iniciativa de extensão universitária focada em "Agir Local, Pensar Global". Fornece uma plataforma para mapear pontos de coleta de resíduos (ODS 12) e educar sobre descarte correto (ODS 4).

## Principais características

- Autenticação JWT para rotas de administrador (/admin).
- Geocoding híbrido no POST /api/ecopontos:
  - Usa lat/lon quando fornecidos.
  - Caso contrário, consulta Nominatim (OpenStreetMap) para obter coordenadas a partir do endereço.
- Consultas geoespaciais com PostGIS (ex.: ST_DWithin) no GET /api/ecopontos para buscar pontos dentro de um raio.
- CRUD completo para ecopontos e gestão de usuários/admin via endpoints protegidos.
- Migrations geridas com golang-migrate/migrate.
- Totalmente containerizado com Docker.

## Stack tecnológico

- Linguagem: Go (Golang)
- Framework: Gin
- Banco de dados: PostgreSQL + PostGIS
- DB access: SQLx
- Migrations: golang-migrate/migrate
- Autenticação: JWT (golang-jwt/jwt)
- Configuração: variáveis de ambiente (sem Viper)
- Containerização: Docker & Docker Compose

## Como executar (ambiente de desenvolvimento)

### Pré‑requisitos

- Docker & Docker Compose
- CLI do golang-migrate (para executar migrations localmente)

Instale a CLI de migration:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Passo 1 — Configurar variáveis de ambiente

Copie `.env.example` para `.env` na raiz do projeto e ajuste conforme necessário.

Exemplo `.env`:

```env
# Configurações da API
API_PORT=8080

# Credenciais do Banco de Dados (usadas pelo docker-compose e pelo script de migration)
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=ecoponto_db
DB_PORT=5432

# String de conexão para a CLI de migration (aponta para localhost)
DATABASE_URL="postgresql://admin:admin@localhost:5432/ecoponto_db?sslmode=disable"

# Segredo JWT (gere um valor forte)
JWT_SECRET="SEU_SEGREDO_FORTE_AQUI"
```

### Passo 2 — Subir containers

Com o Docker em execução:

```bash
docker-compose up --build -d
```

### Passo 3 — Rodar migrations

Com o banco rodando, aplique as migrations:

```bash
migrate -database "postgresql://admin:admin@localhost:5432/ecoponto_db?sslmode=disable" -path internal/database/migrations up
```

### Passo 4 — Criar usuário admin (seed)

Rode o seeder para criar o primeiro administrador:

```bash
go run ./cmd/seeder/main.go
```

(O e‑mail e senha padrão estão definidos em cmd/seeder/main.go.)

A API ficará acessível em http://localhost:8080

## Endpoints principais

Autenticação

- POST /api/auth/login — recebe { email, senha } e retorna JWT.

Ecopontos (público)

- GET /api/ecopontos — lista por proximidade. Query params: lat, lon (obrigatórios), dist (opcional, metros), tipo (opcional).
- GET /api/ecopontos/:id — detalhes de um ecoponto.

Ecopontos (admin — protegido por JWT)

- GET /api/ecopontos/all — lista todos (para gestão/admin).
- POST /api/ecopontos — cria ecoponto (aceita endereço ou lat/lon).
- PUT /api/ecopontos/:id — atualiza ecoponto.
- DELETE /api/ecopontos/:id — apaga ecoponto.

## Observações técnicas

- Geocoding: usa Nominatim (OpenStreetMap) por padrão quando não há lat/lon.
- Proximidade: consultas geoespaciais realizadas via PostGIS (p.ex. ST_DWithin).
- Migrations e seeds facilitam reprodução do ambiente em desenvolvimento.

## Autor

Desenvolvido por Eric Oliveira

- [GitHub: EricOliveiras](https://github.com/EricOliveiras)
- [Linkedin: HeyEriic](https://www.linkedin.com/in/heyeriic/)
