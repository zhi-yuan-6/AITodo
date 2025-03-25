# AITodo 后端服务

提供基于JWT的用户认证系统、阿里云短信服务集成以及待办事项管理的RESTful API服务。

## 📋 先决条件

### 环境要求
- **Go 1.21+** ([安装指南](https://go.dev/doc/install))
- **MySQL 8.0+** ([官方文档](https://dev.mysql.com/doc/))
- **Redis 7.0+** ([快速入门](https://redis.io/docs/install))
- **阿里云账号**（需开通[SMS服务](https://www.aliyun.com/product/sms)）

### 密钥准备
- 阿里云AccessKey ([获取指南](https://help.aliyun.com/zh/ram/user-guide/create-an-accesskey-pair))
- DASHSCOPE_API_KEY ([灵积平台](https://help.aliyun.com/zh/dashscope/developer-reference/activate-dashscope-and-create-an-api-key))

---

## 🚀 快速启动

### 1. 服务配置
```bash
# 复制并重命名配置文件
cp config/config.yaml.example config/config.yaml
cp .env_example .env
```

### 2. 生成RSA私钥
```bash
# 方式一：使用项目工具生成
复制util/pem.go中的GeneratePEM()函数，编写一个main函数调用生成

# 方式二：OpenSSL生成（推荐）
openssl genrsa -out private_key.pem 2048
```

### 3. 数据库初始化
```sql
CREATE DATABASE aitodo CHARACTER SET utf8mb4;
-- 导入SQL文件（项目根目录/schema.sql）
```

---

## 🐳 容器化部署

### 独立运行
```bash
# 后端服务
docker build aitodobackend:latest
docker run -d -p 8080:80 zhhiyuan/aitodobackend

# 前端服务
docker pull zhhiyuan/aitodofrontend:latest
docker run -d -p 5173:80 zhhiyuan/aitodofrontend
```

### Docker Compose 启动指南

按照以下步骤，您可以顺利地拉取前端镜像、构建后端镜像并启动服务：

1. **拉取前端镜像**

   首先，从 Docker Hub 拉取最新的前端镜像：

   ```bash
   docker pull zhhiyuan/aitodofrontend:latest
   ```

2. **构建后端镜像**

   接着，在本地构建后端镜像（因后端镜像依赖配置文件，需在本地构建）：

   ```bash
   docker build -t aitodobackend:latest .
   ```

3. **启动服务**

   最后，使用 Docker Compose 启动服务：

   ```bash
   docker-compose up -d
   ```

   该命令会根据 `docker-compose.yml` 文件配置，启动并运行容器，所有服务将以分离模式在后台运行，端口映射、网络连接等设置均自动完成。

启动完成后，您可以通过浏览器访问 `http://localhost:5173` 查看应用。

## 🛠️ 开发模式

在开发模式下启动后端服务，需要按照以下步骤进行配置和运行：

1. **修改配置文件**

   首先，打开 `config/config.yaml` 文件，根据您的环境修改 MySQL 和 Redis 的连接地址，确保它们指向正确的数据库和缓存服务实例。这一步是必要的，因为后端服务依赖这些配置来建立数据存储的连接。

2. **启动后端服务**

   在配置文件修改完成后，通过以下命令启动后端服务：

   ```bash
   go run main.go
   ```

   该命令会编译并运行项目的主文件，启动后端应用程序。此时，服务将处于开发模式下，您可以根据需要进行调试和开发工作。

---

## 📌 注意事项
1. 阿里云短信服务需完成[资质审核](https://help.aliyun.com/zh/sms/use-cases/apply-for-a-text-message-signature)
2. 生产环境建议：
    - 使用SSL加密数据库连接
    - 配置Redis持久化
    - 定期轮换RSA密钥

---

> 📧 问题反馈：[16655836875@163.com](mailto:zhiyuan@example.com) |  
> 🌐 前端仓库：[AI_Todo_Frontend](https://github.com/zhi-yuan-6/AI_Todo_Frontend)

---