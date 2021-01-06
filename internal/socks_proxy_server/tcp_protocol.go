package socks_proxy_server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"proxy-forward/internal/models"
	"proxy-forward/pkg/utils"
	"time"
)

type AuthHandle func(username, password string) (*models.UserToken, bool)

type TCPProtocol struct {
	cmd    byte
	atyp   byte
	ip     net.IP // when viaUDP is false, this is remote IP otherwise UPD client IP
	port   int    // when viaUDP is false, this is remote port otherwise UDP client port
	domain string
	viaUPD bool

	authHandle AuthHandle
	userToken  *models.UserToken
}

func (p *TCPProtocol) handClientshake(conn *net.TCPConn) error {
	// version
	data, err := p.checkVersion(conn)
	p.writeBuf(conn, data)
	if err != nil {
		return err
	}
	// auth
	data, err = p.checkAuth(conn)
	p.writeBuf(conn, data)
	if err != nil {
		return err
	}
	return nil
}
func (p *TCPProtocol) handRemoteshke(remoteConn, clientConn *net.TCPConn, username, password string) error {
	if err := p.validRemoteVersion(remoteConn, username, password); err != nil {
		return err
	}
	atyp, cmd, addrBytes, port, data, err := p.getAddr(remoteConn, clientConn)
	if err != nil {
		p.writeBuf(clientConn, data)
		return err
	} else {
		p.cmd = cmd
		p.atyp = atyp
	}
	switch p.atyp {
	case ATYPIPv4:
		p.ip = net.IPv4(addrBytes[0], addrBytes[1], addrBytes[2], addrBytes[3])
	case ATYPIPv6:
		p.ip = net.ParseIP(string(addrBytes))
	case ATYPDomain:
		p.domain = string(addrBytes)
		if addr, er := net.ResolveIPAddr("ip", p.domain); er != nil {
			p.writeBuf(clientConn, []byte{Version, 0x04, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			return er
		} else {
			p.ip = addr.IP
		}
	}
	p.port = port

	// check remote addr
	switch p.cmd {
	case CmdConnect:
		if !p.ip.IsGlobalUnicast() || p.port <= 0 {
			p.writeBuf(clientConn, []byte{Version, 0x02, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			return errors.New("remote address error")
		}
	case CmdUdpAssociate:
		p.viaUPD = true
	}
	return nil
}

func (p *TCPProtocol) checkVersion(conn *net.TCPConn) ([]byte, error) {
	var (
		version   byte
		methodLen int
	)
	if buf, err := p.readBuf(conn, 2); err != nil {
		return []byte{Version, MethodNone}, err
	} else {
		version = buf[0]
		methodLen = int(buf[1])
	}
	if version != Version || methodLen <= 0 {
		return []byte{Version, MethodNone}, errors.New("unsupported socks version")
	}

	if _, err := p.readBuf(conn, methodLen); err != nil {
		return []byte{Version, MethodNone}, err
	}

	if p.authHandle != nil {
		return []byte{Version, MethodAuth}, nil
	} else {
		return []byte{Version, MethodNoAuth}, nil
	}
}
func (p *TCPProtocol) validRemoteVersion(conn *net.TCPConn, username, password string) error { // 目标代理 sockts5 服务器 无需校验 {0x05, 0x00}
	var (
		version   byte
		methodLen int
	)
	// 需要鉴权的远端sock服务器
	if username != "" && password != "" {
		p.writeBuf(conn, []byte{Version, 0x02, 0x00, 0x02})
		if buf, err := p.readBuf(conn, 2); err != nil {
			return err
		} else {
			version = buf[0]
			methodLen = int(buf[1])
		}
		if version != Version || methodLen <= 0 {
			return errors.New("unsupported socks version")
		}
		valid := []byte{0x01}
		userLen := utils.B2mMap[len(username)]
		valid = append(valid, userLen)
		valid = append(valid, []byte(username)...)
		passLen := utils.B2mMap[len(password)]
		valid = append(valid, passLen)
		valid = append(valid, []byte(password)...)
		p.writeBuf(conn, valid)
		if buf, err := p.readBuf(conn, 2); err != nil {
			return err
		} else {
			version = buf[0]
			methodLen = int(buf[1])
		}
		if methodLen != 0 {
			return errors.New("unsupported socks version")
		}
		return nil
	} else {
		p.writeBuf(conn, []byte{Version, 0x01, 0x00})
		if buf, err := p.readBuf(conn, 2); err != nil {
			return err
		} else {
			version = buf[0]
			methodLen = int(buf[1])
		}
		if version != Version || methodLen != 0 {
			return errors.New("unsupported socks version")
		}
		return nil
	}
}
func (p *TCPProtocol) checkAuth(conn *net.TCPConn) ([]byte, error) {
	if p.authHandle == nil {
		return nil, nil
	}
	var (
		ver                byte
		userLen, passLen   int
		username, password string
	)
	if buf, err := p.readBuf(conn, 2); err != nil {
		return []byte{0x01, 0x01}, err
	} else {
		ver = buf[0]
		userLen = int(buf[1])
	}

	if ver != 0x01 || userLen <= 0 {
		return []byte{0x01, 0x01}, errors.New("unsupported auth version or username is empty")
	}
	if buf, err := p.readBuf(conn, userLen); err != nil {
		return []byte{0x01, 0x01}, err
	} else {
		username = string(buf)
	}

	if buf, err := p.readBuf(conn, 1); err != nil {
		return []byte{0x01, 0x01}, err
	} else {
		passLen = int(buf[0])
	}
	if passLen <= 0 {
		return []byte{0x01, 0x01}, errors.New("password is empty")
	}

	if buf, err := p.readBuf(conn, passLen); err != nil {
		return []byte{0x01, 0x01}, err
	} else {
		password = string(buf)
	}
	userToken, ok := p.authHandle(username, password)
	if !ok {
		return []byte{0x01, 0x01}, errors.New("username or password invalid")
	} else {
		p.userToken = userToken
		return []byte{0x01, 0x00}, nil
	}
}
func (p *TCPProtocol) getAddr(remoteConn, clientConn *net.TCPConn) (atyp, cmd byte, addrBytes []byte, port int, data []byte, err error) {
	var (
		ver  byte
		abuf []byte
	)
	if buf, er := p.readBuf(clientConn, 4); er != nil {
		err = er
		data = []byte{Version, 0x05, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		return
	} else {
		ver = buf[0]
		cmd = buf[1]
		atyp = buf[3]
		abuf = append(abuf, buf...)
	}
	if ver != Version {
		err = errors.New("unsupported socks version")
		data = []byte{Version, 0x05, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		return
	}
	if bytes.IndexByte([]byte{CmdConnect, CmdUdpAssociate}, cmd) == -1 {
		err = errors.New("unsupported CMD")
		data = []byte{Version, 0x07, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		return
	}
	switch atyp {
	case ATYPIPv4:
		addrBytes, err = p.readBuf(clientConn, 4)
		if err != nil {
			data = []byte{Version, 0x05, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			return
		} else {
			abuf = append(abuf, addrBytes...)
		}
	case ATYPDomain:
		var domainLen int
		if buf, er := p.readBuf(clientConn, 1); er != nil {
			err = er
			data = []byte{Version, 0x05, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			return
		} else {
			domainLen = int(buf[0])
			abuf = append(abuf, buf...)
		}
		if domainLen <= 0 {
			err = errors.New("length of domain is zero")
			data = []byte{Version, 0x05, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			return
		}
		addrBytes, err = p.readBuf(clientConn, domainLen)
		if err != nil {
			data = []byte{Version, 0x05, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			return
		} else {
			abuf = append(abuf, addrBytes...)
		}
	case ATYPIPv6:
		addrBytes, err = p.readBuf(clientConn, 16)
		if err != nil {
			data = []byte{Version, 0x05, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			return
		} else {
			abuf = append(abuf, addrBytes...)
		}
	default:
		err = errors.New("unsupported ATYP")
		data = []byte{Version, 0x08, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		return
	}

	// port
	if buf, er := p.readBuf(clientConn, 2); er != nil {
		err = er
		data = []byte{Version, 0x05, 0x00, ATYPIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		return
	} else {
		port = int(binary.BigEndian.Uint16(buf))
		abuf = append(abuf, buf...)
	}
	// remote handshake
	p.writeBuf(remoteConn, abuf)
	remoteBuf := make([]byte, 256)
	_ = remoteConn.SetReadDeadline(time.Now().Add(time.Second * 10))
	remoteConn.Read(remoteBuf)
	if remoteBuf[0] != Version || remoteBuf[1] != 0x00 {
		data = remoteBuf
		err = errors.New("remote handshake fail.")
		return
	}
	return
}

func (p *TCPProtocol) networkString() (string, string, string, error) {
	return LoadRemoteAddr(p.userToken)
}

func (p *TCPProtocol) readBuf(conn *net.TCPConn, ln int) ([]byte, error) {
	buf := make([]byte, ln)
	curReadLen := 0
	for curReadLen < ln {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		l, err := conn.Read(buf[curReadLen:])
		if err != nil {
			return nil, err
		}
		curReadLen += l
	}
	return buf, nil
}

func (p *TCPProtocol) writeBuf(conn *net.TCPConn, data []byte) {
	if data != nil && len(data) > 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
		_, _ = conn.Write(data)
	}
}
