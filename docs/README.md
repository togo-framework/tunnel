# tunnel — docs

**Tunnel.** Expose a local togo app to the public internet through a pluggable tunnel driver.

## Install

```bash
togo install togo-framework/tunnel
```

Select a driver via **tunnel.provider in togo.yaml (or TUNNEL_DRIVER)**. Drive it from the CLI with **`togo tunnel`**.

## Interface

`Tunnel` — `Start(ctx, addr) -> publicURL`, `Stop`, `Status`.

## Configuration

| Env var | Description |
|---|---|
| `TUNNEL_DRIVER` | Selects the tunnel driver (alternative to togo.yaml `tunnel.provider`). |

## Usage & notes

Powers `togo tunnel`. Drivers self-register; the CLI resolves with `tunnel.Build(name, k)`. Configure in `togo.yaml`:
```yaml
tunnel:
  provider: cloudflare   # cloudflare|ngrok|tailscale|frp
  addr: localhost:8080
```

## Example

```bash
togo tunnel:start --provider ngrok
```

## Links

- [Driver plugins](https://to-go.dev/marketplace)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/tunnel)
