#!/bin/bash
# 智工集市服务器初始化脚本
# 运行后会生成配置文件模板，密钥需单独填入

set -e

BASE=$(cd "$(dirname "$0")" && pwd)

mkdir -p "$BASE/public/uploads" "$BASE/public/resized" \
         "$BASE/../logs" "$BASE/../database"

cat > "$BASE/config/application.yml" << 'ENDOFYML'
datasource:
  driverName: mysql
  host: mysql
  port: "3306"
  database: ssemarket
  username: root
  password: admin
  charset: utf8mb4
redis:
  host: redis
  port: "6379"
  password: ""
cos:
  bucketName: PLACEHOLDER_BUCKET
  appid: "PLACEHOLDER_APPID"
  secretId: PLACEHOLDER_SECRET_ID
  secretKey: PLACEHOLDER_SECRET_KEY
tms:
  secretId: PLACEHOLDER_SECRET_ID
  secretKey: PLACEHOLDER_SECRET_KEY
crypto:
  jwtKey: PLACEHOLDER_JWT_KEY
ENDOFYML

cat > "$BASE/.env" << 'ENDOFENV'
SMTP_EMAIL_USERNAME=PLACEHOLDER_EMAIL
SMTP_EMAIL_PASSWORD=PLACEHOLDER_PASS
SMTP_ADDR=mail.sysu.edu.cn:465
SMTP_HOST=mail.sysu.edu.cn
SMTP_NOTICE_EMAIL_USERNAME=PLACEHOLDER_EMAIL
SMTP_NOTICE_EMAIL_PASSWORD=PLACEHOLDER_PASS
SMTP_NOTICE_ADDR=mail.sysu.edu.cn:465
SMTP_NOTICE_HOST=mail.sysu.edu.cn
ENDOFENV

echo "✅ 配置模板创建完成，请运行 fill_secrets.sh 填入密钥"
