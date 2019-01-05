#!/bin/bash
export GOPATH=/host/gopath
RESEARCH_DIR=`pwd`/../../models/research
export PYTHONPATH=$PYTHONPATH:$RESEARCH_DIR:$RESEARCH_DIR/slim
