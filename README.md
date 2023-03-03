# XIU2/CloudflareSpeedTest

[![Go Version](https://img.shields.io/github/go-mod/go-version/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Go&color=00ADD8&logo=go)](https://github.com/XIU2/CloudflareSpeedTest/)
[![Release Version](https://img.shields.io/github/v/release/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Release&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/releases/latest)
[![GitHub license](https://img.shields.io/github/license/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=License&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)
[![GitHub Star](https://img.shields.io/github/stars/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Star&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)
[![GitHub Fork](https://img.shields.io/github/forks/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Fork&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)

国外很多网站都在使用 Cloudflare CDN，但分配给中国内地访客的 IP 并不友好（延迟高、丢包多、速度慢）。  
虽然 Cloudflare 公开了所有 [IP 段](https://www.cloudflare.com/ips/) ，但想要在这么多 IP 中找到适合自己的，怕是要累死，于是就有了这个软件。

**「自选优选 IP」测试 Cloudflare CDN 延迟和速度，获取最快 IP (IPv4+IPv6)**！好用的话**点个`⭐`鼓励一下叭~**

> _分享我其他开源项目：[**TrackersList.com** - 全网热门 BT Tracker 列表！有效提高 BT 下载速度~](https://github.com/XIU2/TrackersListCollection) <img src="https://img.shields.io/github/stars/XIU2/TrackersListCollection.svg?style=flat-square&label=Star&color=4285dd&logo=github" height="16px" />_  
> _[**UserScript** - 🐵 Github 高速下载、知乎增强、自动无缝翻页、护眼模式 等十几个**油猴脚本**！](https://github.com/XIU2/UserScript)<img src="https://img.shields.io/github/stars/XIU2/UserScript.svg?style=flat-square&label=Star&color=4285dd&logo=github" height="16px" />_

> 本项目也支持对**其他 CDN / 网站 IP** 延迟测速（如：[CloudFront](https://github.com/XIU2/CloudflareSpeedTest/discussions/304)、[Gcore](https://github.com/XIU2/CloudflareSpeedTest/discussions/303) CDN），但下载测速需自行寻找地址

> 对于**代理套 Cloudflare CDN** 的用户，须知这应为**备用方案**，而不应该是**唯一方案**，请勿过度依赖 [#217](https://github.com/XIU2/CloudflareSpeedTest/issues/217) [#188](https://github.com/XIU2/CloudflareSpeedTest/issues/188)

****
## \# 快速使用

### 下载运行

1. 下载编译好的可执行文件 [蓝奏云](https://pan.lanzouf.com/b0742hkxe) / [Github](https://github.com/XIU2/CloudflareSpeedTest/releases) 并解压。  
2. 双击运行 `CloudflareST.exe` 文件（Windows 系统），等待测速完成...

<details>
<summary><code><strong>「 点击查看 Linux 系统下的使用示例 」</strong></code></summary>

****

以下命令仅为示例，版本号和文件名请前往 [**Releases**](https://github.com/XIU2/CloudflareSpeedTest/releases) 查看。

``` yaml
# 如果是第一次使用，则建议创建新文件夹（后续更新时，跳过该步骤）
mkdir CloudflareST

# 进入文件夹（后续更新，只需要从这里重复下面的下载、解压命令即可）
cd CloudflareST

# 下载 CloudflareST 压缩包（自行根据需求替换 URL 中 [版本号] 和 [文件名]）
wget -N https://github.com/XIU2/CloudflareSpeedTest/releases/download/v2.2.2/CloudflareST_linux_amd64.tar.gz
# 如果你是在国内服务器上下载，那么请使用下面这几个镜像加速：
# wget -N https://download.fastgit.org/XIU2/CloudflareSpeedTest/releases/download/v2.2.2/CloudflareST_linux_amd64.tar.gz
# wget -N https://ghproxy.com/https://github.com/XIU2/CloudflareSpeedTest/releases/download/v2.2.2/CloudflareST_linux_amd64.tar.gz
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
> 如果在**路由器**上运行，建议先关闭路由器内的代理（或将其排除），否则测速结果可能会**不准确/无法使用**。

</details>

****

> _在**手机**上独立运行 CloudflareST 测速的简单教程：**[Android](https://github.com/XIU2/CloudflareSpeedTest/discussions/61)、[Android APP](https://github.com/xianshenglu/cloudflare-ip-tester-app)、[IOS](https://github.com/XIU2/CloudflareSpeedTest/discussions/321)**_

### 结果示例

测速完毕后，默认会显示**最快的 10 个 IP**，示例：

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
# 如果在路由器上运行，请先关闭路由器内的代理（或将其排除），否则测速结果可能会不准确/无法使用。

# 因为每次测速都是在每个 IP 段中随机 IP，所以每次的测速结果都不可能相同，这是正常的！

# 注意！我发现电脑开机后第一次测速延迟会明显偏高（手动 TCPing 也一样），后续测速都正常
# 因此建议大家开机后第一次正式测速前，先随便测几个 IP（无需等待延迟测速完成，只要进度条动了就可以直接关了）

# 软件在 默认参数 下的整个流程大概步骤：
# 1. 延迟测速（默认 TCPing 模式，HTTPing 模式需要手动加上参数）
# 2. 延迟排序（延迟从低到高排序，不同丢包率的会分开独立排序，因此可能会有一些延迟低但丢包的 IP 被排到后面）
# 3. 下载测速（从延迟最低的 IP 开始依次下载测速，默认测够 10 个就会停止）
# 4. 速度排序（速度从高到低排序）
# 5. 输出结果（可依靠参数控制是否输出到命令行(-p 0)/文件(-o "")）
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
        延迟测速线程；越多延迟测速越快，性能弱的设备 (如路由器) 请勿太高；(默认 200 最多 1000)
    -t 4
        延迟测速次数；单个 IP 延迟测速次数，为 1 时将过滤丢包的IP；(默认 4 次)
    -dn 10
        下载测速数量；延迟测速并排序后，从最低延迟起下载测速的数量；(默认 10 个)
    -dt 10
        下载测速时间；单个 IP 下载测速最长时间，不能太短；(默认 10 秒)
    -tp 443
        指定测速端口；延迟测速/下载测速时使用的端口；(默认 443 端口)
    -url https://cf.xiu2.xyz/url
        指定测速地址；延迟测速(HTTPing)/下载测速时使用的地址，默认地址不保证可用性，建议自建；

    -httping
        切换测速模式；延迟测速模式改为 HTTP 协议，所用测试地址为 [-url] 参数；(默认 TCPing)
        注意：HTTPing 本质上也算一种 网络扫描 行为，因此如果你在服务器上面运行，需要降低并发(-n)，否则可能会被一些严格的商家暂停服务。
        如果你遇到 HTTPing 首次测速可用 IP 数量正常，后续测速越来越少甚至直接为 0，但停一段时间后又恢复了的情况，那么也可能是被 运营商、Cloudflare CDN 认为你在网络扫描而 触发临时限制机制，因此才会过一会儿就恢复了，建议降低并发(-n)减少这种情况的发生。
    -httping-code 200
        有效状态代码；HTTPing 延迟测速时网页返回的有效 HTTP 状态码，仅限一个；(默认 200 301 302)
    -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD
        匹配指定地区；地区名为当地机场三字码，英文逗号分隔，仅 HTTPing 模式可用；(默认 所有地区)

    -tl 200
        平均延迟上限；只输出低于指定平均延迟的 IP，可与其他上限/下限搭配；(默认 9999 ms)
    -tll 40
        平均延迟下限；只输出高于指定平均延迟的 IP，可与其他上限/下限搭配；(默认 0 ms)
    -sl 5
        下载速度下限；只输出高于指定下载速度的 IP，凑够指定数量 [-dn] 才会停止测速；(默认 0.00 MB/s)

    -p 10
        显示结果数量；测速后直接显示指定数量的结果，为 0 时不显示结果直接退出；(默认 10 个)
    -f ip.txt
        IP段数据文件；如路径含有空格请加上引号；支持其他 CDN IP段；(默认 ip.txt)
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
        指定IP段数据；直接通过参数指定要测速的 IP 段数据，英文逗号分隔；(默认 空)
    -o result.csv
        写入结果文件；如路径含有空格请加上引号；值为空时不写入文件 [-o ""]；(默认 result.csv)

    -dd
        禁用下载测速；禁用后测速结果会按延迟排序 (默认按下载速度排序)；(默认 启用)
    -allip
        测速全部的IP；对 IP 段中的每个 IP (仅支持 IPv4) 进行测速；(默认 每个 /24 段随机测速一个 IP)

    -v
        打印程序版本 + 检查版本更新
    -h
        打印帮助说明
```

### 使用示例

Windows 要指定参数需要在 CMD 中运行，或者把参数添加到快捷方式目标中。

> **注意**：各参数均有**默认值**，使用默认值的参数是可以省略的（**按需选择**），参数**不分前后顺序**。  
> **提示**：Windows **PowerShell** 只需把下面命令中的 `CloudflareST.exe` 改为 `.\CloudflareST.exe` 即可。  
> **提示**：Linux 系统只需要把下面命令中的 `CloudflareST.exe` 改为 `./CloudflareST` 即可。

****

#### \# CMD 带参数运行 CloudflareST

对命令行程序不熟悉的人，可能不知道该如何带参数运行，我就简单说一下。

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

很多人打开 CMD 以**绝对路径**运行 CloudflareST 会报错，这是因为默认的 `-f ip.txt` 参数是相对路径，需要指定绝对路径的 ip.txt 才行，但这样毕竟太麻烦了，因此还是建议进入 CloudflareST 程序目录下，以**相对路径**方式运行：

**方式 一**：
1. 打开 CloudflareST 程序所在目录  
2. 空白处按下 <kbd>Shift + 鼠标右键</kbd> 显示右键菜单  
3. 选择 **\[在此处打开命令窗口\]** 来打开 CMD 窗口，此时默认就位于当前目录下  
4. 输入带参数的命令，如：`CloudflareST.exe -tll 50 -tl 200`即可运行

**方式 二**：
1. 打开 CloudflareST 程序所在目录  
2. 直接在文件夹地址栏中全选并输入 `cmd` 回车来打开 CMD 窗口，此时默认就位于当前目录下  
4. 输入带参数的命令，如：`CloudflareST.exe -tll 50 -tl 200`即可运行

> 当然你也可以随便打开一个 CMD 窗口，然后输入如 `cd /d "D:\Program Files\CloudflareST"` 来进入程序目录

> **提示**：如果用的是 **PowerShell** 只需把命令中的 `CloudflareST.exe` 改为 `.\CloudflareST.exe` 即可。

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
# 测速 IPv4 时，需要指定 IPv4 数据文件（-f 默认值就是 ip.txt，所以该参数可省略）
CloudflareST.exe -f ip.txt

# 测速 IPv6 时，需要指定 IPv6 数据文件（v2.1.0 版本后支持 IPv4+IPv6 混合测速并移除了 -ipv6 参数）
CloudflareST.exe -f ipv6.txt

# 当然你也可以将 IPv4 IPv6 混合在一起测速，也可以直接通过参数指定要测速的 IP
CloudflareST.exe -ip 1.1.1.1,2606:4700::/32
```

> 测速 IPv6 时，可能会注意到每次测速数量都不一样，了解原因： [#120](https://github.com/XIU2/CloudflareSpeedTest/issues/120)  
> 因为 IPv6 太多（以亿为单位），且绝大部分 IP 段压根未启用，所以我只扫了一部分可用的 IPv6 段写到 `ipv6.txt` 文件中，有兴趣的可以自行扫描增删，ASN 数据源来自：[bgp.he.net](https://bgp.he.net/AS13335#_prefixes6)

</details>

****

#### \# HTTPing

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

目前有两种延迟测速模式，分别为 **TCP 协议、HTTP 协议**。  
TCP 协议耗时更短、消耗资源更少，超时时间为 1 秒，该协议为默认模式。  
HTTP 协议适用于快速测试某域名指向某 IP 时是否可以访问，超时时间为 2 秒。  
同一个 IP，各协议去 Ping 得到的延迟一般为：**ICMP < TCP < HTTP**，越靠右对丢包等网络波动越敏感。

> 注意：HTTPing 本质上也算一种**网络扫描**行为，因此如果你在服务器上面运行，需要**降低并发**(`-n`)，否则可能会被一些严格的商家暂停服务。如果你遇到 HTTPing 首次测速可用 IP 数量正常，后续测速越来越少甚至直接为 0，但停一段时间后又恢复了的情况，那么也可能是被 运营商、Cloudflare CDN 认为你在网络扫描而**触发临时限制机制**，因此才会过一会儿就恢复了，建议**降低并发**(`-n`)减少这种情况的发生。

``` bash
# 只需加上 -httping 参数即可切换到 HTTP 协议延迟测速模式
CloudflareST.exe -httping

# 软件会根据访问时网页返回的有效 HTTP 状态码来判断可用性（当然超时也算），默认对返回 200 301 302 这三个 HTTP 状态码的视为有效，可以手动指定认为有效的 HTTP 状态码，但只能指定一个（你需要提前确定测试地址正常情况下会返回哪个状态码）
CloudflareST.exe -httping -httping-code 200

# 通过 -url 参数来指定 HTTPing 测试地址（可以是任意网页 URL，不局限于具体文件地址）
CloudflareST.exe -httping -url https://cf.xiu2.xyz/url

# 注意：如果测速地址为 HTTP 协议，记得加上 -tp 80（这个参数会影响 延迟测速/下载测速 时使用的端口）
# 同理，如果要测速 80 端口，那么也需要加上 -url 参数来指定一个 http:// 协议的地址才行（且该地址不会强制重定向至 HTTPS），如果是非 80 443 端口，那么需要确定该下载测速地址是否支持通过该端口访问。
CloudflareST.exe -httping -tp 80 -url http://cdn.cloudflare.steamstatic.com/steam/apps/5952/movie_max.webm
```

</details>

****

#### \# 匹配指定地区(colo 机场三字码)

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

``` bash
# 指定地区名后，延迟测速后得到的结果就都是指定地区的 IP 了（也可以继续进行下载测速）
# 节点地区名为当地 机场三字码，指定多个时用英文逗号分隔

CloudflareST.exe -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD

# 注意，该参数只有在 HTTPing 延迟测速模式下才可用（因为要访问网页来获得）
```

> Cloudflare 所有节点地区名（机场三字码），请看：https://www.cloudflarestatus.com/

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

#### \# 测速其他端口

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

``` bash
# 如果你想要测速非默认 443 的其他端口，则需要通过 -tp 参数指定（该参数会影响 延迟测速/下载测速 时使用的端口）

# 如果要延迟测速 80 端口+下载测速（如果 -dd 禁用了下载测速则不需要），那么还需要指定 http:// 协议的下载测速地址才行（且该地址不会强制重定向至 HTTPS，因为那样就变成 443 端口了）
CloudflareST.exe -tp 80 -url http://cdn.cloudflare.steamstatic.com/steam/apps/5952/movie_max.webm

# 如果是非 80 443 的其他端口，那么需要确定你使用的下载测速地址是否支持通过该非标端口访问。
```

</details>

****

#### \# 自定义测速地址

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

``` bash
# 该参数适用于下载测速 及 HTTP 协议的延迟测速，对于后者该地址可以是任意网页 URL（不局限于具体文件地址）

# 地址要求：可以直接下载、文件大小超过 200MB、用的是 Cloudflare CDN
CloudflareST.exe -url https://cf.xiu2.xyz/url

# 注意：如果测速地址为 HTTP 协议（该地址不能强制重定向至 HTTPS），记得加上 -tp 80（这个参数会影响 延迟测速/下载测速 时使用的端口），如果是非 80 443 端口，那么需要确定下载测速地址是否支持通过该端口访问。
CloudflareST.exe -tp 80 -url http://cdn.cloudflare.steamstatic.com/steam/apps/5952/movie_max.webm
```

</details>

****

#### \# 自定义测速条件（指定 延迟/下载速度 的目标范围）

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

> 注意：延迟测速进度条右边的**可用数量**，仅指延迟测速过程中**未超时的 IP 数量**，和延迟上下限条件无关。

- 指定 **[平均延迟下限]** 条件

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

> 如果**没有找到一个满足延迟**条件的 IP，那么不会输出任何内容。

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

> 如果**没有找到一个满足速度**条件的 IP，那么会**忽略条件输出所有 IP 测速结果**（方便你下次测速时调整条件）。

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

> 如果**没有找到一个满足延迟**条件的 IP，那么不会输出任何内容。  
> 如果**没有找到一个满足速度**条件的 IP，那么会忽略条件输出所有 IP 测速结果（方便你下次测速时调整条件）。  
> 所以建议先不指定条件测速一遍，看看平均延迟和下载速度大概在什么范围，避免指定条件**过低/过高**！

> 因为Cloudflare 公开的 IP 段是**回源 IP+任播 IP**，而**回源 IP**是无法使用的，所以下载测速是 0.00。  
> 运行时可以加上 `-sl 0.01`（下载速度下限），过滤掉**回源 IP**（下载测速低于 0.01MB/s 的结果）。

</details>

****

#### \# 单独对一个或多个 IP 测速

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

**方式 一**：
直接通过参数指定要测速的 IP 段数据。
``` bash
# 先进入 CloudflareST 所在目录，然后运行：
# Windows 系统（在 CMD 中运行）
CloudflareST.exe -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32

# Linux 系统
./CloudflareST -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
```

****

**方式 二**：
或者把这些 IP 按如下格式写入到任意文本文件中，例如：`1.txt`

```
1.1.1.1
1.1.1.200
1.0.0.1/24
2606:4700::/32
```

> 单个 IP 的话可以省略 `/32` 子网掩码了（即 `1.1.1.1`等同于 `1.1.1.1/32`）。  
> 子网掩码 `/24` 指的是这个 IP 最后一段，即 `1.0.0.1~1.0.0.255`。


然后运行 CloudflareST 时加上启动参数 `-f 1.txt` 来指定 IP 段数据文件。

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

但就如 [**#8**](https://github.com/XIU2/CloudflareSpeedTest/issues/8) 所说，一个个添加域名到 Hosts 实在**太麻烦**了，于是我就找到了个**一劳永逸**的办法！可以看这个 [**还在一个个添加 Hosts？完美本地加速所有使用 Cloudflare CDN 的网站方法来了！**](https://github.com/XIU2/CloudflareSpeedTest/discussions/71) 和另一个[依靠本地 DNS 服务来修改域名解析 IP 为自选 IP](https://github.com/XIU2/CloudflareSpeedTest/discussions/317) 的教程。

****

#### \# 自动更新 Hosts

考虑到很多人获得最快 Cloudflare CDN IP 后，需要替换 Hosts 文件中的 IP。

可以看这个 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/discussions/312) 获取 **Windows/Linux 自动更新 Hosts 脚本**！

****

## 问题反馈

如果你遇到什么问题，可以先去 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues)、[Discussions](https://github.com/XIU2/CloudflareSpeedTest/discussions) 里看看是否有别人问过了（记得去看下  [**Closed**](https://github.com/XIU2/CloudflareSpeedTest/issues?q=is%3Aissue+is%3Aclosed) 的）。  
如果没找到类似问题，请新开个 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues/new) 来告诉我！

> **注意**！_与 `反馈问题、功能建议` 无关的，请前往项目内部 论坛 讨论（上面的 `💬 Discussions`_  

****

## 赞赏支持

![微信赞赏](https://cdn.staticaly.com/gh/XIU2/XIU2/master/img/zs-01.png)![支付宝赞赏](https://cdn.staticaly.com/gh/XIU2/XIU2/master/img/zs-02.png)

****

## 衍生项目

- _https://github.com/xianshenglu/cloudflare-ip-tester-app_  
_**CloudflareST 安卓版 APP [#202](https://github.com/XIU2/CloudflareSpeedTest/discussions/320)**_

- _https://github.com/mingxiaoyu/luci-app-cloudflarespeedtest_  
_**CloudflareST OpenWrt 路由器插件版 [#174](https://github.com/XIU2/CloudflareSpeedTest/discussions/319)**_

- _https://github.com/immortalwrt-collections/openwrt-cdnspeedtest_  
_**CloudflareST OpenWrt 原生编译版本 [#64](https://github.com/XIU2/CloudflareSpeedTest/discussions/64)**_

- _https://github.com/hoseinnikkhah/CloudflareSpeedTest-English_  
_**English language version of CloudflareST (Text language differences only) [#64](https://github.com/XIU2/CloudflareSpeedTest/issues/68)**_

> _此处仅收集了在本项目 Issues 中宣传过的部分 CloudflareST 相关衍生项目，如果有遗漏可以告诉我~_

****

## 感谢项目

- _https://github.com/Spedoske/CloudflareScanner_

> _因为该项目已经很长时间没更新了，而我又产生了很多功能需求，所以我临时学了下 Go 语言就上手了（菜）..._  
> _本软件基于该项目制作，但**已添加大量功能及修复 BUG**，并根据大家的使用反馈积极添加、优化功能（闲）..._

****

## 手动编译

<details>
<summary><code><strong>「 点击展开 查看内容 」</strong></code></summary>

****

为了方便，我是在编译的时候将版本号写入代码中的 version 变量，因此你手动编译时，需要像下面这样在 `go build` 命令后面加上 `-ldflags` 参数来指定版本号：

```bash
go build -ldflags "-s -w -X main.version=v2.3.3"
# 在 CloudflareSpeedTest 目录中通过命令行（例如 CMD、Bat 脚本）运行该命令，即可编译一个可在和当前设备同样系统、位数、架构的环境下运行的二进制程序（Go 会自动检测你的系统位数、架构）且版本号为 v2.3.3
```

如果想要在 Windows 64位系统下编译**其他系统、架构、位数**，那么需要指定 **GOOS** 和 **GOARCH** 变量。

例如在 Windows 系统下编译一个适用于 **Linux 系统 amd 架构 64 位**的二进制程序：

```bat
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w -X main.version=v2.3.3"
```

例如在 Linux 系统下编译一个适用于 **Windows 系统 amd 架构 32 位**的二进制程序：

```bash
GOOS=windows
GOARCH=386
go build -ldflags "-s -w -X main.version=v2.3.3"
```

> 可以运行 `go tool dist list` 来查看当前 Go 版本支持编译哪些组合。

****

当然，为了方便批量编译，我会专门指定一个变量为版本号，后续编译直接调用该版本号变量即可。  
同时，批量编译的话，还需要分开放到不同文件夹才行（或者文件名不同），需要加上 `-o` 参数指定。

```bat
:: Windows 系统下是这样：
SET version=v2.3.3
SET GOOS=linux
SET GOARCH=amd64
go build -o Releases\CloudflareST_linux_amd64\CloudflareST -ldflags "-s -w -X main.version=%version%"
```

```bash
# Linux 系统下是这样：
version=v2.3.3
GOOS=windows
GOARCH=386
go build -o Releases/CloudflareST_windows_386/CloudflareST.exe -ldflags "-s -w -X main.version=${version}"
```

</details>

****

## License

The GPL-3.0 License.