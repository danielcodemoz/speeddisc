# AGENTS.md — speeddisc

<!-- codex-setup -->
## Codex environment
- **Stack:** Go Windows CLI/TUI + static site/
- **Setup:** `go mod download`
- **Check/test/build:** GOOS=windows go build ./... ; go test ./... (non-Windows tests)
- **Env vars:** none (dummy values only in Codex)

## Daniel's rules for Codex / AI agents
- Work on a branch and open a PR. Never push to `main`.
- Never commit `.env`, secrets, signing keys (`*.key`), or `uploads/` / `public/uploads`.
- Never copy GPL code (e.g. OmniGet) — ideas only.
- Keep UI text PT/EN as already done in the repo; keep the existing visual style.
- Every site footer references **danielpro.dev**.
- Support modals: Lemon Squeezy stays hidden until activated; M-Pesa only as a QR code.
- Codex does **not** deploy. Deploys are done separately after review.
- Small, focused PRs with a clear description; include screenshots for UI changes.
