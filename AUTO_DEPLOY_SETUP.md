# 全自动部署到 Kubernetes 集群配置指南

## 🎯 目标

配置完成后，每次推送代码到 `main` 或 `master` 分支时，GitHub Actions 会自动：
1. ✅ 编译代码
2. ✅ 构建 Docker 镜像
3. ✅ 推送到镜像仓库
4. ✅ **自动部署到你的 Kubernetes 集群**

## 📋 前置条件

- ✅ 已准备好 Kubernetes 集群配置文件：`temp_admin.conf`
- ✅ 集群地址：`https://182.42.82.135:6443`
- ✅ CI/CD 文件已配置完成

## 🔧 配置步骤

### 步骤 1: 获取 Base64 编码的 Kubeconfig

已为你生成 base64 编码的 kubeconfig，文件保存在：
```
KUBECONFIG_BASE64.txt
```

**或者手动生成**：
```bash
cat /Volumes/mac_data/code/go_code/101/Allen/01.k8s_config/temp_admin.conf | base64 | tr -d '\n'
```

### 步骤 2: 配置 GitHub Secrets

1. **打开 GitHub 仓库**
   - 访问你的 GitHub 仓库页面

2. **进入 Secrets 配置**
   - 点击 **Settings**（设置）
   - 左侧菜单选择 **Secrets and variables** → **Actions**
   - 点击 **New repository secret**

3. **添加 KUBECONFIG Secret**
   - **Name**: `KUBECONFIG`
   - **Value**: 复制 `KUBECONFIG_BASE64.txt` 文件中的全部内容
     ```bash
     # 查看内容
     cat KUBECONFIG_BASE64.txt
     ```
   - 点击 **Add secret**

### 步骤 3: 验证配置

配置完成后，你的 GitHub Secrets 应该包含：
- ✅ `KUBECONFIG` - Kubernetes 集群访问凭证（base64 编码）

**注意**：使用 GitHub Container Registry (`ghcr.io`) 时，`GITHUB_TOKEN` 会自动提供，无需额外配置。

## 🚀 触发自动部署

### 方式 1: 推送代码（推荐）

```bash
cd /Volumes/mac_data/code/go_code/service_test

# 确保所有文件已提交
git add .
git commit -m "Configure auto deployment to Kubernetes"

# 推送到 main 或 master 分支
git push origin main
# 或
git push origin master
```

### 方式 2: 手动触发

1. 打开 GitHub 仓库 → **Actions** 标签
2. 左侧选择 **Service Test CI/CD**
3. 点击 **Run workflow**
4. 选择分支和服务
5. 点击 **Run workflow**

## 📊 部署流程

当代码推送到 `main`/`master` 分支时，会自动执行：

```
代码 Push
    ↓
┌─────────────────┐
│  Job 1: test    │  ← 测试和编译验证（约 2-3 分钟）
└────────┬────────┘
         ↓
┌─────────────────┐
│  Job 2: build   │  ← 并行构建 4 个 Docker 镜像（约 5-10 分钟）
│  - web          │     - web-service → ghcr.io/{owner}/service-test/web-service:latest
│  - user         │     - user-service → ghcr.io/{owner}/service-test/user-service:latest
│  - product      │     - product-service → ghcr.io/{owner}/service-test/product-service:latest
│  - trade        │     - trade-service → ghcr.io/{owner}/service-test/trade-service:latest
└────────┬────────┘
         ↓
┌─────────────────┐
│  Job 3: deploy  │  ← 自动部署到 Kubernetes（约 2-5 分钟）
│  - 连接集群     │     - 集群：https://182.42.82.135:6443
│  - 创建资源     │     - Namespace: service-test
│  - 部署服务     │     - 部署 4 个服务
│  - 健康检查     │     - 验证部署成功
└─────────────────┘
         ↓
   部署完成 ✅
```

## 🔍 查看部署状态

### 在 GitHub Actions 中查看

1. 打开 GitHub 仓库 → **Actions** 标签
2. 点击最新的 workflow run
3. 查看各个 Job 的执行状态：
   - ✅ **test** - 绿色表示测试通过
   - ✅ **build** - 绿色表示镜像构建成功
   - ✅ **deploy** - 绿色表示部署成功

### 在 Kubernetes 集群中验证

```bash
# 使用你的 kubeconfig
export KUBECONFIG=/Volumes/mac_data/code/go_code/101/Allen/01.k8s_config/temp_admin.conf

# 查看命名空间
kubectl get namespace service-test

# 查看所有 Pod
kubectl get pods -n service-test

# 查看所有 Service
kubectl get svc -n service-test

# 查看部署状态
kubectl get deployments -n service-test

# 查看 Pod 日志
kubectl logs -f <pod-name> -n service-test

# 查看所有资源
kubectl get all -n service-test
```

### 测试服务

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

## 📝 部署配置说明

### 镜像仓库

- **仓库地址**: `ghcr.io` (GitHub Container Registry)
- **镜像格式**: `ghcr.io/{你的GitHub用户名}/service-test/{服务名}-service:{标签}`
- **认证**: 自动使用 `GITHUB_TOKEN`，无需额外配置

### Kubernetes 配置

- **集群地址**: `https://182.42.82.135:6443`
- **命名空间**: `service-test`
- **服务端口**:
  - web: 8888 (HTTP)
  - user: 9001 (gRPC)
  - product: 9002 (gRPC)
  - trade: 9003 (gRPC)

### 部署顺序

自动部署会按以下顺序执行：
1. 创建 namespace
2. 创建 ConfigMap
3. 部署 user 服务
4. 部署 product 服务
5. 部署 trade 服务
6. 部署 web 服务

## 🐛 故障排查

### 1. 部署 Job 失败

**检查 KUBECONFIG Secret**：
- 确保 `KUBECONFIG` Secret 已正确配置
- 确保 base64 编码完整（无换行符）

**检查集群连接**：
```bash
# 本地测试集群连接
export KUBECONFIG=/Volumes/mac_data/code/go_code/101/Allen/01.k8s_config/temp_admin.conf
kubectl cluster-info
kubectl get nodes
```

### 2. 镜像拉取失败

**检查镜像权限**：
- GitHub Container Registry 默认是私有的
- 确保 Kubernetes 集群可以访问 `ghcr.io`
- 如果使用私有仓库，需要配置 `imagePullSecrets`

**检查镜像标签**：
- 在 GitHub Actions 日志中查看构建的镜像标签
- 确保 deployment.yaml 中的镜像标签正确

### 3. Pod 启动失败

**查看 Pod 日志**：
```bash
kubectl logs <pod-name> -n service-test
kubectl describe pod <pod-name> -n service-test
```

**常见原因**：
- 镜像拉取失败
- 配置错误（如 etcd 地址）
- 资源不足

### 4. 服务无法访问

**检查 Service 和 Endpoints**：
```bash
kubectl get svc -n service-test
kubectl get endpoints -n service-test
```

**检查 Pod 标签**：
```bash
kubectl get pods -n service-test --show-labels
```

## 🔄 更新部署

### 自动更新

每次推送代码到 `main`/`master` 分支时，会自动：
- 构建新镜像（使用 commit SHA 作为标签）
- 更新部署（使用新镜像）

### 回滚部署

如果部署出现问题，可以回滚：

```bash
# 查看部署历史
kubectl rollout history deployment/web-service -n service-test

# 回滚到上一个版本
kubectl rollout undo deployment/web-service -n service-test

# 回滚所有服务
for service in user product trade web; do
  kubectl rollout undo deployment/${service}-service -n service-test
done
```

## ✅ 验证清单

部署前检查：
- [ ] GitHub Secrets 中已配置 `KUBECONFIG`
- [ ] Kubeconfig 文件可以正常连接集群
- [ ] 代码已推送到 GitHub
- [ ] GitHub Actions 工作流已触发

部署后检查：
- [ ] GitHub Actions 中所有 Job 显示成功（绿色）
- [ ] Kubernetes 集群中可以看到 Pod 运行
- [ ] 服务可以正常访问
- [ ] 日志中没有错误信息

## 🎉 完成！

配置完成后，你的项目已经实现**全自动 CI/CD 部署**：

1. ✅ **代码推送** → 自动触发
2. ✅ **编译构建** → 在 GitHub Actions 上完成
3. ✅ **镜像推送** → 自动推送到 ghcr.io
4. ✅ **自动部署** → 自动部署到你的 Kubernetes 集群

**下次只需要**：
```bash
git push origin main
```

然后等待 GitHub Actions 自动完成所有工作！🚀

## 📞 需要帮助？

如果遇到问题：
1. 查看 GitHub Actions 日志
2. 检查 Kubernetes 集群状态
3. 参考 `DEPLOYMENT_GUIDE.md` 获取详细故障排查指南

