# 1. Go-Image als Basis verwenden
FROM golang:1.25-alpine

# 2. Arbeitsverzeichnis setzen
WORKDIR /app

# 3. Go-Dateien kopieren
COPY . .

# 4. Build ausführen
RUN go build -o app

# 5. Port freigeben
EXPOSE 8080

# 6. Container starten
CMD ["./app"]
