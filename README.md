# sdoor

适用于边界服务器，当当前主机IP通过防火墙映射到一个公网IP或者本身分配的地址就是公网IP时,可以使用Sdoor在目标主机上打开一个带有socks服务的端口，用于流量的代理

## 免责声明
免责声明：此工具仅限于安全研究，用户承担因使用此工具而导致的所有法律和相关责任！作者不承担任何法律责任！

## 功能

1、在本地机器上监听一个命令行参数输入的端口  
2、该端口支持socks5协议用户名、密码的认证方式，或者没有密码的认证方式  

## 参数说明

port 指定监听哪一个端口  
user 指定用户名（用户名密码必须同时存在）  
pass 指定密码  
