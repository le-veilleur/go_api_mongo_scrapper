# Utiliser l'image officielle de Go
FROM golang:1.20-alpine

# Définir le répertoire de travail pour le serveur
WORKDIR /go_api_mongo_scrapper

# Copier les fichiers nécessaires
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le code source
COPY . .

# Construire le binaire du serveur (dans /go_api_mongo_scrapper)
RUN go build -o server ./main.go

# Construire le binaire du scraper (dans /go_api_mongo_scrapper/scraper)
WORKDIR /go_api_mongo_scrapper/scraper
RUN go build -o scraper ./scraper.go

# Revenir au répertoire principal
WORKDIR /go_api_mongo_scrapper
