#!/usr/bin/python3

import detect_object
import PIL
import numpy as np

class ObjectDetector():
    def __init__(self):
        self.BASE_PATH = "/host/models/research/"
        self.MODEL_PATH = self.BASE_PATH + "ssd_mobilenet_v1_coco_2017_11_17/frozen_inference_graph.pb"
        self.detection_graph = detect_object.load_model(self.MODEL_PATH)
        self.PATH_TO_LABELS = "/host/models/research/object_detection/data/mscoco_label_map.pbtxt"
        self.NUM_CLASSES = 90
        self.category_index = detect_object.get_category_index(self.PATH_TO_LABELS, self.NUM_CLASSES)
        self.MAX_SIZE = [512,512]
        self.IMAGE_SIZE = (12, 8) # Inches

    def detect(self, image):
        image_np = self.load_image_into_numpy_array(image)
        output_dict = detect_object.run_inference(image_np, self.detection_graph)
        detect_object.add_labels_to_image(image_np, output_dict, self.category_index, self.IMAGE_SIZE)
        return PIL.Image.fromarray(np.uint8(image_np))

    def load_image_into_numpy_array(self, image):
        image.thumbnail(self.MAX_SIZE)
        (im_width, im_height) = image.size
        return np.array(image.getdata()).reshape(
            (im_height, im_width, 3)).astype(np.uint8)
