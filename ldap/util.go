// Package ldap contains patched items from github.com/go-ldap/ldap/v3.

// The MIT License (MIT)

// Copyright (c) 2011-2015 Michael Mitton (mmitton@gmail.com)
// Portions copyright (c) 2015-2016 go-ldap Authors

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

//lint:file-ignore ST1005

package ldap

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/go-ldap/ldap/v3"
)

// DialURL connects to the given ldap URL vie TCP using tls.Dial or net.Dial if ldaps://
// or ldap:// specified as protocol. On success a new Conn for the connection
// is returned.
// This patched version supports upgrading unencrypted LDAP connections (with `ldap://` URIs)
// wiht StartTLS if specified. `enableStartTLS` is ignored when not using unencrypted
// LDAP connection.
func DialURL(addr string, tlsConf *tls.Config, enableStartTLS bool) (*ldap.Conn, error) {
	lurl, err := url.Parse(addr)
	if err != nil {
		return nil, NewError(ldap.ErrorNetwork, err)
	}

	host, port, err := net.SplitHostPort(lurl.Host)
	if err != nil {
		// we asume that error is due to missing port
		host = lurl.Host
		port = ""
	}

	switch lurl.Scheme {
	case "ldapi":
		if lurl.Path == "" || lurl.Path == "/" {
			lurl.Path = "/var/run/slapd/ldapi"
		}
		log.Debug("Dialing over Unix socket to ", lurl.Path)
		return ldap.Dial("unix", lurl.Path)
	case "ldap":
		if port == "" {
			port = ldap.DefaultLdapPort
		}
		log.Debug("Dialing unencrypted to ", net.JoinHostPort(host, port))
		l, err := ldap.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			return nil, err
		}
		if enableStartTLS {
			log.Debug("Beginning StartTLS")
			err = l.StartTLS(tlsConf)
		}
		return l, err
	case "ldaps":
		if port == "" {
			port = ldap.DefaultLdapsPort
		}
		if tlsConf == nil {
			tlsConf = &tls.Config{
				ServerName: host,
			}
		}
		log.Debug("Dialing over TLS to ", net.JoinHostPort(host, port))
		return ldap.DialTLS("tcp", net.JoinHostPort(host, port), tlsConf)
	}

	return nil, NewError(ldap.ErrorNetwork, fmt.Errorf("Unknown scheme '%s'", lurl.Scheme))
}

// NewError creates an LDAP error with the given code and underlying error
func NewError(resultCode uint16, err error) error {
	return &Error{ResultCode: resultCode, Err: err}
}

// Error holds LDAP error information
type Error struct {
	// Err is the underlying error
	Err error
	// ResultCode is the LDAP error code
	ResultCode uint16
	// MatchedDN is the matchedDN returned if any
	MatchedDN string
}

func (e *Error) Error() string {
	return fmt.Sprintf("LDAP Result Code %d %q: %s", e.ResultCode, ldap.LDAPResultCodeMap[e.ResultCode], e.Err.Error())
}
