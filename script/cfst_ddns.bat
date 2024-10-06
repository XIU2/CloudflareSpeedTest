:: --------------------------------------------------------------
::	项目: CloudflareSpeedTest 自动更新域名解析记录
::	版本: 1.0.5
::	作者: XIU2
::	项目: https://github.com/XIU2/CloudflareSpeedTest
:: --------------------------------------------------------------
@echo off
Setlocal Enabledelayedexpansion

:: 这里可以自己添加、修改 CloudflareST 的运行参数，echo.| 的作用是自动回车退出程序（不再需要加上 -p 0 参数了）
echo.|CloudflareST.exe -o "result_ddns.txt"

:: 判断结果文件是否存在，如果不存在说明结果为 0
if not exist result_ddns.txt (
    echo.
    echo CloudflareST 测速结果 IP 数量为 0，跳过下面步骤...
    goto :END
)

for /f "tokens=1 delims=," %%i in (result_ddns.txt) do (
    Set /a n+=1 
    If !n!==2 (
        Echo %%i
        if "%%i"=="" (
            echo.
            echo CloudflareST 测速结果 IP 数量为 0，跳过下面步骤...
            goto :END
        )
::      API 密钥方式（全局权限）
        curl -X PUT "https://api.cloudflare.com/client/v4/zones/域名ID/dns_records/域名解析记录ID" ^
                -H "X-Auth-Email: 账号邮箱" ^
                -H "X-Auth-Key: 前面获取的 API 密钥" ^
                -H "Content-Type: application/json" ^
                --data "{\"type\":\"A\",\"name\":\"完整域名\",\"content\":\"%%i\",\"ttl\":1,\"proxied\":true}"
::        API 令牌方式（自定义权限），如果要使用这种方式，可以把上面的删除或注释，然后把下面的行首 "::" 注释符删除即可。
::        curl -X PUT "https://api.cloudflare.com/client/v4/zones/域名ID/dns_records/域名解析记录ID" ^
::                -H "Authorization: Bearer 前面获取的 API 令牌" ^
::                -H "Content-Type: application/json" ^
::                --data "{\"type\":\"A\",\"name\":\"完整域名\",\"content\":\"%%i\",\"ttl\":1,\"proxied\":true}"

        goto :END
    )
)
:END
pause