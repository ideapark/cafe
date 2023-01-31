# Coffee

[English](README.md)

通过 ssh 隧道，将远程网络的 http(s) 服务中转成本地 http 服务。

![](https://user-images.githubusercontent.com/49527198/215134049-1aad187c-1d24-4b00-a006-e08f9278ab55.png)

## 使用场景

### 场景 #1: 打破微服务架构服务依赖环

- 网络拓扑

> 本地：任何网络（工作电脑）
> 远程：云上私有网络

`coffee` 用来将远程网络 http(s) 服务中转成本地 http 服务，这样您可以开
始本地编码和调试，无需在本地运行所有依赖的远程 http(s) 服务，因为这些
服务很难甚至不可能在本地运行起来：

1. 这些远程服务又有自己的服务依赖或者中间件依赖。
2. 需要消耗极大的资源（cpu,内存）才能将他们全部在本地运行起来。
3. 需要耗费极大的精力来确保它们持续可用（试想依赖服务需要持续更新到新的版本）。

总之一句话：你仅仅想这些服务直接稳定的随时被使用，这样你自己才能专注在
自己的服务开发调试等工作上。

### 场景 #2: 极简版自建 VPN 服务

- 网络拓扑

> 本地：中国大陆
> 远程：全球其他网络

假设你有一个服务器跑在 AWS 上，并且它能够访问任何公网服务，例如 Google
（其他的云厂商也同样可行）。你能够 ssh 到这台服务器，在你本地运行
`coffee`，配置它 `https://www.google.com` 中继到本地
`http://www.google.com.127.0.0.1.nip.io`。你就能够访问 Google 了。

反之亦然，你也能把 `coffee` 运行在远程网络服务器上，让通配 DNS 解析到
这个服务器的公网 ip 上，你就有了一个公共托管的 vpn 服务器了。

请明智的使用它！

## 运行效果

> 在以下运行过程，我假设你本地已经起用了 `sshd` 服务，并且能够通过
> `password` 或者 `publickey` 认证方式成功的 `ssh localhost`。

- 启动参数

如果不指定端口，`coffee` 将使用默认 `2046` 端口提供本地中继服务。默认
打印所有的请求响应对象到标准日志输出.

```bash
$ ./coffee --help
Usage of coffee:
  -conf string
    	use customized configuration coffee.json
  -port int
    	use another serving port (default 2046)
  -trace
    	trace every http roundtrip object (default true)
  -version
    	print coffee version
  -view
    	print default builtin configuration coffee.json (as start point of customization)
```

- 运行观察

启动 `coffee` 后，他将开始中继您配置的远程网络 http(s) 服务到本地。

```bash
$ ./coffee
coffee-v0.0.3
#relay  | Remote http(s)           | Local http
------  | --------------           | ----------
1       | http://www.vulnweb.com   | http://www.vulnweb.com.127.0.0.1.nip.io:2046
2       | http://rest.vulnweb.com  | http://rest.vulnweb.com.127.0.0.1.nip.io:2046
3       | https://www.gnu.org      | http://www.gnu.org.127.0.0.1.nip.io:2046
4       | https://kernel.org       | http://kernel.org.127.0.0.1.nip.io:2046
5       | https://go.dev           | http://go.dev.127.0.0.1.nip.io:2046

2023/01/26 20:10:07 #1
GET / HTTP/1.1
Host: www.gnu.org
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate
Accept-Language: en,zh-CN;q=0.9,zh;q=0.8
Connection: keep-alive
Cookie: _ga=GA1.2.630524832.1673162343; _gcl_au=1.1.466449589.1673411024
Purpose: prefetch
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36


2023/01/26 20:10:07 establishing tunnel connection...
2023/01/26 20:10:08 #1
HTTP/1.1 200 OK
Content-Length: 9911
Accept-Ranges: bytes
Access-Control-Allow-Origin: (null)
Cache-Control: max-age=0
Connection: Keep-Alive
Content-Encoding: gzip
Content-Language: en
Content-Location: home.html
Content-Type: text/html
Date: Thu, 26 Jan 2023 12:10:08 GMT
Expires: Thu, 26 Jan 2023 12:10:08 GMT
Keep-Alive: timeout=5, max=100
Server: Apache/2.4.29
Strict-Transport-Security: max-age=63072000
Tcn: choice
Vary: negotiate,accept-language,Accept-Encoding
X-Content-Type-Options: nosniff
X-Frame-Options: sameorigin
```

- curl 请求将被中继到远程网络

以下是一个 `www.gnu.org` 本地中继后的情况。

```bash
$ curl -I http://www.gnu.org.127.0.0.1.nip.io:2046
HTTP/1.1 200 OK
Accept-Ranges: bytes
Access-Control-Allow-Origin: (null)
Cache-Control: max-age=0
Content-Language: en
Content-Location: home.html
Content-Type: text/html
Date: Thu, 26 Jan 2023 12:11:06 GMT
Expires: Thu, 26 Jan 2023 12:11:06 GMT
Server: Apache/2.4.29
Strict-Transport-Security: max-age=63072000
Tcn: choice
Vary: negotiate,accept-language,Accept-Encoding
X-Content-Type-Options: nosniff
X-Frame-Options: sameorigin
```

## Coffee 用户

怎样让 `coffee` 在我的环境运行。

```json
{
  "wild": ".127.0.0.1.nip.io",
  "urls": [
    "http://www.vulnweb.com",
    "http://rest.vulnweb.com",
    "https://www.gnu.org",
    "https://kernel.org"
  ],
  "hops": [
    {
      "host": "127.0.0.1",
      "port": "22",
      "user": "ENV:USER",
      "pass": "ENV:PASS1",
      "key": "FILE:~/.ssh/id_rsa"
    },
    {
      "host": "localhost",
      "port": "22",
      "user": "ENV:USER",
      "pass": "ENV:PASS2",
      "key": "FILE:~/.ssh/id_rsa"
    }
  ]
}
```

- 编辑 `coffee.json`

1. 用你的远程网络 http(s) 地址替换 `urls` 数组。
2. 配置 `coffee` 怎样一跳一跳的到达你的远程网络 (为保护你的隐私，`user`, `pass`, `key` 支持环境变量或者文件读入)。
3. `wild` 支持替换成任何可以解析到 `127.0.0.1` 的通配符域名后缀。

- 编译

编译最终的单二进制可执行文件，祝你使用 `coffee` 有一个愉快的过程！

```bash
make
```

## Coffee 开源贡献者

### 理解 `coffee` 的网络流

```text
              LOCAL NETWORK                                                                                    REMOTE NETWORK
┌────────────────────────────────────────────────┐                                                      ┌───────────────────────────┐
│ [http://www.gnu.org.127.0.0.1.nip.io:2046]     │                                                      │ [https://www.gnu.org]     │
│ [http://kernel.org.127.0.0.1.nip.io:2046]      │                                                      │ [https://kernel.org]      │
│ [http://go.dev.127.0.0.1.nip.io:2046]          │                                                      │ [https://go.dev]          │
│                                                │                     sshd          sshd               │         sshd              │
│             [CURL].......................(coffee:2046)..............(hop.1).......(hop.2)......................(hop.n)            │
│                                                │                                                      │                           │
│ [http://www.vulnweb.com.127.0.0.1.nip.io:2046] │                                                      │ [http://www.vulnweb.com]  │
│ [http://rest.vulnweb.com.127.0.0.1.nip.io:2046]│                                                      │ [http://rest.vulnweb.com] │
└────────────────────────────────────────────────┘                                                      └───────────────────────────┘
                 ↑                               ↑                                                                  ↑
                 └─────────── http ──────────────┴─────────────────────── tcp: ssh tunnel ──────────────────────────┘
```

### 技术参考引用

- [The Secure Shell (SSH) Protocol Architecture](https://www.rfc-editor.org/rfc/rfc4251)
- [The Secure Shell (SSH) Authentication Protocol](https://www.rfc-editor.org/rfc/rfc4252)
- [The Secure Shell (SSH) Transport Layer Protocol](https://www.rfc-editor.org/rfc/rfc4253)
- [The Secure Shell (SSH) Connection Protocol](https://www.rfc-editor.org/rfc/rfc4254)
- [The Secure Shell (SSH) Protocol Assigned Numbers](https://www.rfc-editor.org/rfc/rfc4250)
- [Hypertext Transfer Protocol -- HTTP/1.1](https://www.rfc-editor.org/rfc/rfc2616)
- [Go SSH](https://pkg.go.dev/golang.org/x/crypto/ssh)
