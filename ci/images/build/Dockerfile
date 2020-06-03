FROM golang:1.14.4-stretch
RUN true \
  && apt-get update \
  && apt-get install -y curl git openssl \
  && curl -fLo /usr/local/bin/meta4 https://github.com/dpb587/metalink/releases/download/v0.1.0/meta4-0.1.0-linux-amd64 \
  && echo "235bc60706793977446529830c2cb319e6aaf2da  /usr/local/bin/meta4" | sha1sum -c - \
  && chmod +x /usr/local/bin/meta4 \
  && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
