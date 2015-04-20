package main

import (
	"code.google.com/p/gopass"
	"errors"
	"fmt"
	"github.com/tumdum/gosmbclient"
	"os"
	"strings"
)

var serverAuth = make(map[string]gosmbclient.Auth)

func SimpleAuth(server, share, workgroup, username string) (*gosmbclient.Auth, error) {
	if auth, ok := serverAuth[server]; ok {
		return &auth, nil
	}
	wg := getCredential("workgroup", server)
	if len(wg) == 0 {
		return nil, errors.New("No workgroup credentials")
	}
	un := getCredential("username", server)
	if len(un) == 0 {
		return nil, errors.New("No username credentials")
	}
	ps := getCredential("password", server)
	if len(ps) == 0 {
		return nil, errors.New("No password credentials")
	}
	auth := gosmbclient.Auth{wg, un, ps}
	serverAuth[server] = auth
	return &auth, nil
}

func getCredential(name, server string) string {
	envcred := os.Getenv("SMB_" + strings.ToUpper(name))
	if len(envcred) > 0 {
		return envcred
	}
	cred, err := gopass.GetPass(fmt.Sprintf("%s for '%s': ", strings.ToTitle(name), server))
	if err != nil {
		return ""
	}
	return cred
}
