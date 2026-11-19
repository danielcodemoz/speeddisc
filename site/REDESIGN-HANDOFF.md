# SpeedDisc website redesign 1.0

Local implementation only. No remote writes, PR, deployment or Windows execution performed by this worker.

## Baseline and compatibility

Source repository: `danielcodemoz/speeddisc`, main commit `1214465e0061d18b05a5c8062729f396bf8830c5`.

The complete original site is present for a standalone preview. `comandos.html`, `comandos.js`, `app.js`, `style.css`, fonts, favicon, console screenshot and supporting scripts remain byte-identical. They are **not changes to include in a patch**.

The redesigned homepage does not load the old app.js; this keeps obsolete homepage copy out of the new experience while retaining the existing commands page behaviour.

## New experience

- Problem-first homepage with accent-insensitive search, category filters and useful empty/reset states.
- Five substantive manual guides, each a separate static URL in PT and EN. Content is visible without JavaScript.
- Language links point to the exact counterpart page, not a generic language homepage.
- Shared light/dark state in `sd-theme`, matching legacy storage key; semantic keyboard-operable buttons and focus states.
- Responsive layouts, reduced-motion preference and print stylesheet.
- Existing technician v0.4.0 download and four kit links are version-pinned. All six published release assets (including the unlinked PowerShell launcher) were downloaded as bytes and their SHA-256 compared with GitHub digests, all matched. No downloaded code ran. Record: `release-verification.json`.
- Every new page links to the existing support page at `https://danielpro.dev/apoiar/`; the existing commands-page support modal is unchanged.
- No auto-download, repair, telemetry or fetch on the new pages. No dummy diagnostic download. No instructions to bypass Windows protection.
- Old homepage anchors `#main`, `#fazer`, `#pacotes`, `#correr`, `#smartscreen`, `#limites`, `#obter`, `#kits` remain valid. The package alias points to existing kit explanations; the run alias points to the technician area; SmartScreen points to stop/verify guidance.
- Existing brand assets/fonts retained; new original inline line-art icons avoid changes to the legacy limited icon bundle.

## Build and checks

From the site directory:

```sh
python _source/build.py
node --check redesign.v1.js
node _source/test-ui.cjs
python _source/test-static.py
python -m http.server 8765
```

Static tests compare compatibility assets against the included `_source/legacy-sha256.json`, so no external fixture or checkout is required. UI tests use included guides.json and mocked DOM elements; they cover logic, not actual browser layout.

Passed: twelve generated bilingual HTML pages; one H1 each; no duplicate IDs; local resource, internal-link and guide-fragment checks; legacy aliases; copied asset compatibility; published download digest matches; accent-insensitive search, category/empty/reset flows, theme state/storage, submit and print event logic.

**Not yet verified in a real browser:** screenshots, mobile rendering/overflow, keyboard traversal, browser history and real print output. Installed Chromium was present, but process startup failed with `socket() failed: Operation not permitted`, including the reviewed escalated attempt. The separate cloud browser could not reach the local server (`ERR_CONNECTION_REFUSED`). No visual QA pass or screenshot is claimed. Do not publish as fully QA-approved until the following checks are run in an available cloud browser.

## Browser acceptance checklist

1. Serve this directory as the web root (absolute asset paths retained for compatibility). Check `/` and `/en/` at 1440px, 768px, 390px and 320px; no horizontal overflow.
2. Search `espaco` and `espaço`: disk guide. Search `Wi-Fi`: network guide. Search nonsense: visible empty state; clear restores five guides. Combine a search with a category, then reset.
3. Tab to theme and press Enter/Space; theme changes, focus remains visible, `aria-pressed` changes and reload retains state. Repeat in both languages.
4. Open all guide links directly; switch language and confirm the matching guide. Follow table-of-contents fragments and browser Back/Forward. Check old homepage anchor links.
5. Confirm no automatic download or outgoing request on new pages except local CSS/font/image/script resources. Inspect the technician screenshot and all SVG icons.
6. Open `/comandos.html` and check the existing search/filter/language/copy functions. It is preserved, not redesigned.
7. Print a guide to preview; inspect headings, steps and page breaks. Screenshots should include homepage and a guide in desktop/mobile, plus dark theme.

## Review / deployment instructions

1. Inspect `PATCH-MANIFEST.json`. Apply only the listed new/modified files to `site/` on a new reviewed branch. Do not replace unchanged copied legacy assets.
2. Keep `_source/` and this handoff in the repository for regeneration, or move them to internal documentation if the deployment does not exclude source files. They contain no secrets.
3. Finish browser acceptance above, review both languages and resource cautions, then open a draft PR with the actual check results. This worker did not create a PR.
4. Deployment remains a separate owner decision. If authorised later, deploy the static `site/` contents at the existing domain root, including the new `guias/` and `en/` directories. No backend, package installation or executable rebuild is required.
5. Preserve current assets and route `comandos.html`. Rollback is the previous homepage/sitemap plus removal of the newly added files; no user data or database migration is involved.
