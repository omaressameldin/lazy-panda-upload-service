FROM golang:1.12.4-alpine
ARG APP_SRC
ARG MODS
WORKDIR ${APP_SRC}

RUN apk --update add git ca-certificates

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go test ./...

RUN mkdir -p ${MODS}
RUN cp go.mod ${MODS} && cp go.sum ${MODS}

RUN addgroup -S appuser && adduser -S appuser -G appuser -u 1000
RUN chown -R appuser ${APP_SRC} ${MODS}
USER appuser