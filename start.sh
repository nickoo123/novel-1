#!/bin/bash

set -eu

docker system prune -a
docker build -t novel .
docker stop novels && docker rm novels
docker run -itd --name novels -p 8081:8081 --restart=always -v /data/novel/sitemap:/go/static/sitemap -v /data/novel/up:/go/static/up -v /root/.config/rclone:/go/rclone  novel:latest