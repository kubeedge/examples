# Kubeedge-traffic-light

- build at edge node:

```bash
$ make
```

- create crds at cloud node:

```bash
$ cd crd
$ kubectl apply -f model.yaml
# replace "<your edge node name>" with your edge node name
$ sed -i 's#raspberrypi#<your edge node name>#' instance.yaml
$ kubectl apply -f instance.yaml
```

**Note: instance must be created after model and deleted before model.**

- create demo at cloud node:

```bash
$ kubectl apply -f deploy.yaml
```
