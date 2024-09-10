package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/armon/go-socks5"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// 定义命令行参数
	port := flag.Int("port", 1080, "端口号")
	username := flag.String("user", "", "用户名 (可选)")
	password := flag.String("pass", "", "密码 (可选)")
	flag.Parse()

	// 创建SOCKS5配置
	conf := &socks5.Config{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			log.Printf("尝试连接到: %s", addr)
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			// 解析 IP 地址（支持 IPv4 和 IPv6）
			ips, err := net.LookupIP(host)
			if err != nil {
				return nil, err
			}

			var ip net.IP
			for _, foundIP := range ips {
				ip = foundIP
				break
			}

			if ip == nil {
				return nil, fmt.Errorf("no IP address found for %s", host)
			}

			log.Printf("解析 %s 到 IP 地址: %s", host, ip.String())
			return net.DialTCP(network, nil, &net.TCPAddr{IP: ip, Port: parseInt(port)})
		},
	}

	// 如果提供了用户名和密码，则设置认证
	if *username != "" && *password != "" {
		creds := socks5.StaticCredentials{
			*username: *password,
		}
		auth := socks5.UserPassAuthenticator{Credentials: creds}
		conf.AuthMethods = []socks5.Authenticator{auth}
		log.Printf("启用用户名/密码认证")
	} else {
		log.Printf("未启用认证")
	}

	// 添加自定义的 DNS 解析器
	conf.Resolver = &customResolver{}

	// 创建SOCKS5服务器
	server, err := socks5.New(conf)
	if err != nil {
		log.Fatalf("创建SOCKS5服务器失败: %v", err)
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("SOCKS5服务器正在监听 %s", addr)
	if err := server.ListenAndServe("tcp", addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 自定义 DNS 解析器
type customResolver struct{}

func (r *customResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	log.Printf("正在解析域名: %s", name)
	ips, err := net.LookupIP(name)
	if err != nil {
		log.Printf("解析域名 %s 失败: %v", name, err)
		return ctx, nil, err
	}
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			log.Printf("域名 %s 解析结果: %s", name, ipv4.String())
			return ctx, ipv4, nil
		}
	}
	return ctx, nil, fmt.Errorf("no IPv4 address found for %s", name)
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}
