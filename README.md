# jmfzf

jmfzf 是一个集成了 tmux 和 fzf 的终端工具，旨在简化服务器资源的管理和连接。通过插件系统，jmfzf 支持多种云平台、容器平台等资源的管理，提供快速、高效的终端连接体验。

## 功能特性

- **多资源管理**：支持从云平台、容器平台等多种资源中查询和管理主机。
- **终端连接**：使用 SSH、docker exec、kubectl exec 等方式快速连接到资源。
- **插件系统**：通过插件扩展支持更多的资源类型和平台。
- **tmux 集成**：使用 tmux 实现多窗口、多面板的终端管理。
- **fzf 集成**：通过 fzf 实现快速搜索和筛选资源。
- **高可配置性**：支持多种配置选项，满足不同的使用需求。

## 安装

确保已安装 Go 环境。

```bash
go get -u github.com/raojinlin/jmfzf
```

## 使用前提
在运行 jmfzf 之前，请确保以下命令已正确安装：

- Kubectl：用于管理 Kubernetes 集群。
- Docker：用于管理和运行容器。
- SSH：用于远程连接服务器。
- tmux：用于实现多面板功能（如果需要使用 tmux 的多面板功能）。

可以使用以下命令来安装这些工具：

### 安装 Kubectl
```bash
# 使用包管理器安装，例如 Ubuntu 下使用 apt
sudo apt-get install -y kubectl
```

### 安装docker
```bash
# 使用官方安装脚本
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
```

### 安装 SSH
```bash
# 大多数系统预装了 SSH 客户端，如果没有，可以使用包管理器安装
sudo apt-get install -y openssh-client
```

### 安装 tmux
```bash
# 使用包管理器安装，例如 Ubuntu 下使用 apt
sudo apt-get install -y tmux
```

## 配置文件示例

以下是 `config.yaml` 配置文件的示例，并附上注释：

```yaml
# 默认使用的插件名称
default_plugin: "example"

# 日志配置
logger:
  # 日志级别，可选值: debug, info, warn, error
  level: "info"
  # 是否启用日志
  enabled: true
  # 日志文件路径
  file: "jmfzf.log"

# TMUX 配置
tmux:
  # 是否启用 TMUX
  enabled: true
  # TMUX 会话名称
  session_name: "jmfzf"

# 自定义插件目录
custom_plugins_dir: "/path/to/custom/plugins"

# 插件配置
plugins:
  # 插件名称
  example:
    # 端点地址
    endpoint: "https://api.example.com"
    # 访问密钥
    access_key: "your-access-key"
    # 密钥
    access_key_secret: "your-secret-key"
    # 区域列表
    regions:
      - "us-west-1"
      - "us-east-1"
    # 连接配置
    connections:
      # 是否使用本地 IP
      use_local_ip: true
```

## 插件开发指南

### 插件的作用
插件是 jmfzf 的核心，通过插件系统，用户可以扩展 jmfzf 的功能，以支持更多的云平台、容器平台和其他资源的管理。插件的主要作用包括：

- 资源查询：从不同平台获取资源信息。
- 终端连接：提供连接到资源的功能。
- 插件化开发：用户可以开发自定义插件，扩展 jmfzf 的功能。

### 插件开发流程

1. 创建插件目录
在自定义插件目录中创建一个新的插件目录，例如 `my_plugin`。

2. 实现插件接口
在插件目录中创建一个新的 Go 文件，例如  `my_plugin.go`，并实现插件接口。

```go
type ListOptions struct {
    Order   int    `json:"order" yaml:"order"`
    OrderBy string `json:"order_by" yaml:"order_by"`
}

type Plugin interface {
    List(option *ListOptions) ([]terminal.Host, error)
    Name() string
    Cache() bool
    Init(option interface{}) error
}
```

3. 实现插件逻辑
```go
package my_plugin

import (
    "github.com/raojinlin/jmfzf"
    "github.com/raojinlin/jmfzf/pkg/terminal"
)

type MyPlugin struct {
    option *jmfzf.CloudProviderOption
}

func NewMyPlugin() *MyPlugin {
    return &MyPlugin{option: &jmfzf.CloudProviderOption{}}
}

func (p *MyPlugin) Init(option interface{}) error {
    err := jmfzf.MapToStruct(option, p.option)
    if err != nil {
        return err
    }
    // 初始化云平台客户端
    return nil
}

func (p *MyPlugin) Name() string {
    return "my_plugin"
}

func (p *MyPlugin) List(options *ListOptions) ([]terminal.Host, error) {
    var hosts []terminal.Host
    hosts = append(hosts, terminal.Host{
        Type: terminal.TerminalTypeHost,
        SSHInfo: terminal.SSHInfo{
            Name:       "example-vm",
            PublicIP:   "192.168.1.1",
            LocalIP:    "10.0.0.1",
            Port:       22,
            User:       "root",
            UseLocalIP: p.option.ConnectionOptions.UseLocalIP,
        },
    })
    return hosts, nil
}

func (p *MyPlugin) Cache() bool {
    return true
}
```

4. 配置插件

在 jmfzf 的配置文件中添加自定义插件的配置。
```yaml
custom_plugins_dir: "/path/to/custom/plugins"

plugins:
  my_plugin:
    endpoint: "https://api.example.com"
    access_key: "your-access-key"
    access_key_secret: "your-secret-key"
    regions:
      - "us-west-1"
      - "us-east-1"
    connections:
      use_local_ip: true
```

5. 编译和运行
```go
go build -o jmfzf .
./jmfzf -config ~/.jmfzf.yaml

```


## 贡献
欢迎提交问题和拉取请求。要开始贡献，请克隆此仓库并创建一个新的分支进行更改。

```bash
git clone https://github.com/raojinlin/jmfzf.git
cd jmfzf
git checkout -b new-feature

```