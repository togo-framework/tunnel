// Package tunnel is togo's public-tunnel subsystem: a single Tunnel contract
// over wildly different exposers — Cloudflare Tunnel (cloudflared), ngrok,
// Tailscale Funnel and frp — that publish a local togo app to the public
// internet for dev sharing, webhook testing and demos. The safe dev default
// "log" driver records intent without opening anything. Real drivers ship as
// plugins that call tunnel.RegisterDriver and depend on this package. Select
// one with TUNNEL_DRIVER.
//
// Install: `togo install togo-framework/tunnel` (blank-import registers it),
// then a driver, e.g. `togo install togo-framework/tunnel-cloudflare`.
//
// It pairs with the togo `deploy` and `dns` subsystems: deploy ships the app,
// dns/tunnel make it reachable.
package tunnel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/togo-framework/togo"
)

// ErrUnsupported is returned by a driver for an operation it does not implement.
var ErrUnsupported = errors.New("tunnel: operation not supported by this driver")

// Status describes the current state of a tunnel.
type Status struct {
	Running bool   // a tunnel is active
	URL     string // the public URL, when running
	Driver  string // the active driver name
}

// Tunnel is implemented by driver plugins. Start opens a tunnel from the public
// internet to a local address (host:port, or a bare ":port"/"port") and returns
// the public URL. Stop tears it down. Status reports the current state.
type Tunnel interface {
	Start(ctx context.Context, addr string) (publicURL string, err error)
	Stop(ctx context.Context) error
	Status(ctx context.Context) (Status, error)
}

// DriverFactory builds a Tunnel from the kernel (env-configured).
type DriverFactory func(k *togo.Kernel) (Tunnel, error)

var (
	regMu   sync.RWMutex
	drivers = map[string]DriverFactory{}
)

// RegisterDriver registers a tunnel driver by name (call from a plugin's init()).
func RegisterDriver(name string, f DriverFactory) {
	regMu.Lock()
	drivers[name] = f
	regMu.Unlock()
}

// Drivers lists the registered driver names (unordered).
func Drivers() []string {
	regMu.RLock()
	defer regMu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	return out
}

// Build constructs a Service for the named driver without booting the kernel.
// A togo CLI `tunnel` runner can use it to resolve a provider standalone; the
// stock drivers are env-configured, so a nil kernel is acceptable.
func Build(name string, k *togo.Kernel) (*Service, error) {
	if name == "" {
		name = "log"
	}
	regMu.RLock()
	f, ok := drivers[name]
	regMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("tunnel: unknown driver %q (install its plugin, e.g. togo install togo-framework/tunnel-%s)", name, name)
	}
	t, err := f(k)
	if err != nil {
		return nil, err
	}
	return &Service{tunnel: t, driver: name}, nil
}

func init() {
	RegisterDriver("log", func(k *togo.Kernel) (Tunnel, error) { return &logTunnel{k: k}, nil })

	togo.RegisterProviderFunc("tunnel", togo.PriorityService, func(k *togo.Kernel) error {
		name := os.Getenv("TUNNEL_DRIVER")
		if name == "" {
			name = "log" // safe dev default: record intent, open nothing
		}
		svc, err := Build(name, k)
		if err != nil {
			return err
		}
		k.Set("tunnel", svc)
		return nil
	})
}

// Service is the tunnel runtime stored on the kernel (k.Get("tunnel")).
type Service struct {
	tunnel Tunnel
	driver string
}

// Tunnel returns the active driver implementation.
func (s *Service) Tunnel() Tunnel { return s.tunnel }

// Driver returns the active driver name.
func (s *Service) Driver() string { return s.driver }

// Start opens the tunnel to a local address and returns the public URL.
func (s *Service) Start(ctx context.Context, addr string) (string, error) {
	return s.tunnel.Start(ctx, addr)
}

// Stop tears the tunnel down.
func (s *Service) Stop(ctx context.Context) error { return s.tunnel.Stop(ctx) }

// Status reports the current state (driver name filled in).
func (s *Service) Status(ctx context.Context) (Status, error) {
	st, err := s.tunnel.Status(ctx)
	st.Driver = s.driver
	return st, err
}

// FromKernel fetches the tunnel service from the kernel container.
func FromKernel(k *togo.Kernel) (*Service, bool) {
	v, ok := k.Get("tunnel")
	if !ok {
		return nil, false
	}
	s, ok := v.(*Service)
	return s, ok
}

// logTunnel records intent via the kernel logger — the safe default.
type logTunnel struct {
	k       *togo.Kernel
	addr    string
	running bool
}

func (l *logTunnel) log(op string, kv ...any) {
	if l.k != nil && l.k.Log != nil {
		l.k.Log.Info("tunnel (log driver) "+op, kv...)
	}
}

func (l *logTunnel) Start(_ context.Context, addr string) (string, error) {
	l.addr = addr
	l.running = true
	url := "https://log.tunnel.invalid"
	l.log("start", "addr", addr, "url", url)
	return url, nil
}

func (l *logTunnel) Stop(context.Context) error {
	l.running = false
	l.log("stop", "addr", l.addr)
	return nil
}

func (l *logTunnel) Status(context.Context) (Status, error) {
	url := ""
	if l.running {
		url = "https://log.tunnel.invalid"
	}
	return Status{Running: l.running, URL: url}, nil
}

// NormalizeAddr turns a bare port ("8080"), ":port" (":8080") or host:port into
// a host:port string with a default host of 127.0.0.1 — a small helper drivers
// share so `tunnel.Start(ctx, "8080")` works everywhere.
func NormalizeAddr(addr string) string {
	if addr == "" {
		return "127.0.0.1:80"
	}
	// bare number → :number
	allDigits := true
	for _, r := range addr {
		if r < '0' || r > '9' {
			allDigits = false
			break
		}
	}
	if allDigits {
		return "127.0.0.1:" + addr
	}
	if addr[0] == ':' {
		return "127.0.0.1" + addr
	}
	return addr
}

// PortOf returns just the port from an address accepted by NormalizeAddr.
func PortOf(addr string) string {
	n := NormalizeAddr(addr)
	for i := len(n) - 1; i >= 0; i-- {
		if n[i] == ':' {
			return n[i+1:]
		}
	}
	return ""
}
