# Coffee

[中文](README.zh.md)

Remote network http(s) are relayed as local http over ssh tunnel.

## Use Scenario

### Case #1: Connecting to remote private network services from local

- Network Topology:

> Local: any network (laptop)
> Remote: private network (cloud)

`cafe` was born out to be used to make `remote network` http(s)
services available at local, so you're able to start your local coding
and debugging, without spinning up all the remote dependency services
at local, because these dependency http(s) services are hard or even
impossible to run at local:

1. they have their own dependency services or middlewares
2. too much resources (cpu, memory) are needed to run them all
3. too much effort to ensure them always work (e,g. new versions come
   out with new apis)

In one word, you just want all the dependency http(s) services out
there and always ready to be comsumed. so you could focus on your own
business.

### Case #2: Global network http(s) websites are relayed at local

- Network Topology:

> Local: china network (GFW protected)
> Remote: global network (Internet)

Suppose you have a server running at AWS and it has free access to any
public internet such google (other cloud provider is working as well),
and you have ssh access to this server. Run `cafe` at your local and
configure it with `https://www.google.com` relayed as your local http
`http://www.google.com.local.gd`. now you are free to access google
now.

Vice verse, you could also run `cafe` at the remote network server,
and make the wild dns resolve to this server public ip. Now it's a
public managed vpn like services.

Use it wisely!

## How it looks like

> I would assume you have `sshd` enabled at your localhost, and you
> can successfully `ssh localhost` with either password or publickey
> auth.

- help

The default port `2046` will be used if `-port` not specified, and
every http roundtripping object will be logged out to stdout.

```bash
$ ./cafe --help
Usage of cafe:
  -conf string
    	use customized configuration cafe.json
  -port int
    	use another serving port (default 2046)
  -trace
    	trace every http roundtrip object (default true)
  -version
    	print cafe version
  -view
    	print default builtin configuration cafe.json (as start point of customization)
```

- run it & keep watching

Start `cafe`, it will start relaying the urls you have configured to
the remote network.

```bash
$ ./cafe
cafe-v0.0.5
#relay  | Remote http(s)           | Local http
------  | --------------           | ----------
1       | http://www.vulnweb.com   | http://www.vulnweb.com.local.gd:2046
2       | http://rest.vulnweb.com  | http://rest.vulnweb.com.local.gd:2046
3       | https://www.gnu.org      | http://www.gnu.org.local.gd:2046
4       | https://kernel.org       | http://kernel.org.local.gd:2046
5       | https://go.dev           | http://go.dev.local.gd:2046

2023/03/02 21:19:35 #1
GET / HTTP/1.1
Host: www.gnu.org
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Accept-Encoding: gzip, deflate
Accept-Language: en,zh-CN;q=0.9,zh;q=0.8
Connection: keep-alive
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36


2023/03/02 21:19:35 establishing tunnel connection...
2023/03/02 21:19:37 #1
HTTP/1.1 200 OK
Content-Length: 9653
Accept-Ranges: bytes
Access-Control-Allow-Origin: (null)
Cache-Control: max-age=0
Connection: Keep-Alive
Content-Encoding: gzip
Content-Language: en
Content-Location: home.html
Content-Type: text/html
Date: Thu, 02 Mar 2023 13:19:36 GMT
Expires: Thu, 02 Mar 2023 13:19:36 GMT
Keep-Alive: timeout=5, max=100
Server: Apache/2.4.29
Strict-Transport-Security: max-age=63072000
Tcn: choice
Vary: negotiate,accept-language,Accept-Encoding
X-Content-Type-Options: nosniff
X-Frame-Options: sameorigin
```

- curl request will be relayed to the remote network

Have a try for one of the relayed urls.

```bash
$ curl -I http://www.gnu.org.local.gd:2046
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

## Coffee Users

How to make it work for my environment.

```json
{
  "wild": ".local.gd",
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

- Update `cafe.json`

1. replace `urls` array with your service urls at the remote network.
2. instruct `cafe` how to reach your remote network hop by
   hop. (`user`, `pass`, `key` are kept secure by reading from env or
   file, optional)
3. `wild` suffix can be changed to any public wild dns that resovles
   to `127.0.0.1` (optional)

- Build

Compile the final single binary release. Happy `cafe`

```bash
make
```

## Coffee Contributors

### Understanding Coffee Traffic Flow

```text
              LOCAL NETWORK                                                                REMOTE NETWORK
┌────────────────────────────────────────┐                                          ┌──────────────────────────┐
│ [http://www.gnu.org.local.gd:2046]     │                                          │ [https://www.gnu.org]    │
│ [http://kernel.org.local.gd:2046]      │                                          │ [https://kernel.org]     │
│ [http://go.dev.local.gd:2046]          │                                          │ [https://go.dev]         │
│                                        │              sshd          sshd          │       sshd               │
│             [CURL]..............(cafe:2046)..........(hop.1).......(hop.2)...............(hop.n)             │
│                                        │                                          │                          │
│ [http://www.vulnweb.com.local.gd:2046] │                                          │ [http://www.vulnweb.com] │
│ [http://rest.vulnweb.com.local.gd:2046]│                                          │ [http://rest.vulnweb.com]│
└────────────────────────────────────────┘                                          └──────────────────────────┘
                 ↑                       ↑                                                       ↑
                 └───────── http ────────┴───────────────── tcp: ssh tunnel ─────────────────────┘
```

### Technical references

- [The Secure Shell (SSH) Protocol Architecture](https://www.rfc-editor.org/rfc/rfc4251)
- [The Secure Shell (SSH) Authentication Protocol](https://www.rfc-editor.org/rfc/rfc4252)
- [The Secure Shell (SSH) Transport Layer Protocol](https://www.rfc-editor.org/rfc/rfc4253)
- [The Secure Shell (SSH) Connection Protocol](https://www.rfc-editor.org/rfc/rfc4254)
- [The Secure Shell (SSH) Protocol Assigned Numbers](https://www.rfc-editor.org/rfc/rfc4250)
- [Hypertext Transfer Protocol -- HTTP/1.1](https://www.rfc-editor.org/rfc/rfc2616)
- [Go SSH](https://pkg.go.dev/golang.org/x/crypto/ssh)
