# Stage 1: Build frontend
FROM node:lts AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:latest AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist/ ./web/dist/
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /parcel-tracker ./cmd/parcel-tracker

# Stage 3: Final minimal image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
COPY --from=backend /parcel-tracker /usr/local/bin/parcel-tracker
VOLUME /data
ENV DATABASE_PATH=/data/parcel-tracker.db
EXPOSE 8080
ENTRYPOINT ["parcel-tracker"]
