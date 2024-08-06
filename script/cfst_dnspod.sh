#!/bin/bash

# --------------------------------------------------------------
#	项目: CloudflareSpeedTest 自动更新Dnspod优选解析
#	版本: 1.0.0
#	作者: imashen
# --------------------------------------------------------------

# 清理历史残留
rm -f result4.csv result6.csv
# DNSPod API 凭据
dnspod_token="${API_TOKEN}"
dnspod_domain="${DOMAIN}"
dnspod_record="${SUB_DOMAIN}"

# DNSPod API URL
dnspod_api_url="https://dnsapi.cn"

# 获取记录 ID
get_record_id() {
    local record_type=$1
    local response
    response=$(curl -s -X POST -d "login_token=$dnspod_token&format=json&domain=$dnspod_domain&record_type=$record_type" "$dnspod_api_url/Record.List")
    local record_id
    record_id=$(echo "$response" | jq -r --arg type "$record_type" '.records[] | select(.type == $type) | .id')
    echo "$record_id"
}

# 创建 DNS 记录
create_dns_record() {
    local record_type=$1
    local ip_address=$2
    local response
    response=$(curl -s -X POST -d "login_token=$dnspod_token&format=json&domain=$dnspod_domain&sub_domain=$dnspod_record&record_type=$record_type&record_line=默认&value=$ip_address" "$dnspod_api_url/Record.Create")
    local record_id
    record_id=$(echo "$response" | jq -r '.record.id')
    echo "$record_id"
}

# 更新 DNS 记录
update_dns_record() {
    local record_id=$1
    local record_type=$2
    local ip_address=$3
    curl -s -X POST -d "login_token=$dnspod_token&format=json&domain=$dnspod_domain&record_id=$record_id&sub_domain=$dnspod_record&record_type=$record_type&record_line=默认&value=$ip_address" "$dnspod_api_url/Record.Modify"
}

# 运行 CloudflareST v4
./CloudflareST -f ip.txt -n 500 -o result4.csv

# 读取 CSV 文件并提取优选 IPv4 地址
preferred_ipv4=$(awk -F, 'NR==2 {print $1}' result4.csv)

# 检查是否获取到了 IPv4 地址
if [ -z "$preferred_ipv4" ]; then
  echo "Failed to get the preferred IPv4 address."
else
  echo "BETTER IPv4: $preferred_ipv4"

  # 获取 IPv4 记录 ID
  ipv4_record_id=$(get_record_id "A")

  if [ -n "$ipv4_record_id" ]; then
    # 更新 IPv4 记录
    update_dns_record "$ipv4_record_id" "A" "$preferred_ipv4"
    echo "Updated DNSPod record with IPv4: $preferred_ipv4"
  else
    # 创建 IPv4 记录
    new_ipv4_record_id=$(create_dns_record "A" "$preferred_ipv4")
    if [ -n "$new_ipv4_record_id" ]; then
      echo "Created DNSPod record with IPv4: $preferred_ipv4"
    else
      echo "Failed to create DNSPod record with IPv4."
    fi
  fi
fi

# 运行 CloudflareST v6
./CloudflareST -f ipv6.txt -n 500 -o result6.csv

# 读取 CSV 文件并提取优选 IPv6 地址
preferred_ipv6=$(awk -F, 'NR==2 {print $1}' result6.csv)

# 检查是否获取到了 IPv6 地址
if [ -z "$preferred_ipv6" ]; then
  echo "Failed to get the preferred IPv6 address."
else
  echo "BETTER IPv6: $preferred_ipv6"

  # 获取 IPv6 记录 ID
  ipv6_record_id=$(get_record_id "AAAA")

  if [ -n "$ipv6_record_id" ]; then
    # 更新 IPv6 记录
    update_dns_record "$ipv6_record_id" "AAAA" "$preferred_ipv6"
    echo "Updated DNSPod record with IPv6: $preferred_ipv6"
  else
    # 创建 IPv6 记录
    new_ipv6_record_id=$(create_dns_record "AAAA" "$preferred_ipv6")
    if [ -n "$new_ipv6_record_id" ]; then
      echo "Created DNSPod record with IPv6: $preferred_ipv6"
    else
      echo "Failed to create DNSPod record with IPv6."
    fi
  fi
fi
