#!/bin/bash

# 创建自签名SSL证书用于开发环境
# 生产环境请使用正式的SSL证书

cd "$(dirname "$0")"

if [ ! -f "ssl/server.crt" ] || [ ! -f "ssl/server.key" ]; then
    echo "正在生成自签名SSL证书..."
    
    # 生成私钥
    openssl genrsa -out ssl/server.key 2048
    
    # 生成证书签名请求
    openssl req -new -key ssl/server.key -out ssl/server.csr -subj "/C=CN/ST=Beijing/L=Beijing/O=CloudBP/CN=localhost"
    
    # 生成自签名证书
    openssl x509 -req -days 365 -in ssl/server.csr -signkey ssl/server.key -out ssl/server.crt
    
    # 删除临时文件
    rm ssl/server.csr
    
    echo "SSL证书生成完成"
else
    echo "SSL证书已存在，跳过生成"
fi