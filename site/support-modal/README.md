# Apoiar / Support modal

Framework-free donate dialog used across Daniel’s sites. Copy this folder into a site and add two tags. No build step, no framework, no accounts, no environment variables.

MediaPulse loads this folder (free variant) from `site/index.html`.

## Files

| File | Role |
| --- | --- |
| `partial.html` | Dialog markup, including the platform logos and the M-Pesa QR |
| `support-modal.css` | Layout and theme variables |
| `support-modal.js` | Open/close, PT/EN copy, free/paid wording, focus trap |
| `demo.html` | Local preview of both variants |

The script fetches `partial.html` from the same folder. You can instead paste the contents of `partial.html` just before `</body>`; if `#dp-support` is already in the page, the script does not fetch.

## Drop-in

```html
<link rel="stylesheet" href="support-modal/support-modal.css">
<!-- optional: paste partial.html here -->
<button type="button" data-dp-support-open>Apoiar</button>
<script
  src="support-modal/support-modal.js"
  data-variant="free"
  data-theme="auto"
  defer
></script>
```

Empty slots grow a button on their own (same idea as the MediaPulse header and footer):

```html
<span data-dp-support="nav"></span>
<span data-dp-support></span>
<span data-dp-support="link"></span>
```

`data-apoiar` is an alias of `data-dp-support`, so older markup keeps working. `data-theme` on the script tag accepts `auto`, `dark`, `light`, `default-dark`, or `default-light`. An explicit `data-theme="light"` or `data-theme="dark"` on `<html>` (or `class="dark"` / `class="light"`) wins over the default.

Open from script with `SupportModal.open()` and `SupportModal.close()`. `DPApoiar` is the same object.

## Free and paid

Set `data-variant="free"` or `data-variant="paid"` on the script tag. `SupportModal.setVariant('paid')` switches a live dialog.

| | Português | English |
| --- | --- | --- |
| Free | Qualquer valor ajuda a manter estes projetos gratuitos e atualizados. | Any amount helps keep these projects free and updated. |
| Paid | O seu apoio ajuda a desenvolver novas funcionalidades e a manter o projeto em evolução. | Your support helps develop new features and keep the project evolving. |
| Paid note | Apoio voluntário — não substitui a subscrição de um plano. | Voluntary support — it does not replace a plan subscription. |

Language follows `document.documentElement.lang` (`en` / `en-*` → English, anything else → Portuguese) and updates if `lang` changes.

## Theming

Public variables, set on `:root` or on `.dp-support` / `.dp-support[data-mode="dark"]`:

- `--dp-support-accent` (default `#e11d48`)
- `--dp-support-bg`
- `--dp-support-fg`
- `--dp-support-muted`
- `--dp-support-line`
- `--dp-support-card`
- `--dp-support-glow`

`data-accent="#7c3aed"` on the script tag sets the accent inline. Load `support-modal.css` before your overrides.

## Optional script attributes

`data-name`, `data-page`, `data-paypal-email`, `data-paypal`, `data-kofi`, `data-patreon`, `data-gumroad`, `data-liberapay`, `data-float` (`left` or `right`, only when the page has no slot), `data-open="1"`.

## Accessibility

The panel is a `role="dialog"` with `aria-modal`, a name, and a description. Tab cycles inside the dialog. Escape closes the QR zoom first, then the dialog. Focus returns to the button that opened it. The QR control is a button with an accessible name.

## Rules that stay fixed

- M-Pesa is the QR image plus the account name. Do not add a phone number.
- Lemon Squeezy is not part of this snippet. Do not add it.
- Do not put secrets, API keys, or `.env` values in this folder. The PayPal address is the public support address already shown on the sites.

## Preview

Serve the folder over HTTP (a `file://` page cannot fetch the partial):

```bash
python3 -m http.server
```

Open `demo.html`. Query string for screenshots: `?variant=paid&theme=light&lang=en&shot=1` (`variant` is `free` or `paid`, `theme` is `dark` or `light`, `lang` is `pt` or `en`, `shot=1` hides the toolbar and opens the dialog).

## MediaPulse

`site/index.html` links `support-modal/support-modal.css` and `support-modal/support-modal.js` with `data-variant="free"` and `data-theme="default-dark"`. `site/support-modal` points at this folder, so a site whose document root is `site/` still loads the same files.
