1.在config目录下创建config.yaml并根据config.yaml.example进行配置，mysql数据库和redis的具体安装设置这里不做具体讲解
其中配置文件中sms部分，是阿里云sms短信服务所需要的，可以根据下列链接中进行相应获取配置：https://help.aliyun.com/zh/sms/getting-started/get-started-with-sms
2.创建private_key.pem文件。可以使用util/pem.go中的GeneratePEM函数进行生成，或者使用其他方式，之后将生产的pem私钥复制到private_key.pem文件中即可，可参考config.yaml.example中的示例，。
3.执行go run main.go即可运行