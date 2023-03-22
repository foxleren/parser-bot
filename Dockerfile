FROM golang:1.18-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# Install Python and pip
RUN apt-get -y install python3 python3-pip

# Install required Python packages
COPY requirements.txt .
RUN pip3 install -r requirements.txt

# build go app
RUN go mod download
RUN go build -o parser-bot ./cmd/bot/main.go

CMD ["./parser-bot"]