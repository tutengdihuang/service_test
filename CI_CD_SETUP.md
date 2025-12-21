# GitHub Actions CI/CD 配置说明

## ✅ 已完成的配置

所有 CI/CD 文件已经配置完成，编译和构建镜像都会在 GitHub Actions 上自动完成。

## 📁 项目结构

```
service_test/
├── .github/
│   └── workflows/
│       └── ci-cd.yml          # GitHub Actions 工作流
├── dockerfiles/               # Dockerfile 文件
│   ├── Dockerfile.web
│   ├── Dockerfile.user
│   ├── Dockerfile.product
│   └── Dockerfile.trade
├── k8s/                       # Kubernetes 部署文件
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── web/deployment.yaml
│   ├── user/deployment.yaml
│   ├── product/deployment.yaml
│   └── trade/deployment.yaml
└── CI_CD_SETUP.md            # 本文件
```

## 🚀 使用步骤

### 1. 配置 GitHub Secrets（必需）

在 GitHub 仓库中配置 Kubernetes 访问凭证：

1. 打开 GitHub 仓库页面
2. 点击 **Settings** → **Secrets and variables** → **Actions**
3. 点击 **New repository secret**
4. 添加以下 Secret：

#### KUBECONFIG（必需）
- **Name**: `KUBECONFIG`
- **Value**: Kubernetes 配置文件的 base64 编码
  ```bash
  # 获取并编码
  cat ~/.kube/config | base64 -w 0
  
  # macOS 如果没有 -w 参数
  cat ~/.kube/config | base64 | tr -d '\n'
  ```

### 2. 提交代码到 GitHub

```bash
cd /Volumes/mac_data/code/go_code/service_test

# 添加所有文件
git add .

# 提交
git commit -m "Add GitHub Actions CI/CD configuration"

# 推送到 GitHub（确保推送到 main 或 master 分支）
git push origin main
# 或
git push origin master
```

### 3. 触发 CI/CD 流程

#### 自动触发
当代码推送到以下分支时会自动触发：
- `main`
- `master`
- `develop`

**触发条件**：修改了以下路径的文件
- `api/**`
- `rpc/**`
- `go.mod`
- `go.sum`
- `.github/workflows/**`

#### 手动触发
1. 打开 GitHub 仓库
2. 点击 **Actions** 标签
3. 左侧选择 **Service Test CI/CD**
4. 点击 **Run workflow**
5. 选择分支和服务
6. 点击 **Run workflow**

## 📊 CI/CD 流程说明

### Job 1: test（测试阶段）
- ✅ 拉取代码
- ✅ 安装 Go 1.24
- ✅ 缓存 Go modules（加速）
- ✅ 下载依赖
- ✅ 运行测试
- ✅ 编译验证（验证 4 个服务能否正常编译）

### Job 2: build（构建阶段）
- ✅ 并行构建 4 个服务的 Docker 镜像
  - `web-service` → `ghcr.io/{owner}/service-test/web-service:{tag}`
  - `user-service` → `ghcr.io/{owner}/service-test/user-service:{tag}`
  - `product-service` → `ghcr.io/{owner}/service-test/product-service:{tag}`
  - `trade-service` → `ghcr.io/{owner}/service-test/trade-service:{tag}`
- ✅ 推送到 GitHub Container Registry (ghcr.io)
- ✅ 使用缓存加速构建

**镜像标签策略**：
- `latest` - 仅 main/master 分支
- `{分支名}` - 如 `main`、`develop`
- `{分支名}-{SHA前8位}` - 如 `main-abc12345`
- `pr-{PR号}` - Pull Request 时

### Job 3: deploy（部署阶段，仅 main/master）
- ✅ 安装 kubectl
- ✅ 配置 Kubernetes 连接
- ✅ 验证集群连接
- ✅ 部署到 Kubernetes
  - 创建 namespace: `service-test`
  - 创建 ConfigMap
  - 部署 4 个服务（使用 commit SHA 作为镜像标签）
- ✅ 健康检查

## 🔍 查看执行结果

### 在 GitHub 上查看
1. 打开 GitHub 仓库
2. 点击 **Actions** 标签
3. 点击对应的 workflow run
4. 查看各个 Job 的执行日志

### 查看构建的镜像
访问：`https://github.com/{你的用户名}?tab=packages`

### 在 Kubernetes 集群中验证
```bash
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

# 测试 web 服务（通过 NodePort）
kubectl get svc web-service -n service-test
# 然后访问 http://<node-ip>:30888/api/user/1
```

## ⚙️ 配置说明

### 镜像仓库
- **默认**: GitHub Container Registry (`ghcr.io`)
- **自动认证**: 使用 `GITHUB_TOKEN`，无需额外配置
- **镜像格式**: `ghcr.io/{repository_owner}/service-test/{service}-service:{tag}`

### Kubernetes 命名空间
- **默认**: `service-test`
- 可在 `ci-cd.yml` 中修改 `KUBERNETES_NAMESPACE` 环境变量

### Go 版本
- **当前**: Go 1.24
- 可在 `ci-cd.yml` 中修改 `GO_VERSION` 环境变量

## 🐛 故障排查

### 工作流没有触发
- ✅ 检查文件路径是否匹配 `paths` 配置
- ✅ 检查分支名是否为 `main`、`master` 或 `develop`
- ✅ 检查 `.github/workflows/ci-cd.yml` 文件是否存在

### 镜像构建失败
- ✅ 检查 `dockerfiles/` 目录是否存在
- ✅ 检查 Dockerfile 路径是否正确
- ✅ 查看构建日志中的错误信息

### 部署失败
- ✅ 检查 `KUBECONFIG` Secret 是否正确配置
- ✅ 检查 Kubernetes 集群连接是否正常
- ✅ 检查 `k8s/` 目录下的部署文件是否存在
- ✅ 检查镜像是否成功推送到仓库

### 镜像拉取失败
- ✅ 检查镜像仓库权限
- ✅ 检查镜像标签是否正确
- ✅ 检查 Kubernetes 集群是否能访问 ghcr.io

## 📝 注意事项

1. **首次使用**：确保已配置 `KUBECONFIG` Secret
2. **镜像权限**：GitHub Container Registry 默认是私有的，确保 Kubernetes 集群有拉取权限
3. **etcd 服务**：确保 Kubernetes 集群中有 etcd 服务，或修改 ConfigMap 中的 etcd 地址
4. **资源限制**：根据实际需求调整 `k8s/*/deployment.yaml` 中的资源限制
5. **健康检查**：web 服务使用 HTTP 健康检查，RPC 服务使用 TCP 健康检查

## 🎉 完成！

现在你的项目已经完全配置好 GitHub Actions CI/CD，每次推送代码到 main/master 分支时，会自动：
1. 测试代码
2. 构建 Docker 镜像
3. 推送到镜像仓库
4. 部署到 Kubernetes

祝你使用愉快！

