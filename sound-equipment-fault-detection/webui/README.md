# webui 🎇🎇🎇

## 1 Enter the directory
```shell
cd sound-equipment-fault-detection/webui
```

## 2 Build the image
```shell
docker build -f Dockerfile -t webui .
docker images
```

## 3 Deploy webui ✅
```shell 
kubectl apply -f resource/deployment.yaml
kubectl apply -f resource/role-binding.yaml
```