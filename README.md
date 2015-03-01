# github.com/husobee/dockerspew - Dockerspew

## Introduction

Docker spew is web application that will watch docker containers and spew out 
events and stdout/stderr from containers.

## Environment Variables

This package requires the following environmental variables:

DOCKERSPEW_ENDPOINT -> the api endpoint of docker ("unix:///var/run/docker.sock")
DOCKERSPEW_SERVER_HOST -> the server/port you wish to run this service on (":8080")
DOCKERSPEW_LOG_LEVEL -> the log level you wish to obtain
DOCKERSPEW_GOMAXPROCS -> Max procs overload (defaults to num cpu)

## Usage

Connect to ws://$DOCKERSPEW_SERVER_HOST/spew to start seeing streaming logs from all running containers

to filter by container name

Connect to ws://$DOCKERSPEW_SERVER_HOST/spew?contains=$CONTAINER_NAME_PARTIAL to start seeing streaming logs from filtered running containers

