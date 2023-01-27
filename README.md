# Coffee

[中文](README.zh.md)

Remote network http(s) are relayed as local http over ssh tunnel.

![](https://user-images.githubusercontent.com/49527198/215134049-1aad187c-1d24-4b00-a006-e08f9278ab55.png)

## Use Scenario

### Case #1: Break microservice architecture dependencies cycle

- Network Topology:

> Local: any network (laptop)
> Remote: private network (cloud)

`coffee` was born out to be used to make `remote network` http(s)
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
and you have ssh access to this server. Run `coffee` at your local and
configure it with `https://www.google.com` relayed as your local http
`http://www.google.com.127.0.0.1.nip.io`. now you are free to access
google now.

Vice verse, you could also run `coffee` at the remote network server,
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
$ ./coffee --help
Usage of coffee:
  -conf string
        filepath to coffee.json
  -port int
        use another serving port (default 2046)
  -trace
        trace every http roundtrip object (default true)
  -version
        print coffee version
  -view
        print default coffee.json
```

- run it & keep watching

Start `coffee`, it will start relaying the urls you have configured to
the remote network.

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

- curl request will be relayed to the remote network

Have a try for one of the relayed urls.

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

## Coffee Users

How to make it work for my environment.

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

- Update `coffee.json`

1. replace `urls` array with your service urls at the remote network.
2. instruct `coffee` how to reach your remote network hop by
   hop. (`user`, `pass`, `key` are kept secure by reading from env or
   file, optional)
3. `wild` suffix can be changed to any public wild dns that resovles
   to `127.0.0.1` (optional)

- Build

Compile the final single binary release. Happy `coffee`

```bash
make
```

## Coffee Contributors

### Understanding Coffee Traffic Flow

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

### Technical references

- [The Secure Shell (SSH) Protocol Architecture](https://www.rfc-editor.org/rfc/rfc4251)
- [The Secure Shell (SSH) Authentication Protocol](https://www.rfc-editor.org/rfc/rfc4252)
- [The Secure Shell (SSH) Transport Layer Protocol](https://www.rfc-editor.org/rfc/rfc4253)
- [The Secure Shell (SSH) Connection Protocol](https://www.rfc-editor.org/rfc/rfc4254)
- [The Secure Shell (SSH) Protocol Assigned Numbers](https://www.rfc-editor.org/rfc/rfc4250)
- [Hypertext Transfer Protocol -- HTTP/1.1](https://www.rfc-editor.org/rfc/rfc2616)
- [SSH Protocol by Go](https://pkg.go.dev/golang.org/x/crypto/ssh)
