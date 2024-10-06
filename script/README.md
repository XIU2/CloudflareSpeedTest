# XIU2/CloudflareSpeedTest - Script(脚本)

这里都是一些通过调用 **CloudflareST** 并**扩展实现更多个性化功能**的脚本。  

****
> [!TIP]
> 我之所以将 CloudflareST 制作为一个**命令行程序**，就是考虑到**通用性**，因为毕竟不可能把所有需求都塞到软件内（特别是一些**个性化、小众**的需求），这样增加维护难度和精力不说，还会导致软件异常臃肿（`“变成我讨厌的样子”`），而命令行程序的优势之一就在于**可以很方便的和其他软件、脚本搭配使用**。

比如像下面这些我写的几个脚本，就是把一些需求以外置脚本方式实现。  

> 即脚本调用 CloudflareST 测速并获取结果，然后***按照自己的需求自由决定***如何处理得到的测速结果（比如修改 Hosts 等）。  

总的来说，我写的这几个脚本都比较简单，功能也很单一，除了满足部分用户的需求外，***更像是一个 CloudflareST 与脚本搭配使用的示例参考***，对于一些会写脚本、软件的用户来说，完全可以**自给自足**来实现一些个性化需求。

当然，如果你有一些自用好用的脚本也可以通过 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues)、[**Discussions**](https://github.com/XIU2/CloudflareSpeedTest/discussions) 或 **Pull requests** 发给我添加到这里让更多人用到！

> 小提示：点击↗右上角的三横杠图标按钮即可查看目录~

****
## 📑 cfst_hosts.sh / cfst_hosts.bat (已内置压缩包)

脚本会运行 CloudflareST 获得最快 IP，并替换掉 Hosts 文件中的旧 CDN IP。

> **作者：**[@XIU2](https://github.com/xiu2)  
> **使用说明/问题反馈：https://github.com/XIU2/CloudflareSpeedTest/discussions/312**

<details>
<summary><code><strong>「 更新日志」</strong></code></summary>

****

#### 2021年12月17日，版本 v1.0.6
 - **1. 优化** [找不到满足条件的 IP 就一直循环测速] 功能，在指定下载测速下限时没有重新测速的问题（默认注释）   

#### 2021年12月17日，版本 v1.0.3
 - **1. 新增** 找不到满足条件的 IP 就一直循环测速功能（默认注释）  
 - **2. 优化** 代码  

#### 2021年09月29日，版本 v1.0.2
 - **1. 修复** 当测速结果 IP 数量为 0 时，脚本没有退出的问题  

#### 2021年04月29日，版本 v1.0.1
 - **1. 优化** 不再需要加上 -p 0 参数来避免回车键退出了（现在可以即显示结果，又不用担心回车键退出程序）  

#### 2021年01月28日，版本 v1.0.0
 - **1. 发布** 第一个版本  

</details>

****

## 📑 cfst_3proxy.bat (已内置压缩包)

脚本会运行 CloudflareST 测速后获取最快 IP 并替换 3Proxy 配置文件中的旧 Cloudflare CDN IP。  
可以把所有 Cloudflare CDN IP 都重定向至最快 IP，实现一劳永逸的加速所有使用 Cloudflare CDN 的网站（不需要一个个添加域名到 Hosts 了）。

> **作者：**[@XIU2](https://github.com/xiu2)  
> **使用说明/问题反馈：https://github.com/XIU2/CloudflareSpeedTest/discussions/71**

<details>
<summary><code><strong>「 更新日志」</strong></code></summary>

****

#### 2021年12月17日，版本 v1.0.5
 - **1. 优化** [找不到满足条件的 IP 就一直循环测速] 功能，在指定下载测速下限时没有重新测速的问题（默认注释）   

#### 2021年12月17日，版本 v1.0.4
 - **1. 新增** 找不到满足条件的 IP 就一直循环测速功能（默认注释）  
 - **2. 优化** 代码  

#### 2021年09月29日，版本 v1.0.3
 - **1. 修复** 当测速结果 IP 数量为 0 时，脚本没有退出的问题  

#### 2021年04月29日，版本 v1.0.2
 - **1. 优化** 不再需要加上 -p 0 参数来避免回车键退出了（现在可以即显示结果，又不用担心回车键退出程序）  

#### 2021年03月16日，版本 v1.0.1
 - **1. 优化** 代码及注释内容  

#### 2021年03月13日，版本 v1.0.0
 - **1. 发布** 第一个版本  

</details>

****

## 📑 cfst_dnspod.sh

如果你的域名托管在 **dnspod**，则可以通过 dnspod 官方提供的 API 来自动更新域名解析记录！  
脚本会运行 CloudflareST 测速获得最快 IP，并通过 Cloudflare API 来更新域名解析记录为这个最快 IP。

> **作者：**[@imashen](https://github.com/imashen)  
> **使用说明/问题反馈：https://github.com/XIU2/CloudflareSpeedTest/pull/533**

<details>
<summary><code><strong>「 更新日志」</strong></code></summary>

****

#### 2024年08月06日，版本 v1.0.0
 - **1. 发布** 第一个版本  

</details>

****

## 📑 cfst_ddns.sh / cfst_ddns.bat

如果你的域名托管在 **Cloudflare**，则可以通过 Cloudflare 官方提供的 API 来自动更新域名解析记录！  
脚本会运行 CloudflareST 测速获得最快 IP，并通过 Cloudflare API 来更新域名解析记录为这个最快 IP。

> **作者：**[@XIU2](https://github.com/xiu2)  
> **使用说明/问题反馈：https://github.com/XIU2/CloudflareSpeedTest/discussions/481**

<details>
<summary><code><strong>「 更新日志」</strong></code></summary>

****

#### 2024年10月06日，版本 v1.0.5
 - **1. 新增** 支持 API 令牌方式（相比 API 密钥这种全局权限的，API 令牌可以自由控制权限）   

#### 2021年12月17日，版本 v1.0.4
 - **1. 新增** 找不到满足条件的 IP 就一直循环测速功能（默认注释）  
 - **2. 优化** 代码  

#### 2021年09月29日，版本 v1.0.3
 - **1. 修复** 当测速结果 IP 数量为 0 时，脚本没有退出的问题  

#### 2021年04月29日，版本 v1.0.2
 - **1. 优化** 不再需要加上 -p 0 参数来避免回车键退出了（现在可以即显示结果，又不用担心回车键退出程序）  

#### 2021年01月27日，版本 v1.0.1
 - **1. 优化** 配置从文件中读取  

#### 2021年01月26日，版本 v1.0.0
 - **1. 发布** 第一个版本  

</details>

****

## 功能建议/问题反馈

如果这些脚本使用过程中你遇到了什么问题，可以先去脚本对应的 **`使用说明`** 帖子里看看是否有别人问过了。  
如果没找到类似问题，那么就在脚本对应的 **`使用说明`** 帖子里直接评论问作者吧。