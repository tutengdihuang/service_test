#!/bin/bash
# deploy.sh - 手动部署脚本

set -e

# 配置
NAMESPACE="service-test"
REGISTRY="${REGISTRY:-ghcr.io}"
IMAGE_PREFIX="${IMAGE_PREFIX:-your-org/service-test}"  # 替换为你的实际路径
IMAGE_TAG="${IMAGE_TAG:-latest}"  # 或使用 commit SHA，如 main-abc12345

echo "=========================================="
echo "Deploying services to Kubernetes"
echo "=========================================="
echo "Namespace: $NAMESPACE"
echo "Registry: $REGISTRY"
echo "Image Prefix: $IMAGE_PREFIX"
echo "Image Tag: $IMAGE_TAG"
echo "=========================================="

# 创建命名空间
echo "Creating namespace..."
kubectl apply -f k8s/namespace.yaml

# 创建 ConfigMap
echo "Creating ConfigMap..."
kubectl apply -f k8s/configmap.yaml

# 更新并部署各个服务
for service in user product trade web; do
  echo ""
  echo "Deploying $service service..."
  
  IMAGE="${REGISTRY}/${IMAGE_PREFIX}/${service}-service:${IMAGE_TAG}"
  
  # 更新镜像标签
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s|image:.*${service}-service.*|image: ${IMAGE}|g" k8s/${service}/deployment.yaml
  else
    sed -i "s|image:.*${service}-service.*|image: ${IMAGE}|g" k8s/${service}/deployment.yaml
  fi
  
  # 应用部署文件
  kubectl apply -f k8s/${service}/deployment.yaml
  
  # 等待部署完成
  echo "Waiting for $service service to be ready..."
  kubectl rollout status deployment/${service}-service -n ${NAMESPACE} --timeout=5m
done

echo ""
echo "=========================================="
echo "Deployment completed!"
echo "=========================================="
echo ""
echo "Checking deployment status..."
kubectl get pods -n ${NAMESPACE}
kubectl get svc -n ${NAMESPACE}

