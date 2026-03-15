package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var (
	listenAddr string
	proxyAddr  string
	targetAddr string
)

var rootCmd = &cobra.Command{
	Use:   "proxy-forwarder",
	Short: "HTTP代理转发服务",
	Long:  "通过指定的HTTP代理转发请求到目标HTTPS地址",
	Run:   runProxy,
}

func init() {
	rootCmd.Flags().StringVarP(&listenAddr, "listen", "l", ":8080", "监听地址和端口 (例如 :8080)")
	rootCmd.Flags().StringVarP(&proxyAddr, "proxy", "p", "", "HTTP代理地址 (例如 http://127.0.0.1:7890)")
	rootCmd.Flags().StringVarP(&targetAddr, "target", "t", "", "目标HTTPS地址 (例如 https://api.openai.com)")

	rootCmd.MarkFlagRequired("proxy")
	rootCmd.MarkFlagRequired("target")
}

func runProxy(cmd *cobra.Command, args []string) {
	// 解析目标URL
	targetURL, err := url.Parse(targetAddr)
	if err != nil {
		log.Fatalf("解析目标地址失败: %v", err)
	}

	// 解析代理URL
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		log.Fatalf("解析代理地址失败: %v", err)
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 配置Transport使用HTTP代理
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// 自定义错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("代理错误: %v", err)
		http.Error(w, fmt.Sprintf("代理错误: %v", err), http.StatusBadGateway)
	}

	// 修改请求以正确设置主机头
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
		req.Header.Set("X-Forwarded-Host", req.Host)
	}

	// 设置路由处理
	http.Handle("/", proxy)

	log.Printf("代理服务启动在 %s", listenAddr)
	log.Printf("通过代理 %s 转发到 %s", proxyAddr, targetAddr)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
