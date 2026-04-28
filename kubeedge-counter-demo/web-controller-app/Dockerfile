#Based on ubuntu
FROM ubuntu:22.04
LABEL maintainers="KubeEdge Community Developer"
LABEL description="KubeEdge Counter Web Controller App"

# Copy from build directory
COPY kubeedge-counter-controller /kubeedge-counter-controller
COPY static /static
COPY views /views

# Make sure the binary is executable
RUN chmod +x /kubeedge-counter-controller

# Define default command
ENTRYPOINT ["/kubeedge-counter-controller"]

# Run the executable
CMD ["kubeedge-counter-controller"]
