FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build -o aws-auth-watcher


FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build-env /src/aws-auth-watcher /aws-auth-watcher
ENTRYPOINT ["/aws-auth-watcher"]
