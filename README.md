# evm-listen-transfer
evm系监听地址简单转移eth或铭文的工具

# 执行命令
## 描述
下面俩种方式都行
## -方式1：直接运行
1. 修改etc/deploy.yaml配置文件
2. 右键管理员权限运行 transfer-windows.exe


## -方式2：源码编译
1. 下载golang 安装包,本项目基于golang 1.20版本开发，建议不低于该版本 下载地址：https://go.dev/ 根据系统下相应的版本

2. go env -w GOPROXY=https://goproxy.io,direct

3. go build -o transfer-windows.exe
4. 修改etc/deploy.yaml配置文件
5. 右键管理员权限运行 transfer-windows.exe

