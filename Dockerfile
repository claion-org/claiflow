FROM alpine:latest   as common-build-stage

LABEL version=0.0.1

ENV WORKDIR=/app \
    NAME=claiflow \
    USER=claiflowuser \
    USER_ID=1002 \
    GROUP=claiflowgroup

WORKDIR ${WORKDIR}

RUN apk update && apk add --no-cache bash mariadb-client

RUN mkdir -p ${WORKDIR}/conf

COPY /bin/claiflow ${WORKDIR}/
COPY /conf/claiflow.yaml ${WORKDIR}/conf/
COPY /conf/enigma.yaml ${WORKDIR}/conf/

RUN addgroup ${GROUP} && \
    adduser -D ${USER} -g ${GROUP} -u ${USER_ID} && \
    chown -R ${USER}:${GROUP} ${WORKDIR}/

USER ${USER}

EXPOSE 8099 18099

ENTRYPOINT [ "./claiflow" ]

CMD [ "./conf/claiflow.yaml", "./conf/enigma.yaml" ]