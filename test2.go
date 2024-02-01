package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LoggingTransport struct {
	Transport http.RoundTripper
}

func (lt *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 打印请求信息
	fmt.Printf("转发请求: %s %s\n", req.Method, req.URL)

	// 使用默认的 Transport 执行实际的请求
	resp, err := lt.Transport.RoundTrip(req)

	if err != nil {
		fmt.Println("请求失败:", err)
		return nil, err
	}

	// 打印响应信息
	fmt.Printf("收到响应: %s\n", resp.Status)

	// 打印响应体
	body, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println("无法读取响应体:", err)
	} else {
		fmt.Printf("响应体:\n%s\n", body)
	}

	return resp, nil
}

func main() {
	// 目标服务器的地址
	targetURL, err := url.Parse("http://13.215.101.156:8000/")
	if err != nil {
		fmt.Println("无法解析目标URL:", err)
		return
	}

	// 创建反向代理器
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 使用自定义 Transport
	proxy.Transport = &LoggingTransport{
		Transport: http.DefaultTransport,
	}

	// 使用自定义 Director 修改请求
	proxy.Director = func(req *http.Request) {
		// 在这里可以对请求进行修改
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Host = targetURL.Host
	}

	// 创建HTTP服务器并监听本地端口
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 使用反向代理器处理请求
		proxy.ServeHTTP(w, r)
	})

	// 启动HTTP服务器
	fmt.Println("反向代理服务器正在监听端口 55555")
	err = http.ListenAndServe(":55555", nil)
	if err != nil {
		fmt.Println("启动HTTP服务器失败:", err)
	}
}
