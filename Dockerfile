FROM golang:1.23.4-alpine AS build-stage
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY . /build
RUN templ generate
RUN CGO_ENABLED=1 go build -o ./bin/houston ./cmd/houston/main.go

FROM node:22.13.1-alpine AS static-stage
WORKDIR /static
COPY views ./views
COPY package.json tailwind.config.js ./
RUN npm install
RUN npx tailwindcss -i ./views/tailwind.css -o ./styles.css

FROM golang:1.23.4-alpine AS image-stage
RUN apk add --no-cache tzdata
ENV TZ=Europe/Warsaw
WORKDIR /app
RUN mkdir -p ./data
RUN touch ./data/app.db
COPY static ./public
RUN mkdir -p ./public/css
COPY --from=static-stage /static/styles.css /app/public/css/styles.css
COPY --from=build-stage /build/bin/houston /app/houston
COPY --from=build-stage /build/migrations /app/migrations
CMD ["sh", "-c", "./houston"]
