FROM golang:1.5.1

MAINTAINER st3v "https://github.com/st3v"

# install wget
RUN apt-get update
RUN apt-get install -y --force-yes build-essential wget git

# install latest stable cf cli
RUN wget -O cf_cli.tar.gz https://cli.run.pivotal.io/stable?release=linux64-binary
RUN tar xzvf cf_cli.tar.gz
RUN mv cf /usr/local/bin
RUN rm -rf cf_cli.tar.gz

# add godep
RUN go get github.com/tools/godep

# add local dir as workdir
WORKDIR $GOPATH/src/github.com/st3v/cfkit
ADD . $WORKDIR
ENV GOBIN $GOPATH/bin
ENV GOPATH $GOPATH/src/github.com/st3v/cfkit/Godeps/_workspace:$GOPATH

# install ginkgo
RUN go install github.com/onsi/ginkgo/...

# install necessary tools
RUN go get golang.org/x/tools/cmd/cover
RUN go get golang.org/x/tools/cmd/vet
RUN go get github.com/golang/lint/golint
RUN go get github.com/modocache/gover
RUN go get github.com/mattn/goveralls

