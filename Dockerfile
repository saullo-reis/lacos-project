# Usa uma imagem base oficial do Golang
FROM golang:1.20-alpine

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia o arquivo go.mod e go.sum
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código-fonte da aplicação
COPY . .

# Compila a aplicação
RUN go build -o main .

# Define a porta que será exposta
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]
