#!/bin/bash
OUTPUT_DIR=../assets/css
mkdir -p $OUTPUT_DIR
sass ../assets/scss/sticker.scss:$OUTPUT_DIR/sticker.css --style compressed
