FROM postgres:16-alpine

ENV POSTGRES_DB=boxes
ENV POSTGRES_PASSWORD=adatbazisjelszo
ENV POSTGRES_INITDB_ARGS="--locale hu_HU.iso88592"

ADD init-boxes-db.sh /docker-entrypoint-initdb.d/init-boxes-db.sh 

