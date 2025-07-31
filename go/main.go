package main

import (
	"bytes"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

const (
	TestUser      = "testuser"
	TestPassword  = "password123"
	AllowedPubKey = `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBWiHlrQ6HS7vytwfKb32R70waRKqJ9cZOWx8RDfm4HX jayac@diddy.local`
	PrivateKeyPEM = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAVoh5a0Oh0u78rcHym99ke9MGkSqifXGTlsfEQ35uB1wAAAJgWiaMKFomj
CgAAAAtzc2gtZWQyNTUxOQAAACAVoh5a0Oh0u78rcHym99ke9MGkSqifXGTlsfEQ35uB1w
AAAEBcZQq0UwEx650Q4FFDhsnzC05hFjcOLvSCjL8t6/MKBRWiHlrQ6HS7vytwfKb32R70
waRKqJ9cZOWx8RDfm4HXAAAAEWpheWFjQGRpZGR5LmxvY2FsAQIDBA==
-----END OPENSSH PRIVATE KEY-----
`
)

func passwordHandler(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	if c.User() == TestUser && string(pass) == TestPassword {
		return nil, nil
	}
	return nil, fmt.Errorf("failed to authenticate password")
}

func publicKeyHandler(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	authorizedKey, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(AllowedPubKey))
	if bytes.Equal(key.Marshal(), authorizedKey.Marshal()) {
		return nil, nil
	}
	return nil, fmt.Errorf("failed to authenticate public key")
}

func StartMockSSHServer(addr string) error {
	serverConfig := &ssh.ServerConfig{
		PasswordCallback:  passwordHandler,
		PublicKeyCallback: publicKeyHandler,
	}

	signer, err := ssh.ParsePrivateKey([]byte(PrivateKeyPEM))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	serverConfig.AddHostKey(signer)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			sshConn, chans, reqs, err := ssh.NewServerConn(conn, serverConfig)
			if err != nil {
				return
			}

			defer sshConn.Close()
			go ssh.DiscardRequests(reqs)
			for newChannel := range chans {
				newChannel.Reject(ssh.UnknownChannelType, "no channels supported")
			}
		}()
	}

}

func main() {
	fmt.Println(StartMockSSHServer("127.0.0.1:22"))
}
