# wsnjh007/CloudflareSpeedTest

### 项目介绍
1. 名称: CloudflareSpeedTest 自动更新 Hosts
2. 版本: 自我感觉良好，命名为2.0🤪
3. 作者: wsnjh007（参考XIU2的作品修改）
4. 感谢: 该项目是参考XIU2的作品进行修改，特别是他写的CloudflareST程序，[详见他的作品](https://github.com/XIU2/CloudflareSpeedTest)
 

### 下载基础程序文件到[XIU2的作品](https://github.com/XIU2/CloudflareSpeedTest)那边下载，他有详细的说明，我就不做赘述了。

### 具体说一下我的改进思路和一些其他说明
1. 我的hosts文件中有许多比较常用的域名，用梯子麻烦也有风险
2. 按照XIU2的思路我要每个域名运行一次，很慢麻烦，所以就自己在他的代码做了修改
3. 我提供要测速的域名列表，自动抓取hosts文件对应的IP，然后逐一进行测试，然后自动替换最优IP
4. 这个脚本最大的亮点就是不用手动输入域名IP，也不用管现在的现在hosts文件中的IP，只需在域名列表中添加自己需要替换最优IP的域名即可，其他的都是自动化
5. 下面是我改进的脚本代码，只要替换[XIU2作品](https://github.com/XIU2/CloudflareSpeedTest)中的cfst_hosts.sh文件即可
6. 同时对在终端中显示效果做了一些美化，希望大家喜欢
