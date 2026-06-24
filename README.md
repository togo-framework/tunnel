<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/tunnel</h1>
  <p>Expose a local togo app to the public internet — one contract over Cloudflare Tunnel, ngrok, Tailscale Funnel and frp.</p>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/tunnel"><img src="https://pkg.go.dev/badge/github.com/togo-framework/tunnel.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/tunnel
```
<!-- /togo-header -->

togo's **public-tunnel subsystem**: a single `Tunnel` contract over very different
exposers — **Cloudflare Tunnel** (`cloudflared`), **ngrok**, **Tailscale Funnel** and
**frp** — that publish a local togo app to the public internet for dev sharing,
webhook testing and demos. The safe dev default `log` driver records intent without
opening anything. Real drivers ship as plugins that call `tunnel.RegisterDriver`; pick
one with `TUNNEL_DRIVER`.

It pairs with the togo [`deploy`](https://github.com/togo-framework/deploy) and
[`dns`](https://github.com/togo-framework/dns) subsystems: `deploy` ships the app,
`dns`/`tunnel` make it reachable.

```bash
togo install togo-framework/tunnel             # the base
togo install togo-framework/tunnel-cloudflare  # a driver
```

Drivers: `tunnel-cloudflare`, `tunnel-ngrok`, `tunnel-tailscale`, `tunnel-frp`.

## Configure

```env
TUNNEL_DRIVER=cloudflare   # or ngrok | tailscale | frp | log (default)
# + the selected driver's env (e.g. NGROK_AUTHTOKEN)
```

## Use

```go
svc, _ := tunnel.FromKernel(k)
url, err := svc.Start(ctx, "8080")   // → https://something.trycloudflare.com
defer svc.Stop(ctx)
```

`Start` accepts a bare port (`"8080"`), `":8080"`, or `host:port`. A driver that
doesn't support an operation returns `tunnel.ErrUnsupported`.

## Write a driver

```go
func init() {
    tunnel.RegisterDriver("mytunnel", func(k *togo.Kernel) (tunnel.Tunnel, error) {
        return &myTunnel{ /* env */ }, nil
    })
}
```

Implement `Start`/`Stop`/`Status`. See `tunnel-ngrok` for a pure-Go example and
`tunnel-cloudflare` for a binary-wrapping one.

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
