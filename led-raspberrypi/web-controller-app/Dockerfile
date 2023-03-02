# Based on ubuntu
FROM ubuntu:latest
LABEL maintainers="KubeEdge Community Developer"
LABEL description="KubeEdge Light Controller App"

# Copy from build directory
COPY kubeedge-light-controller /kubeedge-light-controller
COPY static /static
COPY views /views

# Define default command
ENTRYPOINT ["/kubeedge-light-controller"]

# Run the executable
CMD ["kubeedge-light-controller"]
