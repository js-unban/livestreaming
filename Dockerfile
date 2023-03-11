FROM golang:bullseye as builder

ENV WORKDIR /app
WORKDIR ${WORKDIR}

COPY config ${WORKDIR}/config
COPY models ${WORKDIR}/models
COPY *.go ${WORKDIR}/
COPY go.* ${WORKDIR}/

RUN go build

FROM golang:bullseye
ENV WORKDIR /app
WORKDIR ${WORKDIR}
COPY --from=builder ${WORKDIR}/livestreaming ${WORKDIR}
COPY static/ ${WORKDIR}/static/
COPY templates/ ${WORKDIR}/templates/
COPY content/ ${WORKDIR}/content/
COPY config/dev.json ${WORKDIR}/config/dev.json
