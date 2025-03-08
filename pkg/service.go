package pkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ServiceOption configures a Service instance.
type ServiceOption func(*Service) error

// isDisplay validates that the given disp is in the format "x" or "x.y", where
// x and y are both integers.
func isDisplay(disp string) bool {
	ds := strings.Split(disp, ".")
	if len(ds) > 2 {
		return false
	}

	for _, d := range ds {
		if _, err := strconv.Atoi(d); err != nil {
			return false
		}
	}
	return true
}

// FrameBufferOptions describes the options that can be used to create a frame buffer.
type FrameBufferOptions struct {
	// ScreenSize is the option for the frame buffer screen size.
	// This is of the form "{width}x{height}[x{depth}]".  For example: "1024x768x24"
	ScreenSize string
}

// Service controls a locally-running Selenium subprocess.
type Service struct {
	port            int
	addr            string
	cmd             *exec.Cmd
	shutdownURLPath string

	display, xauthPath string

	geckoDriverPath, javaPath string
	chromeDriverPath          string
	htmlUnitPath              string

	output io.Writer
}

// NewSeleniumService starts a Selenium instance in the background.
func NewSeleniumService(jarPath string, port int, opts ...ServiceOption) (*Service, error) {
	s, err := newService(exec.Command("java"), port, opts...)
	if err != nil {
		return nil, err
	}
	if s.javaPath != "" {
		s.cmd.Path = s.javaPath
	}
	if s.geckoDriverPath != "" {
		s.cmd.Args = append([]string{"java", "-Dwebdriver.gecko.driver=" + s.geckoDriverPath}, s.cmd.Args[1:]...)
	}
	if s.chromeDriverPath != "" {
		s.cmd.Args = append([]string{"java", "-Dwebdriver.chrome.driver=" + s.chromeDriverPath}, s.cmd.Args[1:]...)
	}

	var classpath []string
	if s.htmlUnitPath != "" {
		classpath = append(classpath, s.htmlUnitPath)
	}
	classpath = append(classpath, jarPath)
	s.cmd.Args = append(s.cmd.Args, "-cp", strings.Join(classpath, ":"))
	s.cmd.Args = append(s.cmd.Args, "org.openqa.grid.selenium.GridLauncherV3", "-port", strconv.Itoa(port), "-debug")

	if err := s.start(port); err != nil {
		return nil, err
	}
	return s, nil
}

// NewChromeDriverService starts a ChromeDriver instance in the background.
func NewChromeDriverService(path string, port int, opts ...ServiceOption) (*Service, error) {
	cmd := exec.Command(path, "--port="+strconv.Itoa(port), "--verbose")
	s, err := newService(cmd, port, opts...)
	if err != nil {
		return nil, err
	}
	s.shutdownURLPath = "/shutdown"
	if err := s.start(port); err != nil {
		return nil, err
	}
	return s, nil
}

// NewGeckoDriverService starts a GeckoDriver instance in the background.
func NewGeckoDriverService(path string, port int, opts ...ServiceOption) (*Service, error) {
	cmd := exec.Command(path, "--port", strconv.Itoa(port))
	s, err := newService(cmd, port, opts...)
	if err != nil {
		return nil, err
	}
	if err := s.start(port); err != nil {
		return nil, err
	}
	return s, nil
}

func newService(cmd *exec.Cmd, port int, opts ...ServiceOption) (*Service, error) {
	s := &Service{
		port: port,
		addr: fmt.Sprintf("http://localhost:%d", port),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	cmd.Stderr = s.output
	cmd.Stdout = s.output
	cmd.Env = os.Environ()
	// TODO(minusnine): Pdeathsig is only supported on Linux. Somehow, make sure
	// process cleanup happens as gracefully as possible.
	if s.display != "" {
		cmd.Env = append(cmd.Env, "DISPLAY=:"+s.display)
	}
	if s.xauthPath != "" {
		cmd.Env = append(cmd.Env, "XAUTHORITY="+s.xauthPath)
	}
	s.cmd = cmd
	return s, nil
}

func (s *Service) start(port int) error {
	if err := s.cmd.Start(); err != nil {
		return err
	}

	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		resp, err := http.Get(s.addr + "/status")
		if err == nil {
			resp.Body.Close()
			switch resp.StatusCode {
			// Selenium <3 returned Forbidden and BadRequest. ChromeDriver and
			// Selenium 3 return OK.
			case http.StatusForbidden, http.StatusBadRequest, http.StatusOK:
				return nil
			}
		}
	}
	return fmt.Errorf("server did not respond on port %d", port)
}

// Stop shuts down the WebDriver service, and the X virtual frame buffer
// if one was started.
func (s *Service) Stop() error {
	// Selenium 3 stopped supporting the shutdown URL by default.
	// https://github.com/SeleniumHQ/selenium/issues/2852
	if s.shutdownURLPath == "" {
		if err := s.cmd.Process.Kill(); err != nil {
			return err
		}
	} else {
		resp, err := http.Get(s.addr + s.shutdownURLPath)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
	}
	if err := s.cmd.Wait(); err != nil && err.Error() != "signal: killed" {
		return err
	}
	return nil
}
