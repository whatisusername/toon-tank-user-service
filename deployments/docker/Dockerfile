FROM golang:1.22.2-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY api/ ./api/
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN [ -f go.sum ] || go mod tidy
RUN go build -tags lambda.norpc -o main ./cmd/main.go

# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /src/main /main
ENTRYPOINT [ "/main" ]
