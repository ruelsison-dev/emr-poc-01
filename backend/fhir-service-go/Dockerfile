FROM golang:1.21-alpine
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o /bin/fhirservice ./
CMD ["/bin/fhirservice"]
