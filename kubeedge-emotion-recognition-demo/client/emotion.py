#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Tue Jan  7 03:18:26 2020

@author: gaowei
"""

import os
import base64
import json
import requests
import cv2
import numpy as np

faceCascade = cv2.CascadeClassifier('/home/haarcascade_frontalface_default.xml')

env_dist = os.environ

faceemotion_server = env_dist.get('FACEEMOTION_SERVER')
faceemotion_port = env_dist.get('FACEEMOTION_PORT')

cap = cv2.VideoCapture(0)
timeF = 10
count = 0
request_url = "http://%s:%s/model/methods/predict" % (faceemotion_server, faceemotion_port)
init = 0

while(True):
    count = count + 1
    headers = {'accept': 'application/json','content-type': 'application/json'}
    ret, frame = cap.read()
    if (count%timeF != 0):
        continue

    img = cv2.flip(frame, 1)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    faces = faceCascade.detectMultiScale(
        gray,     
        scaleFactor=1.2,
        minNeighbors=5,     
        minSize=(20, 20)
    )

    if len(faces) == 0 and init == 1:
        continue

    init = 1

    small_frame = cv2.resize(frame,(0,0),fx = 0.88,fy = 0.88)
    image = cv2.imencode('.jpg',small_frame)[1]
    img_BASE64 = str(base64.b64encode(image))[2:-1]
    
    post_data = {"img_base64": img_BASE64}   
    data = json.dumps(post_data).encode(encoding = 'utf-8')
    response = requests.post(url=request_url, data=data, headers=headers)
    res = json.loads(response.text)#dict
    res_dict = json.loads(res["value"])
    src = res_dict["img_url"]
    data_processd = src.split(',')[1]
    # base64解码
    image_data = base64.b64decode(data_processd)
    # 转换为np数组
    img_array = np.fromstring(image_data, np.uint8)
    # 转换成opencv可用格式
    img = cv2.imdecode(img_array, cv2.COLOR_RGB2BGR)
    cv2.imshow('img', img)
    # Display the resulting frame
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break
cap.release()
cv2.destroyAllWindows()
