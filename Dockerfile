FROM alpine:latest

# install packages
RUN apk --no-cache add libc6-compat \
		bash

WORKDIR /root/

COPY gou gou
COPY data/ data/

COPY build/version.go version
COPY README.md README.md

COPY .dockerignore .dockerignore
COPY Dockerfile Dockerfile

# 8080 - UI
# 9999 - mesh
EXPOSE 8080 9999

CMD ["./gou"]
