// Jenkins Pipeline for go-zero microservices
// 构建并部署到 K8s 集群
// 代码仓库: https://github.com/tutengdihuang/service_test.git

pipeline {
    agent any
    
    // GitHub Webhook 触发
    triggers {
        githubPush()
    }
    
    environment {
        // Harbor 配置
        HARBOR_URL = '182.42.82.135:30002'
        HARBOR_PROJECT = 'service-test'
        HARBOR_CREDENTIALS = 'harbor-credentials'
        
        // K8s 配置
        K8S_NAMESPACE = 'service-test'
        
        // Go 配置
        GOPROXY = 'https://goproxy.cn,direct'
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo "=== 拉取代码 ==="
                checkout scm
                sh 'git log -1 --oneline'
            }
        }
        
        stage('Build & Push Images') {
            steps {
                script {
                    def services = ['user', 'product', 'trade', 'web']
                    
                    // 登录 Harbor
                    withCredentials([usernamePassword(
                        credentialsId: env.HARBOR_CREDENTIALS,
                        usernameVariable: 'HARBOR_USER',
                        passwordVariable: 'HARBOR_PASS'
                    )]) {
                        sh "docker login -u ${HARBOR_USER} -p ${HARBOR_PASS} ${HARBOR_URL}"
                    }
                    
                    // 构建每个服务
                    services.each { service ->
                        def imageName = "${HARBOR_URL}/${HARBOR_PROJECT}/${service}-service"
                        def imageTag = "${env.BUILD_NUMBER}"
                        
                        echo "=== 构建 ${service} 服务 ==="
                        
                        sh """
                            docker build \
                                -f dockerfiles/Dockerfile.${service} \
                                -t ${imageName}:${imageTag} \
                                -t ${imageName}:latest \
                                .
                        """
                        
                        echo "=== 推送 ${service} 镜像到 Harbor ==="
                        sh """
                            docker push ${imageName}:${imageTag}
                            docker push ${imageName}:latest
                        """
                        
                        sh "docker rmi ${imageName}:${imageTag} || true"
                    }
                }
            }
        }
        
        stage('Deploy to K8s') {
            steps {
                script {
                    def services = ['user', 'product', 'trade', 'web']
                    
                    services.each { service ->
                        def imageName = "${HARBOR_URL}/${HARBOR_PROJECT}/${service}-service"
                        def imageTag = "${env.BUILD_NUMBER}"
                        
                        echo "=== 部署 ${service} 服务 ==="
                        
                        sh """
                            kubectl set image deployment/${service}-service \
                                ${service}=${imageName}:${imageTag} \
                                -n ${K8S_NAMESPACE} || echo "Deployment ${service}-service not found, skipping"
                        """
                    }
                }
            }
        }
        
        stage('Verify') {
            steps {
                echo "=== 验证部署状态 ==="
                sh "kubectl get pods -n ${K8S_NAMESPACE}"
            }
        }
    }
    
    post {
        success {
            echo '=== Pipeline 执行成功 ==='
        }
        failure {
            echo '=== Pipeline 执行失败 ==='
        }
        always {
            sh 'docker logout ${HARBOR_URL} || true'
        }
    }
}
