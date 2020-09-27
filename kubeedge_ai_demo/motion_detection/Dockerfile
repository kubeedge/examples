FROM sixsq/opencv-python:master-arm
WORKDIR /code
COPY . /motion_detection
WORKDIR /motion_detection
ENTRYPOINT ["python3","-u", "/motion_detection/detec_fram.py"]