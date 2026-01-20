// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-20, by liasica

package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"nexis.run/nexa/cmd/nexa/internal/base"
)

// setupTestProject 创建测试项目环境
func setupTestProject(t *testing.T) (rootDir string, cleanup func()) {
	t.Helper()

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "nexa-test-*")
	require.NoError(t, err)

	// 创建 go.mod 文件
	goModContent := `module github.com/test/project

go 1.25.3
`
	goModPath := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	// 返回清理函数
	cleanup = func() {
		_ = os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestRenderDao(t *testing.T) {
	// 设置测试项目
	rootDir, cleanup := setupTestProject(t)
	defer cleanup()

	tests := []struct {
		name       string
		config     *base.Config
		entityName string
		wantErr    bool
		checks     []func(t *testing.T, result []byte)
	}{
		{
			name: "生成基础 Dao 代码",
			config: &base.Config{
				RootDir:   rootDir,
				EntPath:   "internal/infrastructure/ent",
				DaoPath:   "internal/presentation/dao",
				OrmClient: "ent.Database",
			},
			entityName: "User",
			wantErr:    false,
			checks: []func(t *testing.T, result []byte){
				func(t *testing.T, result []byte) {
					code := string(result)
					// 检查包声明
					require.Contains(t, code, "package dao")
					// 检查导入
					require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent"`)
					require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent/user"`)
					// 检查结构体定义
					require.Contains(t, code, "type UserDao struct")
					require.Contains(t, code, "orm *ent.UserClient")
					// 检查构造函数
					require.Contains(t, code, "func NewUser(params ...*ent.UserClient) *UserDao")
					require.Contains(t, code, "orm := ent.Database.User")
					// 检查 return 语句
					require.Contains(t, code, "return &UserDao{")
					require.Contains(t, code, "orm: orm,")
					// 检查版权信息
					require.Contains(t, code, "Copyright")
					require.Contains(t, code, "Created at")
					require.Contains(t, code, "by nexa cli")
				},
			},
		},
		{
			name: "生成自定义 ORM 客户端的 Dao",
			config: &base.Config{
				RootDir:   rootDir,
				EntPath:   "internal/infrastructure/ent",
				DaoPath:   "internal/presentation/dao",
				OrmClient: "db.Client",
			},
			entityName: "Product",
			wantErr:    false,
			checks: []func(t *testing.T, result []byte){
				func(t *testing.T, result []byte) {
					code := string(result)
					require.Contains(t, code, "package dao")
					require.Contains(t, code, "type ProductDao struct")
					require.Contains(t, code, "orm *ent.ProductClient")
					require.Contains(t, code, "orm := db.Client.Product")
					require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent/product"`)
				},
			},
		},
		{
			name: "生成多单词实体名称的 Dao",
			config: &base.Config{
				RootDir:   rootDir,
				EntPath:   "internal/infrastructure/ent",
				DaoPath:   "internal/presentation/dao",
				OrmClient: "ent.Database",
			},
			entityName: "OrderItem",
			wantErr:    false,
			checks: []func(t *testing.T, result []byte){
				func(t *testing.T, result []byte) {
					code := string(result)
					require.Contains(t, code, "type OrderItemDao struct")
					require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent/orderitem"`)
					require.Contains(t, code, "func NewOrderItem(params ...*ent.OrderItemClient) *OrderItemDao")
					require.Contains(t, code, "orm := ent.Database.OrderItem")
					require.Contains(t, code, "orm *ent.OrderItemClient")
				},
			},
		},
		{
			name: "生成带连字符的 DaoPath",
			config: &base.Config{
				RootDir:   rootDir,
				EntPath:   "internal/infrastructure/ent",
				DaoPath:   "internal/data-access",
				OrmClient: "ent.Database",
			},
			entityName: "Customer",
			wantErr:    false,
			checks: []func(t *testing.T, result []byte){
				func(t *testing.T, result []byte) {
					code := string(result)
					// 包名应该将连字符替换为下划线
					require.Contains(t, code, "package data_access")
					require.Contains(t, code, "type CustomerDao struct")
					require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent/customer"`)
				},
			},
		},
		{
			name: "生成深层嵌套路径的 Dao",
			config: &base.Config{
				RootDir:   rootDir,
				EntPath:   "internal/infrastructure/database/ent",
				DaoPath:   "internal/application/data/dao",
				OrmClient: "ent.Database",
			},
			entityName: "Admin",
			wantErr:    false,
			checks: []func(t *testing.T, result []byte){
				func(t *testing.T, result []byte) {
					code := string(result)
					require.Contains(t, code, "package dao")
					require.Contains(t, code, `"github.com/test/project/internal/infrastructure/database/ent"`)
					require.Contains(t, code, `"github.com/test/project/internal/infrastructure/database/ent/admin"`)
					require.Contains(t, code, "type AdminDao struct")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行渲染
			result, err := RenderDao(tt.config, tt.entityName)

			// 检查错误
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, result)

			// 执行自定义检查
			for _, check := range tt.checks {
				check(t, result)
			}
		})
	}
}

// TestRenderDaoOutput 测试实际输出格式
func TestRenderDaoOutput(t *testing.T) {
	rootDir, cleanup := setupTestProject(t)
	defer cleanup()

	cfg := &base.Config{
		RootDir:   rootDir,
		EntPath:   "internal/infrastructure/ent",
		DaoPath:   "internal/presentation/dao",
		OrmClient: "ent.Database",
	}

	result, err := RenderDao(cfg, "User")
	require.NoError(t, err)

	code := string(result)
	t.Logf("Generated code:\n%s", code)

	// 验证代码结构完整性
	lines := strings.Split(code, "\n")
	require.Greater(t, len(lines), 10, "生成的代码应该有足够的行数")

	// 验证没有模板变量未替换
	require.NotContains(t, code, "{{", "不应包含未替换的模板变量")
	require.NotContains(t, code, "}}", "不应包含未替换的模板变量")

	// 验证代码可以编译（基本语法检查）
	require.True(t, strings.HasPrefix(code, "//"), "应该以注释开头")
	require.Contains(t, code, "package ", "应该包含 package 声明")
	require.Contains(t, code, "import (", "应该包含 import 声明")

	// 验证代码格式
	require.Contains(t, code, "type UserDao struct {", "结构体定义应该正确")
	require.Contains(t, code, "func NewUser(params ...*ent.UserClient) *UserDao {", "构造函数签名应该正确")

	// 验证逻辑完整性
	require.Contains(t, code, "if len(params) > 0 {", "应该包含参数检查逻辑")
	require.Contains(t, code, "orm = params[0]", "应该支持自定义 ORM 客户端")
}

// TestRenderDaoEdgeCases 测试边界情况
func TestRenderDaoEdgeCases(t *testing.T) {
	rootDir, cleanup := setupTestProject(t)
	defer cleanup()

	t.Run("实体名称为单字符", func(t *testing.T) {
		cfg := &base.Config{
			RootDir:   rootDir,
			EntPath:   "internal/infrastructure/ent",
			DaoPath:   "internal/presentation/dao",
			OrmClient: "ent.Database",
		}

		result, err := RenderDao(cfg, "A")
		require.NoError(t, err)
		code := string(result)
		require.Contains(t, code, "type ADao struct")
		require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent/a"`)
		require.Contains(t, code, "func NewA(params ...*ent.AClient) *ADao")
	})

	t.Run("实体名称包含数字", func(t *testing.T) {
		cfg := &base.Config{
			RootDir:   rootDir,
			EntPath:   "internal/infrastructure/ent",
			DaoPath:   "internal/presentation/dao",
			OrmClient: "ent.Database",
		}

		result, err := RenderDao(cfg, "User2FA")
		require.NoError(t, err)
		code := string(result)
		require.Contains(t, code, "type User2FADao struct")
		require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent/user2fa"`)
		require.Contains(t, code, "func NewUser2FA(params ...*ent.User2FAClient) *User2FADao")
	})

	t.Run("实体名称全大写", func(t *testing.T) {
		cfg := &base.Config{
			RootDir:   rootDir,
			EntPath:   "internal/infrastructure/ent",
			DaoPath:   "internal/presentation/dao",
			OrmClient: "ent.Database",
		}

		result, err := RenderDao(cfg, "API")
		require.NoError(t, err)
		code := string(result)
		require.Contains(t, code, "type APIDao struct")
		require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent/api"`)
		require.Contains(t, code, "func NewAPI(params ...*ent.APIClient) *APIDao")
	})

	t.Run("实体名称包含下划线（不推荐但应支持）", func(t *testing.T) {
		cfg := &base.Config{
			RootDir:   rootDir,
			EntPath:   "internal/infrastructure/ent",
			DaoPath:   "internal/presentation/dao",
			OrmClient: "ent.Database",
		}

		result, err := RenderDao(cfg, "User_Profile")
		require.NoError(t, err)
		code := string(result)
		require.Contains(t, code, "type User_ProfileDao struct")
		require.Contains(t, code, `"github.com/test/project/internal/infrastructure/ent/user_profile"`)
	})
}

// TestRenderDaoModuleImport 测试模块导入路径正确性
func TestRenderDaoModuleImport(t *testing.T) {
	rootDir, cleanup := setupTestProject(t)
	defer cleanup()

	tests := []struct {
		name       string
		entPath    string
		entityName string
		wantImport string
	}{
		{
			name:       "标准 ent 路径",
			entPath:    "internal/infrastructure/ent",
			entityName: "User",
			wantImport: `"github.com/test/project/internal/infrastructure/ent/user"`,
		},
		{
			name:       "短路径",
			entPath:    "ent",
			entityName: "Product",
			wantImport: `"github.com/test/project/ent/product"`,
		},
		{
			name:       "多层嵌套",
			entPath:    "internal/data/store/ent",
			entityName: "Order",
			wantImport: `"github.com/test/project/internal/data/store/ent/order"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &base.Config{
				RootDir:   rootDir,
				EntPath:   tt.entPath,
				DaoPath:   "internal/presentation/dao",
				OrmClient: "ent.Database",
			}

			result, err := RenderDao(cfg, tt.entityName)
			require.NoError(t, err)
			code := string(result)
			require.Contains(t, code, tt.wantImport)
		})
	}
}
