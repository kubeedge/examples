FROM python:3.6-slim-stretch

RUN  sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list
RUN  sed -i s@/security.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list
RUN  apt-get clean
RUN apt-get -y update && \
    apt-get install -y --fix-missing \
    libglib2.0-dev\
    python-qt4\
    && apt-get clean && rm -rf /tmp/* /var/tmp/*

WORKDIR /code

COPY . /motion_detection

WORKDIR /motion_detection

RUN pip  --default-timeout=10000  install -r requirements.txt -i https://pypi.doubanio.com/simple/

ENTRYPOINT ["python","-u", "/motion_detection/detec_fram.py"]
