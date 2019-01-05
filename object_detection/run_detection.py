#!/usr/bin/python3

import object_detector
from PIL import Image

import io
import picamera

def main():
    stream = io.BytesIO()
    camera = picamera.PiCamera()
    camera.resolution = (1024, 768)
    camera.capture(stream, format='jpeg')
    stream.seek(0)
    image = Image.open(stream)

    IMAGE_DIR = "/host/images/"
    FILENAME = "image1.jpeg"
    OUTPUT_IMAGE = IMAGE_DIR + "od2-" + FILENAME
    od = object_detector.ObjectDetector()
    op = od.detect(image)
    op.save(OUTPUT_IMAGE)

if __name__ == "__main__":
    main()
