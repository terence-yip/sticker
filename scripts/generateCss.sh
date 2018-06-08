#!/bin/bash
OUTPUT_DIR=../static/generated
mkdir -p $OUTPUT_DIR
sass ../static/scss/sticker.scss:$OUTPUT_DIR/sticker.css --no-cache --style compressed
