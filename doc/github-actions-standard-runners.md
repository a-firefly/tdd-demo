# GitHub Actions 标准运行器与 Go Web 测试

整理日期：2026-10-15。运行器规格、系统标签和计费规则可能变化，使用时以文末官方文档为准。

## 1. 运行器是什么

Runner（运行器）是执行 GitHub Actions job 的机器及其执行程序。工作流通过 `runs-on` 选择运行环境：

```yaml
runs-on: ubuntu-latest
```

这表示使用 GitHub 托管的标准 Ubuntu 运行器，无需自己购买、配置服务器。GitHub 为每个 job 分配环境，同一个 job 的所有步骤依次执行，共享该环境中的文件和进程；job 结束后环境会被回收。

标准 Ubuntu、Windows 和 macOS 运行器通常使用虚拟机；`ubuntu-slim` 使用共享虚拟机上的独立容器。

运行器用于测试、构建和部署。需要长期对外提供服务的 Go Web 应用，应部署到服务器或应用托管平台，不能依赖临时运行器持续运行。

## 2. 当前工作流如何触发

本仓库的配置位于 [main.yml](../.github/workflows/main.yml)：

- 推送到远端 `main` 分支时触发，推送其他分支或创建 Pull Request 不会触发。
- `build` job 依次检出代码、按 `go.mod` 配置 Go、执行测试。
- 同一工作流、同一分支出现新的运行时，会取消尚未结束的旧运行。
- 整个 job 最长执行 10 分钟，包含环境准备和测试时间。

本地 `git commit` 不触发工作流；将提交推送到 GitHub 的 `main` 分支后才触发。运行结果可以在仓库的 **Actions → Go CI** 查看。

## 3. 有哪些标准运行器

以下为常用标签及官方列出的资源规格，CPU 数量按逻辑核理解。表中标准运行器的 SSD 均为 14 GB，实际可用空间需在运行时检查。

| 类型 | 常用 `runs-on` 标签举例 | 公开仓库 CPU / 内存 | 私有仓库 CPU / 内存 |
| --- | --- | --- | --- |
| Ubuntu x64 | `ubuntu-latest`、`ubuntu-24.04`、`ubuntu-22.04`、`ubuntu-26.04` | 4 核 / 16 GB | 2 核 / 8 GB |
| Ubuntu ARM64 | `ubuntu-24.04-arm`、`ubuntu-22.04-arm`、`ubuntu-26.04-arm` | 4 核 / 16 GB | 2 核 / 8 GB |
| Windows x64 | `windows-latest`、`windows-2025`、`windows-2022` | 4 核 / 16 GB | 2 核 / 8 GB |
| Windows ARM64 | `windows-11-arm` | 4 核 / 16 GB | 2 核 / 8 GB |
| macOS Intel | `macos-15-intel`、`macos-26-intel` | 4 核 / 14 GB | 4 核 / 14 GB |
| macOS ARM64（M1） | `macos-latest`、`macos-14`、`macos-15`、`macos-26` | 3 核 / 7 GB | 3 核 / 7 GB |
| Ubuntu 轻量容器 | `ubuntu-slim` | 1 核 / 5 GB | 1 核 / 5 GB |

同一个 `ubuntu-latest` 标签，在公开和私有仓库中分配的资源不同。`latest` 指 GitHub 提供的最新稳定镜像，不一定是操作系统厂商发布的最新版本。

若要固定操作系统大版本，可以使用 `ubuntu-24.04`。镜像中的工具和安全补丁仍会更新。

`ubuntu-slim` 适合轻量自动化任务，最长运行 15 分钟，且有容器权限限制。需要 MySQL、Redis 服务容器的常规 Go Web 测试，优先使用 `ubuntu-latest` 等完整 Ubuntu 运行器。

## 4. 免费还是收费

| 场景 | 计费方式 |
| --- | --- |
| 公开仓库使用标准托管运行器 | 执行时间免费 |
| 私有仓库使用标准托管运行器 | 先消耗套餐包含的免费分钟数，超出后按量收费 |
| 使用大型运行器（Larger runners） | 需要付费，公开仓库也不例外 |

私有仓库的常见免费额度为：GitHub Free 每月 2,000 分钟，GitHub Pro 每月 3,000 分钟。额度归仓库所有者账号或组织，由其下的私有仓库共享，不按每个仓库单独提供。

计时包括代码检出、环境准备、测试等 job 执行时间，单个 job 的不足一分钟部分向上取整。不同系统的计费标准不同，不能直接按相同价格估算。

例如，Ubuntu job 每次按 1 分钟计、每天运行 10 次，30 天约消耗 300 分钟。`timeout-minutes: 10` 只是超时上限，不代表每次固定扣除 10 分钟。

制品和缓存存储还有独立额度及计费规则，应与运行分钟数分别查看。

## 5. 每次都重新拉代码、安装 Go 吗

### 代码检出

每个 job 使用新环境，没有上次运行留下的仓库，所以会重新获取代码。当前 `actions/checkout@v7` 的默认行为是：

- 获取触发本次运行的提交，而不是运行开始时再任意选择一个最新提交。
- 获取该提交对应的完整代码文件。
- 默认 `fetch-depth: 1`，只获取一个提交的历史，不拉取所有历史和所有分支。

因此，它会重新获取代码，但不是完整克隆仓库的全部历史，也不是在上次运行的工作目录上增量更新。

`actions/checkout@v7` 是有效的官方用法，跟随维护者移动的 v7 主版本标签接收后续更新。`@v7.0.1` 则指定补丁版本，需要主动升级；如果要求严格固定代码，应使用经过核实的完整提交 SHA。

### Go 环境和缓存

每次都会执行 `actions/setup-go@v7`，但不一定重新下载安装 Go：

1. 根据 `go.mod` 确定所需版本。
2. 优先使用运行器工具缓存中匹配的 Go 版本；没有匹配版本时才下载。
3. 默认启用 Go 依赖和编译产物缓存，在命中缓存时恢复，减少重复下载和编译。

跨运行恢复缓存由 Actions 的缓存机制完成，不代表虚拟机或整个工作目录被保留，也不保证每次都能命中缓存。

当前测试命令：

```bash
go test -v -count=1 -timeout=5m ./...
```

- `-v`：输出详细测试结果。
- `-count=1`：禁用测试结果缓存，确保实际执行测试；依赖和编译缓存仍可使用。
- `-timeout=5m`：每个测试包的测试进程超时为 5 分钟。
- `./...`：测试当前模块下的所有包。

## 6. 哪些项目适合使用

大多数常规项目都适合，主要取决于资源和环境需求：

| 项目 | 常见任务 | 常用系统 |
| --- | --- | --- |
| Go、Java、Python 后端 | 检查代码、编译、单元测试、集成测试 | Ubuntu |
| Vue、React 前端 | 安装依赖、检查代码、测试、打包 | Ubuntu |
| Android 应用 | 编译 APK、运行测试 | Ubuntu |
| iOS、macOS 应用 | 使用 Xcode 编译和测试 | macOS |
| Windows 软件 | 编译和测试 Windows 功能 | Windows |
| 跨平台库、命令行工具 | 验证不同系统的兼容性 | 多系统矩阵 |
| 文档和静态网站 | 生成文档、构建和部署网站 | Ubuntu |

Go Web 项目可以用标准运行器完成测试、编译、构建 Docker 镜像和部署。如果需要 GPU、专用设备、特殊内网访问，或常规资源长期不足，再评估大型运行器或自托管运行器。

## 7. MySQL、Redis 临时服务如何运行

工作流的 `services` 配置会让 GitHub 在同一台 Ubuntu 运行器上启动 Docker 服务容器：

```text
Ubuntu 运行器
├── Go 测试程序
├── MySQL 容器
└── Redis 容器
```

MySQL、Redis 是真实服务，可以供测试读写数据。它们与 Go 测试共享运行器的 CPU、内存和磁盘，不会额外分配独立的数据库服务器。

每个 job 都会创建新容器，结束后销毁；未配置持久化的测试数据不会保留。测试前需要执行项目自己的数据库迁移、建表和测试数据初始化。

以下是可以添加到 job 中的服务配置示例，本仓库当前工作流尚未添加这些服务：

```yaml
jobs:
  build:
    runs-on: ubuntu-latest

    services:
      mysql:
        image: mysql:8.4
        env:
          MYSQL_ROOT_PASSWORD: test_password
          MYSQL_DATABASE: app_test
          MYSQL_USER: app
          MYSQL_PASSWORD: app_test_password
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping -h 127.0.0.1 --silent"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=10

      redis:
        image: redis:7
        ports:
          - 6379:6379
        options: >-
          --health-cmd="redis-cli ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=5

    # 在这里保留检出代码、配置 Go、初始化数据库和运行测试等 steps。
```

示例中的密码仅用于临时测试环境。应用连接 MySQL 时使用 `app` 用户、`app_test_password` 密码和 `app_test` 数据库；具体连接串或环境变量名称由项目决定。

当前工作流直接在运行器上执行 Go 测试，映射端口后可以连接：

- MySQL：`127.0.0.1:3306`。
- Redis：`127.0.0.1:6379`。

如果以后通过 job 的 `container` 配置把测试程序也放进容器，则通常使用服务名 `mysql:3306`、`redis:6379` 连接，不能照搬 `127.0.0.1`。

健康检查用于等待服务可用；数据库表结构和业务测试数据仍需由迁移或测试初始化步骤准备。

## 8. 4 核、16 GB 内存是否够用

对大多数中小型 Go Web 项目，4 核、16 GB 内存同时运行 Go 测试、MySQL 和 Redis，通常足够。是否流畅仍需通过真实测试确认，不能只根据语言或框架判断。

| 测试场景 | 资源判断 |
| --- | --- |
| 单元测试、mock 测试 | 通常负载较轻 |
| 少量测试数据的数据库集成测试 | 通常可以顺畅运行 |
| 大量并行测试、频繁建库建表、导入大量数据 | 可能受 CPU、内存或磁盘 I/O 限制 |
| 压力测试、大数据量性能测试 | 需要单独评估，并保证测试环境具有代表性 |

14 GB 存储需要关注：Docker 镜像、Go 依赖、编译缓存和数据库测试数据都会占空间。测试数据应保持精简，并检查实际磁盘余量。

出现瓶颈时按现象调整：

- 内存或数据库压力大：可用 `go test -p 2 ./...` 限制同时测试的包数；包内并行测试还需结合 `-parallel` 和测试代码控制。
- 数据初始化慢：减少重复建库、导入数据，在保证测试隔离的前提下复用初始化结果。
- 磁盘不足：检查镜像、缓存和测试数据库的占用。
- 运行超时：先定位耗时；确实需要更长时间时，再调整每包 5 分钟和整个 job 10 分钟的限制。

如果仓库是私有的，当前 `ubuntu-latest` 对应的是 2 核、8 GB 内存，应以这一规格评估，而不是按公开仓库的 4 核、16 GB 估算。

## 官方参考

- [GitHub 托管运行器：规格、系统标签和限制](https://docs.github.com/en/actions/reference/runners/github-hosted-runners)
- [GitHub Actions 计费和免费额度](https://docs.github.com/en/billing/concepts/product-billing/github-actions)
- [运行器价格和计时规则](https://docs.github.com/en/billing/reference/actions-runner-pricing)
- [大型运行器及其计费](https://docs.github.com/en/actions/concepts/runners/larger-runners)
- [actions/checkout：默认检出行为](https://github.com/actions/checkout)
- [actions/setup-go：版本选择和缓存](https://github.com/actions/setup-go)
- [Action 版本标签与提交 SHA](https://docs.github.com/en/actions/how-tos/create-and-publish-actions/manage-custom-actions)
- [服务容器的生命周期和连接方式](https://docs.github.com/en/actions/tutorials/use-containerized-services/use-docker-service-containers)
