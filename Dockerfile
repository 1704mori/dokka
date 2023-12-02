FROM node:21-alpine as node

RUN corepack enable

WORKDIR /build

COPY package*.json ./
RUN npm install

COPY frontend ./frontend
COPY index.html .
COPY tsconfig* .
COPY vite* .
COPY postcss* .
COPY tailwind* .
COPY env.d.ts .

RUN npm run build

FROM golang:1.21.4-alpine AS builder

RUN apk add --no-cache ca-certificates file && mkdir /dokka

WORKDIR /dokka

COPY backend ./backend
COPY go.* .
COPY main.go .

RUN go mod download

COPY --from=node /build/dist ./dist
RUN ls -lah

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dokka

RUN mkdir /dist

FROM alpine:latest

RUN apk add curl

ENV PATH /bin
COPY --from=node /build/dist /dist
COPY --from=builder /dokka/dokka /dokka

EXPOSE 7070

ENTRYPOINT ["/dokka"]
