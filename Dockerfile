FROM golang:1.15-alpine3.13 as builder

RUN apk add --update --no-cache make

WORKDIR /app
COPY . .
RUN make

# Final image
FROM alpine:3.13

COPY --from=builder /app/swaggerui-server /usr/local/bin/swaggerui-server
CMD swaggerui-server -docs PetStore=https://petstore.swagger.io/v2/swagger.json
