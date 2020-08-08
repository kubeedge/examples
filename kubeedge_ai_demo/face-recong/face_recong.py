# -!- coding: utf-8 -!-
import json

import cv2
import numpy as np
import align.detect_face
import tensorflow as tf
import facenet
from scipy import misc
import time
import os
import sys
import argparse

minsize = 20
threshold = [ 0.6, 0.7, 0.7 ]
factor = 0.709
def load_and_align_data(img, image_size, margin,pnet, rnet, onet):

    minsize = 20 # minimum size of face
    threshold = [ 0.6, 0.7, 0.7 ]  # three steps's threshold
    factor = 0.709 # scale factor

    img_size = np.asarray(img.shape)[0:2]

    # bounding_boxes shape:(1,5)  type:np.ndarray
    bounding_boxes, _ = align.detect_face.detect_face(img, minsize, pnet, rnet, onet, threshold, factor)
    if len(bounding_boxes) < 1:
        return 0,0,0

    # det = np.squeeze(bounding_boxes[:,0:4])
    det=bounding_boxes

    # print('det shape type')
    # print(det.shape)
    # print(type(det))

    det[:,0] = np.maximum(det[:,0]-margin/2, 0)
    det[:,1] = np.maximum(det[:,1]-margin/2, 0)
    det[:,2] = np.minimum(det[:,2]+margin/2, img_size[1]-1)
    det[:,3] = np.minimum(det[:,3]+margin/2, img_size[0]-1)

    det=det.astype(int)
    crop=[]
    for i in range(len(bounding_boxes)):
        temp_crop=img[det[i,1]:det[i,3],det[i,0]:det[i,2],:]
        aligned=misc.imresize(temp_crop, (image_size, image_size), interp='bilinear')
        prewhitened = facenet.prewhiten(aligned)
        crop.append(prewhitened)

    crop_image=np.stack(crop)

    return 1,det,crop_image


def  main(args,model_dir):
    if args.camera=="0":
        camera=0
    else:
        camera=args.camera
    frame_start_time=0
    fps=0
    font = cv2.FONT_ITALIC
    with tf.Graph().as_default():
        gpu_options = tf.GPUOptions(per_process_gpu_memory_fraction=1.0)
        sess = tf.Session(config=tf.ConfigProto(gpu_options=gpu_options, log_device_placement=False))
        with sess.as_default():
            pnet, rnet, onet = align.detect_face.create_mtcnn(sess, None)
    with tf.Graph().as_default():
        with tf.Session() as sess:
            # 加载模型
            facenet.load_model(model_dir)
            # 返回给定名称的tensor
            images_placeholder = tf.get_default_graph().get_tensor_by_name("input:0")
            embeddings = tf.get_default_graph().get_tensor_by_name("embeddings:0")
            phase_train_placeholder = tf.get_default_graph().get_tensor_by_name("phase_train:0")
            #读取people中的数据
            dir="people/"
            compare_emb=[]
            compare_list=[]
            for i in os.listdir(dir):
                compare_list.append(i.split(".")[0])
                compare_emb.append(np.load(dir+i))
            compare_emb=np.array(compare_emb)
            print(compare_emb)
            cap = cv2.VideoCapture(camera)
            cap.set(3,160)
            while cap.isOpened():
                fin_obj=[]
                ok,img=cap.read()
                image=cv2.resize(img, (0, 0), fx=0.5, fy=0.5, interpolation=cv2.INTER_NEAREST)
                kk = cv2.waitKey(1)
                # 按下 q 键退出 / Press 'q' to quit
                if kk == ord('q'):
                    break
                start1=time.time()
                gray = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
                mark, bounding_box, crop_image = load_and_align_data(gray, 160, 44, pnet, rnet, onet)
                # if(count==len(bounding_box)):
                #     print(cv2.absdiff(image,frame))
                # else:
                #     count=len(bounding_box)
                #     frame=image
                end1=time.time()
                print("提取特征时间："+str(end1-start1))
                if mark:
                    start2 = time.time()
                    feed_dict = {images_placeholder: crop_image,phase_train_placeholder:False}
                    emb = sess.run(embeddings, feed_dict=feed_dict)
                    end2=time.time()
                    print("特征计算时间："+str(end2-start2))
                    if kk == ord('s'):
                        name=input("please input your name:")
                        np.save("people/"+name,emb)
                        compare_emb=np.insert(compare_emb,len(compare_emb),emb,axis=0)
                        print(len(compare_emb))
                        compare_list.append(name)
                        continue
                    fin_obj = []
                    temp_num=len(emb)
                    if (len(compare_emb) != 0):
                        for i in range(temp_num):
                            dist_list = []
                            for j in range(len(compare_emb)):
                                dist = np.sqrt(np.sum(np.square(np.subtract(emb[i, :], compare_emb[j, :]))))
                                print(str(dist))
                                dist_list.append(dist)
                            min_value = min(dist_list)
                            if (min_value > 1):
                                fin_obj.append('unknow')
                                print('unknown')
                            else:
                                fin_obj.append(compare_list[dist_list.index(min_value)])
                                print(compare_list[dist_list.index(min_value)])
                    else:
                        for i in range(temp_num):
                            fin_obj.append('unknow')
                            print('unknown')

                    for rec_position in range(temp_num):
                        cv2.rectangle(img, (bounding_box[rec_position, 0]*2, bounding_box[rec_position, 1]*2),
                                      (bounding_box[rec_position, 2]*2, bounding_box[rec_position, 3]*2), (0, 255, 0),
                                      2, 8, 0)

                        cv2.putText(
                            img,
                            fin_obj[rec_position],
                            (bounding_box[rec_position, 0]*2, bounding_box[rec_position, 1]*2),
                            cv2.FONT_HERSHEY_COMPLEX_SMALL,
                            0.8,
                            (0, 0, 255),
                            thickness=2,
                            lineType=2)

                cv2.putText(img, "Face Recognizer", (20, 40), font, 1, (255, 255, 255), 1, cv2.LINE_AA)
                cv2.putText(img, "FPS:   " + str(fps.__round__(2)), (20, 100), font, 0.8, (0, 255, 0), 1,
                            cv2.LINE_AA)
                cv2.putText(img, "Faces: " + str(len(fin_obj)), (20, 140), font, 0.8, (0, 255, 0), 1,
                            cv2.LINE_AA)
                cv2.putText(img,"Q: Quit", (20, 450), font, 0.8, (255, 255, 255), 1, cv2.LINE_AA)
                if(args.show):
                  cv2.imshow('camera', img)
                now = time.time()
                frame_time = now - frame_start_time
                fps = 1.0 / frame_time
                frame_start_time = now
def parse_arguments(argv):
    parser = argparse.ArgumentParser()
    parser.add_argument('--camera', type=str,
                        help='camera/ip camera',
                        default=0)
    parser.add_argument('--show', action='store_true',
                        help='display or not ',
                        default=False)
    return parser.parse_args(argv)

if __name__ == '__main__':
    main(parse_arguments(sys.argv[1:]),"model/facenet/facenet.pb")