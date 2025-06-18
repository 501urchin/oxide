package ssh

import (
	"golang.org/x/crypto/ssh"
)

func ConnectSshViaPassword(user, host string, password string) (client *ssh.Client, err error) {
	conf := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", host, &conf)
}

func ConnectSshViakeys(user, host string, password string) (client *sshClient, err error) {
	conf := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", host, &conf)
	if err != nil {
		return nil, err
	}

	return &sshClient{
		user:   user,
		host:   host,
		client: conn,
	}, nil
}
