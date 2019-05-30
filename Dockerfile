FROM golang:1.12.4-alpine AS builder
# Setup Args
ARG APP_SRC
ARG MODS
ARG USER_SVC
ARG USER_SRC
WORKDIR ${APP_SRC}

# Download required packages
RUN apk --update add git ca-certificates

# Dowmload modules
COPY go.mod go.sum* ./
RUN go mod download

# Run tests and update mod fiels
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go test ./...

# Copy new mod files
RUN mkdir -p ${MODS}
RUN cp go.mod ${MODS} && cp go.sum ${MODS}

# Build User Service
COPY ${USER_SRC}/. /${USER_SRC}/.
# RUN cp -R ${MODS}/* ${USER_SRC}/.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${USER_SVC} ${USER_SRC}/.

# Adjusting ownership
RUN addgroup -S appuser && adduser -S appuser -G appuser -u 1000
RUN chown -R appuser ${APP_SRC} ${MODS}
USER appuser