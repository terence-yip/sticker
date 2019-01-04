#!/bin/bash
export GOPATH=/app/gopath
RESEARCH_DIR=`pwd`/../../models/research
export PYTHONPATH=$PYTHONPATH:$RESEARCH_DIR:$RESEARCH_DIR/slim
