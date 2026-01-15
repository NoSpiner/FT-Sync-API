FROM golang:latest AS build
ENV CGO_ENABLED=0
WORKDIR /FT
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go get github.com/gin-gonic/gin
RUN go get github.com/glebarez/go-sqlite 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ft-sync-api .
# Stage 2 - Runtime
FROM scratch
WORKDIR /app
COPY --from=build /FT/ft-sync-api .
EXPOSE 9191
ENTRYPOINT ["./ft-sync-api", "-addr", "0.0.0.0"] 
