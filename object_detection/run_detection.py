#!/usr/bin/python3

import object_detector
import PIL

def main():
    IMAGE_DIR = "/host/images/"
    FILENAME = "image1.jpeg"
    IMAGE_PATH = IMAGE_DIR + FILENAME
    OUTPUT_IMAGE = IMAGE_DIR + "od-" + FILENAME
    image = PIL.Image.open(IMAGE_PATH)
    od = object_detector.ObjectDetector()
    op = od.detect(image)
    op.save(OUTPUT_IMAGE)

if __name__ == "__main__":
    main()
