FROM golang as builder
RUN export GO111MODULE=on
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build  -mod=mod -o /ksooo-study
CMD ["/ksooo-study"]

FROM alpine

COPY --from=builder /ksooo-study /ksooo-echo
COPY /config.yaml /config.yaml
ENTRYPOINT ["/ksooo-echo"]