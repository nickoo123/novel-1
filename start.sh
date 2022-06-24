#!/bin/bash

set -eu

docker stop novels && docker rm novels
docker system prune -a
docker build -t novel .
docker run -itd --name novels -p 8081:8081 --restart=always -v /data/novel/up:/go/static/up novel:latest