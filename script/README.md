# XIU2/CloudflareSpeedTest - Script

这里都是一些基于 **XIU2/CloudflareSpeedTest** 并**扩展更多功能**的脚本。  
有什么现有脚本功能上的建议可以告诉我，如果你有一些自用好用的脚本也可以通过 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues) 或 Pull requests 发给我添加到这里让更多人用到（会标注作者的~

> 小提示：点击↖左上角的三横杠图标按钮即可查看目录~

****
## 📑 cfst_hosts.sh / cfst_hosts.bat (已内置)

运行 CloudflareST 获得最快 IP 后，脚本会替换 Hosts 文件中的旧 CDN IP。

> **使用说明：https://github.com/XIU2/CloudflareSpeedTest/issues/42**

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

## 📑 cfst_3proxy.bat (已内置)

该脚本的作用为 CloudflareST 测速后获取最快 IP 并替换 3Proxy 配置文件中的 Cloudflare CDN IP。  
可以把所有 Cloudflare CDN IP 都重定向至最快 IP，实现一劳永逸的加速所有使用 Cloudflare CDN 的网站（不需要一个个添加域名到 Hosts 了）。

> **使用说明：https://github.com/XIU2/CloudflareSpeedTest/discussions/71**

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

## 📑 cfst_ddns.sh / cfst_ddns.bat

如果你的域名托管在 Cloudflare，则可以通过 Cloudflare 官方提供的 API 来自动更新域名解析记录！

> **使用说明：https://github.com/XIU2/CloudflareSpeedTest/issues/40**

<details>
<summary><code><strong>「 更新日志」</strong></code></summary>

****

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

如果你遇到什么问题，可以先去 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues) 里看看是否有别人问过了（记得去看下  [**Closed**](https://github.com/XIU2/CloudflareSpeedTest/issues?q=is%3Aissue+is%3Aclosed) 的）。  
如果没找到类似问题，请新开个 [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues/new) 来告诉我！