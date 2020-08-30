# CloudflareSpeedTest

国外很多网站都在使用 Cloudflare CDN，但分配给中国访客的 IP 并不友好。  
虽然 Cloudflare 公开了所有 [IP 段](https://www.cloudflare.com/ips/) ，但想要在这么多 IP 中找到适合自己的，怕是要累死，所以就有了这个软件。  

该软件可以**测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最佳 IP**！  
你可以将 IP 添加到 `Hosts` 文件中，来帮你提高访问使用 Cloudflare CDN 服务的网站速度！  

****
### 快速使用

1. 下载编译好的可执行文件 [蓝奏云](https://www.lanzoux.com/b0742hkxe) / [Github](https://github.com/XIU2/CloudflareSpeedTest/releases) 并解压。  
2. 双击运行 `CloudflareST.exe`文件（Windows系统），等待测速...  

测速完毕后，会把结果保存在当前目录下的 `result.csv` 文件中，用记事本打开，按照**延迟由低到高排序**。  

****
### 进阶使用

直接双击运行使用的是默认参数，如果想要测试速度更快、测试结果更全面，可以自定义参数。  
``` cmd
C:\>CloudflareST.exe -h
CloudflareSpeedTest
    测试 Cloudflare CDN 所有 IP 的延迟和速度，获取最佳 IP！
    https://github.com/XIU2/CloudflareSpeedTest

参数：
    -n 500
        测速线程数量；请勿超过1000 (default 500)
    -t 4
        延迟测速次数；单个 IP (default 4)
    -dn 10
        下载测速数量；延迟测速后，从最低延迟起测试下载速度的数量，请勿太多 (default 10)
    -dt 10
        下载测试时间；单个 IP 测速最长时间，单位：秒 (default 10)
    -v
        打印程序版本
    -h
        打印帮助说明

示例：
    Windows：CloudflareST.exe -n 800 -t 4 -dn 10 -dt 10
    Linux：CloudflareST -n 800 -t 4 -dn 10 -dt 10
```

## 使用示例

在 CMD 中运行，或者把启动参数添加到快捷方式中。  
> 注意：不需要四个参数都加上，如果你认为某个参数默认就很合适，那就跳过。
``` cmd
# CMD 示例
CloudflareST.exe -n 800 -t 4 -dn 10 -dt 10
```
``` cmd
# 快捷方式示例（右键快捷方式 - 目标）
## 如果有引号就放在引号外面，记得引号和 - 之间有空格。
"D:\Program Files\CloudflareST\CloudflareST.exe" -n 800 -t 4 -dn 10 -dt 10
```

****
### 感谢项目
https://github.com/Spedoske/CloudflareScanner

****
### 许可证
The GPL-3.0 License.
