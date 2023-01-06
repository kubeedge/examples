FROM ubuntu:latest

CMD mkdir -p light_mapper

COPY . light_mapper/

WORKDIR light_mapper

ENTRYPOINT ["/light_mapper/light_mapper","-logtostderr=true"]
