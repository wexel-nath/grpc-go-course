#!/bin/bash

protoc "$1" --go_out=plugins=grpc:.
