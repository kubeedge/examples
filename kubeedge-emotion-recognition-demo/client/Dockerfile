From kuramal/onnxruntime_opencv:arm32

RUN pip install requests

ADD ./emotion.py /home

ADD ./haarcascade_frontalface_default.xml /home

ENTRYPOINT [ "python", "/home/emotion.py" ]
