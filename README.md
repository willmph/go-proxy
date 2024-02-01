TCP流量转发工具

## 使用方法

Termux安装环境
pkg install golang

手机端运行
go run main.go --bind 0.0.0.0:55555 --backend 127.0.0.1:45889	
#55555 为本地端口
#45889为盾的端口
先从盾获取接口后在运行此命令


