# hostnav

Language： English | [中文](./README_CN.md)

hostnav is a terminal tool integrated with tmux and fzf, designed to simplify the management and connection of server resources. Through its plugin system, hostnav supports managing resources from various cloud platforms, container platforms, and more, providing a fast and efficient terminal connection experience.

https://github.com/raojinlin/hostnav/assets/19492031/707557d9-8032-40a6-9fff-6ea3b8d0b800

## Features

- **Multi-resource Management**: Supports querying and managing hosts from multiple resources such as cloud platforms and container platforms.
- **Terminal Connection**: Quickly connect to resources using SSH, docker exec, kubectl exec, etc.
- **Plugin System**: Extend support for more resource types and platforms through plugins.
- **tmux Integration**: Manage multiple windows and panels in the terminal using tmux.
- **fzf Integration**: Quickly search and filter resources using fzf.
- **High Configurability**: Supports various configuration options to meet different usage needs.

## Installation


### Install from Source
Ensure that the Go environment is installed.

```bash
git clone github.com/raojinlin/hostnav
cd hostnav
go mod tidy
go install ./cmd/hostnav

hostnav
```

### Binary Download
Use the following command to download.

```bash
# Linux amd64
curl -O https://github.com/raojinlin/hostnav/releases/latest/download/hostnav-linux-amd64

# MacOS amd64
curl -O https://github.com/raojinlin/hostnav/releases/latest/download/hostnav-darwin-amd64
```

### windows install
Download the latest version of hostnav.exe

[hostnav.exe](https://github.com/raojinlin/hostnav/releases//latest/download/hostnav.exe)


## Using tmux
To bind keys in tmux to open hostnav, you can add the following configurations to your ~/.tmux.conf file:

```bash
# Bind Prefix-G to open hostnav in a new window
tmux bind g new-window 'hostnav'

# Bind Prefix-G to split the current window and open hostnav in a new pane
tmux bind g split-window 'hostnav'
```

For a more specific example related to splitting windows and opening hostnav in different ways:

```bash
# Bind Prefix-G to split the window and open hostnav in a new vertical pane
tmux bind g split-window -v 'hostnav'

# Bind Prefix-G to split the window and open hostnav in a new horizontal pane
tmux bind g split-window -h 'hostnav'

```

These commands bind the key G (when pressed after the tmux prefix key, usually Ctrl-b) to either create a new window or split the current window, and then execute hostnav in the new pane or window​ ([Koen Woortman](https://koenwoortman.com/tmux-remap-split-window-keys/))​​ ([Stack Overflow](https://stackoverflow.com/questions/38967247/tmux-how-do-i-bind-function-keys-to-commands))​.




## Prerequisites

Before running hostnav, make sure the following commands are properly installed:

- Kubectl: For managing Kubernetes clusters.
- Docker: For managing and running containers.
- SSH: For remote server connections.
- tmux: For multi-panel functionality (if using tmux's multi-panel feature).


## Configuration File Example

Below is an example of the config.yaml configuration file with comments:

```yaml
# Cache configuration
cache: 
  # Cache save path
  directory: ./cache
  # Cache duration (minutes)
  duration: 30

# List of default plugins to use
default_plugins:
- docker
- kubernetes
- cvm
- jumpserver

# Plugin configurations
plugins:
  # Plugin name
  example:
    # Endpoint address
    endpoint: "https://api.example.com"
    # Access key
    access_key: "your-access-key"
    # Secret key
    access_key_secret: "your-secret-key"
    # Region list
    regions:
      - "us-west-1"
      - "us-east-1"
    # Connection configuration
    connections:
      # Use local IP
      use_local_ip: true

```


## Plugins

Below are the plugins supported by hostnav.

### ECS Alibaba Cloud plugin
```bash
hostnav -plugins ecs
```

Configuration：
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

### BCC Baidu Cloud Server Plugin

```bash
hostnav -plugins bcC
```

Configuration:

```yaml
plugins:
  bcc:
    # Access key
    access_key: "your-access-key"
    # Secret key
    access_key_secret: "your-secret-key"
    # Region list
    regions:
      - cn-bj
      - cn-bd
      - cn-gz
      - cn-su
      - cn-hkg
      - cn-fwh
      - cn-cd
      - cn-fsh
    # Connection configuration
    connections:
      # Use local IP
      use_local_ip: true 
```

### CVM Tencent Cloud Server Plugin

```bash
hostnav -plugins cvm
```

Configuration:

```yaml
plugins:
  cvm:
    # Access key
    access_key: "your-access-key"
    # Secret key
    access_key_secret: "your-secret-key"
    # Region list
    regions:
      - "ap-guangzhou"
      - "ap-shanghai"
    # Connection configuration
    connections:
      # Use local IP
      use_local_ip: true 

```

### Docker Container Plugin

```bash
hostnav -plugins docker
```

Configuration:
```yaml
plugins:
  docker:
    host: unix:///var/run/docker.sock
    version: "1.43"

```

### Kubernetes Container Plugin

```bash
hostnav -plugins kubernetes
```

Configuration:
```yaml
plugins:
  kubernetes:
    # Namespaces to query
    namespaces:
    - default
    - kube_system
    # Kubeconfig path
    kubeconfig: ~/.kube/config

```

### JumpServer Plugin

```bash
hostnav -plugins jumpserver
```


Configuration:
```yaml
plugins:
  jumpserver:
    # JumpServer address
    url: http://your-jumpserver-url
    access_key: your-api-access-key
    access_key_secret: your-api-access-key-secret
    # Asset filter
    search: xxxx
```

## Plugin Development Guide
### Role of Plugins
Plugins are the core of hostnav. Through the plugin system, users can extend the functionality of hostnav to support more cloud platforms, container platforms, and other resources. The main roles of plugins include:

- Resource Query: Retrieve resource information from different platforms.
- Terminal Connection: Provide functionality to connect to resources.
- Plugin Development: Users can develop custom plugins to extend the functionality of hostnav.


### Plugin Development Process

1. Create Plugin Directory

    Create a new plugin directory in the custom plugins directory, e.g., my_plugin.

2. Implement Plugin Interface

    Create a new Go file in the plugin directory, e.g., my_plugin.go, and implement the plugin interface.

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
3. Implement Plugin Logic

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
        // Initialize cloud platform client
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

4. Configure Plugin

    Add the custom plugin configuration in the hostnav configuration file.

    ```yaml
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

5. Compile and Run

```bash
go build -o hostnav .
./hostnav -config ~/.hostnav.yaml
```

## Contributing
Issues and pull requests are welcome. To start contributing, clone this repository and create a new branch for your changes.

```bash
git clone https://github.com/raojinlin/hostnav.git
cd hostnav
git checkout -b new-feature
```
