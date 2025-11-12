package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oodzchen/dizkaz/config"
	"github.com/oodzchen/dizkaz/model"
	"github.com/oodzchen/dizkaz/service"
)

// UserConfig 定义用户配置结构
type UserConfig struct {
	Email        string `yaml:"email"`
	Name         string `yaml:"name"`
	Password     string `yaml:"password"`
	Role         string `yaml:"role"`
	SuperAdmin   bool   `yaml:"super_admin"`
	Introduction string `yaml:"introduction"`
}

// UsersConfig 定义用户配置文件的根结构
type UsersConfig struct {
	Users []UserConfig `yaml:"users"`
}

// setSuperAdmin 设置用户的超级管理员状态
func setSuperAdmin(userId int, superAdmin bool) error {
	// 创建数据库连接
	cfg := config.Config
	pool, err := pgxpool.New(context.Background(), cfg.DB.GetDSN())
	if err != nil {
		return err
	}
	defer pool.Close()

	// 更新super_admin字段
	_, err = pool.Exec(context.Background(), "UPDATE users SET super_admin = $1 WHERE id = $2", superAdmin, userId)
	return err
}

// seedUsers 从YAML文件读取用户配置并插入到数据库
func seedUsers(userSrv *service.User, usersFile string) {
	// 读取YAML文件
	data, err := os.ReadFile(usersFile)
	if err != nil {
		log.Fatalf("Failed to read users file %s: %v", usersFile, err)
	}

	// 解析YAML配置
	var config UsersConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse users YAML file: %v", err)
	}

	if len(config.Users) == 0 {
		fmt.Println("No users found in the configuration file")
		return
	}

	fmt.Printf("Found %d users in configuration\n", len(config.Users))

	// 创建用户
	successCount := 0
	for i, userConfig := range config.Users {
		fmt.Printf("Creating user %d/%d: %s (%s)\n", i+1, len(config.Users), userConfig.Name, userConfig.Email)

		// 检查用户是否已存在
		existingUserId, err := userSrv.Store.User.Exists(userConfig.Email, userConfig.Name)
		if err != nil && err.Error() != "no rows in result set" {
			fmt.Printf("Failed to check user existence %s: %v\n", userConfig.Email, err)
			continue
		}
		if existingUserId > 0 {
			fmt.Printf("User %s already exists with ID: %d\n", userConfig.Email, existingUserId)
			continue
		}

		// 创建用户模型并进行数据校验
		user := &model.User{
			Email:        userConfig.Email,
			Name:         userConfig.Name,
			Password:     userConfig.Password,
			Introduction: userConfig.Introduction,
		}

		// 清理数据
		user.TrimSpace()
		user.Sanitize(userSrv.SantizePolicy)

		// 数据校验
		err = user.Valid(true)
		if err != nil {
			fmt.Printf("User data validation failed for %s: %v\n", userConfig.Email, err)
			continue
		}

		// 密码加密
		err = user.EncryptPassword()
		if err != nil {
			fmt.Printf("Password encryption failed for user %s: %v\n", userConfig.Email, err)
			continue
		}

		// 使用store层创建用户，直接传入加密后的密码
		userId, err := userSrv.Store.User.Create(user.Email, user.Password, user.Name, string(model.DefaultUserRoleCommon))
		if err != nil {
			fmt.Printf("Failed to create user %s: %v\n", userConfig.Email, err)
			continue
		}

		// 如果指定了角色，更新用户角色
		if userConfig.Role != "" {
			_, err = userSrv.Store.User.SetRole(userId, userConfig.Role)
			if err != nil {
				fmt.Printf("Failed to set role %s for user %s: %v\n", userConfig.Role, userConfig.Email, err)
				// 继续处理，角色设置失败不影响用户创建
			}
		}

		// 如果指定了简介，更新用户简介
		if userConfig.Introduction != "" {
			// 先获取用户名
			user, err := userSrv.Store.User.Item(userId)
			if err != nil {
				fmt.Printf("Failed to get user info for setting introduction: %v\n", err)
			} else {
				err = userSrv.Store.User.UpdateIntroduction(user.Name, userConfig.Introduction)
				if err != nil {
					fmt.Printf("Failed to set introduction for user %s: %v\n", userConfig.Email, err)
					// 继续处理，简介设置失败不影响用户创建
				}
			}
		}

		// 设置超级管理员状态
		if userConfig.SuperAdmin {
			err = setSuperAdmin(userId, true)
			if err != nil {
				fmt.Printf("Failed to set super_admin for user %s: %v\n", userConfig.Email, err)
				// 继续处理，super_admin设置失败不影响用户创建
			} else {
				fmt.Printf("Set super_admin = true for user %s\n", userConfig.Email)
			}
		}

		successCount++
		fmt.Printf("Successfully created user %s with ID: %d\n", userConfig.Email, userId)
	}

	fmt.Printf("User seeding completed. Success: %d/%d\n", successCount, len(config.Users))
}