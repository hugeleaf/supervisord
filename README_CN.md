[English](README.md) | [中文](README_CN.md)

# 为什么有这个项目？

Python 版本的 supervisord 是一个被广泛使用的进程管理工具。

但它需要在目标系统中安装 Python 环境，在某些场景下，例如在 Docker 环境中，Python 对我们来说太大了。

这个项目使用 Go 语言重新实现了 supervisord。编译后的 supervisord 非常适合在没有安装 Python 的环境中使用。

# 构建 supervisord

在编译 supervisord 之前，请确保您的环境中已安装 Go 1.11+。

要为 **Linux** 编译 supervisord，请运行以下命令：

```bash
cd supervisord
GOOS=linux go build -a -ldflags "-linkmode external -extldflags -static" -o supervisord
```

或者从项目根目录构建：

```bash
go build -o supervisord ./supervisord
```

# 运行 supervisord

生成 supervisord 二进制文件后，创建 supervisord 配置文件并按如下方式启动 supervisord：

```Shell
$ cat supervisor.conf
[program:test]
command = /your/program args
$ supervisord -c supervisor.conf
```

请注意，配置文件位置按以下顺序自动检测：

1. $CWD/supervisord.conf
2. $CWD/etc/supervisord.conf
3. /etc/supervisord.conf
4. /etc/supervisor/supervisord.conf (since Supervisor 3.3.0)
5. ../etc/supervisord.conf (相对于可执行文件)
6. ../supervisord.conf (相对于可执行文件)


## 守护进程方式运行

启用 Web UI，在配置中添加 inet 接口，并以守护进程方式运行

```ini
$ cat supervisor.conf
[inet_http_server]
port=127.0.0.1:9001
$ supervisord -c supervisor.conf -d
```

为了管理守护进程，您可以使用 `supervisord ctl` 子命令，可用子命令有：`status`、`start`、`stop`、`shutdown`、`reload`。

```shell
$ supervisord ctl status
$ supervisord ctl status program-1 program-2...
$ supervisord ctl status group:*
$ supervisord ctl stop program-1 program-2...
$ supervisord ctl stop group:*
$ supervisord ctl stop all
$ supervisord ctl start program-1 program-2...
$ supervisord ctl start group:*
$ supervisord ctl start all
$ supervisord ctl shutdown
$ supervisord ctl reload
$ supervisord ctl signal <signal_name> <process_name> <process_name> ...
$ supervisord ctl signal all
$ supervisord ctl pid <process_name>
$ supervisord ctl fg <process_name>
```

请注意，`supervisor ctl` 子命令只有在 [inet_http_server] 中启用了 http 服务器并且正确设置了 **serverurl** 时才能正常工作。目前不支持 Unix 域套接字用于此目的。

Serverurl 参数按以下顺序检测：

- 检查是否存在选项 -s 或 --serverurl，使用此 url
- 检查是否存在 -c 选项，并且 "supervisorctl" 部分中存在 "serverurl"，使用 "supervisorctl" 部分中的 "serverurl"
- 检查在自动检测的 supervisord.conf 文件位置中是否定义了 "supervisorctl" 部分的 "serverurl"，如果是则使用找到的值
- 使用 http://localhost:9001

# 特性说明

## HTTP 服务配置

Http 服务器可以通过 Unix 域套接字和 TCP 工作。也支持可选的基本身份验证。

Unix 域套接字设置在 "unix_http_server" 部分，TCP http 服务器设置在 "inet_http_server" 部分。

如果配置文件中未设置 "inet_http_server" 和 "unix_http_server"，则不会启动 http 服务器。

以下参数可以在 `[inet_http_server]` 或 `[unix_http_server]` 部分中配置：

- **port** (仅 inet_http_server)。HTTP 服务器绑定的地址，例如 `:9001` 或 `127.0.0.1:9001`。
- **username**。HTTP 基本认证的用户名。可选。
- **password**。HTTP 基本认证的密码。可选。支持明文或 SHA1 哈希格式（前缀为 `{SHA}`）。
- **path_prefix**。所有 HTTP 端点的 URL 路径前缀。在反向代理后运行时很有用。默认为空（无前缀）。

带身份验证的示例：

```ini
[inet_http_server]
port=127.0.0.1:9001
username=admin
password=secret
```

使用 SHA1 哈希密码：

```ini
[inet_http_server]
port=127.0.0.1:9001
username=admin
password={SHA}5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8
```

配置身份验证后，访问 Web GUI 或 API 时浏览器会弹出认证窗口。

带路径前缀的示例：

```ini
[inet_http_server]
port=127.0.0.1:9001
path_prefix=/supervisord
```

使用上述配置，所有端点都将以 `/supervisord` 为前缀：
- Web GUI: `http://127.0.0.1:9001/supervisord/webgui/`
- REST API: `http://127.0.0.1:9001/supervisord/program/list`
- XML-RPC: `http://127.0.0.1:9001/supervisord/RPC2`
- Metrics: `http://127.0.0.1:9001/supervisord/metrics`

## 守护进程配置

在 "supervisord" 部分配置以下参数：

- **logfile**。supervisord 自身日志的存放位置。
- **logfile_maxbytes**。日志文件超过此长度后轮换。
- **logfile_backups**。保留的轮换日志文件数量。
- **loglevel**。日志详细程度，可以是 trace、debug、info、warning、error、fatal 和 panic（根据用于此功能的模块文档）。默认为 info。
- **pidfile**。包含当前 supervisord 实例进程 ID 的文件的完整路径。
- **minfds**。在 supervisord 启动时保留至少此数量的文件描述符。（Rlimit nofiles）。
- **minprocs**。在 supervisord 启动时保留至少此数量的进程资源。（Rlimit noproc）。
- **identifier**。此 supervisord 实例的标识符。如果在一台机器上的同一命名空间中运行多个 supervisord，则需要此参数。

## 被监督程序配置

被监督程序设置在 [program:programName] 部分中配置，包括以下选项：

- **command**。要监督的命令。可以作为可执行文件的完整路径给出，或者通过 PATH 变量计算。命令行参数也应该在此字符串中提供。
- **process_name**。进程名称
- **numprocs**。进程数量
- **numprocs_start**。??
- **autostart**。被监督命令是否应在 supervisord 启动时运行？默认为 **true**。
- **startsecs**。程序在启动后需要保持运行的总秒数，以认为启动成功（将进程从 STARTING 状态移动到 RUNNING 状态）。设置为 0 表示程序不需要保持运行任何特定时间。
- **startretries**。supervisord 在放弃并将进程置于 FATAL 状态之前，允许尝试启动程序的连续失败次数。有关 FATAL 状态的解释，请参阅进程状态。
- **autorestart**。如果被监督命令死亡，自动重新运行它。
- **exitcodes**。与此程序一起使用的 "预期" 退出代码列表。如果 autorestart 参数设置为 unexpected，并且进程以除 supervisor 停止请求之外的任何方式退出，如果进程以未在此列表中定义的退出代码退出，supervisord 将重新启动进程。
- **stopsignal**。发送给命令以优雅停止它的信号。如果配置了多个 stopsignal，停止程序时，supervisor 将按间隔 "stopwaitsecs" 逐个向程序发送信号。如果在向程序发送所有信号后程序仍未退出，supervisord 将终止程序。
- **stopwaitsecs**。在发送 SIGKILL 给被监督命令以使其不优雅停止之前等待的时间量。
- **stdout_logfile**。被监督命令的 STDOUT 应该重定向到哪里。（特定值在本文档后面描述）。
- **stdout_logfile_maxbytes**。超过此大小的日志将被轮换。
- **stdout_logfile_backups**。保留的轮换日志文件数量。
- **redirect_stderr**。是否应将 STDERR 重定向到 STDOUT。
- **stderr_logfile**。被监督命令的 STDERR 应该重定向到哪里。（特定值在本文档后面描述）。
- **stderr_logfile_maxbytes**。超过此大小的日志将被轮换。
- **stderr_logfile_backups**。保留的轮换日志文件数量。
- **environment**。要传递给被监督程序的 VARIABLE=value 列表。
- **priority**。程序在启动和关闭顺序中的相对优先级
- **user**。在执行被监督命令之前 sudo 到此 USER 或 USER:GROUP。
- **directory**。跳转到此路径并在那里执行被监督命令。
- **stopasgroup**。停止此程序所在程序组时，也停止此程序。
- **killasgroup**。停止此程序所在程序组时，也终止此程序。
- **restartpause**。在重新启动被监督程序之前等待（至少）这么多秒。
- **restart_when_binary_changed**。布尔值（false 或 true），控制当其可执行二进制文件更改时是否应重新启动被监督命令。默认为 false。
- **restart_cmd_when_binary_changed**。如果程序二进制文件本身更改，用于重新启动程序的命令。
- **restart_signal_when_binary_changed**。如果程序二进制文件更改，发送给程序以重新启动的信号。
- **restart_directory_monitor**。用于重新启动目的的监控路径。
- **restart_file_pattern**。如果 **restart_directory_monitor** 下的文件更改且文件名匹配此模式，被监督命令将重新启动。
- **restart_cmd_when_file_changed**。如果 **restart_directory_monitor** 下具有模式 **restart_file_pattern** 的任何监控文件更改，用于重新启动程序的命令。
- **restart_signal_when_file_changed**。如果 **restart_directory_monitor** 下具有模式 **restart_file_pattern** 的任何监控文件更改，将发送给程序的信号，例如 Nginx，用于重新启动。
- **depends_on**。定义被监督命令启动依赖关系。如果程序 A 依赖于程序 B、C，则程序 B、C 将在程序 A 之前启动。示例：

```ini
[program:A]
depends_on = B, C

[program:B]
...
[program:C]
...
```

### 为所有被监督程序设置默认参数

所有对被监督程序都相同的通用参数可以在 "program-default" 部分中定义一次，并在所有其他程序部分中省略。

在下面的示例中，VAR1 和 VAR2 环境变量适用于 test1 和 test2 被监督程序：

```ini
[program-default]
environment=VAR1="value1",VAR2="value2"

[program:test1]
...

[program:test2]
...

```

## 组

支持 "group" 部分，您可以设置 "programs" 项

## 事件

部分支持 Supervisord 3.x 定义的事件。现在它支持以下事件：

- 所有与进程状态相关的事件
- 进程通信事件
- 远程通信事件
- 与 tick 相关的事件
- 与进程日志相关的事件

## 日志

Supervisord 可以将被监督程序的 stdout 和 stderr（字段 stdout_logfile、stderr_logfile）重定向到：

- **/dev/null**：忽略日志 - 将其发送到 /dev/null。
- **/dev/stdout**：将日志写入 STDOUT。
- **/dev/stderr**：将日志写入 STDERR。
- **syslog**：将日志发送到本地 syslog 服务。
- **syslog @[protocol:]host[:port]**：将日志事件发送到远程 syslog 服务器。协议必须是 "tcp" 或 "udp"，如果缺失，则假定为 "udp"。如果端口缺失，对于 "udp" 协议，默认为 514，对于 "tcp" 协议，其值为 6514。
- **file name**：将日志写入指定文件。

可以为 stdout_logfile 和 stderr_logfile 配置多个日志文件，使用 ',' 作为分隔符。例如：

```ini
stdout_logfile = test.log, /dev/stdout
```

# Web GUI

Supervisord 有内置的 Web GUI：您可以从 GUI 启动、停止和检查程序状态。Web GUI 提供：

- **仪表板**：查看所有被管理程序及其当前状态（运行中/已停止）
- **进程控制**：一键启动/停止单个或多个程序
- **统计信息**：显示程序总数、运行中数量和已停止数量
- **Supervisor 管理**：重新加载配置或关闭 supervisor

请注意，要查看/使用 Web GUI，您应该在 /etc/supervisord.conf 中的 [inet_http_server]（如果您更喜欢 Unix 域套接字，则为 [unix_http_server]）和 [supervisorctl] 中配置它：

```ini
[inet_http_server]
port=127.0.0.1:9001
;username=test1
;password=thepassword

[supervisorctl]
serverurl=http://127.0.0.1:9001
```

# 从 Docker 容器中使用

supervisord 在 Docker 镜像中编译，可以直接在另一个镜像中使用，来自 Docker Hub 版本。

```Dockerfile
FROM debian:latest
COPY --from=ochinchina/supervisord:latest /usr/local/bin/supervisord /usr/local/bin/supervisord
CMD ["/usr/local/bin/supervisord"]
```

# 与 Prometheus 集成

Prometheus 节点导出器支持的 supervisord 指标现在已集成到 supervisor 中。因此不需要部署额外的 node_exporter 来收集 supervisord 指标。要收集指标，必须在 "inet_http_server" 部分中配置 port 参数，指标服务器在 supervisor http 服务器的 /metrics 路径上启动。

例如，如果 "inet_http_server" 中的 port 参数是 "127.0.0.1:9001"，则指标服务器应该在 url "http://127.0.0.1:9001/metrics" 中访问

# 检查版本

命令 "version" 将显示当前 supervisord 二进制文件的版本。

```shell
$ supervisord version
```

# 注册服务

操作系统启动后自动启动 supervisord。查看 [kardianos/service](https://github.com/kardianos/service) 支持的平台。

```Shell
# 安装
sudo supervisord service install -c full_path_to_conf_file
# 卸载
sudo supervisord service uninstall
# 启动
supervisord service start
# 停止
supervisord service stop
```
