package ssh

import (
	"fmt"
	"testing"

	taurinetesting "github.com/scott-mescudi/taurine/pkg/testing"
)

func TestMain(m *testing.M) {
	lis, err := taurinetesting.StartMockSSHServer("127.0.0.1:3098")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lis.Close()

	m.Run()
}

func TestPasswordAuth(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		user        string
		host        string
		expectError bool
	}{
		{name: "valid password login", password: taurinetesting.TestPassword, user: taurinetesting.TestUser, host: "127.0.0.1:3098", expectError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := ConnectSshViaPassword(taurinetesting.TestUser, "127.0.0.1:3098", taurinetesting.TestPassword)
			if tt.expectError && err == nil {
				client.Close()
				t.Fatal("failed to throw err")
			}

			if !tt.expectError && err != nil {
				client.Close()
				t.Fatal(err)
			}

			client.Close()
		})
	}

}
