# Based on centos
FROM centos:7.6.1810
LABEL maintainers="KubeEdge Authors"
LABEL description="KubeEdge WeChat App"

# Copy from build directory
COPY kubeedge-wechat-app /kubeedge-wechat-app

# Update
RUN yum -y update

# Define default command
ENTRYPOINT ["/kubeedge-wechat-app"]

# Run the executable
CMD ["kubeedge-wechat-app"]
