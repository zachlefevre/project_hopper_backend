FROM golang:latest as builder
WORKDIR /go/src/github.com/zachlefevre/project_hopper_backend/algorithmSync
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest
LABEL maintainer="@zachlefevre"
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/zachlefevre/project_hopper_backend/algorithmSync .
CMD [ "./algorithmSync" ]