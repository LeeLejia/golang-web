
DROP DATABASE IF EXISTS MY_SITE;
CREATE DATABASE MY_SITE;
USE MY_SITE;
# 日志表
DROP TABLE IF EXISTS "t_log";
CREATE TABLE "t_log" (
  "id"         SERIAL PRIMARY KEY,
  "type"       VARCHAR(10) COLLATE "default",
  "tag"        VARCHAR(50) COLLATE "default",
  "operator"      VARCHAR(50) COLLATE "default",
  "content"    VARCHAR(255) COLLATE "default",
  "created_at" TIMESTAMP
)
WITH (OIDS = FALSE);

# 用户表
DROP TABLE IF EXISTS "t_user";
CREATE TABLE "t_user" (
  "id"         SERIAL PRIMARY KEY,
  "role"        VARCHAR (10),
  "developer"   INT,
  "nick"       VARCHAR(50) COLLATE "default",
  "pwd"        VARCHAR(50) COLLATE "default",
  "avatar"      VARCHAR(225) COLLATE "default",
  "phone"       VARCHAR(11) COLLATE "default",
  "email"      VARCHAR(50) COLLATE "default",
  "qq"      VARCHAR(50) COLLATE "default",
  "status"  INT,
  "expend"  JSONB NOT NULL,
  "created_at" TIMESTAMP,
  "updated_at" TIMESTAMP,
)
WITH (OIDS = FALSE);

# 邀请码表
DROP TABLE IF EXISTS "t_code";
CREATE TABLE "t_code" (
  "id"         SERIAL PRIMARY KEY,
  "code"       VARCHAR(10) COLLATE "default",
  "appId"      VARCHAR(16) COLLATE "default",
  "developer"  INT,
  "user"       JSONB NOT NULL,
  "desc"       VARCHAR(255) COLLATE "default",
  "valid"      BOOL,

  "machineCount"  INT,
  "mostCount"    INT,
  "enableTime"  BOOL,
  "startTime"   TIMESTAMP,
  "endTime"     TIMESTAMP,
  "created_at" TIMESTAMP
)
WITH (OIDS = FALSE);

# 应用列表
DROP TABLE IF EXISTS "t_app";
CREATE TABLE "t_app" (
  "id"         SERIAL PRIMARY KEY,
  "icon"       text COLLATE "default",
  "appId"      VARCHAR(16) COLLATE "default",
  "version"    VARCHAR(16) COLLATE "default",
  "desc"       VARCHAR(255) COLLATE "default",
  "developer"  INT,
  "valid"      BOOL,
  "file"       INT,
  "src"        INT,
  "expend"     JSONB NOT NULL,

  "downloadCount"  INT,
  "created_at" TIMESTAMP
)
WITH (OIDS = FALSE);

# 文件列表
DROP TABLE IF EXISTS "t_file";
CREATE TABLE "t_file" (
  "id"         SERIAL PRIMARY KEY,
  "file_key"       VARCHAR(16) COLLATE "default",
  "file_type"       VARCHAR(16) COLLATE "default",
  "file_name"        VARCHAR(32) COLLATE "default",
  "owner"           INT ,
  "path"            VARCHAR(255) COLLATE "default",
  "created_at" TIMESTAMP
)
WITH (OIDS = FALSE);

