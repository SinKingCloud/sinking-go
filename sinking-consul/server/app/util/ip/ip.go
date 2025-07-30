package ip

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// GetLocalIP 获取本机IP地址，优先选择公网IP，超时时间5秒
func GetLocalIP() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	type result struct {
		ip  string
		err error
	}
	ch := make(chan result, 1)
	go func() {
		// 策略1: 通过连接外网确定IP
		conn, err := net.DialTimeout("udp", "8.8.8.8:53", 2*time.Second)
		if err == nil {
			defer conn.Close()
			if localAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
				ip := localAddr.IP.String()
				if net.ParseIP(ip) != nil && !strings.Contains(ip, ":") {
					ch <- result{ip, nil}
					return
				}
			}
		}
		// 策略2: 遍历网络接口查找非回环IP
		faces, err := net.Interfaces()
		if err == nil {
			for _, face := range faces {
				if face.Flags&net.FlagUp == 0 {
					continue
				}
				if face.Flags&net.FlagLoopback != 0 {
					continue
				}

				adders, err := face.Addrs()
				if err != nil {
					continue
				}

				for _, addr := range adders {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					default:
						continue
					}

					if ip == nil || ip.To4() == nil {
						continue
					}

					ipStr := ip.String()
					if net.ParseIP(ipStr) != nil && !strings.Contains(ipStr, ":") {
						ch <- result{ipStr, nil}
						return
					}
				}
			}
		}
		// 策略3: 使用默认路由接口
		conn, err = net.DialTimeout("udp", "google.com:80", 2*time.Second)
		if err == nil {
			defer conn.Close()
			if localAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
				ip := localAddr.IP.String()
				if net.ParseIP(ip) != nil && !strings.Contains(ip, ":") {
					ch <- result{ip, nil}
					return
				}
			}
		}
		ch <- result{"", errors.New("无法获取有效的本地IP地址")}
	}()
	var (
		ip  string
		err error
	)
	select {
	case res := <-ch:
		ip, err = res.ip, res.err
	case <-ctx.Done():
		ip, err = "", ctx.Err()
	}
	if err != nil || ip == "" {
		return "127.0.0.1"
	} else {
		return ip
	}
}

// IsValidIP 判断IP是否有效(非空且为IPv4或IPv6)
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidIPv4 判断IP是否为有效的IPv4
func IsValidIPv4(ip string) bool {
	return net.ParseIP(ip) != nil && strings.Count(ip, ":") == 0
}

// IsValidIPv6 判断IP是否为有效的IPv6
func IsValidIPv6(ip string) bool {
	return net.ParseIP(ip) != nil && strings.Count(ip, ":") > 0
}

// IsPrivateIP 判断IP是否为私有地址
func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// IPv4私有地址范围
	if ip4 := parsedIP.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}

	// IPv6链路本地地址
	return strings.HasPrefix(ip, "fe80:") || strings.HasPrefix(ip, "fc00:")
}

// IsPublicIP 判断IP是否为公网地址
func IsPublicIP(ip string) bool {
	return IsValidIP(ip) && !IsPrivateIP(ip) && !IsLoopBack(ip)
}

// IsLoopBack 判断IP是否为回环地址(127.0.0.1/8 或 ::1)
func IsLoopBack(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.IsLoopback()
}

// ToInt 将IPv4地址转换为整数
func ToInt(ip string) (uint32, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0, errors.New("无效的IP地址")
	}

	ip4 := parsedIP.To4()
	if ip4 == nil {
		return 0, errors.New("不是有效的IPv4地址")
	}

	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3]), nil
}

// ToIP 将整数转换为IPv4地址
func ToIP(ipInt uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		ipInt>>24,
		(ipInt>>16)&0xFF,
		(ipInt>>8)&0xFF,
		ipInt&0xFF)
}

// GetAllLocalIPs 获取所有非回环IP地址
func GetAllLocalIPs() ([]string, error) {
	faces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, face := range faces {
		if face.Flags&net.FlagUp == 0 {
			continue
		}
		if face.Flags&net.FlagLoopback != 0 {
			continue
		}

		adders, err := face.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range adders {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}

			if ip == nil {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	if len(ips) == 0 {
		return nil, errors.New("未找到非回环IP地址")
	}

	return ips, nil
}

// GetLocalIPv4s 获取所有非回环IPv4地址
func GetLocalIPv4s() ([]string, error) {
	allIPs, err := GetAllLocalIPs()
	if err != nil {
		return nil, err
	}
	var ipv4s []string
	for _, ip := range allIPs {
		if IsValidIPv4(ip) {
			ipv4s = append(ipv4s, ip)
		}
	}
	if len(ipv4s) == 0 {
		return nil, errors.New("未找到非回环IPv4地址")
	}
	return ipv4s, nil
}

// GetLocalIPv6s 获取所有非回环IPv6地址
func GetLocalIPv6s() ([]string, error) {
	allIPs, err := GetAllLocalIPs()
	if err != nil {
		return nil, err
	}
	var ipv6s []string
	for _, ip := range allIPs {
		if IsValidIPv6(ip) {
			ipv6s = append(ipv6s, ip)
		}
	}
	if len(ipv6s) == 0 {
		return nil, errors.New("未找到非回环IPv6地址")
	}
	return ipv6s, nil
}
