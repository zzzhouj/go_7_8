package main

import (
	"crypto/sha256"
	"encoding/hex"
	"filesys/model_def"
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// 清空数据库所有数据（保留表结构）
func clearAllData(db *gorm.DB) error {
	// 按依赖关系逆序删除数据，避免外键约束问题
	if err := db.Where("1 = 1").Delete(&model_def.Session{}).Error; err != nil {
		return fmt.Errorf("清空 sessions 表失败: %v", err)
	}
	if err := db.Where("1 = 1").Delete(&model_def.StoreRef{}).Error; err != nil {
		return fmt.Errorf("清空 store_refs 表失败: %v", err)
	}
	if err := db.Where("1 = 1").Delete(&model_def.Version{}).Error; err != nil {
		return fmt.Errorf("清空 versions 表失败: %v", err)
	}
	if err := db.Where("1 = 1").Delete(&model_def.File{}).Error; err != nil {
		return fmt.Errorf("清空 files 表失败: %v", err)
	}
	if err := db.Where("1 = 1").Delete(&model_def.User{}).Error; err != nil {
		return fmt.Errorf("清空 users 表失败: %v", err)
	}
	return nil
}

func main() {
	// db, err := gorm.Open(sqlite.Open("filesys.db"), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	// dao.SetDefault(db) // 让 gen 生成的代码使用这个 db

	// // 新增用户
	// user := &model_def.User{Name: "xxxxxx"}
	// dao.Q.User.WithContext(context.Background()).Create(user)

	// // 查询所有用户
	// users, err := dao.Q.User.WithContext(context.Background()).Find()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, u := range users {
	// 	fmt.Println(u.ID, u.Name)
	// }

	// // 分页查询
	// users, _, err = dao.Q.User.WithContext(context.Background()).FindByPage(2, 10)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, u := range users {
	// 	fmt.Println(u.ID, u.Name)
	// }

	// // 更新用户
	// _, err = dao.Q.User.WithContext(context.Background()).Where(dao.Q.User.ID.Eq(1)).Or(dao.Q.User.Name.Eq("Alice")).UpdateSimple(
	// 	dao.Q.User.Name.Value("AliceNew"),
	// )

	// // 删除用户
	// _, err = dao.Q.User.WithContext(context.Background()).Where(dao.Q.User.ID.Eq(1)).Delete()

	// 使用无 CGO 的 SQLite 驱动
	db, err := gorm.Open(sqlite.Open("../../filesys.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("打开数据库失败:", err)
	}

	// 可选的数据库清空操作 - 根据需要取消注释使用

	// 方式1: 只清空数据，保留表结构
	// if err := clearAllData(db); err != nil {
	// 	log.Fatal("清空数据失败:", err)
	// }

	// 方式2: 删除所有表
	// if err := dropAllTables(db); err != nil {
	// 	log.Fatal("删除表失败:", err)
	// }

	// 方式3: 完全重置数据库（删除表后重新创建）
	// if err := resetDatabase(db); err != nil {
	// 	log.Fatal("重置数据库失败:", err)
	// 	return
	// }

	// 自动建表
	if err := db.AutoMigrate(&model_def.User{}); err != nil {
		log.Fatal("建表失败:", err)
	}

	if err := db.AutoMigrate(&model_def.File{}); err != nil {
		log.Fatal("建表失败:", err)
	}

	if err := db.AutoMigrate(&model_def.Version{}); err != nil {
		log.Fatal("建表失败:", err)
	}

	if err := db.AutoMigrate(&model_def.StoreRef{}); err != nil {
		log.Fatal("建表失败:", err)
	}

	if err := db.AutoMigrate(&model_def.Session{}); err != nil {
		log.Fatal("建表失败:", err)
	}

	// 演示清空数据功能
	log.Println("开始清空所有表数据...")
	if err := clearAllData(db); err != nil {
		log.Printf("清空数据失败: %v", err)
	}
	// 默认带一个admin用户，默认密码是123456，数据库里存储的密码是经过sha256加密的，只有这个用户才能创建其它用户
	adminUser := &model_def.User{Name: "admin", Password: "dfjaidfaf"}
	// 计算 SHA-256 哈希
	hash := sha256.Sum256([]byte(adminUser.Password))

	// 转为十六进制字符串
	hashStr := hex.EncodeToString(hash[:])
	adminUser.Password = hashStr
	adminUser.Ctime = time.Now().Unix() // 创建时间
	adminUser.Mtime = adminUser.Ctime   // 修改时间
	if err := db.Create(adminUser).Error; err != nil {
		log.Printf("创建默认用户失败: %v", err)
	}

	log.Println("所有表数据已清空")
}
