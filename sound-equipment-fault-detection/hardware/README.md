# hardware 🎇🎇🎇
The pymodbus library has been updated significantly, please use the specified version pymodbus==2.5.3

## 1 Enter the directory
```shell
cd sound-equipment-fault-detection/hardware
```

## 2 Build images and Send to the edge
```shell
# Build the image
docker build -t hardware .
# Exporting Docker images
docker save -o hardware_image.tar hardware
# Transferring the image file to edge machine
rsync -avz --progress hardware_image.tar root@192.168.1.201:/root/hardware/
rm hardware_image.tar
```

## 3 Edge loading images
```shell
# Importing Docker images
docker load -i /root/hardware/hardware_image.tar
rm /root/hardware/hardware_image.tar
# Confirm that the image was loaded successfully
docker images
```

## 4 Deploy hardware ✅
```shell
# Deployment
kubectl apply -f resource/config.yaml
kubectl apply -f resource/deployment.yaml
```


## Extensions

If your hardware uses Modbus RTU transmission equipment, please replace the Modbus TCP client in the code with the Modbus RTU client, and keep the rest unchanged.

```shell
# If you don't have an RTU device, you can use virtual RTU communication, but the speed is very slow.
sudo apt-get install socat
# Create a virtual serial port pair ttyVIRT0 and ttyVIRT1
socat -d -d PTY,link=/dev/ttyVIRT0,raw,echo=0 PTY,link=/dev/ttyVIRT1,raw,echo=0
stty -F /dev/ttyVIRT0 115200 raw
stty -F /dev/ttyVIRT1 115200 raw
stty -F /dev/ttyVIRT0
stty -F /dev/ttyVIRT1
```
