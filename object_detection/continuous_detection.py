#!/usr/bin/python3

import object_detector
from PIL import Image

import io
import picamera
import argparse
import time

def main(duration):
    stream = io.BytesIO()
    camera = picamera.PiCamera()
    camera.resolution = (400, 300)

    IMAGE_DIR = "/host/sticker/staging/images/"
    FILENAME = "image.jpg"
    OUTPUT_IMAGE = IMAGE_DIR + FILENAME
    od = object_detector.ObjectDetector()

    while(1):
        time.sleep(duration)
        camera.capture(stream, format='jpeg')
        stream.seek(0)
        image = Image.open(stream)
        t1 = time.time()
        op = od.detect(image)
        op.save(OUTPUT_IMAGE)
        #image.save(OUTPUT_IMAGE)
        t2 = time.time()
        tdiff = t2 - t1
        print(str(tdiff) + " - Wrote to " + OUTPUT_IMAGE)
        stream.seek(0)
        stream.truncate()

if __name__ == "__main__":
    p = argparse.ArgumentParser("Auto picture + image detection")
    p.add_argument('-d', '--duration', type=float, help="Duration between pictures in seconds", default=0.0)
    args = p.parse_args()
    main(args.duration)
