package main

import (
	"crypto/sha1"
	"encoding/base64"
	"log"

	auth "github.com/abbot/go-http-auth"
)

// SimpleSecurityProvider is a go-http-auth security provider which only
// allows a single set of simple credentials, for a single realm.
func SimpleSecurityProvider(username, password string) auth.SecretProvider {
	return func(user, realm string) string {
		log.Print("SimpleSecurityProvider called for user " + user + " realm " + realm)
		if user != username {
			return ""
		}
		d := sha1.New()
		d.Write([]byte(password))
		pw := "{SHA}" + base64.StdEncoding.EncodeToString(d.Sum(nil))
		log.Print("For user " + username + " created password " + pw)
		return pw
	}
}
