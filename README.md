# XIU2/CloudflareSpeedTest

[![Go Version](https://img.shields.io/github/go-mod/go-version/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Go&color=00ADD8)](https://github.com/XIU2/CloudflareSpeedTest/blob/master/go.mod)
[![Release Version](https://img.shields.io/github/v/release/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Release&color=1784ff)](https://github.com/XIU2/CloudflareSpeedTest/releases/latest)
[![GitHub license](https://img.shields.io/github/license/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=License&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/blob/master/LICENSE)
[![GitHub Star](https://img.shields.io/github/stars/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Star&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/stargazers)
[![GitHub Fork](https://img.shields.io/github/forks/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Fork&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/network/members)

国外很多网站都在使用 Cloudflare CDN，但分配给中国访客的 IP 并不友好。  
虽然 Cloudflare 公开了所有 [IP 段](https://www.cloudflare.com/ips/) ，但想要在这么多 IP 中找到适合自己的，怕是要累死，所以就有了这个软件。  

该软件可以**测试 Cloudflare CDN 所有 IP 的延迟和速度，获得最快 IP**！  
你可以将 IP 添加到 `Hosts` 文件中，以提高访问使用 Cloudflare CDN 服务的国外网站速度！  

****
### 快速使用

1. 下载编译好的可执行文件 [蓝奏云](https://www.lanzoux.com/b0742hkxe) / [Github](https://github.com/XIU2/CloudflareSpeedTest/releases) 并解压。  
2. 双击运行 `CloudflareST.exe`文件（Windows系统），等待测速...  

测速完毕后，会把结果保存在当前目录下的 `result.csv` 文件中，用记事本打开，排序为**延迟由低到高**，每一列用逗号分隔，分别是：  
```
IP 地址, 测试次数, 成功次数, 成功比率, 平均延迟, 下载速度 (MB/s)
104.27.70.18, 4, 4, 1.00, 150.79, 12.89
```
选择一个平均延迟与下载速度都不错的 IP 放到 `Hosts` 文件中（指向域名）。  

****
### 进阶使用

直接双击运行使用的是默认参数，如果想要测试速度更快、测试结果更全面，可以自定义参数。  
``` cmd
C:\>CloudflareST.exe -h

CloudflareSpeedTest
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最快 IP！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 500
        测速线程数量；数值越大速度越快，请勿超过1000(结果误差大)；(默认 500)
    -t 4
        延迟测速次数；单个 IP 测速次数，为 1 时将过滤丢包的IP，TCP协议；(默认 4)
    -dn 20
        下载测速数量；延迟测速并排序后，从最低延迟起测试下载速度的数量，请勿太多(速度慢)；(默认 20)
    -dt 10
        下载测速时间；单个 IP 测速最长时间，单位：秒；(默认 10)
    -f ip.txt
        IP 数据文件；相对/绝对路径，如包含空格请加上引号；支持其他 CDN IP段，记得禁用下载测试；(默认 ip.txt)
    -o result.csv
        输出结果文件；相对/绝对路径，如包含空格请加上引号；允许 .txt 等后缀；(默认 result.csv)
    -dd
        禁用下载测速；如果带上该参数就是禁用下载测速；(默认 启用)
    -v
        打印程序版本
    -h
        打印帮助说明

示例：
	CloudflareST.exe -n 500 -t 4 -dn 20 -dt 10
    CloudflareST.exe -n 500 -t 4 -dn 20 -dt 10 -f "ip.txt" -o "result.csv" -dd
    CloudflareST.exe -n 500 -t 4 -dn 20 -dt 10 -f "C:\abc\ip.txt" -o "C:\abc\result.csv" -dd
```

#### 使用示例

在 CMD 中运行，或者把启动参数添加到快捷方式中。  
> **注意：** 不需要通顺加上所有参数，按需选择，参数前后顺序随意。  

``` cmd
# CMD 示例
CloudflareST.exe -n 500 -t 4 -dn 20 -dt 10
# 指定 IP数据文件 及 输出结果文件（相对路径，即当前目录下）
CloudflareST.exe -n 500 -t 4 -dn 20 -dt 10 -f "ip.txt" -o "result.csv" -dd
# 指定 IP数据文件 及 输出结果文件（绝对路径，即 C:\abc\ 目录下）
CloudflareST.exe -n 500 -t 4 -dn 20 -dt 10 -f "C:\abc\ip.txt" -o "C:\abc\result.csv" -dd
```

``` cmd
# 快捷方式示例（右键快捷方式 - 目标）
## 如果有引号就放在引号外面，记得引号和 - 之间有空格。
"D:\Program Files\CloudflareST\CloudflareST.exe" -n 500 -t 4 -dn 20 -dt 10
```

****
### 感谢项目
* https://github.com/Spedoske/CloudflareScanner

意外发现了这个项目，看了之后发现正好解决了我的问题，但是我更喜欢用户命令行方式运行，这样会更方便、有更多使用姿势，于是我临时学了下 Golang 并 Fork 修改了一份命令行方式交互的版本，如果有什么问题可以告诉我，虽然我不一定会~

****
### 许可证
The GPL-3.0 License.
