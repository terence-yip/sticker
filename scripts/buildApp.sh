#!/bin/bash
  
ROOT_DIR=../
OUTPUT_DIR=$ROOT_DIR/staging

./setupStaging.sh

GOARCH=arm GOARM=7 GOOS=linux go build -o $OUTPUT_DIR/app $ROOT_DIR/app/pi 
