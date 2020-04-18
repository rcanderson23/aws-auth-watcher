FROM alpine
RUN apk --no-cache add ca-certificates
COPY aws-auth-watcher /aws-auth-watcher
ENTRYPOINT ["/aws-auth-watcher"]
