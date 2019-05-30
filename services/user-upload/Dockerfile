# User service
FROM scratch

ARG APP_SRC
ARG SVC
ARG CONFIG

COPY --from=lazy-panda-upload-core /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ${CONFIG} ${CONFIG}
COPY --from=lazy-panda-upload-core ${APP_SRC}/${SVC} /${SVC}