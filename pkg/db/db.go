package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

type User struct {
	gorm.Model         // 内含ID, CreatedAt, UpdatedAt, DeletedAt四个字段
	UserName    string `gorm:"column:user_name;size:100;not null;unique"` // 用户名长度最大100，不能为空，且唯一
	PassWord    string `gorm:"column:pass_word;size:100;not null"`        // 密码长度最大100，不能为空
	Email       string `gorm:"column:email;size:20"`
	PhoneNumber string `gorm:"column:phone_number;size:20"`
	UserID      int64  `gorm:"column:user_id;unique"`
	Delete      bool   `gorm:"column:delete"`
}

type File struct {
	gorm.Model
	FileKey   string `gorm:"column:file_key;size:100;not null"`   // 文件键长度最大100，不能为空
	ChunkNum  int    `gorm:"column:chunk_num;not null"`           // 块编号，不能为空
	Chunk     string `gorm:"column:chunk;type:longtext;not null"` // 文件块，使用longtext类型，不能为空
	FileType  string `gorm:"column:file_type;size:100;not null"`  // 文件类型长度最大100，不能为空
	OwnUserID int64  `gorm:"column:user_id;not null"`
}

type FileShare struct {
	gorm.Model
	FileKey string `gorm:"column:file_key;size:100;not null"`
	UserID  int64  `gorm:"column:user_id;not null"`
}

func InitDB() error {
	var err error
	dsn := "root:123456@tcp(localhost:3306)/app?charset=utf8mb4&parseTime=True&loc=Local"
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
  `user_name` varchar(100) NOT NULL,
  `pass_word` varchar(100) NOT NULL,
  `email` varchar(20) NOT NULL,
  `phone_number` varchar(20) NOT NULL,
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
  `file_key` varchar(100) NOT NULL,
  `chunk_num` int NOT NULL,
  `chunk` longtext NOT NULL,
  `file_type` varchar(100) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
   PRIMARY KEY (`id`),
  UNIQUE KEY `idx_files_file_key` (`file_key`, `chunk_num`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `file_shares` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `file_key` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_file_shares_user_id` (`user_id`, `file_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
*/
