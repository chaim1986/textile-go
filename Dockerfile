FROM golang:1.10

# replace shell with bash so we can source files
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

# update the repository sources list
# and install dependencies
RUN apt-get update \
    && apt-get install -y curl \
    && apt-get -y autoclean

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# install gx
RUN go get -u github.com/whyrusleeping/gx \
    && go get -u github.com/whyrusleeping/gx-go

# install go-ipfs@0.4.14 source
RUN echo '{"language": "go", "gxVersion": "0.12.1", "gxDependencies": [{"hash": "QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn", "name": "go-ipfs", "version": "0.4.14"}]}' >package.json
RUN gx install --global
