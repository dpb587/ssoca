FROM debian:stretch-slim
RUN true \
  && apt-get update \
  && apt-get install -y curl git openssl \
  && curl -fLo /usr/local/bin/meta4-repo https://github.com/dpb587/metalink/releases/download/v0.3.0/meta4-repo-0.3.0-linux-amd64 \
  && echo "dce06d852585d5a5201b2d21cd864cbdf81dff093747ab368acb8d4fa1e11216  /usr/local/bin/meta4-repo" | sha256sum -c - \
  && chmod +x /usr/local/bin/meta4-repo \
  && curl -fLo /tmp/hugo.tar.gz https://github.com/gohugoio/hugo/releases/download/v0.53/hugo_extended_0.53_Linux-64bit.tar.gz \
  && tar -xzf /tmp/hugo.tar.gz -C /usr/local/bin hugo \
  && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
