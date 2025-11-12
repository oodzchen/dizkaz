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

// CategoryConfig 定义分类配置结构
type CategoryConfig struct {
	FrontId  string `yaml:"front_id"`
	Name     string `yaml:"name"`
	Describe string `yaml:"describe"`
	Approved bool   `yaml:"approved"`
}

// CategoriesConfig 定义分类配置文件的根结构
type CategoriesConfig struct {
	Categories []CategoryConfig `yaml:"categories"`
}

// getSuperAdminId 获取超级管理员用户ID
func getSuperAdminId() (int, error) {
	// 创建数据库连接
	cfg := config.Config
	pool, err := pgxpool.New(context.Background(), cfg.DB.GetDSN())
	if err != nil {
		return 0, err
	}
	defer pool.Close()

	// 查询超级管理员用户ID
	var adminId int
	err = pool.QueryRow(context.Background(), "SELECT id FROM users WHERE super_admin = true LIMIT 1").Scan(&adminId)
	if err != nil {
		return 0, err
	}
	return adminId, nil
}

// setCategoryApproval 设置分类的审核状态
func setCategoryApproval(frontId string, approved bool) error {
	// 创建数据库连接
	cfg := config.Config
	pool, err := pgxpool.New(context.Background(), cfg.DB.GetDSN())
	if err != nil {
		return err
	}
	defer pool.Close()

	// 更新approved字段
	_, err = pool.Exec(context.Background(), "UPDATE categories SET approved = $1 WHERE front_id = $2", approved, frontId)
	return err
}

// seedCategories 从YAML文件读取分类配置并插入到数据库
func seedCategories(categorySrv *service.Category, categoriesFile string) {
	// 读取YAML文件
	data, err := os.ReadFile(categoriesFile)
	if err != nil {
		log.Fatalf("Failed to read categories file %s: %v", categoriesFile, err)
	}

	// 解析YAML配置
	var config CategoriesConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse categories YAML file: %v", err)
	}

	if len(config.Categories) == 0 {
		fmt.Println("No categories found in the configuration file")
		return
	}

	// 获取超级管理员用户ID
	adminId, err := getSuperAdminId()
	if err != nil {
		log.Fatalf("Failed to get super admin user ID: %v", err)
	}
	fmt.Printf("Using super admin user ID: %d as category author\n", adminId)

	fmt.Printf("Found %d categories in configuration\n", len(config.Categories))

	// 创建分类
	successCount := 0
	for i, categoryConfig := range config.Categories {
		fmt.Printf("Creating category %d/%d: %s (%s)\n", i+1, len(config.Categories), categoryConfig.Name, categoryConfig.FrontId)

		// 检查分类是否已存在
		existingCategory, err := categorySrv.Store.Category.Item(categoryConfig.FrontId, 0)
		if err != nil && err.Error() != "no rows in result set" {
			fmt.Printf("Failed to check category existence %s: %v\n", categoryConfig.FrontId, err)
			continue
		}
		if existingCategory != nil {
			fmt.Printf("Category %s already exists with ID: %d\n", categoryConfig.FrontId, existingCategory.Id)

			// 如果分类已存在，更新分类信息
			_, err = categorySrv.Store.Category.Update(categoryConfig.FrontId, categoryConfig.Name, categoryConfig.Describe)
			if err != nil {
				fmt.Printf("Failed to update category %s: %v\n", categoryConfig.FrontId, err)
			} else {
				fmt.Printf("Successfully updated category %s\n", categoryConfig.FrontId)
				successCount++
			}

			// 设置审核状态
			if categoryConfig.Approved {
				err = setCategoryApproval(categoryConfig.FrontId, true)
				if err != nil {
					fmt.Printf("Failed to set approval status for category %s: %v\n", categoryConfig.FrontId, err)
				} else {
					fmt.Printf("Set approved = true for category %s\n", categoryConfig.FrontId)
				}
			}
			continue
		}

		// 创建分类模型并进行数据校验
		category := &model.Category{
			FrontId:  categoryConfig.FrontId,
			Name:     categoryConfig.Name,
			Describe: categoryConfig.Describe,
			Approved: categoryConfig.Approved,
		}

		// 清理数据
		category.TrimSpace()

		// 数据校验
		err = category.Valid()
		if err != nil {
			fmt.Printf("Category data validation failed for %s: %v\n", categoryConfig.FrontId, err)
			continue
		}

		// 使用store层创建分类，使用超级管理员用户作为作者
		categoryId, err := categorySrv.Store.Category.Create(
			category.FrontId,
			category.Name,
			category.Describe,
			adminId, // 使用超级管理员用户ID
		)
		if err != nil {
			fmt.Printf("Failed to create category %s: %v\n", categoryConfig.FrontId, err)
			continue
		}

		// 设置审核状态
		if categoryConfig.Approved {
			err = setCategoryApproval(categoryConfig.FrontId, true)
			if err != nil {
				fmt.Printf("Failed to set approval status for category %s: %v\n", categoryConfig.FrontId, err)
			} else {
				fmt.Printf("Set approved = true for category %s\n", categoryConfig.FrontId)
			}
		}

		successCount++
		fmt.Printf("Successfully created category %s with ID: %d\n", categoryConfig.FrontId, categoryId)
	}

	fmt.Printf("Category seeding completed. Success: %d/%d\n", successCount, len(config.Categories))
}