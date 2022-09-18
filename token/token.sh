#!/bin/bash

docker run -itd --name refreshToken -p 8010:8010 -v /usr/local/var/www/novel-1/rclone:/app python:3.7