# Reverse Proxy

This Program helps you set up a server that provides reverse proxy service.
Typically, this can be used to bypass the ICP filling requirements in China.

这个程序可以方便地提供反向代理服务，特别适合提供用于在中国绕过工信部备案的服务。

In brief, a reverse proxy service is that, you visit a reverse proxy server,
and the reverse proxy server will response by retrieving the data you requested from the real destination.
To bypass ICP filling requirements, just let the reverse proxy server retrieves data by ip and port,
instead of domain name.

大概地讲，反向代理服务就是你访问反向代理服务器，反向代理服务器帮你向你要访问的网站请求数据，
再将数据返回给你。要绕过备案，只要让反向代理服务器请求的时候不是用域名，而是用 IP 地址和端口就可以了。

Using this program, you can provide a public service that users are able to register by themselves on your website.
All they need to do is to fill in the host-ip, port, domain name, and set up a DNS record to the reverse proxy server.

用户可以利用这个网站提供的服务，自行注册反向代理。
只需要填入要被代理的网站的 IP、端口、域名，然后添加一条 DNS 记录解析到反向代理服务器即可。

On our server we hold a demo service, feel free to [try](http://cool2645.pub)!

在我们的服务器上有这个程序，你可以[直接尝试它](http://cool2645.pub)！


## Build

```
go get github.com/2645Corp/reverse-proxy
cd $GOPATH/bin
cp $GOPATH/src/github.com/2645Corp/reverse-proxy/config.toml.example ./config.toml
# Now you modify the config.toml
cp -r $GOPATH/src/github.com/2645Corp/reverse-proxy/tmpl ./
nohup ./reverse-proxy &
```

其中 `proxy_port` 为反向代理服务器所在的端口，通常为 80。
`manage_port` 为服务注册网站所在的端口，当然了，你可以使用反向代理服务来为它分配一个域名。

## Tips

+ 把一个一级域名泛解析到这个服务器，然后把反向代理服务开在 80 端口，就可以为你的用户分配二级域名了。
+ 用户注册服务的逻辑类似 Github，只是网站不是 host 在我们这里，而是反向代理了一个用户指定的网站。用户注册服务后会给他分一个二级域名，可以使用这个域名来访问。也可以自定义域名，自定义域名需要将这个域名在服务管理页面中填写上，并且添加域名到反向代理服务器的解析。