package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

const (
	VideoColumnValue = "video"
)

type User struct {
	gorm.Model         // 内含ID, CreatedAt, UpdatedAt, DeletedAt四个字段
	UserName    string `gorm:"column:user_name;size:100;not null;unique"` // 用户名长度最大100，不能为空，且唯一
	PassWord    string `gorm:"column:pass_word;size:100;not null"`        // 密码长度最大100，不能为空
	Email       string `gorm:"column:email;size:20"`
	PhoneNumber string `gorm:"column:phone_number;size:20"`
	UserID      int64  `gorm:"column:user_id;unique"`
	Delete      bool   `gorm:"column:delete;not null"`
}

type File struct {
	gorm.Model
	FileKey      string `gorm:"column:file_key;size:100;not null"`   // 文件键长度最大100，不能为空
	ChunkNum     int    `gorm:"column:chunk_num;not null"`           // 块编号，不能为空
	ChunkSize    int    `gorm:"column:chunk_size;not null"`          // 块大小，不能为空
	Chunk        []byte `gorm:"column:chunk;type:longblob;not null"` // 文件块，blob，不能为空
	FileType     string `gorm:"column:file_type;size:100;not null"`  // 文件类型长度最大100，不能为空
	HasMore      bool   `gorm:"column:has_more;not null"`
	OwnUserID    int64  `gorm:"column:user_id;not null"`
	Delete       bool   `gorm:"column:delete;not null"`
	IsCompressed bool   `gorm:"column:is_compressed;not null"`
}

type FileShare struct {
	gorm.Model
	FileKey    string `gorm:"column:file_key;size:100;not null"`
	UserID     int64  `gorm:"column:user_id;not null"`
	Permission int    `gorm:"column:permission;not null"`
	Delete     bool   `gorm:"column:delete;not null"`
}

type FileInfo struct {
	gorm.Model
	FileKey         string `gorm:"column:file_key;size:100;not null"`
	OwnUserID       int64  `gorm:"column:user_id;not null"`
	UploadEnd       bool   `gorm:"column:upload_end;not null"`
	Delete          bool   `gorm:"column:delete;not null"`
	CompressionType string `gorm:"compression_type;not null"`
}

func InitDB() error {
	var err error
	dsn := "root:QmlGls08280709@tcp(bj-cynosdbmysql-grp-4u4v5eag.sql.tencentcdb.com:21216)/app?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func GetDB() *gorm.DB {
	return _db
}

/*
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `delete` bool DEFAULT false,
  `user_name` varchar(100) NOT NULL,
  `pass_word` varchar(100) NOT NULL,
  `email` varchar(20),
  `phone_number` varchar(20),
  `user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`user_name`),
  UNIQUE KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `files` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `delete` bool DEFAULT false,
  `file_key` varchar(100) NOT NULL,
  `chunk_num` int NOT NULL,
  `chunk` LONGBLOB NOT NULL,
  `chunk_size` int NOT NULL,
  `has_more` bool NOT NULL default TRUE,
  `file_type` varchar(100) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `is_compressed` bool DEFAULT false,
   PRIMARY KEY (`id`),
  UNIQUE KEY `idx_files_file_key` (`file_key`, `chunk_num`),
  FOREIGN KEY (file_key) REFERENCES file_infos(file_key),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `file_shares` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `delete` bool DEFAULT false,
  `user_id` bigint unsigned NOT NULL,
  `file_key` varchar(100) NOT NULL,
  `has_more` bool NOT NULL default TRUE,
  `permission` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_shares_user_id` (`user_id`, `file_key`),
  FOREIGN KEY (file_key) REFERENCES file_infos(file_key),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `file_infos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `delete` bool DEFAULT false,
  `user_id` bigint unsigned NOT NULL,
  `file_key` varchar(100) NOT NULL,
  `upload_end` bool NOT NULL default FALSE,
  `compression_type` varchar(100) NOT NULL default 'av1',
  PRIMARY KEY (`id`),
  UNIQUE KEY(`file_key`),
  UNIQUE KEY `idx_file_shares_user_id` (`user_id`, `file_key`),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
*/
