FROM golang:1.16 as builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM jrottenberg/ffmpeg:4.4-alpine as runner

WORKDIR /app

COPY --from=builder /workspace/agqr-toshitai-recording /app/

ENTRYPOINT [ "" ]
