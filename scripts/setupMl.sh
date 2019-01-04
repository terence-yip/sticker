#!/bin/bash

RESEARCH_DIR=`pwd`/../../models/research
export PYTHONPATH=$PYTHONPATH:$RESEARCH_DIR:$RESEARCH_DIR/slim

cd $RESEARCH_DIR
protoc object_detection/protos/*.proto --python_out=.

MODEL_BASE=ssd_mobilenet_v1_coco_2017_11_17
MODEL_FILENAME=$MODEL_BASE.tar.gz
wget http://download.tensorflow.org/models/object_detection/$MODEL_FILENAME
tar -xzvf $MODEL_FILENAME

cd -
