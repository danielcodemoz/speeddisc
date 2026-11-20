# SpeedDisc Linux · pré-visualização / preview

Edição portátil, sem instalador. O executável abre uma interface no navegador
local. Não precisa de Internet e não escuta na rede: liga-se apenas a
`127.0.0.1`, numa sessão temporária. Execute com a sua conta normal, **sem sudo**.
Descarregue o arquivo para a arquitectura certa, extraia-o e abra o executável.
Se o gestor de ficheiros não abrir o programa, execute-o no terminal, a partir
da pasta extraída: `./SpeedDisc-Linux`.

The portable executable opens a local browser workbench. It needs no installer
or Internet connection. Run as your normal user, not as root. Use
`./SpeedDisc-Linux --no-browser` to print the workbench URL instead. Do not share
that URL: it is specific to the running session. Close the session in the
workbench or by pressing Ctrl+C in the terminal. Idle sessions stop
automatically.

## O que esta edição inclui / Scope

- Unprivileged, read-only Linux diagnostics: disk space, memory, processes,
  local interfaces, large files and an estimate for temporary directories.
- The scans are bounded. Their results are not a complete disk inventory.
  No caches, processes, services or other resources are deleted or changed.
- Linux command examples are reference only, not execution buttons. TXT and
  branded PDF reports can contain local paths and names; review before sharing.
- PT/EN and light/dark. Use `--demo --no-browser` for fictional data only.

## Build

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o SpeedDisc-Linux ./cmd/speeddisc-linux
go test ./...
```

Use `GOARCH=arm64` for 64-bit ARM. A desktop web browser is required to open the
workbench. No WebView is embedded and no Windows maintenance command is used.
