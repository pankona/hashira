#!/bin/bash -e

rm -rf ./protoc
mkdir ./protoc
cd ./protoc
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.18.0/protoc-3.18.0-linux-x86_64.zip
unzip ./protoc-3.18.0-linux-x86_64.zip
./bin/protoc --version
