package common

import (
	"fmt"
)

// ProxyType represents the type of proxy configuration
type ProxyType string

const (
	// DirectProxy represents a direct connection, no proxy (default on Windows)
	DirectProxy ProxyType = "DIRECT"
	// ManualProxy represents manual proxy settings (e.g., for httpProxy)
	ManualProxy ProxyType = "MANUAL"
	// PacProxy represents proxy autoconfiguration from URL
	PacProxy ProxyType = "PAC"
	// AutodetectProxy represents proxy autodetection (presumably with WPAD)
	AutodetectProxy ProxyType = "AUTODETECT"
	// SystemProxy represents system proxy settings (default on Linux)
	SystemProxy ProxyType = "SYSTEM"
	// UnspecifiedProxy represents not initialized proxy settings (for internal use)
	UnspecifiedProxy ProxyType = "UNSPECIFIED"
)

// Proxy contains information about proxy type and necessary proxy settings
type Proxy struct {
	proxyType          ProxyType
	autodetect         bool
	ftpProxy           string
	httpProxy          string
	noProxy            string
	proxyAutoconfigUrl string
	sslProxy           string
	socksProxy         string
	socksUsername      string
	socksPassword      string
	socksVersion       int
}

// NewProxy creates a new Proxy instance
func NewProxy() *Proxy {
	return &Proxy{
		proxyType:    UnspecifiedProxy,
		autodetect:   false,
		socksVersion: 0,
	}
}

// GetProxyType returns the proxy type
func (p *Proxy) GetProxyType() ProxyType {
	return p.proxyType
}

// SetProxyType sets the proxy type
func (p *Proxy) SetProxyType(proxyType ProxyType) error {
	p.proxyType = proxyType
	return nil
}

// GetAutodetect returns whether proxy autodetection is enabled
func (p *Proxy) GetAutodetect() bool {
	return p.autodetect
}

// SetAutodetect sets whether proxy autodetection is enabled
func (p *Proxy) SetAutodetect(autodetect bool) error {
	if err := p.verifyProxyTypeCompatibility(AutodetectProxy); err != nil {
		return err
	}
	p.proxyType = AutodetectProxy
	p.autodetect = autodetect
	return nil
}

// GetFtpProxy returns the FTP proxy
func (p *Proxy) GetFtpProxy() string {
	return p.ftpProxy
}

// SetFtpProxy sets the FTP proxy
func (p *Proxy) SetFtpProxy(ftpProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}
	p.proxyType = ManualProxy
	p.ftpProxy = ftpProxy
	return nil
}

// GetHttpProxy returns the HTTP proxy
func (p *Proxy) GetHttpProxy() string {
	return p.httpProxy
}

// SetHttpProxy sets the HTTP proxy
func (p *Proxy) SetHttpProxy(httpProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}
	p.proxyType = ManualProxy
	p.httpProxy = httpProxy
	return nil
}

// GetNoProxy returns the no proxy list
func (p *Proxy) GetNoProxy() string {
	return p.noProxy
}

// SetNoProxy sets the no proxy list
func (p *Proxy) SetNoProxy(noProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}
	p.proxyType = ManualProxy
	p.noProxy = noProxy
	return nil
}

// GetProxyAutoconfigUrl returns the proxy autoconfig URL
func (p *Proxy) GetProxyAutoconfigUrl() string {
	return p.proxyAutoconfigUrl
}

// SetProxyAutoconfigUrl sets the proxy autoconfig URL
func (p *Proxy) SetProxyAutoconfigUrl(proxyAutoconfigUrl string) error {
	if err := p.verifyProxyTypeCompatibility(PacProxy); err != nil {
		return err
	}
	p.proxyType = PacProxy
	p.proxyAutoconfigUrl = proxyAutoconfigUrl
	return nil
}

// GetSslProxy returns the SSL proxy
func (p *Proxy) GetSslProxy() string {
	return p.sslProxy
}

// SetSslProxy sets the SSL proxy
func (p *Proxy) SetSslProxy(sslProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}
	p.proxyType = ManualProxy
	p.sslProxy = sslProxy
	return nil
}

// GetSocksProxy returns the SOCKS proxy
func (p *Proxy) GetSocksProxy() string {
	return p.socksProxy
}

// SetSocksProxy sets the SOCKS proxy
func (p *Proxy) SetSocksProxy(socksProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}
	p.proxyType = ManualProxy
	p.socksProxy = socksProxy
	return nil
}

// GetSocksUsername returns the SOCKS username
func (p *Proxy) GetSocksUsername() string {
	return p.socksUsername
}

// SetSocksUsername sets the SOCKS username
func (p *Proxy) SetSocksUsername(socksUsername string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}
	p.proxyType = ManualProxy
	p.socksUsername = socksUsername
	return nil
}

// GetSocksPassword returns the SOCKS password
func (p *Proxy) GetSocksPassword() string {
	return p.socksPassword
}

// SetSocksPassword sets the SOCKS password
func (p *Proxy) SetSocksPassword(socksPassword string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}
	p.proxyType = ManualProxy
	p.socksPassword = socksPassword
	return nil
}

// GetSocksVersion returns the SOCKS version
func (p *Proxy) GetSocksVersion() int {
	return p.socksVersion
}

// SetSocksVersion sets the SOCKS version
func (p *Proxy) SetSocksVersion(socksVersion int) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}
	p.proxyType = ManualProxy
	p.socksVersion = socksVersion
	return nil
}

// verifyProxyTypeCompatibility verifies that the proxy type is compatible with the requested operation
func (p *Proxy) verifyProxyTypeCompatibility(requiredType ProxyType) error {
	if p.proxyType != UnspecifiedProxy && p.proxyType != requiredType {
		return fmt.Errorf("proxy type is %s, requested operation requires %s", p.proxyType, requiredType)
	}
	return nil
}

// ToCapabilities returns a map of proxy capabilities
func (p *Proxy) ToCapabilities() map[string]interface{} {
	caps := make(map[string]interface{})

	caps["proxyType"] = p.proxyType

	if p.autodetect {
		caps["autodetect"] = p.autodetect
	}
	if p.ftpProxy != "" {
		caps["ftpProxy"] = p.ftpProxy
	}
	if p.httpProxy != "" {
		caps["httpProxy"] = p.httpProxy
	}
	if p.noProxy != "" {
		caps["noProxy"] = p.noProxy
	}
	if p.proxyAutoconfigUrl != "" {
		caps["proxyAutoconfigUrl"] = p.proxyAutoconfigUrl
	}
	if p.sslProxy != "" {
		caps["sslProxy"] = p.sslProxy
	}
	if p.socksProxy != "" {
		caps["socksProxy"] = p.socksProxy
	}
	if p.socksUsername != "" {
		caps["socksUsername"] = p.socksUsername
	}
	if p.socksPassword != "" {
		caps["socksPassword"] = p.socksPassword
	}
	if p.socksVersion != 0 {
		caps["socksVersion"] = p.socksVersion
	}

	return caps
}
