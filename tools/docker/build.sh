#!/bin/bash -e

docker build . -t pankona/golangci:`cat version.txt`
