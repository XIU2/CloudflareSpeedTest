:: --------------------------------------------------------------
::	项目: CloudflareSpeedTest 自动更新 Hosts
::	版本: 1.0.4
::	作者: XIU2
::	项目: https://github.com/XIU2/CloudflareSpeedTest
:: --------------------------------------------------------------
@echo off
Setlocal Enabledelayedexpansion

::判断是否已获得管理员权限

>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system" 

if '%errorlevel%' NEQ '0' (  
    goto UACPrompt  
) else ( goto gotAdmin )  

::写出 vbs 脚本以管理员身份运行本脚本（bat）

:UACPrompt  
    echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs" 
    echo UAC.ShellExecute "%~s0", "", "", "runas", 1 >> "%temp%\getadmin.vbs" 
    "%temp%\getadmin.vbs" 
    exit /B  

::如果临时 vbs 脚本存在，则删除
  
:gotAdmin  
    if exist "%temp%\getadmin.vbs" ( del "%temp%\getadmin.vbs" )  
    pushd "%CD%" 
    CD /D "%~dp0" 


::上面是判断是否以获得管理员权限，如果没有就去获取，下面才是本脚本主要代码


::如果 nowip_hosts.txt 文件不存在，说明是第一次运行该脚本
if not exist "nowip_hosts.txt" (
    echo 该脚本的作用为 CloudflareST 测速后获取最快 IP 并替换 Hosts 中的 Cloudflare CDN IP。
    echo 使用前请先阅读：https://github.com/XIU2/CloudflareSpeedTest/issues/42#issuecomment-768273768
    echo.
    echo 第一次使用，请先将 Hosts 中所有 Cloudflare CDN IP 统一改为一个 IP。
    set /p nowip="输入该 Cloudflare CDN IP 并回车（后续不再需要该步骤）:"
    echo !nowip!>nowip_hosts.txt
    echo.
)  

::从 nowip_hosts.txt 文件获取当前 Hosts 中使用的 Cloudflare CDN IP
set /p nowip=<nowip_hosts.txt
echo 开始测速...


:: 这个 RESET 是给需要 "找不到满足条件的 IP 就一直循环测速下去" 功能的人准备的
:: 如果需要这个功能就把下面 3 个 goto :STOP 改为 goto :RESET 即可
:RESET


:: 这里可以自己添加、修改 CloudflareST 的运行参数，echo.| 的作用是自动回车退出程序（不再需要加上 -p 0 参数了）
echo.|CloudflareST.exe -o "result_hosts.txt"


:: 判断结果文件是否存在，如果不存在说明结果为 0
if not exist result_hosts.txt (
    echo.
    echo CloudflareST 测速结果 IP 数量为 0，跳过下面步骤...
    goto :STOP
)

:: 获取第一行的最快 IP
for /f "tokens=1 delims=," %%i in (result_hosts.txt) do (
    SET /a n+=1 
    If !n!==2 (
        SET bestip=%%i
        goto :END
    )
)
:END

:: 判断刚刚获取的最快 IP 是否为空，以及是否和旧 IP 一样
if "%bestip%"=="" (
    echo.
    echo CloudflareST 测速结果 IP 数量为 0，跳过下面步骤...
    goto :STOP
)
if "%bestip%"=="%nowip%" (
    echo.
    echo CloudflareST 测速结果 IP 数量为 0，跳过下面步骤...
    goto :STOP
)


:: 下面这段代码是 "找不到满足条件的 IP 就一直循环测速下去" 才需要的代码
:: 考虑到当指定了下载速度下限，但一个满足全部条件的 IP 都没找到时，CloudflareST 就会输出所有 IP 结果
:: 因此当你指定 -sl 参数时，需要移除下面这段代码开头的这个 :: 冒号注释符，来做文件行数判断（比如下载测速数量：10 个，那么下面的值就设在为 11）
::set /a v=0
::for /f %%a in ('type result_hosts.txt') do set /a v+=1
::if %v% GTR 11 (
::    echo.
::    echo CloudflareST 测速结果没有找到一个完全满足条件的 IP，重新测速...
::    goto :RESET
::)


echo %bestip%>nowip_hosts.txt
echo.
echo 旧 IP 为 %nowip%
echo 新 IP 为 %bestip%

CD /d "C:\Windows\System32\drivers\etc"
echo.
echo 开始备份 Hosts 文件（hosts_backup）...
copy hosts hosts_backup
echo.
echo 开始替换...
(
    for /f "tokens=*" %%i in (hosts_backup) do (
        set s=%%i
        set s=!s:%nowip%=%bestip%!
        echo !s!
        )
)>hosts

echo 完成...
echo.
:STOP
pause 