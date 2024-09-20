# Initial stage: download modules
FROM golang:1.20 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.20 as builder

ARG NAME_ENV=hmb_bot

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 ${NAME_ENV}

RUN mkdir -p /${NAME_ENV}
ADD . /${NAME_ENV}
WORKDIR /${NAME_ENV}

RUN go get -u github.com/swaggo/swag/cmd/swag

RUN swag init -g ./app/cmd/main.go

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
     go build -o /go/bin/${NAME_ENV} ./app/cmd


# Final stage: Run the binary
FROM scratch

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER ${NAME_ENV}

# and finally the binary
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/${NAME_ENV} /${NAME_ENV}

CMD ["/hmb_bot"]
