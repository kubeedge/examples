FROM sixsq/opencv-python:master-arm

WORKDIR /code
COPY . /face-recongition
WORKDIR /face-recongition

RUN pip  --default-timeout=10000  install -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple


RUN pip install https://www.piwheels.org/simple/scikit-learn/scikit_learn-0.21.3-cp35-cp35m-linux_armv7l.whl
ENTRYPOINT ["python3","-u","/face-recongition/face_recong.py"]
