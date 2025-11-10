# Estágio 1: Build (Compilação)
# Usamos a imagem completa do Go para compilar
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copia os arquivos de módulo e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do código-fonte
COPY . .

# Compila o binário da aplicação
RUN go build -o /main ./cmd/api/main.go

# Estágio 2: Execução (imagem final leve)
# Usamos uma imagem 'alpine' pura, que é minúscula
FROM alpine:latest

WORKDIR /app

# Copia SOMENTE o binário compilado do estágio 'builder'
COPY --from=builder /main .

# Expõe a porta que o Gin vai usar
EXPOSE 8080

# O comando para iniciar a aplicação
# ENTRYPOINT ["/app/main"]
CMD ["/app/main"]