FROM golang:1.8.1-alpine
RUN true \
  && apk --no-cache add bash ca-certificates curl git openssl \
  && curl -fLo /usr/local/bin/meta4 https://github.com/dpb587/metalink/releases/download/v0.1.0/meta4-0.1.0-linux-amd64 \
  && echo "235bc60706793977446529830c2cb319e6aaf2da  /usr/local/bin/meta4" | sha1sum -c - \
  && chmod +x /usr/local/bin/meta4
