package ftp

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Host struct {
	Host, User, Password string
	Port                 int
}

type Session struct {
	Host     *Host
	FromFile string
	FromPath string
	ToFile   string
	ToPath   string
	sc       *sftp.Client
}

//NewConn Create a new SFTP connection by given parameters
func NewConn(new *Session) (session *Session, err error) {
	switch {
	case `` == strings.TrimSpace(new.Host.Host),
		`` == strings.TrimSpace(new.Host.User),
		`` == strings.TrimSpace(new.Host.Password),
		0 >= new.Host.Port || new.Host.Port > 65535:
		return nil, errors.New("invalid parameters")
	}
	if err = new.connect(); nil != err {
		return nil, err
	}
	return new, nil
}

//Upload file to sftp server
func (inst *Session) Upload() (err error) {
	localFile := fmt.Sprintf("%s/%s", inst.FromPath, inst.FromFile)
	remoteFile := fmt.Sprintf("%s/%s", inst.ToPath, inst.ToFile)
	log.Infoln("try and upload file:", localFile, " to:", remoteFile)
	srcFile, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		_ = inst.sc.Mkdir(path)
	}

	dstFile, err := inst.sc.Create(remoteFile)
	if err != nil {
		log.Errorln("upload file error:", err)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return
}

func (inst *Session) connect() (err error) {

	config := &ssh.ClientConfig{
		User:            inst.Host.User,
		Auth:            []ssh.AuthMethod{ssh.Password(inst.Host.Password)},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connect to ssh
	addr := fmt.Sprintf("%s:%d", inst.Host.Host, inst.Host.Port)
	log.Infoln("try connect to:", addr)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}

	// create sftp client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Errorln("connect error:", err)
		return err
	}
	inst.sc = client
	return nil
}
