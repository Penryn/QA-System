# QA-System
一个简单的问卷系统


### 项目目录
```
QA-System/
├── main.go                   # 应用程序的入口点，包含启动服务器的代码
├── conf                      # 存放配置文件，如 YAML、JSON 格式的配置
├── docs                      # 项目文档，可能包括 API 文档、开发者指南等
│   └── README.md             # 项目 README 文档
├── go.mod                    # Go Modules 模块依赖文件
├── go.sum                    # Go Modules 模块依赖的校验和
├── hack                      # 构建脚本、CI 配置和辅助工具
│   └── docker                # Docker 相关配置，如 Dockerfile 和 Docker Compose 文件
├── internal                  # 项目内部包，包含服务器、模型、配置等
│   ├── global                # 全局可用的配置和初始化代码
│   │   └── config            # 配置加载和解析
│   ├── router                # 路由注册，定义应用程序的路由结构
│   ├── middleware            # 中间件逻辑，处理跨请求的任务如日志、认证等
│   ├── handle                # 放置 handler 函数，可能用于处理具体的业务逻辑
│   ├── service               # 业务逻辑服务层，实现应用程序的核心业务逻辑
│   ├── dao                   # 数据访问对象层，与数据库交互，执行增删改查等操作
│   ├── model                 # 数据模型定义
│   │   ├── user.go           # 用户模型定义
│   │   └── ...
│   └── pkg                  # 内部工具包
│       ├── code             # 错误码定义，用于标准化错误处理
│       ├── utils            # 内部使用的工具函数
│       ├── log              # 日志配置和管理，封装日志记录器的配置和使用
│       ├── database         # 数据库连接和初始化，管理数据库连接池
│       ├── session          # 会话管理，处理用户会话和状态
│       └── redis            # Redis 配置和管理，封装 Redis 缓存操作
├── logs                     # 日志文件输出目录，存放应用程序生成的日志文件
├── LICENSE                   # 项目许可证文件
├── Makefile                  # 根 Makefile 文件，包含构建和编译项目的指令
├── pkg                       # 可被外部引用的全局工具包
│   └── util                  # 通用工具代码
├── README.md                  # 项目 README 文档，通常提供项目概览和快速开始指南
├── public                    # 公共静态资源，如未构建的前端资源或可直接访问的静态文件
└── .gitignore                # Git 忽略文件配置
```