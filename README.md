# Brasil CEP APIs

Exemplo prático de race condition em Go, onde duas APIs de consulta de CEP (BrasilAPI e ViaCEP) competem para retornar o resultado mais rápido.

## Como funciona

O programa faz requisições simultâneas para:
- BrasilAPI
- ViaCEP

A primeira API que responder terá seu resultado exibido, enquanto a resposta mais lenta é descartada. Se nenhuma API responder em 1 segundo, um timeout é retornado.

## Instalação

```bash
git clone https://github.com/lcidral/goExpertBrasilCep.git
cd goExpertBrasilCep
```

## Uso

### Local
```bash
go run main.go -cep 01153000
```

### Docker
```bash
# Build
docker-compose build

# Execução
CEP=01153000 docker-compose up

# Ou use o CEP default
docker-compose up
```

## Resposta

O programa retorna:
- API utilizada
- CEP consultado
- Estado
- Cidade
- Rua
- Bairro

## Tecnologias

- Go 1.21+
- Docker
- Docker Compose
- Multithreading com goroutines
- Channels
- APIs REST

## Features

- ✨ Race condition entre APIs
- 🔄 Consultas paralelas
- ⏱️ Timeout de 1 segundo
- 🐳 Containerização
- 🔍 Retorno mais rápido
