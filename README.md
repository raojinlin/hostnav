# hostnav

语言：中文 | [English](./README_en.md)

hostnav 是一个集成了 tmux 和 fzf 的终端工具，旨在简化服务器资源的管理和连接。通过插件系统，hostnav 支持多种云平台、容器平台等资源的管理，提供快速、高效的终端连接体验。



https://github.com/raojinlin/hostnav/assets/19492031/707557d9-8032-40a6-9fff-6ea3b8d0b800


## 功能特性

- **多资源管理**：支持从云平台、容器平台等多种资源中查询和管理主机。
- **终端连接**：使用 SSH、docker exec、kubectl exec 等方式快速连接到资源。
- **插件系统**：通过插件扩展支持更多的资源类型和平台。
- **tmux 集成**：使用 tmux 实现多窗口、多面板的终端管理。
- **fzf 集成**：通过 fzf 实现快速搜索和筛选资源。
- **高可配置性**：支持多种配置选项，满足不同的使用需求。
- **多平台支持**: 支持Widnows、Linux、MacOS

## 安装

确保已安装 Go 环境。

```bash
git clone github.com/raojinlin/hostnav
cd hostnav
go mod tidy
go install ./cmd/hostnav

hostnav
```


### 二进制文件下载
```bash
# Linux amd64
curl -O https://github.com/raojinlin/hostnav/releases/latest/download/hostnav-linux-amd64

# MacOS amd64
curl -O https://github.com/raojinlin/hostnav/releases/latest/download/hostnav-darwin-amd64
```

### Windows下载
点击 [hostnav.exe](https://github.com/raojinlin/hostnav/releases/latest/download/hostnav.exe) 下载


## 使用前提
在运行 hostnav 之前，请确保以下命令已正确安装：

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
# 缓存配置
cache: 
  # 缓存保存路径
  directory: ./cache
  # 缓冲时间（分钟）
  duration: 30
# 默认使用的插件名称列表
default_plugins:
- docker
- kubernetes
- cvm
- jumpserver

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

## 插件

下面是hostnav支持的插件。

### 阿里云ECS插件
使用：
```bash
hostnav -plugins ecs
```

配置：
```yaml
plugins:
  ecs:
    access_key: xxxx
    access_key_secret: xxxx
    regions:
    - cn-qingdao
    - cn-beijing
    - cn-zhangjiakou
    - cn-hongkong
```

### 文件插件
从文件加载自定义的服务器列表

使用:
```bash
hostnav -plugins file
```

配置：

```yaml
- type: host
  ssh_info:
    name: my test server
    public_ip: 192.168.5.43
    local_ip: 192.168.5.43
    user: root
    port: 22
    user_local_ip: true
```

### BCC百度云服务器插件
使用:
```bash
hostnav -plugins bcc
```

配置:
```yaml
plugins:
  bcc:
    # 访问密钥
    access_key: "your-access-key"
    # 密钥
    access_key_secret: "your-secret-key"
    # 区域列表
    regions:
      - cn-bj
      - cn-bd
      - cn-gz
      - cn-su
      - cn-hkg
      - cn-fwh
      - cn-cd
      - cn-fsh
    # 连接配置
    connections:
      # 是否使用本地 IP
      use_local_ip: true 
```

### CVM腾讯云服务器插件
使用
```bash
hostnav -plugins cvm
```

配置：

```yaml
plugins:
  cvm:
    # 访问密钥
    access_key: "your-access-key"
    # 密钥
    access_key_secret: "your-secret-key"
    # 区域列表
    regions:
      - "ap-guangzhou"
      - "ap-shanghai"
    # 连接配置
    connections:
      # 是否使用本地 IP
      use_local_ip: true 
```

### Docker容器插件
使用：
```bash
hostnav -plugins docker
```

配置：
```yaml
plugins:
  docker:
    host: unix:///var/run/docker.sock
    version: "1.43"
```

### Kubernetes容器插件
使用:
```bash
hsotnav -plugins kubernetes
```

配置：
```yaml
plugins:
  kubernetes:
    # 要查询的namespace
    namespaces:
    - default
    - kube_system
    # kubeconfig路径
    kubeconfig: ~/.kube/config
    # kubectl命令的路径
    kubectl: /snap/bin/microk8s.kubectl
```

### JumpServer插件
使用:
```bash
hostnav -plugins jumpserver
```

配置：
```yaml
plugins:
  jumpserver:
    # Jumpserver服务器地址
    url: http://your-jumpserver-url
    access_key: your api access key
    access_key_secret: your api access key secret
    # 资产过滤
    search: xxxx
```


### 使用tmux
要在 tmux 中绑定按键以打开 hostnav，你可以将以下配置添加到您的 ~/.tmux.conf 文件中：

```bash
# 绑定 Prefix-G 在新窗口中打开 hostnav
tmux bind g split-window 'hostnav -new-window'

# 绑定 Prefix-G 在当前窗口分隔出一个面板并打开 hostnav
tmux bind g split-window 'hostnav'
```

关于分隔窗口并以不同方式打开 hostnav 的更具体示例：

```bash
# 绑定 Prefix-G 分隔出一个垂直面板并打开 hostnav
tmux bind g split-window -v 'hostnav'

# 绑定 Prefix-G 分隔出一个水平面板并打开 hostnav
tmux bind g split-window -h 'hostnav'
```

这些命令绑定键 G（在按下 tmux 前缀键后，通常是 Ctrl-b），以创建一个新窗口或分隔当前窗口，并在新面板或窗口中执行 hostnav​ ([Koen Woortman](https://koenwoortman.com/tmux-remap-split-window-keys/))​​ ([Stack Overflow](https://stackoverflow.com/questions/38967247/tmux-how-do-i-bind-function-keys-to-commands))​。








## 插件开发指南

### 插件的作用
插件是 hostnav 的核心，通过插件系统，用户可以扩展 hostnav 的功能，以支持更多的云平台、容器平台和其他资源的管理。插件的主要作用包括：

- 资源查询：从不同平台获取资源信息。
- 终端连接：提供连接到资源的功能。
- 插件化开发：用户可以开发自定义插件，扩展 hostnav 的功能。

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
    "github.com/raojinlin/hostnav"
    "github.com/raojinlin/hostnav/pkg/terminal"
)

type MyPlugin struct {
    option *hostnav.CloudProviderOption
}

func NewMyPlugin() *MyPlugin {
    return &MyPlugin{option: &hostnav.CloudProviderOption{}}
}

func (p *MyPlugin) Init(option interface{}) error {
    err := hostnav.MapToStruct(option, p.option)
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

在 hostnav 的配置文件中添加自定义插件的配置。
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
go build -o hostnav .
./hostnav -config ~/.hostnav.yaml

```


## 贡献
欢迎提交问题和拉取请求。要开始贡献，请克隆此仓库并创建一个新的分支进行更改。

```bash
git clone https://github.com/raojinlin/hostnav.git
cd hostnav
git checkout -b new-feature

```
