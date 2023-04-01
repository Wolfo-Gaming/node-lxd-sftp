package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared"
	"golang.org/x/crypto/ssh"
)

func main() {
	args := os.Args
	//fmt.Println(args)

	randString := func(length int) string {
		var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
		randStr := make([]rune, length)
		for i := range randStr {
			randStr[i] = chars[rand.Intn(len(chars))]
		}

		return string(randStr)
	}

	con, lxdErr := lxd.ConnectLXD(args[1], &lxd.ConnectionArgs{
		TLSServerCert: args[2],
		TLSClientCert: args[3],
		TLSClientKey:  args[4],
	})
	if lxdErr != nil {
		fmt.Printf("{\"type\": \"error\",\"data\": %q}\n", lxdErr)
		return
	}
	config := &ssh.ServerConfig{}
	var authUser, authPass string
	authUser = randString(10)
	authPass = randString(20)
	instanceName := args[5]
	_, privKey, err := shared.GenerateMemCert(false, false)
	if err != nil {
		fmt.Printf("{\"type\": \"error\",\"data\": %q}\n", err)
		return
	}
	config.PasswordCallback = func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		if c.User() == authUser && string(pass) == authPass {
			return nil, nil
		}

		return nil, fmt.Errorf("Password rejected for %q", c.User())
	}
	private, err := ssh.ParsePrivateKey(privKey)
	if err != nil {
		fmt.Printf("{\"type\": \"error\",\"data\": %q}\n", err)
		return
	}
	config.AddHostKey(private)
	listenAddr := args[6]
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("{\"type\": \"error\",\"data\": %q}\n", err)
		return
	}
	//fmt.Printf("SSH SFTP listening on %v\n", listener.Addr())

	fmt.Printf("{\"type\": \"auth\",\"user\": %q, \"password\": %q, \"address\": %q}\n", authUser, authPass, listener.Addr())
	for {
		// Wait for new SSH connections.
		nConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("{\"type\": \"error\",\"data\": %q}\n", err)
			return
		}

		// Handle each SSH connection in its own go routine.
		go func() {
			fmt.Printf("{\"type\": \"connect\",\"data\": %q}\n", nConn.RemoteAddr().String())
			defer fmt.Printf("{\"type\": \"disconnect\",\"data\": %q}\n", nConn.RemoteAddr().String())
			defer func() { _ = nConn.Close() }()

			// Before use, a handshake must be performed on the incoming net.Conn.
			_, chans, reqs, err := ssh.NewServerConn(nConn, config)
			if err != nil {
				fmt.Printf("{\"type\": \"error-withclient\",\"data\": %q,\"remote\": %q}\n", err, nConn.RemoteAddr().String())
				return
			}

			// The incoming Request channel must be serviced.
			go ssh.DiscardRequests(reqs)

			// Service the incoming Channel requests.
			for newChannel := range chans {
				localChannel := newChannel

				// Channels have a type, depending on the application level protocol intended.
				// In the case of an SFTP session, this is "subsystem" with a payload string of
				// "<length=4>sftp"
				if localChannel.ChannelType() != "session" {
					_ = localChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
					fmt.Printf("{\"type\": \"error-withclient\",\"data\": \"unknown channel type\",\"remote\": %q}\n", nConn.RemoteAddr().String())
					continue
				}

				// Accept incoming channel request.
				channel, requests, err := localChannel.Accept()
				if err != nil {
					fmt.Printf("{\"type\": \"error-withclient\",\"data\": %q,\"remote\": %q}\n", err, nConn.RemoteAddr().String())
					return
				}

				// Sessions have out-of-band requests such as "shell", "pty-req" and "env".
				// Here we handle only the "subsystem" request.
				go func(in <-chan *ssh.Request) {
					for req := range in {
						ok := false
						switch req.Type {
						case "subsystem":
							if string(req.Payload[4:]) == "sftp" {
								ok = true
							}
						}

						_ = req.Reply(ok, nil)
					}
				}(requests)

				// Handle each channel in its own go routine.
				go func() {
					defer func() { _ = channel.Close() }()

					// Connect to the instance's SFTP server.
					sftpConn, err := con.GetInstanceFileSFTPConn(instanceName)
					if err != nil {
						fmt.Printf("{\"type\": \"error-withclient\",\"data\": %q,\"remote\": %q}\n", err, nConn.RemoteAddr().String())
						return
					}

					defer func() { _ = sftpConn.Close() }()

					// Copy SFTP data between client and remote instance.

					go func() {
						_, err := io.Copy(channel, sftpConn)
						if err != nil {
							fmt.Printf("{\"type\": \"error-withclient\",\"data\": %q,\"remote\": %q}\n", err, nConn.RemoteAddr().String())
						} else {
							fmt.Printf("{\"type\": \"info\",\"data\": \"disconnect\",\"remote\": %q}\n", nConn.RemoteAddr().String())
						}
						_ = channel.Close()
					}()

					_, err = io.Copy(sftpConn, channel)
					if err != nil {
						fmt.Printf("{\"type\": \"error-withclient\",\"data\": %q,\"remote\": %q}\n", err, nConn.RemoteAddr().String())
					}

					_ = sftpConn.Close()
				}()
			}
		}()
	}
}
