# GoForPT

GoForPT 是一个基于 Go 语言开发的项目，主要用于处理种子文件和提供相关的 API 服务。

## 项目结构

- `api/`: 包含 API 相关的代码和逻辑。
- `model/`: 包含项目所使用的相关文件
- `pkg/cfg/`: 配置文件加载和管理模块。
- `pkg/database/`: 数据库初始化和数据操作模块。
- `pkg/ptcaches/`: 缓存初始化和管理模块。
- 

## 功能特性

- 加载配置文件。
- 检查并创建种子文件和静态文件夹。
- 初始化缓存（邮件缓存、公告缓存、对等节点缓存）。
- 初始化数据库并加载初始数据。
- 启动 API 服务器。

## 如何运行

1. 确保你已经安装了 Go 环境，PostgreSQL数据库环境。
2. 克隆项目到本地
```bash
git clone https://github.com/Jlan45/GoForPT
```
3. 修改配置文件并重命名为config.yaml
4. 运行即可