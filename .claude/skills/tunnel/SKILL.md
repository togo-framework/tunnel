---
name: tunnel
description: Expose a local togo app to the public internet — pick a tunnel driver (Cloudflare/ngrok/Tailscale/frp) and call tunnel.Start for dev sharing & webhook testing
---

# togo tunnel

The `tunnel` plugin publishes a local togo app to a public URL for dev sharing,
webhook testing and demos, over a single `Tunnel` contract. The safe default
`log` driver records intent without opening anything; real drivers ship as
plugins selected by `TUNNEL_DRIVER`.

## Setup

```bash
togo install togo-framework/tunnel
# then ONE driver:
togo install togo-framework/tunnel-cloudflare   # *.trycloudflare.com (no account)
togo install togo-framework/tunnel-ngrok        # ngrok
togo install togo-framework/tunnel-tailscale    # Tailscale Funnel
togo install togo-framework/tunnel-frp          # self-hosted frp
```

`.env`:

```bash
TUNNEL_DRIVER=cloudflare   # or ngrok | tailscale | frp | log
```

## Use

```go
import (
	_ "github.com/togo-framework/tunnel"
	_ "github.com/togo-framework/tunnel-cloudflare"
	"github.com/togo-framework/tunnel"
)

if tn, ok := tunnel.FromKernel(k); ok {
	url, err := tn.Start(ctx, "8080")   // public URL → local :8080
	defer tn.Stop(ctx)
	st, _ := tn.Status(ctx)             // {Running, URL, Driver}
}
```

- `Start(ctx, addr)` accepts a bare port (`"8080"`), `":8080"`, or `host:port`.
- The default `log` driver is safe for tests — it opens nothing.
- Pairs with the `deploy` + `dns` subsystems: deploy ships the app, dns/tunnel
  make it reachable.

## Rules
- One driver is active per process (`TUNNEL_DRIVER`).
- Tunnels expose your local app publicly — only run them intentionally, never by
  default in production.
- Always `tn.Stop(ctx)` when done; most drivers tear down the tunnel on stop.
