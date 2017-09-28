/*
Navicat PGSQL Data Transfer

Source Server         : postgres_zpdb
Source Server Version : 90603
Source Host           : localhost:5432
Source Database       : my_site
Source Schema         : public

Target Server Type    : PGSQL
Target Server Version : 90603
File Encoding         : 65001

Date: 2017-09-17 11:13:37
*/


-- ----------------------------
-- Sequence structure for t_app_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."t_app_id_seq";
CREATE SEQUENCE "public"."t_app_id_seq"
 INCREMENT 1
 MINVALUE 1
 MAXVALUE 9223372036854775807
 START 1
 CACHE 1;

-- ----------------------------
-- Sequence structure for t_code_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."t_code_id_seq";
CREATE SEQUENCE "public"."t_code_id_seq"
 INCREMENT 1
 MINVALUE 1
 MAXVALUE 9223372036854775807
 START 1
 CACHE 1;

-- ----------------------------
-- Sequence structure for t_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."t_log_id_seq";
CREATE SEQUENCE "public"."t_log_id_seq"
 INCREMENT 1
 MINVALUE 1
 MAXVALUE 9223372036854775807
 START 4
 CACHE 1;
SELECT setval('"public"."t_log_id_seq"', 4, true);

-- ----------------------------
-- Sequence structure for t_user_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."t_user_id_seq";
CREATE SEQUENCE "public"."t_user_id_seq"
 INCREMENT 1
 MINVALUE 1
 MAXVALUE 9223372036854775807
 START 8455
 CACHE 1;
SELECT setval('"public"."t_user_id_seq"', 8455, true);

-- ----------------------------
-- Table structure for t_app
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_app";
CREATE TABLE "public"."t_app" (
"id" int4 DEFAULT nextval('t_app_id_seq'::regclass) NOT NULL,
"icon" text COLLATE "default",
"appId" varchar(16) COLLATE "default",
"version" varchar(16) COLLATE "default",
"developer" int4,
"desc" varchar(255) COLLATE "default",
"valid" int4,
"expend" jsonb NOT NULL,
"downloadCount" int4,
"created_at" timestamp(6)
)
WITH (OIDS=FALSE)

;

-- ----------------------------
-- Table structure for t_code
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_code";
CREATE TABLE "public"."t_code" (
"id" int4 DEFAULT nextval('t_code_id_seq'::regclass) NOT NULL,
"code" varchar(10) COLLATE "default",
"appId" varchar(16) COLLATE "default",
"developer" int4,
"user" jsonb NOT NULL,
"desc" varchar(255) COLLATE "default",
"valid" int4,
"machineCount" int4,
"mostCount" int4,
"startTime" timestamp(6),
"endTime" timestamp(6),
"last_time" timestamp(6),
"created_at" timestamp(6)
)
WITH (OIDS=FALSE)

;

-- ----------------------------
-- Table structure for t_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_log";
CREATE TABLE "public"."t_log" (
"id" int4 DEFAULT nextval('t_log_id_seq'::regclass) NOT NULL,
"type" varchar(10) COLLATE "default",
"tag" varchar(50) COLLATE "default",
"operator" varchar(50) COLLATE "default",
"content" varchar(255) COLLATE "default",
"created_at" timestamp(6)
)
WITH (OIDS=FALSE)

;

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_user";
CREATE TABLE "public"."t_user" (
"id" int4 DEFAULT nextval('t_user_id_seq'::regclass) NOT NULL,
"role" int4 NOT NULL,
"nick" varchar(50) COLLATE "default" NOT NULL,
"pwd" varchar(50) COLLATE "default" NOT NULL,
"avator" varchar(255) COLLATE "default" NOT NULL,
"phone" varchar(11) COLLATE "default" NOT NULL,
"email" varchar(50) COLLATE "default" NOT NULL,
"qq" varchar(50) COLLATE "default" NOT NULL,
"expend" jsonb NOT NULL,
"created_at" timestamp(6) NOT NULL,
"updated_at" timestamp(6) NOT NULL
)
WITH (OIDS=FALSE)

;

-- ----------------------------
-- Alter Sequences Owned By 
-- ----------------------------
ALTER SEQUENCE "public"."t_app_id_seq" OWNED BY "t_app"."id";
ALTER SEQUENCE "public"."t_code_id_seq" OWNED BY "t_code"."id";
ALTER SEQUENCE "public"."t_log_id_seq" OWNED BY "t_log"."id";
ALTER SEQUENCE "public"."t_user_id_seq" OWNED BY "t_user"."id";

-- ----------------------------
-- Primary Key structure for table t_app
-- ----------------------------
ALTER TABLE "public"."t_app" ADD PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table t_code
-- ----------------------------
ALTER TABLE "public"."t_code" ADD PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table t_log
-- ----------------------------
ALTER TABLE "public"."t_log" ADD PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table t_user
-- ----------------------------
ALTER TABLE "public"."t_user" ADD PRIMARY KEY ("id");
