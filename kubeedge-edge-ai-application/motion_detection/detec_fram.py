# -!- coding: utf-8 -!-
import argparse
import itertools
import sys
import threading
import time
import numpy as np
import cv2


def parse_arguments(argv):
    parser = argparse.ArgumentParser()
    parser.add_argument('--camera', type=str,
                        help='camera/ip camera',
                        default=0)
    parser.add_argument('--show', action='store_true',
                        help='display or not ',
                        default=False)
    return parser.parse_args(argv)


def main(args):
    num = 0
    global flag
    global over
    global image
    over = True
    flag = False
    if args.camera == "0":
        camera = 0
    else:
        camera = args.camera
    camera = cv2.VideoCapture(camera)
    lastFrame1 = None
    lastFrame2 = None
    count = 0
    frame_start_time=0
    fps=0
    while camera.isOpened():
        (ret, frame) = camera.read()
        image=frame
        frame = cv2.resize(frame, (640, 480), interpolation=cv2.INTER_CUBIC)
        if lastFrame2 is None:
            if lastFrame1 is None:
                lastFrame1 = frame
            else:
                lastFrame2 = frame
                global frameDelta1
                frameDelta1 = cv2.absdiff(lastFrame1, lastFrame2)
            continue
        frameDelta2 = cv2.absdiff(lastFrame2, frame)  # 帧差�?
        thresh = cv2.bitwise_and(frameDelta1, frameDelta2)  # 图像与运�?
        thresh2 = thresh.copy()
        # 当前帧设为下一帧的前帧,前帧设为下一帧的前前�?帧差二设为帧差一
        lastFrame1 = lastFrame2
        lastFrame2 = frame.copy()
        frameDelta1 = frameDelta2
        # 结果转为灰度�?
        thresh = cv2.cvtColor(thresh, cv2.COLOR_BGR2GRAY)
        # 图像二值化
        thresh = cv2.threshold(thresh, 25, 255, cv2.THRESH_BINARY)[1]
        # 去除图像噪声,先腐蚀再膨胀(形态学开运算)
        thresh = cv2.dilate(thresh, None, iterations=3)
        thresh = cv2.erode(thresh, None, iterations=1)
        List = list(itertools.chain.from_iterable(thresh.tolist()))
        # 通过像素数量判断是否存在运动物体
        if (List.count(255) > 1000):
            if (count == 5):
                print("Moving object detected")
                flag = True
                if (over):
                    print("Start the video")
                    num += 1
                    threading.Thread(target=storeVideo, args=(num,)).start()
            else:
                count = count + 1
        else:
            print("Still picture")
            flag = False
            count = 0
        cv2.putText(frame, "FPS:   " + str(fps.__round__(2)), (20, 100),cv2.FONT_ITALIC, 0.8, (0, 255, 0), 1,
                    cv2.LINE_AA)
        isShow = args.show
        if (isShow):
            cv2.imshow("frame", frame)
            cv2.imshow("thresh", thresh)
            cv2.imshow("threst2", thresh2)
        # 如果q键被按下，跳出循�?
        if cv2.waitKey(200) & 0xFF == ord('q'):
            break
    # 清理资源并关闭打开的窗�?
        now = time.time()
        frame_time = now - frame_start_time
        fps = 1.0 / frame_time
        frame_start_time = now
    camera.release()
    cv2.destroyAllWindows()


def storeVideo(num):
    global image
    global over
    over = False
    global flag
    out_fps = 20.0  # 输出文件的帧�?
    fourcc = cv2.VideoWriter_fourcc(*'XVID')
    out1 = cv2.VideoWriter('./data/video/' + str(num) + ".avi", fourcc, out_fps, (640,480))
    start = time.time()
    count = 5
    image0=np.zeros(image.shape)
    while (True):
        if(image!=image0).all():
            out1.write(image)
            image=image0
        else:continue
        end = time.time()
        if(flag==False):
            if (end - start >= 10):
                if (flag):
                    continue
                    start = end
                    count = 0
                else:
                    if (count == 5):
                        over = True
                        print("End of the video")
                        break
                    else:
                        count += 1
                        continue
        else:
            start = end
    out1.release()


if __name__ == '__main__':
    main(parse_arguments(sys.argv[1:]))
