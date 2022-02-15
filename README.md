
## Deploy all resources

1. Create EKS cluster

Ex.
```
$ eksctl create cluster --with-oidc --vpc-nat-mode=Disable
```

2. Deploy App Mesh CRD & Controller <br>
https://github.com/aws/eks-charts/tree/master/stable/appmesh-controller
   
You should set `tracing.enabled` and `tracing.provider` options if you want to use X-Ray. <br>
Ex.
```
helm upgrade -i appmesh-controller eks/appmesh-controller \
    --namespace appmesh-system \
    --set region=$AWS_REGION \
    --set serviceAccount.create=false \
    --set serviceAccount.name=appmesh-controller \
    --set tracing.enabled=true \
    --set tracing.provider=x-ray    
```

3. Deploy gRPC client/server as Pod

```
$ git clone https://github.com/a2ush/sample-grpc-server-with-appmesh-xray.git
$ cd sample-grpc-server-with-appmesh-xray

$ kubectl apply -f manifests/WithoutECRImage/
```