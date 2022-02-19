# XIU2/CloudflareSpeedTest

[![Go Version](https://img.shields.io/github/go-mod/go-version/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Go&color=00ADD8&logo=go)](https://github.com/XIU2/CloudflareSpeedTest/)
[![Release Version](https://img.shields.io/github/v/release/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Release&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/releases/latest)
[![GitHub license](https://img.shields.io/github/license/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=License&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)
[![GitHub Star](https://img.shields.io/github/stars/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Star&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)
[![GitHub Fork](https://img.shields.io/github/forks/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Fork&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)

## 前排提醒：[关于下载测速不稳定/不可用的 情况说明 及 解决方法...](https://github.com/XIU2/CloudflareSpeedTest/issues/168)

国外很多网站都在使用 Cloudflare CDN，但分配给中国内地访客的 IP 并不友好（延迟高、丢包多、速度慢）。  
虽然 Cloudflare 公开了所有 [IP 段](https://www.cloudflare.com/ips/) ，但想要在这么多 IP 中找到适合自己的，怕是要累死，于是就有了这个软件。  

**「自选优选 IP」测试 Cloudflare CDN 延迟和速度，获取最快 IP (IPv4+IPv6)**！好用的话**点个`⭐`鼓励一下叭~**  

> _分享我其他开源项目：[**TrackersList.com** - 全网热门 BT Tracker 列表！有效提高 BT 下载速度~](https://github.com/XIU2/TrackersListCollection) <img src="https://img.shields.io/github/stars/XIU2/TrackersListCollection.svg?style=flat-square&label=Star&color=4285dd&logo=github" height="16px" />_   
> _[**UserScript** - 🐵 Github 高速下载、知乎增强、自动无缝翻页、护眼模式 等十几个**油猴脚本**！](https://github.com/XIU2/UserScript)<img src="https://img.shields.io/github/stars/XIU2/UserScript.svg?style=flat-square&label=Star&color=4285dd&logo=github" height="16px" />_   

> 本项目也支持对**其他 CDN / 网站 IP** 延迟测速（如：[CloudFront IP 段](https://github.com/XIU2/CloudflareSpeedTest/issues/180)），但下载测速需自行寻找地址。

****
## \# 快速使用

### 下载运行

1. 下载编译好的可执行文件 [蓝奏云](https://pan.lanzouo.com/b0742hkxe) / [Github](https://github.com/XIU2/CloudflareSpeedTest/releases) 并解压。  
2. 双击运行 `CloudflareST.exe` 文件（Windows 系统），等待测速完成...  

<details>
<summary><code><strong>「 点击查看 Linux 系统下的使用示例 」</strong></code></summary>

****

以下命令仅为示例，版本号和文件名请前往 [**Releases**](https://github.com/XIU2/CloudflareSpeedTest/releases) 查看。

``` yaml
# 如果是第一次使用，则建议创建新文件夹（后续更新请跳过该步骤）
mkdir CloudflareST

# 进入文件夹（后续更新，只需要从这里重复下面的下载、解压命令即可）
cd CloudflareST

# 下载 CloudflareST 压缩包（自行根据需求替换 URL 中 [版本号] 和 [文件名]）
wget -N https://download.fastgit.org/XIU2/CloudflareSpeedTest/releases/download/v2.0.3/CloudflareST_linux_amd64.tar.gz
# 考虑到国内直接从 Github 下载速度很慢，这里替换为镜像站了，如果还是下载很慢/无法下载，那就试试下面这几个镜像：
# wget -N https://gh.xiu.workers.dev/https://github.com/XIU2/CloudflareSpeedTest/releases/download/v2.0.3/CloudflareST_linux_amd64.tar.gz
# wget -N https://ghproxy.com/https://github.com/XIU2/CloudflareSpeedTest/releases/download/v2.0.3/CloudflareST_linux_amd64.tar.gz
# 如果下载失败的话，尝试删除 -N 参数（如果是为了更新，则记得提前删除旧压缩包 rm CloudflareST_linux_amd64.tar.gz ）

# 解压（不需要删除旧文件，会直接覆盖，自行根据需求替换 文件名）
tar -zxf CloudflareST_linux_amd64.tar.gz

# 赋予执行权限
chmod +x CloudflareST

# 运行（不带参数）
./CloudflareST

# 运行（带参数示例）
./CloudflareST -dd -tll 90
```

> 如果平**均延迟非常低**（如 0.xx），则说明 CloudflareST **测速时走了代理**，请先关闭代理软件后再测速。  
> 如果在**路由器**上运行（如 OpenWrt），请先关闭路由器内的代理，否则测速结果会**不准确且无法使用**。

</details>

****

> _在**手机**上独立运行 CloudflareST 测速的简单教程：**[Android](https://github.com/XIU2/CloudflareSpeedTest/discussions/61)、[IOS](https://github.com/XIU2/CloudflareSpeedTest/issues/151)**_  
> _**建议测速时避开晚上高峰期（20:00~24:00）**，否则测速结果会与其他时间**相差很大...**_  

### 结果示例

测速完毕后，默认会显示**最快的 10 个 IP**，示例（我联通白天测速结果）：  

``` bash
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

# 如果平均延迟非常低（如 0.xx），则说明 CloudflareST 测速时走了代理，请先关闭代理软件后再测速。
# 如果延迟很低 (几十ms)，且你也不是移动 (香港直连)，那么你就是遇到假墙 IP 了，记得加上 -tll 参数。
# 如果在路由器上运行（如 OpenWrt），请先关闭路由器内的代理，否则测速结果会不准确且无法使用。

# 因为每次测速都是在每个 IP 段中随机 IP，所以每次的测速结果都不可能相同，这是正常的！

# 软件是先 延迟测速并按从低到高排序后，再从 最低延迟的 IP 开始下载测速的，所以：
```

测速结果第一行就是**既下载速度最快、又平均延迟最低的最快 IP**！至于拿来干嘛？取决于你~  

完整结果保存在当前目录下的 `result.csv` 文件中，用**记事本/表格软件**打开，格式如下：  

```
IP 地址, 已发送, 已接收, 丢包率, 平均延迟, 下载速度 (MB/s)
104.27.200.69, 4, 4, 0.00, 146.23, 28.64
```

> _大家可以按自己需求，对完整结果**进一步筛选处理**，或者去看一看进阶使用**指定过滤条件**！_

****
## \# 进阶使用

直接运行使用的是默认参数，如果想要测速结果更全面、更符合自己的要求，可以自定义参数。  

``` cmd
C:\>CloudflareST.exe -h

CloudflareSpeedTest vX.X.X
测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最快 IP (IPv4+IPv6)！
https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 200
        测速线程数量；越多测速越快，性能弱的设备 (如路由器) 请勿太高；(默认 200 最多 1000)
    -t 4
        延迟测速次数；单个 IP 延迟测速次数，为 1 时将过滤丢包的IP，TCP协议；(默认 4 次)
    -tp 443
        指定测速端口；延迟测速/下载测速时使用的端口；(默认 443 端口)
    -dn 10
        下载测速数量；延迟测速并排序后，从最低延迟起下载测速的数量；(默认 10 个)
    -dt 10
        下载测速时间；单个 IP 下载测速最长时间，不能太短；(默认 10 秒)
    -url https://cf.xiu2.xyz/url
        下载测速地址；用来下载测速的 Cloudflare CDN 文件地址，默认地址不保证可用性，建议自建；
    -tl 200
        平均延迟上限；只输出低于指定平均延迟的 IP，可与其他上限/下限搭配；(默认 9999 ms)
    -tll 40
        平均延迟下限；只输出高于指定平均延迟的 IP，可与其他上限/下限搭配、过滤假墙 IP；(默认 0 ms)
    -sl 5
        下载速度下限；只输出高于指定下载速度的 IP，凑够指定数量 [-dn] 才会停止测速；(默认 0.00 MB/s)
    -p 10
        显示结果数量；测速后直接显示指定数量的结果，为 0 时不显示结果直接退出；(默认 10 个)
    -f ip.txt
        IP段数据文件；如路径含有空格请加上引号；支持其他 CDN IP段；(默认 ip.txt)
    -o result.csv
        写入结果文件；如路径含有空格请加上引号；值为空时不写入文件 [-o ""]；(默认 result.csv)
    -dd
        禁用下载测速；禁用后测速结果会按延迟排序 (默认按下载速度排序)；(默认 启用)
    -ipv6
        IPv6测速模式；确保 IP 段数据文件内只包含 IPv6 IP段，软件不支持同时测速 IPv4+IPv6；(默认 IPv4)
    -allip
        测速全部的IP；对 IP 段中的每个 IP (仅支持 IPv4) 进行测速；(默认 每个 IP 段随机测速一个 IP)
    -v
        打印程序版本+检查版本更新
    -h
        打印帮助说明
```

### 使用示例

Windows 要指定参数需要在 CMD 中运行，或者把参数添加到快捷方式目标中。  

> **注意**：各参数均有**默认值**，使用默认值的参数是可以省略的（**按需选择**），参数**不分前后顺序**。  
> **提示**：Linux 系统只需要把下面命令中的 `CloudflareST.exe` 改为 `./CloudflareST` 即可。  

****

#### \# CMD 带参数运行 CloudflareST

对命令行程序不熟悉的人，可能不知道该如何带参数运行，我就简单说一下。

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

很多人打开 CMD 以**绝对路径**运行 CloudflareST 会报错，这是因为默认的 `-f ip.txt` 参数是相对路径，需要指定绝对路径的 ip.txt 才行，但这样毕竟太麻烦了，因此还是建议进入 CloudflareST 程序目录下，以**相对路径**方式运行：

1. 打开 CloudflareST 程序所在目录
2. 空白处按下 <kbd>Shift + 鼠标右键</kbd> 显示右键菜单
3. 选择 **\[在此处打开命令窗口\]** 来打开 CMD 窗口，此时默认就位于当前目录下
4. 输入带参数的命令，如：`CloudflareST.exe -tll 50 -tl 200`即可运行

> 当然你也可以随便打开一个 CMD 窗口，然后输入如 `cd /d "D:\Program Files\CloudflareST"` 来进入程序目录

</details>

****

#### \# Windows 快捷方式带参数运行 CloudflareST

如果不经常修改运行参数（比如平时都是直接双击运行）的人，建议使用快捷方式，更方便点。

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

右键 `CloudflareST.exe` 文件 - **\[创建快捷方式\]**，然后右键该快捷方式 - **\[属性\]**，修改其**目标**：

``` bash
# 如果要不输出结果文件，那么请加上 -o " "，引号里的是空格（没有空格会导致该参数被省略）。
D:\ABC\CloudflareST\CloudflareST.exe -n 500 -t 4 -dn 20 -dt 5 -o " "

# 如果文件路径包含引号，则需要把启动参数放在引号外面，记得引号和 - 之间有空格。
"D:\Program Files\CloudflareST\CloudflareST.exe" -n 500 -t 4 -dn 20 -dt 5 -o " "

# 注意！快捷方式 - 起始位置 不能是空的，否则就会因为绝对路径而找不到 ip.txt 文件
```

</details>

****

#### \# IPv4/IPv6

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****
``` bash
# 测速 IPv4 时，需要指定 IPv4 数据文件（-f 默认值就是 ip.txt，所以该参数可以省略）
CloudflareST.exe -f ip.txt

# 测速 IPv6 时，需要指定 IPv6 数据文件( ipv6.txt ) 的同时再加上 -ipv6 参数
CloudflareST.exe -f ipv6.txt -ipv6
```

</details>

****

#### \# 文件相对/绝对路径

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

``` bash
# 指定 IPv4 数据文件，不显示结果直接退出，输出结果到文件（-p 值为 0）
CloudflareST.exe -f 1.txt -p 0 -dd

# 指定 IPv4 数据文件，不输出结果到文件，直接显示结果（-p 值为 10 条，-o 值为空但引号不能少）
CloudflareST.exe -f 2.txt -o "" -p 10 -dd

# 指定 IPv4 数据文件 及 输出结果到文件（相对路径，即当前目录下，如含空格请加上引号）
CloudflareST.exe -f 3.txt -o result.txt -dd


# 指定 IPv4 数据文件 及 输出结果到文件（相对路径，即当前目录内的 abc 文件夹下，如含空格请加上引号）
# Linux（CloudflareST 程序所在目录内的 abc 文件夹下）
./CloudflareST -f abc/3.txt -o abc/result.txt -dd

# Windows（注意是反斜杠）
CloudflareST.exe -f abc\3.txt -o abc\result.txt -dd


# 指定 IPv4 数据文件 及 输出结果到文件（绝对路径，即 C:\abc\ 目录下，如含空格请加上引号）
# Linux（/abc/ 目录下）
./CloudflareST -f /abc/4.txt -o /abc/result.csv -dd

# Windows（注意是反斜杠）
CloudflareST.exe -f C:\abc\4.txt -o C:\abc\result.csv -dd


# 如果要以【绝对路径】运行 CloudflareST，那么 -f / -o 参数中的文件名也必须是【绝对路径】，否则会报错找不到文件！
# Linux（/abc/ 目录下）
/abc/CloudflareST -f /abc/4.txt -o /abc/result.csv -dd

# Windows（注意是反斜杠）
C:\abc\CloudflareST.exe -f C:\abc\4.txt -o C:\abc\result.csv -dd
```
</details>

****

#### \# 自定义下载测速地址

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

因为目前默认下载测速地址流量太大被 Cloudflare 限速，因此建议大家**改用其他**下载测速地址（如下面的 Cloudflare 官方下载测速地址），更多请见： [#168](https://github.com/XIU2/CloudflareSpeedTest/issues/168) 

``` bash
# 地址要求：可以直接下载、文件大小超过 200MB、用的是 Cloudflare CDN
CloudflareST.exe -url https://cf.xiu2.xyz/url

# 因为默认下载测速地址的文件大小只有 300MB，如果你速度太快的话，测速结果可能会低于实际速度。
# 因此推荐使用 Cloudflare CDN 官方下载测速地址（300MB 且可自定义大小，即末尾数字）：
CloudflareST.exe -url https://speed.cloudflare.com/__down?bytes=300000000

# 注意：如果下载测速地址为 HTTP 协议，记得加上 -tp 80（这个参数会影响 延迟测速/下载测速 时使用的端口）
CloudflareST.exe -tp 80 -url http://xxx/xxx
```

</details>

****

#### \# 自定义测速条件

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

- 指定 **[平均延迟下限]** 条件（用于过滤**被假蔷的 IP**，这类 IP 都被 TCP 劫持，因此延迟很低只有几十ms）

``` bash
# 平均延迟下限：40 ms （一般除了移动直连香港外，几乎不存在低于 100ms 的，自行测试适合的下限延迟）
# 平均延迟下限和其他的上下限参数一样，都可以单独使用、互相搭配使用！
CloudflareST.exe -tll 40
```

- 仅指定 **[平均延迟上限]** 条件

``` bash
# 平均延迟上限：200 ms，下载速度下限：0 MB/s，数量：10 个（可选）
# 即找到平均延迟低于 200 ms 的 IP，然后再按延迟从低到高进行 10 次下载测速
CloudflareST.exe -tl 200 -dn 10
```

> 如果没有一个 IP **平均延迟低于 200ms**，那么不会输出任何内容。  

****

- 仅指定 **[平均延迟上限]** 条件，且**只延迟测速，不下载测速**

``` bash
# 平均延迟上限：200 ms，下载速度下限：0 MB/s，数量：不知道多少 个
# 即只输出低于 200ms 的 IP，且不再下载测速（因为不再下载测速，所以 -dn 参数就无效了）
CloudflareST.exe -tl 200 -dd
```

****

- 仅指定 **[下载速度下限]** 条件

``` bash
# 平均延迟上限：9999 ms，下载速度下限：5 MB/s，数量：10 个（可选）
# 即需要找到 10 个平均延迟低于 9999 ms 且下载速度高于 5 MB/s 的 IP 才会停止测速
CloudflareST.exe -sl 5 -dn 10
```

> 没有指定平均延迟上限时，如果一直**凑不够**满足条件的 IP 数量，就会**一直测速**下去。    
> 所以建议**同时指定 [下载速度下限] + [平均延迟上限]**，这样测速到指定延迟上限还没凑够数量，就会终止测速。

****

- 同时指定 **[平均延迟上限] + [下载速度下限]** 条件

``` bash
# 平均延迟上限、下载速度下限均支持小数（如 -sl 0.5）
# 平均延迟上限：200 ms，下载速度下限：5.6 MB/s，数量：10 个（可选）
# 即需要找到 10 个平均延迟低于 200 ms 且下载速度高于 5 .6MB/s 的 IP 才会停止测速
CloudflareST.exe -tl 200 -sl 5.6 -dn 10
```

> 如果没有一个 IP **平均延迟低于 200ms**，那么不会输出任何内容。  
> 如果没有一个 IP **下载速度高于 5.6 MB/s**，那么就会**和不指定 [下载速度下限] 条件一样**输出结果。  
> 所以建议先不指定条件测速一遍，看看平均延迟和下载速度大概在什么范围，避免指定条件**过低/过高**！

> 因为Cloudflare 公开的 IP 段是**回源 IP+任播 IP**，而**回源 IP**是无法使用的，所以下载测速是 0.00。  
> 运行时可以加上 `-sl 0.01`（下载速度下限），过滤掉**回源 IP**（下载测速低于 0.01MB/s 的结果）。

</details>

****

#### \# 单独对一个或多个 IP 测速

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

如果要单独**对一个或多个 IP 进行测速**，只需要把这些 IP 按如下格式写入到任意文本文件中，例如：`1.txt`

```
1.1.1.1
1.1.1.200
1.0.0.1/24
```

> 单个 IP 的话可以省略 `/32` 子网掩码了（即 `1.1.1.1`等同于 `1.1.1.1/32`）。  
> 子网掩码 `/24` 指的是这个 IP 最后一段，即 `1.0.0.1~1.0.0.255`。


然后运行 CloudflareST 时加上启动参数 `-f 1.txt` 即可。

``` bash
# 先进入 CloudflareST 所在目录，然后运行：
# Windows 系统（在 CMD 中运行）
CloudflareST.exe -f 1.txt

# Linux 系统
./CloudflareST -f 1.txt

# 对于 1.0.0.1/24 这样的 IP 段只会随机最后一段（1.0.0.1~255），如果要测速该 IP 段中的所有 IP，请加上 -allip 参数。
```

</details>

****

#### \# 一劳永逸加速所有使用 Cloudflare CDN 的网站（不需要再一个个添加域名到 Hosts 了）

我以前说过，开发该软件项目的目的就是为了通过**改 Hosts 的方式来加速访问使用 Cloudflare CDN 的网站**。

但就如 [**#8**](https://github.com/XIU2/CloudflareSpeedTest/issues/8) 所说，一个个添加域名到 Hosts 实在**太麻烦**了，于是我就找到了个**一劳永逸**的办法！

可以看这个 [**还在一个个添加 Hosts？完美本地加速所有使用 Cloudflare CDN 的网站方法来了！**](https://github.com/XIU2/CloudflareSpeedTest/discussions/71) 

****

#### \# 自动更新 Hosts

考虑到很多人获得最快 Cloudflare CDN IP 后，需要替换 Hosts 文件中的 IP。

可以看这个 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues/42) 获取 **Windows/Linux 自动更新 Hosts 脚本**！

****

## 问题反馈

如果你遇到什么问题，可以先去 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues) 里看看是否有别人问过了（记得去看下  [**Closed**](https://github.com/XIU2/CloudflareSpeedTest/issues?q=is%3Aissue+is%3Aclosed) 的）。  
如果没找到类似问题，请新开个 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues/new) 来告诉我！

> _有问题请**大胆告诉我**，描述越详细越好（必要时可远程协助），如果不说那我怎么去完善功能或~~修复 BUG~~ 呢？！_

****

## 感谢项目

* https://github.com/Spedoske/CloudflareScanner

> _因为该项目已经很长时间没更新了，而我又产生了很多功能需求，所以我临时学了下 Go 语言就上手了（菜）..._  
> _本软件基于该项目制作，但**已添加大量功能及修复 BUG**，并根据大家的使用反馈积极添加、优化功能（闲）..._

****

## License
The GPL-3.0 License.
