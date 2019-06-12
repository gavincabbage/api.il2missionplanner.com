ARG GO_VERSION=1.12

FROM golang:${GO_VERSION}-alpine AS build

# Create the user and group files that will be used in the running container to
# run the process an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -mod 'vendor' \
    -o /api ./cmd/...

FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=build /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the first stage.
COPY --from=build /api /api

EXPOSE 8080
EXPOSE 6060

# Perform any further action as an unprivileged user.
USER nobody:nobody

ENTRYPOINT ["/api"]
