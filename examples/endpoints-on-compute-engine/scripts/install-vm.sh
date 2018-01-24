#!/bin/bash

sudo apt-get install git -y
curl https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz | tar xzf -
wget https://raw.githubusercontent.com/GoogleCloudPlatform/golang-samples/master/endpoints/getting-started/app.go
GOPATH=$PWD GOROOT=$PWD/go go/bin/go get ./... 2> /tmp/goget_log.txt > /tmp/goget_outlog.txt
echo "deb http://packages.cloud.google.com/apt google-cloud-endpoints-jessie main" | sudo tee /etc/apt/sources.list.d/google-cloud-endpoints.list
curl --silent https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
sudo apt-get update && sudo apt-get install google-cloud-sdk
sudo apt-get install endpoints-runtime
sudo echo "PORT=80" >> /etc/default/nginx
PORT=8081 GOPATH=$PWD GOROOT=$PWD/go go/bin/go run app.go > /tmp/gorun_outlog.txt 2> /tmp/gorun_log.txt &
sleep 10
sudo service nginx restart
