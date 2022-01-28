/*
 Navicat Premium Data Transfer
 
 Source Server         : cachee
 Source Server Type    : SQLite
 Source Server Version : 3035005
 Source Schema         : main
 
 Target Server Type    : SQLite
 Target Server Version : 3035005
 File Encoding         : 65001
 
 Date: 20/01/2022 12:54:26
 */
PRAGMA foreign_keys = false;
-- ----------------------------
-- Table structure for records
-- ----------------------------
CREATE TABLE IF NOT EXISTS "cmd_cache" (
  "body" TEXT NOT NULL,
  "command" TEXT NOT NULL,
  "exit" integer NOT NULL DEFAULT 0,
  "go_version" text NOT NULL,
  "src" TEXT NOT NULL,
  "sum" TEXT NOT NULL
);
PRAGMA foreign_keys = true;