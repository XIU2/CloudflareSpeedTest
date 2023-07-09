# wsnjh007/CloudflareSpeedTest

### 项目介绍
1. 名称: CloudflareSpeedTest 自动更新 Hosts
2. 版本: 自我感觉良好，命名为2.0🤪
3. 作者: wsnjh007（参考XIU2的作品修改）
4. 感谢: 该项目是参考XIU2的作品进行修改，特别是他写的CloudflareST程序，[详见他的作品](https://github.com/XIU2/CloudflareSpeedTest)
 

### 下载基础程序文件到[XIU2的作品](https://github.com/XIU2/CloudflareSpeedTest)那边下载，他有详细的说明，我就不做赘述了。

### 具体说一下我的改进思路和一些其他说明
1. 我的hosts文件中有许多比较常用的域名，用梯子麻烦也有风险
2. 按照原来的方法我要每个域名运行一次，很慢麻烦，所以就自己在XIU2的代码做了修改
3. 我提供要测速的域名列表，自动抓取hosts文件对应的IP，然后逐一进行测试，然后自动替换最优IP
4. 这个脚本最大的亮点就是不用手动输入域名IP，也不用管现在的现在hosts文件中的IP，只需在域名列表中添加自己需要替换最优IP的域名即可，其他的都是自动化
5. 下面是我改进的脚本代码，只要替换[XIU2作品](https://github.com/XIU2/CloudflareSpeedTest)中的cfst_hosts.sh文件即可
6. 同时对在终端中显示效果做了一些美化，希望大家喜欢


### 具体代码如下：
``` bash
#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
set -euo pipefail
# --------------------------------------------------------------
#   项目: CloudflareSpeedTest 自动更新 Hosts
#   版本: 2.0
#   作者: wsnjh007（参考XIU2的作品修改）
#   感谢: 该项目是参考XIU2的作品进行修改，特别是他写的CloudflareST程序，他的作品详见https://github.com/XIU2/CloudflareSpeedTest
# --------------------------------------------------------------

# 定义文件名
# CloudflareST 程序的文件名
CLOUDFLAREST="CloudflareST"
# 存储当前 IP 地址的文件名
NOWIP_FILE="nowip_hosts.txt"
# 存储测速结果的文件名
RESULT_FILE="result_hosts.txt"
# 你要测速的域名列表
DOMAINS=("www.themoviedb.org" "image.tmdb.org" "api.thetvdb.com" "api.tmdb.org" "api.themoviedb.org" "tmdb.nastool.org" "t.nastool.workers.dev" "webservice.fanart.tv" "api.telegram.org" "www.opensubtitles.org" "dicmusic.club")

# 用于输出的颜色代码
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
NC="\033[0m" # No Color

# 定义三个关联数组来分别存储每个域名的IP和替换状态
declare -A domain_ips
declare -A domain_old_ips
declare -A domain_status

# 定义一个函数_BACKUP_HOSTS，用于备份当前的hosts文件。
_BACKUP_HOSTS() {
    echo -e "\n${RED}----------------------------------------------------------------------${NC}"
    echo -e "${RED}该脚本将在 hosts 文件中查找每个指定域名的当前 IP，并进行测速和替换操作${NC}"
    echo -e "${RED}----------------------------------------------------------------------${NC}"
    if [ -w "/etc/hosts" ]; then
        TIMESTAMP=$(date +"%Y%m%d%H%M%S")
        echo -e "\n开始备份 Hosts 文件..."
        # 将"/volume1/Tool/04_code/hostsupdate"改成自己的路径
        cp -f /etc/hosts "/volume1/Tool/04_code/hostsupdate/hosts_backup_${TIMESTAMP}"
        echo -e "${GREEN}完成...${NC}"
    else
        echo -e "${YELLOW}无法备份 Hosts 文件，检查你的权限${NC}"
        exit 1
    fi
}

# 获取所有域名的当前IP的函数
_INIT_IP() {
    for DOMAIN in "${DOMAINS[@]}"; do
        # 通过对 grep 命令的结果进行进一步处理，确保只取出一个 IP。
        NOWIP=$(grep ${DOMAIN} /etc/hosts | head -n 1 | awk '{print $1}')
        if [[ ! -z "${NOWIP}" ]]; then
            domain_ips["${DOMAIN}"]=${NOWIP}  # 存储当前IP到关联数组中
            domain_old_ips["${DOMAIN}"]=${NOWIP}  # 存储旧的IP到关联数组中
        else
            echo -e "${RED}在 hosts 文件中未找到 ${DOMAIN} 的 IP 地址，跳过此域名的测速和替换操作${NC}"
        fi
    done
}


# 进行测速并更新hosts文件
_UPDATE() {
    # 删除之前运行的结果文件，保证每次运行都是新的结果
    if [ -f "${RESULT_FILE}" ]; then
        rm -f "${RESULT_FILE}"
    fi
    if [ -f "${NOWIP_FILE}" ]; then
        rm -f "${NOWIP_FILE}"
    fi

    # 对每个在DOMAINS列表中的域名进行处理
    for DOMAIN in "${DOMAINS[@]}"; do
        # 从我们之前存储的关联数组中获取当前域名的IP
        NOWIP=${domain_ips["${DOMAIN}"]}

        # 输出开始测速的信息
        echo -e "\n------------------------------------------------------------"
        echo -e "${GREEN}【${DOMAIN}】${NC}开始测速..."

        # 使用CloudflareST工具进行测速，将结果存储在临时文件中
        ./"${CLOUDFLAREST}" -o "temp_result_hosts.txt"


        # 从测速结果中获取最快的IP地址
        BESTIP=$(sed -n "2,1p" temp_result_hosts.txt | awk -F, '{print $1}')


        # 如果没有获取到最快的IP，则跳过后续的步骤，进行下一个域名的处理
        if [[ -z "${BESTIP}" ]]; then
            echo "CloudflareST 测速结果 IP 数量为 0，跳过下面步骤..."
            continue
        fi

        # 输出旧的IP和新的IP
        echo -e "\n旧 IP 为 ${NOWIP}\n新 IP 为 ${BESTIP}\n"

        # 输出开始替换的信息
        echo -e "开始替换..."

        # 在hosts文件中将旧的IP替换为新的IP
        sed -i "s#${NOWIP}#${BESTIP}#g" /etc/hosts

        # 检查替换是否成功
        if grep -q "${BESTIP}" /etc/hosts; then
            # 如果替换成功，将新的IP和域名写入到结果文件中，并存储替换状态到关联数组中
            echo ${BESTIP} ${DOMAIN} >> "${NOWIP_FILE}"
            echo -e "${DOMAIN} 替换成功\n旧 IP 为 ${NOWIP}\n新 IP 为 ${BESTIP}" >> "${RESULT_FILE}"
            domain_status["${DOMAIN}"]="替换成功"
            domain_ips["${DOMAIN}"]=${BESTIP}  # 更新当前域名的 IP 地址
            echo -e "\n\n" >> "${RESULT_FILE}"
            echo -e "${GREEN}完成...${NC}\n"
        else
            # 如果替换失败，将旧的IP和域名写入到结果文件中，并存储替换状态到关联数组中
            NOWIP=$(awk '/'${DOMAIN}'/{print $1}' /etc/hosts)
            echo ${NOWIP} ${DOMAIN} >> "${NOWIP_FILE}"
            echo -e "${DOMAIN} 替换失败 \n当前 IP 为 ${NOWIP}" >> "${RESULT_FILE}"
            domain_status["${DOMAIN}"]="替换失败"
            echo -e "\n\n" >> "${RESULT_FILE}"
            echo -e "${RED}失败...${NC}${YELLOW}  可能原因包括没有正确的文件权限，或者hosts文件不存在。${NC}\n"
        fi

        # 删除临时的测速结果文件
        if [ -f "temp_result_hosts.txt" ]; then
            echo -e "开始删除 temp_result_hosts.txt 文件..."
            rm -f "temp_result_hosts.txt"
            echo -e "${GREEN}完成...${NC}\n \n"
        else
            echo -e "${YELLOW}未找到 temp_result_hosts.txt 文件${NC}\n \n"
        fi
    done
}

# 显示结果的函数
_DISPLAY_RESULT() {
    echo -e "\n${YELLOW}下面是本次运行的结果汇总：${NC}"
    echo -e "-----------------------------------------------------------"
    for DOMAIN in "${DOMAINS[@]}"; do
        OLDIP=${domain_old_ips["${DOMAIN}"]}
        NOWIP=${domain_ips["${DOMAIN}"]}
        echo -e "域名：${GREEN}${DOMAIN}${NC}"
        echo -e "旧 IP：${OLDIP}"
        if [[ ${domain_status["${DOMAIN}"]} == "替换成功" ]]; then
            echo -e "新 IP：${NOWIP}"  # 引用新的 IP 地址
            echo -e "${GREEN}替换成功${NC}\n"
        elif [[ ${domain_status["${DOMAIN}"]} == "替换失败" ]]; then
            echo -e "${RED}替换失败${NC}\n"
        fi
    done
    echo -e "-----------------------------------------------------------"
    echo -e "${NC}"  # 结束时重置颜色。
}


# 定义一个函数MAIN，用于执行上面定义的所有函数。
_MAIN() {
    _BACKUP_HOSTS
    _INIT_IP
    _UPDATE
    _DISPLAY_RESULT
}

# 执行_MAIN函数
_MAIN

```
