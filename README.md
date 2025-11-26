# MyNoteBook - 简易笔记管理系统

一个基于 Go 语言开发的笔记管理系统，支持用户注册、登录、笔记的增删改查及标签管理，适合学习 Go Web 开发的初学者参考。

## 前端页面访问
本项目包含完整的 HTML 静态页面，部署后可直接通过浏览器访问
static/index.html

## 功能特点
- 用户认证：注册、登录（基于 JWT 令牌）
- 笔记管理：创建、查询（列表/详情）、更新、删除笔记
- 标签功能：为笔记添加标签，方便分类
- 分类筛选：支持按分类筛选笔记
- 分页查询：笔记列表支持分页加载


## 技术栈
- 后端框架：Gin v1.11.0（轻量级 Web 框架）
- 数据库：MySQL（主数据存储）、Redis v9（缓存）
- 认证：JWT（JSON Web Token）
- 配置管理：Viper
- ORM：GORM v1.31.1
- 日志：Zap
- 参数校验：go-playground/validator


## 项目结构

MyNoteBook/
├── cmd/ # 入口文件
│ └── main.go # 程序入口
├── internal/ # 内部代码（不对外暴露）
│ ├── config/ # 配置相关
│ ├── model/ # 数据库模型
│ ├── service/ # 业务逻辑层
│ └── middlewares/ # 中间件（如认证）
├── pkg/ # 公共工具
│ ├── jwt/ # JWT 工具
│ ├── db/ # 数据库工具
│ ├── errcode/ # 统一错误码
│ ├── redis/ # 连接redis
│ ├── validator/ # 参数校验
│ └── response/ # 统一响应
├── api/ # API 接口层
├── router/ # 路由注册
├── config.yaml # 配置文件（本地开发用，不上传 Git）
├── config.yaml.example # 配置示例（上传 Git，供参考）
└── README.md # 项目说明


## 快速开始

### 前置条件
- Go 1.25+ 环境
- MySQL 5.7+
- Redis 6.0+


### 安装步骤
1. **克隆仓库**
   ```bash
   git clone https://github.com/JokerYuan-lang/MyNoteBook.git
   cd MyNoteBook


配置文件
复制示例配置文件并修改为自己的环境信息：
bash
运行
cp config_example.yaml config.yaml
项目启动时会自动创建数据表（基于 GORM 自动迁移）
启动项目
bash
运行
# 安装依赖
go mod tidy
# 启动服务
go run cmd/main.go
服务会在 http://localhost:8080 启动