#!/bin/bash

ROOT_DIR=../
OUTPUT_DIR=$ROOT_DIR/staging

./setupStaging.sh

go build -o $OUTPUT_DIR/app $ROOT_DIR/app/test 
