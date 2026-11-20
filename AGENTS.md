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
- Changes still go through a PR from a branch; merging to `main` is Daniel's approval (unless an auto-merge rule applies to this repo).
- Small, focused PRs with a clear description; include screenshots for UI changes.
- Handoff: keep your branch pushed often, and keep a section "## Estado / Próximos passos" in the PR description up to date (what is done, what is left, how to test). Another agent (Codex or Grok Bot) may continue the same branch mid-way.

## Publishing (https://speeddisc.danielpro.dev)
- This is a registered static site on the server (name `speeddisc`, published from `site/` only; the Go CLI ships via GitHub Releases).
- **Auto-publish:** every merge to `main` goes live within ~2 minutes via the server's `dots-static-sync` timer (a backup is taken first). No manual deploy needed.
- **Manual (Dots only):** publish now or a specific commit that is already on `main`: `ssh contabo-codex "sudo -n dots-deploy speeddisc [commit]"`; undo the last publish: `ssh contabo-codex "sudo -n dots-rollback speeddisc"`. Never copy files into `/var/www` yourself.
- `download/` on the server is preserved: deploys never delete it. Large binaries stay out of git (link GitHub Release assets instead).
- Never published: `*.md`, `.git*`, `*.php`, `.env`, `*secret*`, `uploads/`.
