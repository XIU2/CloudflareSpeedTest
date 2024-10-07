#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
# --------------------------------------------------------------
#	项目: CloudflareSpeedTest 自动更新域名解析记录
#	版本: 1.0.6
#	作者: XIU2(大部分代码) & yyyyt(少部分)
#	项目: https://github.com/XIU2/CloudflareSpeedTest
# --------------------------------------------------------------
# 如果你想要指定自定义配置文件的话请将下面的"cfst_ddns.conf"换成对应的文件名
# 多子域名的已经默认注释,如果你想更新多子域名的话请在配置文件指定NAME2和DNS_RECORDS_ID2并取消掉对应的注释(_READ的和下面curl的)
# 如果你想要v4 v6都优选的话请将脚本文件复制一份并修改为指定配置文件~~(其实一个脚本也是可以的,但是我懒得改2333)~~,配置文件填上对应的内容 测速时候分别去运行两个脚本

_READ() {
	[[ ! -e "cfst_ddns.conf" ]] && echo -e "[错误] 配置文件不存在 [cfst_ddns.conf] !" && exit 1
	CONFIG=$(cat "cfst_ddns.conf")
	FOLDER=$(echo "${CONFIG}"|grep 'FOLDER='|awk -F '=' '{print $NF}')
	[[ -z "${FOLDER}" ]] && echo -e "[错误] 缺少配置项 [FOLDER] !" && exit 1
	ZONE_ID=$(echo "${CONFIG}"|grep 'ZONE_ID='|awk -F '=' '{print $NF}')
	[[ -z "${ZONE_ID}" ]] && echo -e "[错误] 缺少配置项 [ZONE_ID] !" && exit 1
	DNS_RECORDS_ID=$(echo "${CONFIG}"|grep 'DNS_RECORDS_ID='|awk -F '=' '{print $NF}')
	[[ -z "${DNS_RECORDS_ID}" ]] && echo -e "[错误] 缺少配置项 [DNS_RECORDS_ID] !" && exit 1
#	DNS_RECORDS_ID2=$(echo "${CONFIG}"|grep 'DNS_RECORDS_ID2='|awk -F '=' '{print $NF}')
#	[[ -z "${DNS_RECORDS_ID2}" ]] && echo -e "[错误] 缺少配置项 [DNS_RECORDS_ID2] !" && exit 1
	KEY=$(echo "${CONFIG}"|grep 'KEY='|awk -F '=' '{print $NF}')
	[[ -z "${KEY}" ]] && echo -e "[错误] 缺少配置项 [KEY] !" && exit 1
	EMAIL=$(echo "${CONFIG}"|grep 'EMAIL='|awk -F '=' '{print $NF}')
	[[ -z "${EMAIL}" ]] && echo -e "[信息] 缺少配置项 [EMAIL]，由 [API 密钥] 方式转为 [API 令牌] 方式!"
	TYPE=$(echo "${CONFIG}"|grep 'TYPE='|awk -F '=' '{print $NF}')
	[[ -z "${TYPE}" ]] && echo -e "[错误] 缺少配置项 [TYPE] !" && exit 1
	NAME=$(echo "${CONFIG}"|grep 'NAME='|awk -F '=' '{print $NF}')
	[[ -z "${NAME}" ]] && echo -e "[错误] 缺少配置项 [NAME] !" && exit 1
#	NAME2=$(echo "${CONFIG}"|grep 'NAME2='|awk -F '=' '{print $NF}')
#	[[ -z "${NAME2}" ]] && echo -e "[错误] 缺少配置项 [NAME2] !" && exit 1
	TTL=$(echo "${CONFIG}"|grep 'TTL='|awk -F '=' '{print $NF}')
	[[ -z "${TTL}" ]] && echo -e "[错误] 缺少配置项 [TTL] !" && exit 1
	PROXIED=$(echo "${CONFIG}"|grep 'PROXIED='|awk -F '=' '{print $NF}')
	[[ -z "${PROXIED}" ]] && echo -e "[错误] 缺少配置项 [PROXIED] !" && exit 1
    RESULT=$(echo "${CONFIG}"|grep 'RESULT='|awk -F '=' '{print $NF}')
    [[ -z "${RESULT}" ]] && echo -e "[错误] 缺少配置项 [RESULT] !" && exit 1
}

_UPDATE() {
	# 执行前删除原有的测速结果
	rm -rf "${RESULT}"

	# 这里可以自己添加、修改 CloudflareST 的运行参数.如果你想要测速非cfcdn的节点/IPv6的话请指定-f参数
	./CloudflareST -o "${RESULT}" -url https://download.parallels.com/desktop/v17/17.1.1-51537/ParallelsDesktop-17.1.1-51537.dmg

	# 判断结果文件是否存在，如果不存在说明结果为 0
	[[ ! -e "${RESULT}" ]] && echo "CloudflareST 测速结果 IP 数量为 0，跳过下面步骤..." && exit 0

	CONTENT=$(sed -n "2,1p" "${RESULT}" | awk -F, '{print $1}')
	if [[ -z "${CONTENT}" ]]; then
		echo "CloudflareST 测速结果 IP 数量为 0，跳过下面步骤..."
		exit 0
	fi
	# 如果 EMAIL 变量是空的，那么就代表要使用 API 令牌方式
	if [[ -n "${EMAIL}" ]]; then
		# API 密钥方式（全局权限）
		curl -X PUT "https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records/${DNS_RECORDS_ID}" \
			-H "X-Auth-Email: ${EMAIL}" \
			-H "X-Auth-Key: ${KEY}" \
			-H "Content-Type: application/json" \
			--data "{\"type\":\"${TYPE}\",\"name\":\"${NAME}\",\"content\":\"${CONTENT}\",\"ttl\":${TTL},\"proxied\":${PROXIED}}"
#		curl -X PUT "https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records/${DNS_RECORDS_ID2}" \
#			-H "X-Auth-Email: ${EMAIL}" \
#			-H "X-Auth-Key: ${KEY}" \
#			-H "Content-Type: application/json" \
#			--data "{\"type\":\"${TYPE}\",\"name\":\"${NAME2}\",\"content\":\"${CONTENT}\",\"ttl\":${TTL},\"proxied\":${PROXIED}}"
	else
		# API 令牌方式（自定义权限）
		curl -X PUT "https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records/${DNS_RECORDS_ID}" \
			-H "Authorization: Bearer ${KEY}" \
			-H "Content-Type: application/json" \
			--data "{\"type\":\"${TYPE}\",\"name\":\"${NAME}\",\"content\":\"${CONTENT}\",\"ttl\":${TTL},\"proxied\":${PROXIED}}"
#		curl -X PUT "https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records/${DNS_RECORDS_ID2}" \
#			-H "Authorization: Bearer ${KEY}" \
#			-H "Content-Type: application/json" \
#			--data "{\"type\":\"${TYPE}\",\"name\":\"${NAME2}\",\"content\":\"${CONTENT}\",\"ttl\":${TTL},\"proxied\":${PROXIED}}"
	fi
}

_READ
cd "${FOLDER}"
_UPDATE
