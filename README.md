# XIU2/CloudflareSpeedTest

[![Go Version](https://img.shields.io/github/go-mod/go-version/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Go&color=00ADD8)](https://github.com/XIU2/CloudflareSpeedTest/blob/master/go.mod)
[![Release Version](https://img.shields.io/github/v/release/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Release&color=1784ff)](https://github.com/XIU2/CloudflareSpeedTest/releases/latest)
[![GitHub license](https://img.shields.io/github/license/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=License&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/blob/master/LICENSE)
[![GitHub Star](https://img.shields.io/github/stars/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Star&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/stargazers)
[![GitHub Fork](https://img.shields.io/github/forks/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Fork&color=f38020)](https://github.com/XIU2/CloudflareSpeedTest/network/members)

国外很多网站都在使用 Cloudflare CDN，但分配给中国访客的 IP 并不友好。  
虽然 Cloudflare 公开了所有 [IP 段](https://www.cloudflare.com/ips/) ，但想要在这么多 IP 中找到适合自己的，怕是要累死，所以就有了这个软件。  

该软件可以**测试 Cloudflare CDN 延迟和速度，获取最快 IP (IPv4+IPv6)**！觉得好用请**点个⭐鼓励一下下~**  

> 本项目也**适用于其他 CDN**，但是需要自行寻找 **CDN IP 段及下载测速地址**（否则只能延迟测速）。

****
## 快速使用

### 下载运行

1. 下载编译好的可执行文件 [蓝奏云](https://xiu.lanzoux.com/b0742hkxe) / [Github](https://github.com/XIU2/CloudflareSpeedTest/releases) 并解压。  
2. 双击运行 `CloudflareST.exe`文件（Windows），等待测速完成...  

>  **提示：Linux 系统**请先赋予执行权限 `chmod +x CloudflareST` ，然后再执行 `./CloudflareST`   

> **注意：建议测速时避开高峰期（晚上~凌晨）**，否则测速结果会与其他时间**差距很大...**  

### 结果示例

测速完毕后，默认会显示**最快的 20 个 IP**，示例（我的白天测速结果）：  

```
IP 地址           已发送  已接收  丢包率  平均延迟  下载速度 (MB/s)
104.27.200.69     4       4       0.00    146.23    28.64
172.67.60.78      4       4       0.00    139.82    15.02
104.25.140.153    4       4       0.00    146.49    14.90
104.27.192.65     4       4       0.00    140.28    14.07
172.67.62.214     4       4       0.00    139.29    12.71
104.27.207.5      4       4       0.00    145.92    11.95
172.67.54.193     4       4       0.00    146.71    11.55
104.22.66.8       4       4       0.00    147.42    11.11
104.27.197.63     4       4       0.00    131.29    10.26
172.67.58.91      4       4       0.00    140.19    9.14
...
```

测速结果第一行是**兼顾平均延迟与下载速度的最快 IP**，至于拿来干嘛？取决于你~  

完整结果保存在当前目录下的 `result.csv` 文件中，用**记事本/表格软件**打开，格式如下：  

```
IP 地址, 已发送, 已接收, 丢包率, 平均延迟, 下载速度 (MB/s)
104.27.200.69, 4, 4, 0.00, 146.23, 28.64
```

> 大家可以按自己需求，对完整结果**进一步筛选处理**，或者去看一看进阶使用**指定过滤条件**！

****
## 进阶使用

直接运行使用的是默认参数，如果想要测速结果更全面、更符合自己的要求，可以自定义参数。  

``` cmd
C:\>CloudflareST.exe -h

CloudflareSpeedTest vX.X.X
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最快 IP！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 500
        测速线程数量；线程数量越多延迟测速越快，请勿超过 1000 (误差大)；(默认 500)
    -t 4
        延迟测速次数；单个 IP 延迟测速次数，为 1 时将过滤丢包的IP，TCP协议；(默认 4)
    -tp 443
        延迟测速端口；延迟测速 TCP 协议的端口；(默认 443)
    -dn 20
        下载测速数量；延迟测速并排序后，从最低延迟起下载测速的数量；(默认 20)
    -dt 10
        下载测速时间；单个 IP 下载测速最长时间，单位：秒；(默认 10)
    -url https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png
        下载测速地址；用来下载测速的 Cloudflare CDN 文件地址，如地址含有空格请加上引号；
    -tl 200
        平均延迟上限；只输出低于指定平均延迟的 IP，与下载速度下限搭配使用；(默认 9999 ms)
    -sl 5
        下载速度下限；只输出高于指定下载速度的 IP，凑够指定数量 [-dn] 才会停止测速；(默认 0 MB/s)
    -p 20
        显示结果数量；测速后直接显示指定数量的结果，为 0 时不显示结果直接退出；(默认 20)
    -f ip.txt
        IP段数据文件；如路径含有空格请加上引号；支持其他 CDN IP段；(默认 ip.txt)
    -o result.csv
        输出结果文件；如路径含有空格请加上引号；值为空格时不输出 [-o " "]；(默认 result.csv)
    -dd
        禁用下载测速；禁用后测速结果会按延迟排序（默认按下载速度排序）；(默认 启用)
    -ipv6
        IPv6测速模式；确保 IP 段数据文件内只包含 IPv6 IP段，软件不支持同时测速 IPv4+IPv6；(默认 IPv4)
    -allip
        测速全部的IP；对 IP 段中的每个 IP (仅支持 IPv4) 进行测速；(默认 每个 IP 段随机测速一个 IP)
    -v
        打印程序版本+检查版本更新
    -h
        打印帮助说明
```

> 如果**下载速度都是 0.00**，那可能默认的**下载测速地址**用的人太多到上限了，**请去这个 [Issues](https://github.com/XIU2/CloudflareSpeedTest/issues/6) 获得解决方法！**  

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

# ——————————————————————

# 指定 IPv4 数据文件，不输出结果到文件，直接显示结果（-p 值为 10 条）
CloudflareST.exe -p 10 -f ip.txt -o " " -dd

# 指定 IPv4 数据文件 及 输出结果到文件（相对路径，即当前目录下，如含空格请加上引号）
CloudflareST.exe -f ip.txt -o result.csv -dd

# 指定 IPv4 数据文件 及 输出结果到文件（绝对路径，即 C:\abc\ 目录下，如含空格请加上引号）
CloudflareST.exe -f C:\abc\ip.txt -o C:\abc\result.csv -dd

# ——————————————————————

# 指定下载测速地址（要求：可以直接下载、文件大小超过 200MB、用的是 Cloudflare CDN），如含空格请加上引号
CloudflareST.exe -url https://cf.xiu2.xyz/Github/CloudflareSpeedTest.png

# ——————————————————————

# 指定测速条件（只有同时满足三个条件时才会停止测速）：

# 平均延迟上限：9999 ms，下载速度下限：5 MB/s，数量：10 个
# 即需要找到 10 个平均延迟低于 9999 ms 且 下载速度高于 5 MB/s 的 IP 才会停止测速。
CloudflareST.exe -sl 5 -dn 10

# 没有指定平均延迟上限时，如果一直凑不够满足条件的 IP 数量，会一直测速下去。  
# 所以建议同时指定 下载速度下限 和 平均延迟上限，这样测试到指定延迟还没凑够数量，就会终止测速。

# 平均延迟上限：200 ms，下载速度下限：5 MB/s，数量：10 个
# 即需要找到 10 个平均延迟低于 200 ms 且 下载速度高于 5 MB/s 的 IP 才会停止测速。
CloudflareST.exe -tl 200 -sl 5 -dn 10

# 如果一个满足条件的 IP 都没找到，那么就会和不指定条件一样输出完整结果。
# 所以建议先不指定条件测速一遍，看看平均延迟和下载速度大概在什么范围，避免指定条件过低/过高！
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

****
## 许可证
The GPL-3.0 License.
