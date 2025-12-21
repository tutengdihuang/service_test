# 部署指南

构建完成后，有两种方式可以部署到 Kubernetes：

## 方式一：自动部署（推荐）🚀

### 自动部署流程

当代码推送到 `main` 或 `master` 分支时，GitHub Actions 会自动完成部署：

1. **构建镜像** → 推送到 `ghcr.io`
2. **自动部署** → 部署到 Kubernetes 集群

### 触发自动部署

```bash
# 推送到 main 分支
git push origin main

# 或推送到 master 分支
git push origin master
```

### 查看部署状态

1. **在 GitHub Actions 中查看**
   - 打开 GitHub 仓库 → **Actions** 标签
   - 查看最新的 workflow run
   - 点击 **deploy** job 查看部署日志

2. **在 Kubernetes 集群中验证**
   ```bash
   # 查看所有 Pod
   kubectl get pods -n service-test
   
   # 查看部署状态
   kubectl get deployments -n service-test
   
   # 查看服务
   kubectl get svc -n service-test
   ```

---

## 方式二：手动部署 🔧

如果构建完成但需要手动部署，可以使用以下方法：

### 步骤 1: 获取镜像标签

首先需要知道构建的镜像标签。可以通过以下方式获取：

**方法 A：从 GitHub Actions 日志中获取**
- 打开 GitHub Actions → 查看 build job 的日志
- 找到镜像推送信息，例如：`ghcr.io/your-org/service-test/web-service:main-abc12345`

**方法 B：从 GitHub Packages 中查看**
- 访问：`https://github.com/{你的用户名}?tab=packages`
- 查看 `service-test` 下的镜像标签

**方法 C：使用 commit SHA**
```bash
# 获取最新的 commit SHA（前8位）
git rev-parse --short HEAD
# 输出：abc12345
```

### 步骤 2: 更新部署文件中的镜像标签

```bash
cd /Volumes/mac_data/code/go_code/service_test

# 设置镜像标签（替换为实际的标签）
IMAGE_TAG="main-abc12345"  # 或使用 latest
REGISTRY="ghcr.io"
IMAGE_PREFIX="your-org/service-test"  # 替换为你的实际路径

# 更新各个服务的镜像标签
for service in user product trade web; do
  IMAGE="${REGISTRY}/${IMAGE_PREFIX}/${service}-service:${IMAGE_TAG}"
  
  # macOS
  sed -i '' "s|image:.*${service}-service.*|image: ${IMAGE}|g" k8s/${service}/deployment.yaml
  
  # Linux
  # sed -i "s|image:.*${service}-service.*|image: ${IMAGE}|g" k8s/${service}/deployment.yaml
done
```

### 步骤 3: 部署到 Kubernetes

```bash
# 1. 创建命名空间
kubectl apply -f k8s/namespace.yaml

# 2. 创建 ConfigMap
kubectl apply -f k8s/configmap.yaml

# 3. 部署各个服务（按依赖顺序）
kubectl apply -f k8s/user/deployment.yaml
kubectl apply -f k8s/product/deployment.yaml
kubectl apply -f k8s/trade/deployment.yaml
kubectl apply -f k8s/web/deployment.yaml

# 4. 等待部署完成
kubectl rollout status deployment/user-service -n service-test --timeout=5m
kubectl rollout status deployment/product-service -n service-test --timeout=5m
kubectl rollout status deployment/trade-service -n service-test --timeout=5m
kubectl rollout status deployment/web-service -n service-test --timeout=5m
```

---

## 部署脚本（一键部署）

创建一个部署脚本，方便手动部署：

```bash
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
```

**使用方法**：

```bash
# 使用 latest 标签
./deploy.sh

# 使用特定标签
IMAGE_TAG="main-abc12345" ./deploy.sh

# 使用自定义镜像仓库
REGISTRY="registry.example.com" IMAGE_PREFIX="my-org/service-test" IMAGE_TAG="v1.0.0" ./deploy.sh
```

---

## 验证部署

### 1. 检查 Pod 状态

```bash
# 查看所有 Pod
kubectl get pods -n service-test

# 查看 Pod 详细信息
kubectl get pods -n service-test -o wide

# 查看 Pod 日志
kubectl logs -f <pod-name> -n service-test

# 查看所有 Pod 的日志
kubectl logs -f -l app=web-service -n service-test
```

### 2. 检查服务状态

```bash
# 查看所有 Service
kubectl get svc -n service-test

# 查看 Service 详细信息
kubectl describe svc web-service -n service-test

# 查看 Endpoints（确保 Pod 已关联）
kubectl get endpoints -n service-test
```

### 3. 检查部署状态

```bash
# 查看所有 Deployment
kubectl get deployments -n service-test

# 查看 Deployment 详细信息
kubectl describe deployment web-service -n service-test

# 查看部署历史
kubectl rollout history deployment/web-service -n service-test
```

### 4. 测试服务

```bash
# 获取 web 服务的 NodePort
kubectl get svc web-service -n service-test

# 测试 API（替换 <node-ip> 为实际节点 IP，30888 为 NodePort）
curl http://<node-ip>:30888/api/user/1
curl http://<node-ip>:30888/api/product/1

# 或使用 port-forward（本地测试）
kubectl port-forward svc/web-service 8888:8888 -n service-test
# 然后访问 http://localhost:8888/api/user/1
```

### 5. 健康检查

```bash
# 在集群内部测试（使用临时 Pod）
kubectl run curl-test --image=curlimages/curl:latest --rm -i --restart=Never -n service-test -- \
  curl -f http://web-service:8888/health || echo "Health check failed"

# 测试各个服务
for service in user product trade web; do
  echo "Testing $service service..."
  kubectl run curl-test-$service --image=curlimages/curl:latest --rm -i --restart=Never -n service-test -- \
    curl -f http://${service}-service:$(kubectl get svc ${service}-service -n service-test -o jsonpath='{.spec.ports[0].port}') || echo "$service service check failed"
done
```

---

## 回滚部署

如果部署出现问题，可以回滚到之前的版本：

```bash
# 查看部署历史
kubectl rollout history deployment/web-service -n service-test

# 回滚到上一个版本
kubectl rollout undo deployment/web-service -n service-test

# 回滚到特定版本
kubectl rollout undo deployment/web-service --to-revision=2 -n service-test

# 回滚所有服务
for service in user product trade web; do
  kubectl rollout undo deployment/${service}-service -n service-test
done
```

---

## 更新部署

### 更新镜像版本

```bash
# 方法 1: 使用 kubectl set image
kubectl set image deployment/web-service web=ghcr.io/your-org/service-test/web-service:new-tag -n service-test

# 方法 2: 修改 deployment.yaml 后重新应用
sed -i '' 's|image:.*web-service.*|image: ghcr.io/your-org/service-test/web-service:new-tag|g' k8s/web/deployment.yaml
kubectl apply -f k8s/web/deployment.yaml

# 方法 3: 使用部署脚本
IMAGE_TAG="new-tag" ./deploy.sh
```

### 扩缩容

```bash
# 扩容到 3 个副本
kubectl scale deployment/web-service --replicas=3 -n service-test

# 缩容到 1 个副本
kubectl scale deployment/web-service --replicas=1 -n service-test
```

---

## 常见问题

### 1. Pod 一直处于 Pending 状态

```bash
# 查看 Pod 事件
kubectl describe pod <pod-name> -n service-test

# 常见原因：
# - 资源不足
# - NodeSelector 不匹配
# - 没有可用的节点
```

### 2. Pod 一直处于 CrashLoopBackOff 状态

```bash
# 查看 Pod 日志
kubectl logs <pod-name> -n service-test

# 查看 Pod 事件
kubectl describe pod <pod-name> -n service-test

# 常见原因：
# - 应用启动失败
# - 配置错误
# - 依赖服务不可用
```

### 3. 镜像拉取失败

```bash
# 检查镜像是否存在
docker pull ghcr.io/your-org/service-test/web-service:tag

# 检查镜像仓库权限
kubectl describe pod <pod-name> -n service-test | grep -i image

# 如果使用私有仓库，需要配置 imagePullSecrets
```

### 4. 服务无法访问

```bash
# 检查 Service 配置
kubectl get svc web-service -n service-test -o yaml

# 检查 Endpoints
kubectl get endpoints web-service -n service-test

# 检查 Pod 标签是否匹配
kubectl get pods -n service-test --show-labels
kubectl get svc web-service -n service-test -o jsonpath='{.spec.selector}'
```

---

## 最佳实践

1. **部署顺序**：先部署基础服务（user, product），再部署依赖服务（trade），最后部署网关（web）
2. **健康检查**：确保配置了正确的 livenessProbe 和 readinessProbe
3. **资源限制**：根据实际需求设置合理的资源请求和限制
4. **版本管理**：使用有意义的镜像标签（如 commit SHA），避免使用 latest
5. **监控告警**：配置监控和告警，及时发现问题
6. **备份配置**：定期备份 Kubernetes 配置和部署文件

---

## 快速参考

```bash
# 一键查看所有资源状态
kubectl get all -n service-test

# 查看所有资源
kubectl get pods,svc,deployments,ingress -n service-test

# 删除所有部署（谨慎使用）
kubectl delete namespace service-test

# 重新部署
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/user/deployment.yaml
kubectl apply -f k8s/product/deployment.yaml
kubectl apply -f k8s/trade/deployment.yaml
kubectl apply -f k8s/web/deployment.yaml
```

祝你部署顺利！🎉

