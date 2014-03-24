package main

import (
	auth "github.com/abbot/go-http-auth"
)

// SimpleSecurityProvider is a go-http-auth security provider which only
// allows a single set of simple credentials, for a single realm.
func SimpleSecurityProvider(username, password string) auth.SecretProvider {
	return func(user, realm string) string {
		if user != username {
			return ""
		}
		return password
	}
}
