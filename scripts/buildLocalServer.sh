#!/bin/bash

ROOT_DIR=../
ABSOLUTE_ROOT_DIR=$(pwd)/$ROOT_DIR
OUTPUT_DIR=$ROOT_DIR/staging

./setupStaging.sh

go build -o $OUTPUT_DIR/app $ROOT_DIR/app/test 

rm $OUTPUT_DIR/static/sticker.css
rm $OUTPUT_DIR/static/app.js

ln -s $ABSOLUTE_ROOT_DIR/static/generated/sticker.css $OUTPUT_DIR/static/sticker.css
ln -s $ABSOLUTE_ROOT_DIR/static/js/app.js $OUTPUT_DIR/static/app.js