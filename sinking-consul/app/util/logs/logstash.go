package logs

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Logstash struct {
	Hostname   string
	Port       int
	Connection *net.TCPConn
	Timeout    int
}

func New(hostname string, port int, timeout int) *Logstash {
	l := Logstash{}
	l.Hostname = hostname
	l.Port = port
	l.Connection = nil
	l.Timeout = timeout
	return &l
}

func (l *Logstash) Dump() {
	fmt.Println("Hostname:   ", l.Hostname)
	fmt.Println("Port:       ", l.Port)
	fmt.Println("Connection: ", l.Connection)
	fmt.Println("Timeout:    ", l.Timeout)
}

func (l *Logstash) SetTimeouts() {
	deadline := time.Now().Add(time.Duration(l.Timeout) * time.Millisecond)
	err := l.Connection.SetDeadline(deadline)
	if err != nil {
		return
	}
	err = l.Connection.SetWriteDeadline(deadline)
	if err != nil {
		return
	}
	err = l.Connection.SetReadDeadline(deadline)
	if err != nil {
		return
	}
}

func (l *Logstash) Connect() (*net.TCPConn, error) {
	var connection *net.TCPConn
	service := fmt.Sprintf("%s:%d", l.Hostname, l.Port)
	addr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		return connection, err
	}
	connection, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return connection, err
	}
	if connection != nil {
		l.Connection = connection
		err := l.Connection.SetLinger(0)
		if err != nil {
			return nil, err
		} // default -1
		err = l.Connection.SetNoDelay(true)
		if err != nil {
			return nil, err
		}
		err = l.Connection.SetKeepAlive(true)
		if err != nil {
			return nil, err
		}
		err = l.Connection.SetKeepAlivePeriod(time.Duration(5) * time.Second)
		if err != nil {
			return nil, err
		}
		l.SetTimeouts()
	}
	return connection, err
}

func (l *Logstash) Writeln(message string) error {
	err := errors.New("TCP Connection is nil.")
	message = fmt.Sprintf("%s\n", message)
	if l.Connection == nil {
		l.Connection, err = l.Connect()
		if err != nil {
			log.Println("重连日志服务器失败")
		}
	}
	if l.Connection != nil {
		_, err = l.Connection.Write([]byte(message))
		if err != nil {
			l.Connection, err = l.Connect()
			_, err = l.Connection.Write([]byte(message))
		}
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				err = l.Connection.Close()
				if err != nil {
					return err
				}
				l.Connection = nil
				if err != nil {
					return err
				}
			} else {
				err = l.Connection.Close()
				if err != nil {
					return err
				}
				l.Connection = nil
				return err
			}
		} else {
			// Successful write! Let's extend the timeoul.
			l.SetTimeouts()
			return nil
		}
	}
	return err
}
