FROM golang:1.22-alpine as build

ARG MAIN_FILE_PATH
ARG OUT_FILE

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd ${MAIN_FILE_PATH} && CGO_ENABLED=0 go build -o ${OUT_FILE}

ENV BIN_FILE=${OUT_FILE}

ENTRYPOINT $BIN_FILE
