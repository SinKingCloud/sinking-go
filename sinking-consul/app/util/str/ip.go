package str

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

func GetInternalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", errors.New("internal IP fetch failed, detail:" + err.Error())
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	res := conn.LocalAddr().String()
	res = strings.Split(res, ":")[0]
	return res, nil
}

func GetExternalIP() (string, error) {
	response, err := http.Get("http://ip.dhcp.cn/?ip")
	if err != nil {
		return "", errors.New("external IP fetch failed, detail:" + err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	res := ""
	for {
		tmp := make([]byte, 32)
		n, err := response.Body.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return "", errors.New("external IP fetch failed, detail:" + err.Error())
			}
			res += string(tmp[:n])
			break
		}
		res += string(tmp[:n])
	}
	return strings.TrimSpace(res), nil
}
