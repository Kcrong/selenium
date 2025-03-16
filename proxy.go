package selenium

import (
	"errors"
	"fmt"
	"reflect"
)

// ProxyType represents the type of proxy configuration.
type ProxyType string

const (
	// DirectProxy represents a direct connection, no proxy (default on Windows).
	DirectProxy ProxyType = "DIRECT"
	// ManualProxy represents manual proxy settings (e.g., for httpProxy).
	ManualProxy ProxyType = "MANUAL"
	// PacProxy represents proxy autoconfiguration from URL.
	PacProxy ProxyType = "PAC"
	// AutodetectProxy represents proxy autodetection (presumably with WPAD).
	AutodetectProxy ProxyType = "AUTODETECT"
	// SystemProxy represents system proxy settings (default on Linux).
	SystemProxy ProxyType = "SYSTEM"
	// UnspecifiedProxy represents not initialized proxy settings (for internal use).
	UnspecifiedProxy ProxyType = "UNSPECIFIED"
)

type Proxy interface {
	GetProxyType() ProxyType
	SetProxyType(proxyType ProxyType) error
	GetAutodetect() bool
	SetAutodetect(d bool) error
	GetFTPProxy() string
	SetFTPProxy(p string) error
	GetHTTPProxy() string
	SetHTTPProxy(p string) error
	GetNoProxy() string
	SetNoProxy(p string) error
	GetProxyAutoConfigURL() string
	SetProxyAutoConfigURL(url string) error
	GetSSLProxy() string
	SetSSLProxy(ssl string) error
	GetSocksProxy() string
	SetSocksProxy(socks string) error
	GetSocksUsername() string
	SetSocksUsername(name string) error
	GetSocksPassword() string
	SetSocksPassword(pw string) error
	GetSocksVersion() int
	SetSocksVersion(version int) error
	verifyProxyTypeCompatibility(proxyType ProxyType) error
	ToCapabilities() map[string]interface{}
}

// proxy contains information about proxy type and necessary proxy settings.
type proxy struct {
	ProxyType          ProxyType `capabilities:"proxyType"`
	FTPProxy           string    `capabilities:"ftpProxy"`
	HTTPProxy          string    `capabilities:"httpProxy"`
	NoProxy            string    `capabilities:"noProxy"`
	ProxyAutoConfigURL string    `capabilities:"proxyAutoconfigUrl"`
	SSLProxy           string    `capabilities:"sslProxy"`
	SocksProxy         string    `capabilities:"socksProxy"`
	SocksUsername      string    `capabilities:"socksUsername"`
	SocksPassword      string    `capabilities:"socksPassword"`
	SocksVersion       int       `capabilities:"socksVersion"`
	Autodetect         bool      `capabilities:"autodetect"`
}

var _ Convertible = (*proxy)(nil)

// NewProxy creates a new Proxy instance.
func NewProxy() Proxy {
	return &proxy{
		ProxyType:          UnspecifiedProxy,
		FTPProxy:           "",
		HTTPProxy:          "",
		NoProxy:            "",
		ProxyAutoConfigURL: "",
		SSLProxy:           "",
		SocksProxy:         "",
		SocksUsername:      "",
		SocksPassword:      "",
		SocksVersion:       0,
		Autodetect:         false,
	}
}

// GetProxyType returns the proxy type.
func (p *proxy) GetProxyType() ProxyType {
	return p.ProxyType
}

// SetProxyType sets the proxy type.
func (p *proxy) SetProxyType(proxyType ProxyType) error {
	p.ProxyType = proxyType

	return nil
}

// GetAutodetect returns whether proxy autodetection is enabled.
func (p *proxy) GetAutodetect() bool {
	return p.Autodetect
}

// SetAutodetect sets whether proxy autodetection is enabled.
func (p *proxy) SetAutodetect(autodetect bool) error {
	if err := p.verifyProxyTypeCompatibility(AutodetectProxy); err != nil {
		return err
	}

	p.ProxyType = AutodetectProxy
	p.Autodetect = autodetect

	return nil
}

// GetFTPProxy returns the FTP proxy.
func (p *proxy) GetFTPProxy() string {
	return p.FTPProxy
}

// SetFTPProxy sets the FTP proxy.
func (p *proxy) SetFTPProxy(ftpProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}

	p.ProxyType = ManualProxy
	p.FTPProxy = ftpProxy

	return nil
}

// GetHTTPProxy returns the HTTP proxy.
func (p *proxy) GetHTTPProxy() string {
	return p.HTTPProxy
}

// SetHTTPProxy sets the HTTP proxy.
func (p *proxy) SetHTTPProxy(httpProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}

	p.ProxyType = ManualProxy
	p.HTTPProxy = httpProxy

	return nil
}

// GetNoProxy returns the no proxy list.
func (p *proxy) GetNoProxy() string {
	return p.NoProxy
}

// SetNoProxy sets the no proxy list.
func (p *proxy) SetNoProxy(noProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}

	p.ProxyType = ManualProxy
	p.NoProxy = noProxy

	return nil
}

// GetProxyAutoConfigURL returns the proxy autoconfig URL.
func (p *proxy) GetProxyAutoConfigURL() string {
	return p.ProxyAutoConfigURL
}

// SetProxyAutoConfigURL sets the proxy autoconfig URL.
func (p *proxy) SetProxyAutoConfigURL(url string) error {
	if err := p.verifyProxyTypeCompatibility(PacProxy); err != nil {
		return err
	}

	p.ProxyType = PacProxy
	p.ProxyAutoConfigURL = url

	return nil
}

// GetSSLProxy returns the SSL proxy.
func (p *proxy) GetSSLProxy() string {
	return p.SSLProxy
}

// SetSSLProxy sets the SSL proxy.
func (p *proxy) SetSSLProxy(sslProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}

	p.ProxyType = ManualProxy
	p.SSLProxy = sslProxy

	return nil
}

// GetSocksProxy returns the SOCKS proxy.
func (p *proxy) GetSocksProxy() string {
	return p.SocksProxy
}

// SetSocksProxy sets the SOCKS proxy.
func (p *proxy) SetSocksProxy(socksProxy string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}

	p.ProxyType = ManualProxy
	p.SocksProxy = socksProxy

	return nil
}

// GetSocksUsername returns the SOCKS username.
func (p *proxy) GetSocksUsername() string {
	return p.SocksUsername
}

// SetSocksUsername sets the SOCKS username.
func (p *proxy) SetSocksUsername(name string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}

	p.ProxyType = ManualProxy
	p.SocksUsername = name

	return nil
}

// GetSocksPassword returns the SOCKS password.
func (p *proxy) GetSocksPassword() string {
	return p.SocksPassword
}

// SetSocksPassword sets the SOCKS password.
func (p *proxy) SetSocksPassword(password string) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}

	p.ProxyType = ManualProxy
	p.SocksPassword = password

	return nil
}

// GetSocksVersion returns the SOCKS version.
func (p *proxy) GetSocksVersion() int {
	return p.SocksVersion
}

// SetSocksVersion sets the SOCKS version.
func (p *proxy) SetSocksVersion(socksVersion int) error {
	if err := p.verifyProxyTypeCompatibility(ManualProxy); err != nil {
		return err
	}

	p.ProxyType = ManualProxy
	p.SocksVersion = socksVersion

	return nil
}

var ErrIncompatibleProxyType = errors.New("proxy type is incompatible with the requested operation")

// verifyProxyTypeCompatibility verifies that the proxy type is compatible with the requested operation.
func (p *proxy) verifyProxyTypeCompatibility(requiredType ProxyType) error {
	if p.ProxyType != UnspecifiedProxy && p.ProxyType != requiredType {
		return fmt.Errorf("%w: proxy type %s is not supported", ErrIncompatibleProxyType, requiredType)
	}

	return nil
}

// ToCapabilities returns a map of proxy capabilities.
func (p *proxy) ToCapabilities() map[string]interface{} {
	const fieldTag = "capabilities"

	caps := make(map[string]interface{})
	val := reflect.ValueOf(p).Elem()
	typ := val.Type()

	for i := range val.NumField() {
		field := typ.Field(i)

		tagValue := field.Tag.Get(fieldTag)
		if tagValue == "-" {
			continue // Skip fields marked "-"
		}

		key := field.Name
		if tagValue != "" {
			key = tagValue
		}

		caps[key] = val.Field(i).Interface()
	}

	return caps
}
