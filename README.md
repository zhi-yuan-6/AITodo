# 项目名称

这是一个后端服务项目，旨在[简要描述项目目标和功能，例如“提供用户认证和短信通知功能”]。本指南将帮助您快速配置和运行项目。

## 先决条件

在开始之前，请确保您的环境中已安装以下组件：
- **Go**：1.16 或更高版本
- **MySQL**：5.7 或更高版本
- **Redis**：6.0 或更高版本
- **阿里云SMS服务账号**（可选）：如果需要使用短信功能，请提前注册并获取相关配置

## 快速启动

按照以下步骤配置和运行项目：

### 1. 配置数据库和Redis
- 安装并启动MySQL和Redis服务。请参考各自的官方文档完成安装和配置，确保服务可正常访问。

### 2. 创建配置文件
- 在项目的`config`目录下创建`config.yaml`文件。
- 参考`config/config.yaml.example`中的示例配置，填写您的MySQL和Redis连接信息。
- **SMS配置（可选）**：
    - 如果需要使用阿里云SMS短信服务，请在`config.yaml`中填写`sms`部分的配置。
    - 通过[阿里云SMS快速入门](https://help.aliyun.com/zh/sms/getting-started/get-started-with-sms)获取以下信息：
        - AccessKey ID
        - AccessKey Secret
        - SignName
        - TemplateCode
    - 将这些信息填入`config.yaml`的相应字段。

### 3. 生成私钥文件
- 创建一个名为`private_key.pem`的文件，并将其放置在项目根目录下。
- **生成方法**：
    - **使用项目工具**：运行`go run util/pem.go`，调用其中的`GeneratePEM`函数生成私钥，然后将生成的私钥内容复制到`private_key.pem`文件中。
    - **其他方式**：您也可以使用OpenSSL等工具生成PEM格式的私钥，例如运行命令：
      ```bash
      openssl genrsa -out private_key.pem 2048
确保private_key.pem的内容格式正确，可参考config/config.yaml.example中的示例。
4. 运行项目
   在项目根目录下执行以下命令启动服务：
    ```bash
    go run main.go

