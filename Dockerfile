# Stage 1 Build
From golang:latest AS build
WORKDIR /FT
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o FT-Sync-API . 
# Stage 2 - Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=build /FT/FT-Sync-API .
EXPOSE 9191
ENTRYPOINT ["./FT-Sync-API", "-addr", "0.0.0.0"] 
