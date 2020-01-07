## My Kubernetes Custom Controller

> For now, custom controller for monitoring changes of kubernetes workload.

### Usage

```
Usage:
  k8s-custom-controller [flags]

Flags:
      --enable-deployment          enable watching deployment
      --enable-job                 enable watching job
      --enable-persistent-volume   enable watching persistent volume
      --enable-pod                 enable watching pod
      --enable-service             enable watching service
  -h, --help                       help for k8s-custom-controller
      --version                    version for k8s-custom-controller
```

### Testing

**TODO**

### Deploy

build application, create and push docker image to registry.

```
$ make push
```

