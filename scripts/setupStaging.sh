#!/bin/bash

ROOT_DIR=../

#Setup staging
OUTPUT_DIR=$ROOT_DIR/staging
STATIC_OUTPUT=$OUTPUT_DIR/static
mkdir -p $STATIC_OUTPUT

#Generate static files
./generateCss.sh

#Copy static files
cp $ROOT_DIR/static/js/* $STATIC_OUTPUT
cp $ROOT_DIR/static/generated/*.css $STATIC_OUTPUT
mkdir -p $OUTPUT_DIR/images
