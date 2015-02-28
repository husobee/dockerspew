# github.com/husobee/dockerspew - Dockerspew

## Introduction

Docker spew is web application that will watch docker containers and spew out 
events and stdout/stderr from containers.

## Code Organization

main.go - entry point start server
controllers - web controllers that house the handler functions
context - globals for context
middlewares - negroni middlewares needed for the application
models - data access and manipulation
content - content based code for request and response
