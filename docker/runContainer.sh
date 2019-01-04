#!/bin/bash
SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
BASEDIR=$SCRIPTPATH/../..
docker run -it --privileged -p 8080:8080 --rm -v $BASEDIR:/host sticker bash
