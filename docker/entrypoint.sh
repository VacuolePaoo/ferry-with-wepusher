#!/bin/sh
set -e

# 等待MySQL和Redis就绪
echo "Waiting for MySQL and Redis to be ready..."
while ! nc -z ferry_mysql 3306; do
  sleep 1
done
while ! nc -z ferry_redis 6379; do
  sleep 1
done

# 复制配置文件
if [[ ! -f /opt/workflow/ferry/config/settings.yml ]]
then
    cp -r /opt/workflow/ferry/default_config/* /opt/workflow/ferry/config/
fi

# 初始化数据库
if [[ -f /opt/workflow/ferry/config/needinit ]]
then
    echo "Initializing database..."
    /opt/workflow/ferry/ferry init -c=/opt/workflow/ferry/config/settings.yml
    rm -f /opt/workflow/ferry/config/needinit
fi

# 启动服务
echo "Starting ferry service..."
exec /opt/workflow/ferry/ferry server -c=/opt/workflow/ferry/config/settings.yml
