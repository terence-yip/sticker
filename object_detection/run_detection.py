#!/usr/bin/python3

import detect_object
import PIL
import numpy as np

BASE_PATH = "/host/models/research/"
MODEL_PATH = BASE_PATH + "ssd_mobilenet_v1_coco_2017_11_17/frozen_inference_graph.pb"
detection_graph = detect_object.load_model(MODEL_PATH)

PATH_TO_LABELS = "/host/models/research/object_detection/data/mscoco_label_map.pbtxt"
#PATH_TO_LABELS = BASE_PATH + "ssd_mobilenet_v1_coco_2017_11_17/mscoco_label_map.pbtxt"
NUM_CLASSES = 90
category_index = detect_object.get_category_index(PATH_TO_LABELS, NUM_CLASSES)

IMAGE_DIR = "/host/images/"
FILENAME = "image1.jpeg"
IMAGE_PATH = IMAGE_DIR + FILENAME
OUTPUT_IMAGE = IMAGE_DIR + "od-" + FILENAME
IMAGE_SIZE = (12, 8) # Inches
MAX_SIZE = [512,512]
image_np = detect_object.load_image_into_numpy_array(IMAGE_PATH, MAX_SIZE)
output_dict = detect_object.run_inference(image_np, detection_graph)
detect_object.add_labels_to_image(image_np, output_dict, category_index, IMAGE_SIZE)
im = PIL.Image.fromarray(np.uint8(image_np))
im.save(OUTPUT_IMAGE)
