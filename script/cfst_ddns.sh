#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
# --------------------------------------------------------------
#	项目: CloudflareSpeedTest 自动更新域名解析记录
#	版本: 1.0.4
#	作者: XIU2
#	项目: https://github.com/XIU2/CloudflareSpeedTest
# --------------------------------------------------------------

_READ() {
	[[ ! -e "cfst_ddns.conf" ]] && echo -e "[错误] 配置文件不存在 [cfst_ddns.conf] !" && exit 1
	CONFIG=$(cat "cfst_ddns.conf")
	FOLDER=$(echo "${CONFIG}"|grep 'FOLDER='|awk -F '=' '{print $NF}')
	[[ -z "${FOLDER}" ]] && echo -e "[错误] 缺少配置项 [FOLDER] !" && exit 1
	ZONE_ID=$(echo "${CONFIG}"|grep 'ZONE_ID='|awk -F '=' '{print $NF}')
	[[ -z "${ZONE_ID}" ]] && echo -e "[错误] 缺少配置项 [ZONE_ID] !" && exit 1
	DNS_RECORDS_ID=$(echo "${CONFIG}"|grep 'DNS_RECORDS_ID='|awk -F '=' '{print $NF}')
	[[ -z "${DNS_RECORDS_ID}" ]] && echo -e "[错误] 缺少配置项 [DNS_RECORDS_ID] !" && exit 1
	EMAIL=$(echo "${CONFIG}"|grep 'EMAIL='|awk -F '=' '{print $NF}')
	[[ -z "${EMAIL}" ]] && echo -e "[错误] 缺少配置项 [EMAIL] !" && exit 1
	KEY=$(echo "${CONFIG}"|grep 'KEY='|awk -F '=' '{print $NF}')
	[[ -z "${KEY}" ]] && echo -e "[错误] 缺少配置项 [KEY] !" && exit 1
	TYPE=$(echo "${CONFIG}"|grep 'TYPE='|awk -F '=' '{print $NF}')
	[[ -z "${TYPE}" ]] && echo -e "[错误] 缺少配置项 [TYPE] !" && exit 1
	NAME=$(echo "${CONFIG}"|grep 'NAME='|awk -F '=' '{print $NF}')
	[[ -z "${NAME}" ]] && echo -e "[错误] 缺少配置项 [NAME] !" && exit 1
	TTL=$(echo "${CONFIG}"|grep 'TTL='|awk -F '=' '{print $NF}')
	[[ -z "${TTL}" ]] && echo -e "[错误] 缺少配置项 [TTL] !" && exit 1
	PROXIED=$(echo "${CONFIG}"|grep 'PROXIED='|awk -F '=' '{print $NF}')
	[[ -z "${PROXIED}" ]] && echo -e "[错误] 缺少配置项 [PROXIED] !" && exit 1
}

_UPDATE() {
	# 这里可以自己添加、修改 CloudflareST 的运行参数
	./CloudflareST -o "result_ddns.txt"

	# 判断结果文件是否存在，如果不存在说明结果为 0
	[[ ! -e "result_ddns.txt" ]] && echo "CloudflareST 测速结果 IP 数量为 0，跳过下面步骤..." && exit 0

	CONTENT=$(sed -n "2,1p" result_ddns.txt | awk -F, '{print $1}')
	if [[ -z "${CONTENT}" ]]; then
		echo "CloudflareST 测速结果 IP 数量为 0，跳过下面步骤..."
		exit 0
	fi
	curl -X PUT "https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records/${DNS_RECORDS_ID}" \
		-H "X-Auth-Email: ${EMAIL}" \
		-H "X-Auth-Key: ${KEY}" \
		-H "Content-Type: application/json" \
		--data "{\"type\":\"${TYPE}\",\"name\":\"${NAME}\",\"content\":\"${CONTENT}\",\"ttl\":${TTL},\"proxied\":${PROXIED}}"
}

_READ
cd "${FOLDER}"
_UPDATE