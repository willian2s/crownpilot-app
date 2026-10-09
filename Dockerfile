# Estágio 1: compila um binário estático. Só o código de produção entra no
# contexto (veja .dockerignore); testes, tools/ e frontend ficam de fora.
FROM golang:1.27.2-trixie AS build
WORKDIR /src

# Usa exatamente o toolchain da imagem, sem baixar outro em tempo de build.
ENV GOTOOLCHAIN=local

# go.mod antes do código: a camada de dependências fica em cache enquanto
# go.mod/go.sum não mudam.
COPY go.mod go.sum* ./
RUN go mod download

COPY api ./api
COPY cmd ./cmd
COPY internal ./internal

# CGO desligado gera binário estático, que roda numa imagem sem libc.
# -trimpath remove caminhos da máquina de build do binário.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/crownpilot-api ./cmd/crownpilot-api

# Estágio 2: imagem final mínima, sem shell nem gerenciador de pacotes, rodando
# como usuário sem privilégios.
FROM gcr.io/distroless/static-debian13:nonroot AS final
COPY --from=build /out/crownpilot-api /crownpilot-api

# Porta padrão do container; hosting sobrescreve com PORT. Não há default para
# CROWNPILOT_ENVIRONMENT: sem ele o processo encerra antes do listener.
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/crownpilot-api"]
