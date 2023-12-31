package socks_proxy_server

import (
	"net"
	"proxy-forward/internal/service/user_token_service"
	"proxy-forward/pkg/logging"
	"time"
)

type Request struct {
	tcpGram    TCPProtocol
	remoteAddr string

	ClientConn *net.TCPConn
	RemoteConn *net.TCPConn

	server *Server
}

type RequestList struct {
	Prev *RequestList
	Data Request
	Next *RequestList
}

func (r *Request) Close() error {
	var err error
	if r.ClientConn != nil {
		er := r.ClientConn.Close()
		if er != nil {
			err = er
		}
	}
	if r.RemoteConn != nil {
		er := r.RemoteConn.Close()
		if er != nil {
			err = er
		}
	}

	return err
}

func (r *Request) Process() {
	r.tcpGram.authHandle = Auth
	if err := r.tcpGram.handClientshake(r.ClientConn); err != nil {
		return
	}

	remoteAddr, username, password, err := r.tcpGram.networkString()
	if err != nil {
		_, _ = r.ClientConn.Write([]byte{Version, 0x01, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}
	if conn, err := net.DialTimeout("tcp", remoteAddr, time.Second*time.Duration(r.server.writeTimeout)); err != nil {
		_, _ = r.ClientConn.Write([]byte{Version, 0x01, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	} else {
		r.RemoteConn = conn.(*net.TCPConn)
	}
	if err := r.tcpGram.handRemoteshke(r.RemoteConn, r.ClientConn, username, password, remoteAddr); err != nil {
		return
	}
	r.remoteAddr = remoteAddr
	r.transformTCP()

	/*
		if !r.tcpGram.viaUPD { // tcp
			bindIP := r.ClientConn.LocalAddr().(*net.TCPAddr).IP
			res := make([]byte, 0, 22)
			if ip := bindIP.To4(); ip != nil {
				// IPv4, len is 4
				res = append(res, []byte{Version, 0x00, 0x00, ATYPIPv4}...)
				res = append(res, ip...)
			} else {
				// IPv6, len is 16
				res = append(res, []byte{Version, 0x00, 0x00, ATYPIPv6}...)
				res = append(res, bindIP...)
			}
			portByte := [2]byte{}
			binary.BigEndian.PutUint16(portByte[:], uint16(r.ClientConn.LocalAddr().(*net.TCPAddr).Port))
			res = append(res, portByte[:]...)
			if _, err := r.ClientConn.Write(res); err != nil {
				return
			}
			r.transformTCP()
		} else {
			// bind UDP addr and answer
			if !r.server.enableUPD {
				_, _ = r.ClientConn.Write([]byte{Version, 0x07, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
				return
			}
			// 暂时不支持
		}
	*/
}

func (r *Request) transformTCP() {
	logging.Log.Infof("[%s]connect to: %s", "tcp", r.remoteAddr)

	userTokenService := user_token_service.UserToken{
		ID:              r.tcpGram.userToken.ID,
		ReqUsageAmount:  r.tcpGram.userToken.ReqUsageAmount,
		RespUsageAmount: r.tcpGram.userToken.RespUsageAmount,
		Uid:             r.tcpGram.userToken.Uid,
		//LaID:            r.tcpGram.userToken.LaID,
		IsStatic:   r.tcpGram.userToken.IsStatic,
		DataCenter: r.tcpGram.userToken.DataCenter,
		//PsID:            r.tcpGram.userToken.PsID,
	}
	done := make(chan int)
	go func() {
		defer func() { _ = r.Close(); done <- 0 }()
		buf := make([]byte, 1024*8)
		for {
			_ = r.RemoteConn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(r.server.readTimeout)))
			if ln, err := r.RemoteConn.Read(buf); err != nil {
				break
			} else {
				_ = r.ClientConn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(r.server.writeTimeout)))
				if _, err := r.ClientConn.Write(buf[:ln]); err != nil {
					break
				}
				userTokenService.IncrRespBytes(r.remoteAddr, ln)
				// userTokenService.SetRespUsageKey(r.tcpGram.userToken.ID)
				userTokenService.SetReqUsageKey(r.tcpGram.userToken.ID)
			}
		}
	}()

	buf := make([]byte, 1024*8)
	for {
		_ = r.ClientConn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(r.server.readTimeout)))
		if ln, err := r.ClientConn.Read(buf); err != nil {
			break
		} else {
			_ = r.RemoteConn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(r.server.writeTimeout)))
			if _, err := r.RemoteConn.Write(buf[:ln]); err != nil {
				break
			}
			userTokenService.IncrReqBytes(r.remoteAddr, ln)
			userTokenService.SetReqUsageKey(r.tcpGram.userToken.ID)
		}
	}
	_ = r.Close()
	<-done
}
