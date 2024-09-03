# vled 🎇🎇🎇

```markdown
*** MODIFIED FILE ***
----------------------------
config.yaml
resource/configmap.yaml
resource/vled-model.yaml
resource/vled-instance.yaml
resource/deployment.yaml

driver/devicetype.go
driver/driver.go
————————————————————————————
```

## 1 Enter the directory
```shell
cd sound-equipment-fault-detection/vled
```

## 2 Build images and Send to the edge
```shell
# Build the image
docker build -f Dockerfile_nostream -t vled .
# Exporting Docker images
docker save -o vled_image.tar vled
# Transferring the image file to edge machine
rsync -avz --progress vled_image.tar root@192.168.1.201:/root/vled/
rm vled_image.tar
```

## 3 Edge loading images
```shell
# Importing Docker images
docker load -i /root/vled/vled_image.tar
rm /root/vled/vled_image.tar
# Confirm that the image was loaded successfully
docker images
```

## 4 Deploy vled ✅
```shell 
kubectl apply -f resource/vled-model.yaml
kubectl apply -f resource/vled-instance.yaml
kubectl apply -f resource/deployment.yaml
```

## Debug
```shell
docker run -it --name vled-container --network host -v /etc/kubeedge:/etc/kubeedge vled /bin/bash
```