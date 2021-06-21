FROM golang:latest as goserver

WORKDIR /forum/server

ADD ./server/go.mod .
ADD ./server/go.sum .
RUN go mod download

ADD ./server .
RUN go build server_main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

ADD ./postgres/DB_BACKUP.sql .

ENV PGPASSWORD ohhibitchitsme
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER mikhail WITH SUPERUSER PASSWORD 'ohhibitchitsme';" &&\
    createdb -O mikhail forum &&\
    psql -h localhost -d forum -U mikhail -p 5432 -a -q -f ./DB_BACKUP.sql &&\
    /etc/init.d/postgresql stop

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

COPY --from=goserver /forum/server .

EXPOSE 5000/tcp

CMD service postgresql start &&  ./server_main