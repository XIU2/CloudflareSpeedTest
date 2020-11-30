# XIU2/CloudflareSpeedTest

[![Go Version](https://img.shields.io/github/go-mod/go-version/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Go&color=00ADD8)](https://github.com/XIU2/CloudflareSpeedTest/blob/master/go.mod)
[![Release Version](https://img.shields.io/github/v/release/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Release&color=1784ff)](https://github.com/XIU2/CloudflareSpeedTest/releases/latest)
[![GitHub license](https://img.shields.io/github/license/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=License&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/blob/master/LICENSE)
[![GitHub Star](https://img.shields.io/github/stars/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Star&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/stargazers)
[![GitHub Fork](https://img.shields.io/github/forks/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Fork&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/network/members)

国外很多网站都在使用 Cloudflare CDN，但分配给中国访客的 IP 并不友好。  
虽然 Cloudflare 公开了所有 [IP 段](https://www.cloudflare.com/ips/) ，但想要在这么多 IP 中找到适合自己的，怕是要累死，所以就有了这个软件。  

该软件可以**测试 Cloudflare CDN 所有 IP 的延迟和速度，获得最快 IP**！觉得好用请**点个⭐鼓励一下下~**  
将 IP 添加到 `Hosts` 文件或 DNS 程序中，以提高访问使用 Cloudflare CDN 的网站速度！  

> 本项目也**适用于其他 CDN**，但是需要自行寻找 **CDN IP 段及下载测速地址**（否则只能延迟测速）！

****
## 快速使用

### 下载运行

1. 下载编译好的可执行文件 [蓝奏云](https://xiu.lanzoux.com/b0742hkxe) / [Github](https://github.com/XIU2/CloudflareSpeedTest/releases) 并解压。  
2. 双击运行 `CloudflareST.exe`文件（Windows），等待测速...  

>  **注意：Linux 系统**请先赋予执行权限 `chmod +x CloudflareST` ，然后再执行 `./CloudflareST` 。  

### 结果示例

测速完毕后，会直接显示**最快的 20 个 IP**，示例：  

```
IP 地址           已发送  已接收  丢包率  平均延迟  下载速度 (MB/s)
104.27.198.101    4       4       0.00    126.52    12.71
104.22.43.157     4       4       0.00    129.38    16.74
104.27.214.140    4       4       0.00    132.02    4.65
104.22.42.165     4       4       0.00    133.63    12.00
104.22.35.177     4       4       0.00    135.75    3.92
104.22.87.44      4       4       0.00    136.00    5.86
104.22.67.122     4       4       0.00    136.50    9.47
104.22.88.154     4       4       0.00    140.75    13.00
104.22.69.218     4       4       0.00    142.00    19.07
104.27.184.10     4       4       0.00    148.02    21.05
...
```

完整结果保存在当前目录下的 `result.csv` 文件中，用**记事本/表格软件**打开，排序为**延迟由低到高**，分别是：  

```
IP 地址, 已发送, 已接收, 丢包率, 平均延迟, 下载速度 (MB/s)
104.27.199.141, 4, 4, 0.00, 139.52, 11.71
```
> 大家可以按照自己的需求，对完整测速数据**进一步筛选处理**！

选择一个平均延迟与下载速度都不错的 IP 放到 `Hosts` 文件中（指向使用 Cloudflare CDN 的网站域名）。  

****
## 进阶使用

直接双击运行使用的是默认参数，如果想要测试速度更快、测试结果更全面，可以自定义参数。  

``` cmd
C:\>CloudflareST.exe -h

CloudflareSpeedTest vX.X.X
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最快 IP！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 500
        测速线程数量；数值越大速度越快，请勿超过 1000(结果误差大)；(默认 500)
    -t 4
        延迟测速次数；单个 IP 测速次数，为 1 时将过滤丢包的IP，TCP协议；(默认 4)
    -tp 443
        延迟测速端口；延迟测速 TCP 协议的端口；(默认 443)
    -dn 20
        下载测速数量；延迟测速并排序后，从最低延迟起下载测速数量，请勿太多(速度慢)；(默认 20)
    -dt 5
        下载测速时间；单个 IP 测速最长时间，单位：秒；(默认 5)
    -url https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png
        下载测速地址；用来 Cloudflare CDN 测速的文件地址，如含有空格请加上引号；
    -tl 200
        延迟时间上限；只输出指定延迟时间以下的结果，数量为 -dn 参数的值，单位：ms；
    -sl 5
        下载速度下限；只输出指定下载速度以上的结果，数量为 -dn 参数的值，单位：MB/s；
    -p 20
        显示结果数量；测速后直接显示指定数量的结果，为 0 时不显示结果直接退出；(默认 20)
    -f ip.txt
        IP 数据文件；如含有空格请加上引号；支持其他 CDN IP段，记得禁用下载测速；(默认 ip.txt)
    -o result.csv
        输出结果文件；如含有空格请加上引号；为空格时不输出结果文件(-o " ")；允许其他后缀；(默认 result.csv)
    -dd
        禁用下载测速；如果带上该参数就是禁用下载测速；(默认 启用)
    -ipv6
        IPv6 测速模式；请确保 IP 数据文件内只包含 IPv6 IP段，软件不支持同时测速 IPv4+IPv6；(默认 IPv4)
    -v
        打印程序版本
    -h
        打印帮助说明
```

> 如果**下载速度都是 0.00**，那说明默认的**下载测速地址**用的人太多又到上限了，**请去这个 [Issues](https://github.com/XIU2/CloudflareSpeedTest/issues/6) 获得解决方法！**  

### 使用示例

在 CMD 中运行，或者把启动参数添加到快捷方式中。  

``` bash
# 命令行示例
# 注意：各参数均有默认值，只有不使用默认值时，才需要手动指定参数的值（按需选择），参数不分前后顺序。  
# 提示： Linux 系统只需要把下面命令中的 CloudflareST.exe 改为 ./CloudflareST 即可。  

# 指定 IPv4 数据文件，不显示结果直接退出（-p 值为 0）
CloudflareST.exe -p 0 -f ip.txt -dd

# 指定 IPv6 数据文件( ipv6.txt )，不显示结果直接退出（-p 值为 0）
CloudflareST.exe -p 0 -f ipv6.txt -dd -ipv6

# 指定 IPv4 数据文件，不输出结果到文件，直接显示结果（-p 值为 10 条）
CloudflareST.exe -p 10 -f ip.txt -o " " -dd

# 指定 IPv4 数据文件 及 输出结果到文件（相对路径，即当前目录下，如果包含空格请加上引号）
CloudflareST.exe -f ip.txt -o result.csv -dd

# 指定 IPv4 数据文件 及 输出结果到文件（绝对路径，即 C:\abc\ 目录下，如果包含空格请加上引号）
CloudflareST.exe -f C:\abc\ip.txt -o C:\abc\result.csv -dd

# 指定下载测速地址（要求：可以直接下载、文件大小超过 200MB、用的是 Cloudflare CDN），如果包含空格请加上引号
CloudflareST.exe -url https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png

# 指定测速条件（只有同时满足三个条件时才会停止测速）：
# 延迟时间上限：200 ms，下载速度下限：0 MB/s，数量：10 个
CloudflareST.exe -tl 200 -dn 10

# 延迟时间上限：0 ms，下载速度下限：5 MB/s，数量：10 个
CloudflareST.exe -sl 5 -dn 10

# 延迟时间上限：200 ms，下载速度下限：5 MB/s，数量：10 个
CloudflareST.exe -tl 200 -sl 5 -dn 10

# 如果一直凑不够指定数量，会一直测速下去。  
# 建议指定下载速度下限时，同时指定延迟时间上限，如果测试到指定延迟还没凑够数，就会终止测速。
# 如果一个满足条件的 IP 都没有，那么就会正常输出结果（和不指定条件一样）。
# 如果你需要通过外部程序进一步筛选处理，那么只需要判断测速结果数量，如果上千个说明一个满足条件的 IP 都没有。
```

``` cmd
# Windows 快捷方式示例（右键快捷方式 - 目标）
## 如果有引号就放在引号外面，记得引号和 - 之间有空格。
### 如果要不输出结果文件，那么请加上 -o " "，引号里的是空格。
"D:\Program Files\CloudflareST\CloudflareST.exe" -n 500 -t 4 -dn 20 -dt 5
```

****
## 感谢项目
* https://github.com/Spedoske/CloudflareScanner

意外发现了这个项目，看了之后发现正好解决了我的问题，但是我更喜欢用户命令行方式运行，这样会更方便、有更多使用姿势，于是我临时学了下 Golang 并 Fork 按照我自己的需求修改了一下（包括但不限于命令行方式交互、直接输出结果等），如果有什么问题可以告诉我，虽然我不一定会~

****
## 许可证
The GPL-3.0 License.
