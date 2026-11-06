/*! Apoiar v2.6 \u2014 bot\u00e3o + janela de apoio (Daniel Marcos). Paleta dark/light + variantes free/paid. Sem depend\u00eancias, sem servi\u00e7os externos. */
(function () {
  'use strict';
  if (window.__dpApoiar) return; window.__dpApoiar = 1;

  var script = document.currentScript || {};
  var ds = script.dataset || {};
  var CFG = {
    variant: ds.variant === 'paid' ? 'paid' : 'free',
    theme: ds.theme || 'auto',          // auto | dark | light | default-dark | default-light
    accent: ds.accent || '#e11d48',
    page: 'https://danielpro.dev/apoiar/',
    name: 'Daniel Marcos Domingos Manganhe',
    paypalEmail: 'danielmarcos173@gmail.com',
    paypal: 'https://www.paypal.com/myaccount/transfer/homepage/pay',
    kofi: 'https://ko-fi.com/danielpro',
    patreon: 'https://www.patreon.com/cw/DanielMarcos958',
    gumroad: 'https://danielprodev.gumroad.com/',
    liberapay: 'https://liberapay.com/danielmarcos'
  };

  function accentRGB(hex) {
    hex = String(hex || '#e11d48').replace('#', '').trim();
    if (hex.length === 3) hex = hex[0]+hex[0]+hex[1]+hex[1]+hex[2]+hex[2];
    var n = parseInt(hex, 16);
    if (isNaN(n)) n = 0xe11d48;
    return [(n >> 16) & 255, (n >> 8) & 255, n & 255];
  }
  var AR = accentRGB(CFG.accent);
  var A = AR[0] + ',' + AR[1] + ',' + AR[2];
  var glowSoft = 'rgba(' + A + ',.12)';
  var glowMid = 'rgba(' + A + ',.22)';
  var glowHot = 'rgba(' + A + ',.18)';
  var hoverTint = 'rgba(' + A + ',.10)';
  var hoverBorder = 'rgba(' + A + ',.55)';
  var hoverBorderLt = 'rgba(' + A + ',.45)';
  var hoverBgLt = 'rgba(' + A + ',.08)';

  var QR = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 25 25" class="qr" role="img" aria-label="QR M-Pesa"><path fill="#fff" d="M0 0h25v25h-25z"/><path class="qrline" stroke="#111" d="M2 2.5h7m2 0h4m1 0h7m-21 1h1m5 0h1m4 0h1m2 0h1m5 0h1m-21 1h1m1 0h3m1 0h1m1 0h2m2 0h1m1 0h1m1 0h3m1 0h1m-21 1h1m1 0h3m1 0h1m4 0h1m2 0h1m1 0h3m1 0h1m-21 1h1m1 0h3m1 0h1m1 0h2m4 0h1m1 0h3m1 0h1m-21 1h1m5 0h1m1 0h1m2 0h2m1 0h1m5 0h1m-21 1h7m1 0h1m1 0h1m1 0h1m1 0h7m-12 1h2m-10 1h1m2 0h1m1 0h1m1 0h1m2 0h3m1 0h2m1 0h1m-16 1h1m4 0h2m2 0h2m2 0h1m1 0h2m-19 1h11m3 0h2m3 0h1m-21 1h1m2 0h2m2 0h3m1 0h2m2 0h2m1 0h1m-19 1h5m1 0h4m1 0h1m1 0h1m1 0h1m1 0h4m-13 1h2m1 0h1m1 0h3m2 0h1m-19 1h7m2 0h1m3 0h8m-21 1h1m5 0h1m4 0h1m4 0h1m3 0h1m-21 1h1m1 0h3m1 0h1m1 0h4m1 0h2m2 0h2m1 0h1m-21 1h1m1 0h3m1 0h1m5 0h5m2 0h2m-21 1h1m1 0h3m1 0h1m2 0h2m1 0h1m2 0h2m1 0h1m-19 1h1m5 0h1m1 0h1m1 0h3m4 0h1m2 0h1m-21 1h7m2 0h1m1 0h1m1 0h3m2 0h3"/></svg>';

  var TX = {en: {
    btn: 'Support', title: 'Enjoyed my work? Support the project',
    subFree: 'Any amount helps keep these projects free and updated.',
    subPaid: 'Your support helps develop new features and keep the project evolving.',
    notePaid: 'Voluntary support \u2014 it does not replace a plan subscription.',
    mpesaHint: 'Scan the QR with the M-Pesa app', copied: 'Copied!',
    paypal: 'PayPal', paypalHint: 'PayPal \u2192 Send \u2192 this email', copyMail: 'Copy email', openPP: 'Open PayPal',
    kofi: 'Ko-fi', kofiHint: 'Card or PayPal',
    patreon: 'Patreon', patreonHint: 'Monthly or one-time',
    gumroad: 'Gumroad', gumroadHint: 'Digital products',
    liberapay: 'Liberapay', liberapayHint: 'Weekly, monthly or yearly',
    close: 'Close', more: 'Support page', thanks: 'Thank you!', enlargeQr: 'Enlarge QR'
  }, pt: {
    btn: 'Apoiar', title: 'Gostou do que fiz? Apoie o projeto',
    subFree: 'Qualquer valor ajuda a manter estes projetos gratuitos e atualizados.',
    subPaid: 'O seu apoio ajuda a desenvolver novas funcionalidades e a manter o projeto em evolu\u00e7\u00e3o.',
    notePaid: 'Apoio volunt\u00e1rio \u2014 n\u00e3o substitui a subscri\u00e7\u00e3o de um plano.',
    mpesaHint: 'Leia o QR com a app M-Pesa', copied: 'Copiado!',
    paypal: 'PayPal', paypalHint: 'PayPal \u2192 Enviar \u2192 este email', copyMail: 'Copiar email', openPP: 'Abrir PayPal',
    kofi: 'Ko-fi', kofiHint: 'Cart\u00e3o ou PayPal',
    patreon: 'Patreon', patreonHint: 'Mensal ou \u00fanico',
    gumroad: 'Gumroad', gumroadHint: 'Produtos digitais',
    liberapay: 'Liberapay', liberapayHint: 'Semanal, mensal ou anual',
    close: 'Fechar', more: 'P\u00e1gina de apoio', thanks: 'Obrigado!', enlargeQr: 'Ampliar QR'
  }};
  var T = TX.pt;
  function lang() { T = /^en/i.test(document.documentElement.lang || '') ? TX.en : TX.pt; return T; }
  lang();


  var HEART = '<svg viewBox="0 0 24 24" width="1em" height="1em" aria-hidden="true" fill="currentColor"><path d="M12 21s-7.5-4.6-10-9.3C.4 8.4 2.2 4.5 5.9 4.1 8 3.9 9.8 5 12 7.4 14.2 5 16 3.9 18.1 4.1c3.7.4 5.5 4.3 3.9 7.6C19.5 16.4 12 21 12 21z"/></svg>';

  /* Log\u00f3tipos oficiais das marcas (embutidos, sem hotlink): M-Pesa (Vodacom), Ko-fi, PayPal, Patreon, Gumroad, Liberapay. */
  var LOGOS = {
    mpesa: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" aria-hidden="true"><rect width="48" height="48" rx="11" fill="#e60000"/><g transform="translate(8 8) scale(.08)"><path transform="translate(107 0)" fill="#fff" d="M0 0 C1 0 2 0 2 0 C3 0 4 -0 5 -0 C7 -0 10 -0 12 -0 C14 -0 16 -0 18 0 C23 0 28 0 34 0 C37 0 40 0 43 0 C54 0 66 0 77 0 C87 0 98 0 108 0 C117 0 126 0 135 0 C140 0 145 0 151 0 C156 0 161 0 166 0 C168 0 170 0 171 0 C174 0 176 0 179 0 C180 0 180 0 181 0 C185 0 187 0 190 3 C190 4 191 4 192 5 C192 6 193 7 194 8 C199 14 201 21 201 30 C201 30 201 31 201 32 C201 34 201 36 201 38 C201 40 201 41 201 43 C201 47 201 51 201 55 C201 59 201 63 201 67 C201 74 201 82 201 90 C191 88 184 84 175 79 C175 62 175 44 175 26 C120 26 64 26 8 26 C8 70 8 114 9 159 C17 150 17 150 20 145 C24 147 26 150 28 154 C37 167 48 177 62 183 C68 182 72 178 76 174 C77 174 78 173 78 172 C80 170 82 169 83 167 C85 166 87 165 88 163 C90 162 92 160 94 158 C97 154 100 151 102 148 C121 128 135 104 147 79 C151 81 155 82 160 85 C160 85 161 85 162 86 C164 87 166 88 168 89 C170 90 172 91 175 92 C176 93 178 93 179 94 C186 97 192 100 198 103 C207 108 216 112 225 116 C225 116 225 116 227 117 C233 120 239 123 246 126 C272 139 272 139 284 144 C278 160 261 174 249 184 C248 185 247 186 246 186 C191 233 110 238 42 234 C31 233 20 232 9 231 C8 259 8 287 8 316 C63 316 118 316 175 316 C175 290 175 263 175 235 C181 233 187 230 193 228 C194 228 195 228 196 227 C197 227 199 226 201 225 C201 245 201 264 201 284 C201 293 201 302 201 312 C201 320 201 327 201 335 C201 340 201 344 201 348 C201 352 201 356 201 360 C201 361 201 363 201 364 C201 376 200 385 193 394 C192 395 192 395 191 397 C189 399 189 399 187 400 C184 400 182 400 179 400 C178 400 178 400 177 400 C174 400 172 400 169 400 C167 400 165 400 163 401 C146 401 129 401 112 401 C99 401 86 401 74 401 C65 401 56 401 47 401 C42 401 36 401 31 401 C26 401 21 401 16 401 C14 401 12 401 10 401 C-4 402 -4 402 -8 398 C-10 396 -11 395 -12 392 C-13 391 -13 390 -14 390 C-14 389 -15 388 -15 387 C-16 386 -16 385 -16 385 C-18 380 -18 377 -18 372 C-18 371 -18 370 -18 368 C-18 367 -18 365 -18 364 C-18 362 -18 361 -18 360 C-18 356 -18 351 -19 347 C-19 343 -19 339 -19 335 C-19 327 -19 319 -19 311 C-19 302 -19 292 -19 283 C-19 265 -19 246 -19 227 C-20 227 -20 227 -22 227 C-60 220 -60 220 -78 215 C-80 215 -81 214 -83 214 C-87 213 -90 212 -94 211 C-94 211 -94 211 -94 210 C-94 210 -93 210 -92 210 C-67 210 -46 204 -25 190 C-24 190 -23 189 -22 189 C-20 187 -19 185 -19 182 C-18 178 -18 174 -19 170 C-18 168 -18 167 -18 165 C-18 162 -18 158 -19 155 C-19 150 -19 144 -19 139 C-19 128 -19 118 -19 107 C-19 94 -19 80 -19 67 C-19 62 -19 56 -19 51 C-19 48 -19 44 -19 41 C-19 39 -19 38 -19 36 C-19 25 -18 15 -11 6 C-10 6 -10 5 -9 4 C-6 1 -5 0 0 0 Z M74 341 C73 346 73 350 73 354 C73 357 73 359 73 361 C73 362 73 364 73 365 C73 366 73 368 73 369 C74 372 74 372 75 374 C78 376 80 376 84 376 C84 376 84 376 86 376 C87 376 88 376 90 376 C92 376 94 376 96 376 C97 376 99 376 100 376 C101 376 102 376 104 376 C107 375 107 375 109 374 C110 371 110 369 110 365 C110 365 110 364 110 363 C110 362 110 360 110 359 C110 357 110 354 110 352 C110 342 110 342 108 339 C105 339 105 339 101 339 C100 339 100 339 98 339 C97 339 96 339 94 339 C92 339 90 339 87 339 C86 339 84 339 83 339 C82 339 80 339 79 339 C76 339 76 339 74 341 Z M32 348 C32 350 32 351 32 352 C32 353 32 354 32 355 C32 356 32 356 32 357 C32 358 32 359 32 360 C32 361 32 362 32 362 C32 363 32 364 32 364 C32 366 32 366 33 367 C36 367 38 367 41 367 C42 367 43 367 44 367 C50 367 50 367 51 366 C51 365 51 364 51 362 C51 362 51 361 51 360 C51 359 51 359 51 357 C51 356 51 356 51 355 C51 354 51 354 51 352 C51 352 51 351 51 350 C51 348 51 348 50 347 C47 347 44 347 41 347 C40 347 40 347 39 347 C33 347 33 347 32 348 Z M133 349 C133 352 132 355 133 357 C133 358 133 359 133 360 C133 365 133 365 134 367 C137 367 140 367 143 367 C144 367 144 367 145 367 C146 367 147 367 148 367 C148 367 149 367 150 367 C152 367 152 367 153 366 C153 365 153 364 153 362 C153 362 153 361 153 360 C153 359 153 359 153 357 C153 356 153 356 153 355 C153 354 153 354 153 352 C153 352 153 351 153 350 C153 348 153 348 152 347 C149 347 146 347 143 347 C143 347 143 347 141 347 C135 347 135 347 133 349 Z "/></g></svg>',
    kofi: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" aria-hidden="true"><rect width="48" height="48" rx="11" fill="#fff"/><svg x="8" y="8" width="32" height="32" viewBox="0 0 24 24"><path fill="#202020" d="M11.351 2.715c-2.7 0-4.986.025-6.83.26C2.078 3.285 0 5.154 0 8.61c0 3.506.182 6.13 1.585 8.493 1.584 2.701 4.233 4.182 7.662 4.182h.83c4.209 0 6.494-2.234 7.637-4a9.5 9.5 0 0 0 1.091-2.338C21.792 14.688 24 12.22 24 9.208v-.415c0-3.247-2.13-5.507-5.792-5.87-1.558-.156-2.65-.208-6.857-.208M11.351 4.662c4.208 0 5.09.052 6.571.182 2.624.311 4.13 1.584 4.13 4v.39c0 2.156-1.792 3.844-3.87 3.844h-.935l-.156.649c-.208 1.013-.597 1.818-1.039 2.546-.909 1.428-2.545 3.064-5.922 3.064h-.805c-2.571 0-4.831-.883-6.078-3.195-1.09-2-1.298-4.155-1.298-7.506 0-2.181.857-3.402 3.012-3.714 1.533-.233 3.559-.26 6.39-.26M17.898 6.949c-.416 0-.65.234-.65.546v2.935c0 .311.234.545.65.545 1.324 0 2.051-.754 2.051-2s-.727-2.026-2.052-2.026"/><path fill="#ff5e5b" d="M7.507 7.131c-1.818 0-3.013 1.48-3.013 3.142 0 1.533.858 2.857 1.949 3.897.727.701 1.87 1.429 2.649 1.896a1.47 1.47 0 0 0 1.507 0c.78-.467 1.922-1.195 2.623-1.896 1.117-1.039 1.974-2.364 1.974-3.897 0-1.662-1.247-3.142-3.039-3.142-1.065 0-1.792.545-2.338 1.298-.493-.753-1.246-1.298-2.312-1.298"/></svg></svg>',
    paypal: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" aria-hidden="true"><rect width="48" height="48" rx="11" fill="#fff"/><svg x="10" y="8" width="28" height="32" viewBox="0 0 154.728 190.5"><g transform="translate(898.192 276.071)"><path d="M-837.663-237.968a5.49 5.49 0 0 0-5.423 4.633l-9.013 57.15-8.281 52.514-.005.044.01-.044 8.281-52.514c.421-2.669 2.719-4.633 5.42-4.633h26.404c26.573 0 49.127-19.387 53.246-45.658.314-1.996.482-3.973.52-5.924v-.003h-.003c-6.753-3.543-14.683-5.565-23.372-5.565z" fill="#001c64"/><path d="M-766.506-232.402c-.037 1.951-.207 3.93-.52 5.926-4.119 26.271-26.673 45.658-53.246 45.658h-26.404c-2.701 0-4.999 1.964-5.42 4.633l-8.281 52.514-5.197 32.947a4.46 4.46 0 0 0 4.405 5.153h28.66a5.49 5.49 0 0 0 5.423-4.633l7.55-47.881c.423-2.669 2.722-4.636 5.423-4.636h16.876c26.573 0 49.124-19.386 53.243-45.655 2.924-18.649-6.46-35.614-22.511-44.026z" fill="#0070e0"/><path d="M-870.225-276.071a5.49 5.49 0 0 0-5.423 4.636l-22.489 142.608a4.46 4.46 0 0 0 4.405 5.156h33.351l8.281-52.514 9.013-57.15a5.49 5.49 0 0 1 5.423-4.633h47.782c8.691 0 16.621 2.025 23.375 5.563.46-23.917-19.275-43.666-46.412-43.666z" fill="#003087"/></g></svg></svg>',
    patreon: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" aria-hidden="true"><rect width="48" height="48" rx="11" fill="#000"/><svg x="11" y="11" width="26" height="26" viewBox="0 0 1080 1080"><path fill="#fff" d="M1033.05,324.45c-0.19-137.9-107.59-250.92-233.6-291.7c-156.48-50.64-362.86-43.3-512.28,27.2 C106.07,145.41,49.18,332.61,47.06,519.31c-1.74,153.5,13.58,557.79,241.62,560.67c169.44,2.15,194.67-216.18,273.07-321.33 c55.78-74.81,127.6-95.94,216.01-117.82C929.71,603.22,1033.27,483.3,1033.05,324.45z"/></svg></svg>',
    gumroad: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" aria-hidden="true"><rect width="48" height="48" rx="11" fill="#ff90e8"/><svg x="9" y="9" width="30" height="30" viewBox="0 0 24 24"><path fill="#000" d="M12 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0Zm-.007 5.12c4.48 0 5.995 3.025 6.064 4.744h-3.239c-.069-.962-.897-2.406-2.896-2.406-2.136 0-3.514 1.857-3.514 4.126 0 2.27 1.378 4.125 3.514 4.125 1.93 0 2.758-1.512 3.103-3.025h-3.103v-1.238h6.509v6.327h-2.855v-3.989c-.207 1.444-1.102 4.264-4.617 4.264-3.516 0-5.584-2.82-5.584-6.326 0-3.645 2.276-6.602 6.618-6.602z"/></svg></svg>',
    liberapay: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 80" aria-hidden="true"><g transform="matrix(.83012 0 0 .83012-135.4-247.7)"><path d="m259.55 385.57c0 5.145-4.169 9.318-9.318 9.318h-77.74c-5.144 0-9.318-4.174-9.318-9.318v-77.74c0-5.145 4.174-9.318 9.318-9.318h77.74c5.149 0 9.318 4.173 9.318 9.318v77.74" fill="#f6c915"/><g fill="#fff"><path d="m202.45 366.03c-3.104 0-5.541-.405-7.311-1.213-1.77-.809-3.039-1.912-3.803-3.313-.766-1.398-1.137-3-1.115-4.818.021-1.814.272-3.748.754-5.803l8.327-34.817 10.164-1.573-9.114 37.768c-.175.786-.273 1.508-.295 2.163-.023.655.098 1.235.36 1.737.262.504.71.908 1.344 1.213.633.307 1.519.504 2.656.591l-1.967 8.06"/><path d="m239.16 344.33c0 3.19-.525 6.108-1.574 8.753-1.049 2.646-2.503 4.929-4.36 6.852-1.858 1.925-4.087 3.421-6.688 4.491-2.601 1.07-5.432 1.607-8.49 1.607-1.487 0-2.973-.132-4.459-.395l-2.951 11.869h-9.704l10.884-45.37c1.748-.524 3.748-.994 5.999-1.41 2.252-.415 4.689-.622 7.312-.622 2.448 0 4.558.371 6.327 1.114 1.771.743 3.224 1.76 4.361 3.049 1.136 1.29 1.977 2.798 2.523 4.524.546 1.726.82 3.574.82 5.542m-23.802 13.442c.743.175 1.661.262 2.754.262 1.704 0 3.256-.316 4.655-.951 1.398-.633 2.59-1.518 3.574-2.655.982-1.136 1.747-2.501 2.294-4.098.546-1.595.819-3.354.819-5.278 0-1.879-.416-3.475-1.245-4.787-.831-1.311-2.273-1.967-4.327-1.967-1.4 0-2.711.131-3.935.394l-4.589 19.08"/></g></g></svg>'
  };

  function isDark() {
    if (CFG.theme === 'dark') return true;
    if (CFG.theme === 'light') return false;
    var els = [document.documentElement, document.body];
    for (var i = 0; i < els.length; i++) {
      var e = els[i]; if (!e) continue;
      var c = ' ' + (e.className && e.className.baseVal === undefined ? e.className : '') + ' ';
      if (/ (dark|dark-mode|dark-theme|theme-dark)\s/.test(c)) return true;
      var a = (e.getAttribute('data-theme') || e.getAttribute('data-bs-theme') || e.getAttribute('data-mode') || '').toLowerCase();
      if (a === 'dark') return true;
      if (a === 'light' || / (light|light-mode|light-theme|theme-light)\s/.test(c)) return false;
    }
    if (CFG.theme === 'default-dark') return true;
    if (CFG.theme === 'default-light') return false;
    return !!(window.matchMedia && matchMedia('(prefers-color-scheme: dark)').matches);
  }

  /* ---------- bot\u00e3o (light DOM, herda a fonte do site) ---------- */
  var css = document.createElement('style');
  css.textContent =
    'button.dp-apoiar-btn{display:inline-flex!important;align-items:center;gap:.4em;font:inherit;font-size:.9em;font-weight:600;line-height:1.2;' +
    'padding:.45em .95em!important;border-radius:999px!important;border:1px solid currentColor!important;background:transparent;color:inherit;cursor:pointer;' +
    'opacity:.92;transition:opacity .15s,transform .15s,background .15s;white-space:nowrap;text-decoration:none;margin:0;vertical-align:middle;box-shadow:none}' +
    'button.dp-apoiar-btn:hover{opacity:1;transform:translateY(-1px);background:rgba(127,127,127,.12)}' +
    '.dp-apoiar-btn:focus-visible{outline:2px solid ' + CFG.accent + ';outline-offset:2px}' +
    'button.dp-apoiar-btn svg{display:inline-block!important;color:' + CFG.accent + ';flex:none;width:1em;height:1em}' +
    'button.dp-apoiar-btn.dp-link{border:0!important;padding:0!important;font-size:inherit;font-weight:inherit;background:none}' +
    'button.dp-apoiar-btn.dp-link:hover{background:none;text-decoration:underline}' +
    'button.dp-apoiar-btn.dp-nav{font-size:.85em;padding:.4em .8em!important;gap:.35em}' +
    '@media (max-width:560px){button.dp-apoiar-btn.dp-nav{padding:0!important;width:34px;height:34px;justify-content:center}button.dp-apoiar-btn.dp-nav .dp-l{display:none}}' +
    'button.dp-apoiar-float{position:fixed;left:16px;bottom:16px;z-index:2147483000;background:' + CFG.accent + ';color:#fff;border-color:transparent;box-shadow:0 6px 20px rgba(0,0,0,.25)}' +
    'button.dp-apoiar-float svg{color:#fff}button.dp-apoiar-float:hover{background:' + CFG.accent + '}';
  (document.head || document.documentElement).appendChild(css);

  function makeBtn(variant) {
    var b = document.createElement('button');
    b.type = 'button';
    b.className = 'dp-apoiar-btn' + (variant === 'link' ? ' dp-link' : variant === 'nav' ? ' dp-nav' : '');
    b.setAttribute('aria-label', T.btn);
    b.setAttribute('aria-haspopup', 'dialog');
    b.innerHTML = HEART + '<span class="dp-l">' + T.btn + '</span>';
    b.addEventListener('click', function (e) { e.preventDefault(); open(b); });
    return b;
  }

  /* ---------- painel (shadow DOM, isolado do CSS do site) ---------- */
  var PANEL_CSS =
    ':host{all:initial}' +
    '.p{--bg:#fefefe;--fg:#0d0d0f;--mut:#5a5a68;--line:#e8e3e9;--card:#faf8fb;--ac:' + CFG.accent + ';--glow:' + glowSoft + ';' +
    'font:14px/1.4 system-ui,-apple-system,"Segoe UI",Roboto,Inter,Arial,sans-serif;color:var(--fg);' +
    'background:radial-gradient(120% 60% at 0% 0%,' + glowSoft.replace(',.12)', ',.08)') + ',transparent 60%),linear-gradient(180deg,#ffffff,#f7f5f8);' +
    'border:1px solid rgba(15,13,21,.06);box-shadow:inset 0 1px 0 rgba(255,255,255,.9);' +
    'border-radius:16px;box-sizing:border-box;width:100%;max-width:390px;padding:18px;position:relative;text-align:left}' +
    '.p:not(.dk) .card,.p:not(.dk) a.b{background:linear-gradient(180deg,rgba(255,255,255,.95),rgba(250,248,251,.9));border-color:rgba(15,13,21,.08);box-shadow:0 1px 3px rgba(15,13,21,.04),inset 0 1px 0 rgba(255,255,255,.8)}' +
    '.p:not(.dk) a.b:hover{border-color:' + hoverBorderLt + ';background:linear-gradient(180deg,rgba(255,255,255,.98),rgba(255,255,255,.9));box-shadow:0 6px 18px ' + glowSoft + '}' +
    '.p:not(.dk) .cp{background:rgba(255,255,255,.7);border-color:rgba(15,13,21,.1)}.p:not(.dk) .cp:hover{border-color:' + hoverBorder + ';background:' + hoverBgLt + '}' +
    '.p.dk{--bg:#131316;--fg:#f4f4f5;--mut:#a1a1aa;--line:rgba(255,255,255,.08);--card:#1c1c20;--glow:' + glowMid + ';' +
    'background:radial-gradient(120% 60% at 0% 0%,' + glowSoft + ',transparent 60%),linear-gradient(180deg,#18181c,#111114);border:1px solid rgba(255,255,255,.07);box-shadow:inset 0 1px 0 rgba(255,255,255,.05)}' +
    '.p.dk .card,.p.dk a.b{background:linear-gradient(180deg,rgba(255,255,255,.045),rgba(255,255,255,.015));border-color:rgba(255,255,255,.08);box-shadow:inset 0 1px 0 rgba(255,255,255,.05),0 2px 6px rgba(0,0,0,.25)}' +
    '.p.dk a.b:hover{border-color:' + hoverBorder + ';background:linear-gradient(180deg,' + hoverTint + ',rgba(255,255,255,.02));box-shadow:inset 0 1px 0 rgba(255,255,255,.06),0 6px 18px ' + glowHot + '}' +
    '.p.dk .cp{background:rgba(255,255,255,.04);border-color:rgba(255,255,255,.12)}.p.dk .cp:hover{border-color:' + hoverBorder + ';background:' + hoverBgLt + '}' +
    '.p.dk h2 svg{color:var(--ac)}.p.dk .x:hover{background:rgba(255,255,255,.08)}' +
    '.p *{box-sizing:border-box}.p:focus{outline:none}' +
    'h2{margin:0 30px 4px 0;font-size:17.5px;line-height:1.2;font-weight:750;display:flex;gap:7px;align-items:center}' +
    'h2 svg{color:var(--ac);flex:none;filter:drop-shadow(0 0 8px var(--glow))}' +
    '.sub{margin:0 0 10px;color:var(--mut);font-size:12.5px}' +
    '.note{margin:0 0 14px;color:var(--mut);font-size:11.5px;padding:8px 10px;border-radius:10px;border:1px dashed var(--line);background:transparent}' +
    '.card{border:1px solid var(--line);background:var(--card);border-radius:12px;padding:12px;margin-bottom:8px;' +
    'box-shadow:0 1px 3px rgba(0,0,0,.04)}' +
    
    '.mp{display:flex;gap:12px;align-items:center}' +
    '.qr{width:80px;height:80px;flex:none;border-radius:8px;background:#fff;padding:0;display:block;cursor:pointer;' +
    'transition:transform .15s,box-shadow .15s;border:2px solid transparent}' +
    '.qr:hover{transform:scale(1.02);box-shadow:0 4px 12px var(--glow);border-color:var(--ac)}' +
    '.qr:focus-visible{outline:2px solid var(--ac);outline-offset:2px}' +
    '.lab{display:flex;align-items:center;gap:7px;font-size:13px;font-weight:750;letter-spacing:.01em;color:#e60000;margin-bottom:2px}' +
    '.lab .ic{width:30px;height:30px;border-radius:8px}.p.dk .lab{color:#ff4d4d}' +
    '.num{font-size:18px;font-weight:750;letter-spacing:.02em;margin:2px 0 0;white-space:nowrap}' +
    '.nm{font-size:13px;font-weight:650;margin:3px 0 5px}' +
    '.scan{font-size:12px;color:var(--mut)}' +
    '.hint{font-size:11.5px;color:var(--mut);margin-top:4px}' +
    'button,a.b{font:inherit;cursor:pointer;text-decoration:none}' +
    '.grid{display:grid;grid-template-columns:1fr 1fr;gap:8px;margin-bottom:8px}' +
    '.cp{display:inline-flex;align-items:center;gap:6px;border:1px solid var(--line);background:var(--bg);color:var(--fg);' +
    'border-radius:999px;padding:6px 11px;font-size:12.5px;font-weight:600;transition:border-color .15s}' +
    '.cp:hover{border-color:var(--ac)}.cp.ok{border-color:#16a34a;color:#16a34a}' +
    'a.b{display:flex;flex-direction:column;align-items:center;justify-content:center;gap:6px;border:1px solid var(--line);background:var(--card);color:var(--fg);border-radius:12px;padding:10px 8px;transition:border-color .15s,transform .15s,box-shadow .15s;text-align:center}' +
    'a.b:hover{border-color:var(--ac);transform:translateY(-1px);box-shadow:0 4px 12px var(--glow)}' +
    '.ic{width:34px;height:34px;border-radius:9px;display:block;flex:none;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,.18),inset 0 0 0 1px rgba(0,0,0,.06)}' +
    '.p.dk .ic{box-shadow:0 1px 3px rgba(0,0,0,.4),0 0 0 1px rgba(255,255,255,.14)}' +
    '.ic svg{display:block;width:100%;height:100%}' +
    'a.b strong{display:block;font-size:13px;font-weight:650}a.b small{display:block;color:var(--mut);font-size:11px;margin-top:2px}' +
    '.pp{display:flex;gap:11px;align-items:flex-start}' +
    '.ppb strong{display:block;font-size:13px;font-weight:650}.mail{font-weight:650;font-size:13px;margin:2px 0 7px;word-break:break-all}.row{display:flex;gap:7px;flex-wrap:wrap}a.cp{text-decoration:none}' +
    '.x{position:absolute;top:12px;right:12px;width:30px;height:30px;border-radius:50%;border:0;background:transparent;color:var(--mut);font-size:20px;line-height:1;transition:background .15s,color .15s}' +
    '.x:hover{background:var(--card);color:var(--fg)}' +
    '.ft{text-align:center;font-size:12px;margin-top:6px}.ft a{color:var(--mut);transition:color .15s}.ft a:hover{color:var(--ac)}' +
    'button:focus-visible,a:focus-visible{outline:2px solid var(--ac);outline-offset:2px}' +
    '.ov{position:fixed;inset:0;z-index:2147483646;background:rgba(2,6,23,.62);display:flex;align-items:center;justify-content:center;padding:16px;' +
    'opacity:0;transition:opacity .16s;-webkit-backdrop-filter:blur(3px);backdrop-filter:blur(3px);overflow-y:auto}' +
    '.ov.on{opacity:1}.ov .p{box-shadow:0 24px 60px rgba(0,0,0,.35);transform:translateY(8px) scale(.98);transition:transform .16s;margin:auto}.ov.on .p{transform:none}' +
    '.lb{position:fixed;inset:0;z-index:2147483647;background:rgba(0,0,0,.85);display:flex;align-items:center;justify-content:center;padding:20px;opacity:0;transition:opacity .2s}' +
    '.lb.on{opacity:1}.lb-qr{max-width:min(90vw,400px);max-height:90vh;border-radius:12px;background:#fff;padding:16px;box-shadow:0 25px 50px rgba(0,0,0,.5);' +
    'transform:scale(.95);transition:transform .2s}.lb.on .lb-qr{transform:scale(1)}.lb-qr svg{display:block;width:100%;height:auto}' +
    '.lb-x{position:absolute;top:12px;right:12px;width:40px;height:40px;border-radius:50%;border:0;background:rgba(255,255,255,.15);color:#fff;font-size:28px;' +
    'line-height:1;cursor:pointer;transition:background .15s;-webkit-backdrop-filter:blur(10px);backdrop-filter:blur(10px)}' +
    '.lb-x:hover{background:rgba(255,255,255,.25)}.lb-x:focus-visible{outline:2px solid #fff;outline-offset:2px}' +
    '@media (max-width:380px){.p{padding:16px}.qr{width:72px;height:72px}}' +
    '@media (prefers-reduced-motion:reduce){.ov,.ov .p,a.b,.qr,.lb,.lb-qr{transition:none}}';

  function panelHTML(modal) {
    var qrBtn = '<button class="qr" type="button" aria-label="' + T.enlargeQr + '" title="' + T.enlargeQr + '">' + QR + '</button>';
    return '<div class="p" role="' + (modal ? 'dialog' : 'region') + '" aria-modal="' + (modal ? 'true' : 'false') + '" aria-labelledby="dpa-t">' +
      (modal ? '<button class="x" type="button" aria-label="' + T.close + '">&times;</button>' : '') +
      '<h2 id="dpa-t">' + HEART + '<span>' + T.title + '</span></h2>' +
      '<p class="sub">' + (CFG.variant === 'paid' ? T.subPaid : T.subFree) + '</p>' +
      (CFG.variant === 'paid' ? '<p class="note">' + T.notePaid + '</p>' : '') +
      '<div class="card mp">' + qrBtn +
        '<div><div class="lab"><span class="ic">' + LOGOS.mpesa + '</span><span>M-Pesa</span></div><div class="nm">' + CFG.name + '</div>' +
        '<div class="scan">' + T.mpesaHint + '</div></div></div>' +
      '<div class="grid">' +
      '<a class="b" href="' + CFG.patreon + '" target="_blank" rel="noopener"><span class="ic">' + LOGOS.patreon + '</span>' +
        '<strong>' + T.patreon + '</strong><small>' + T.patreonHint + '</small></a>' +
      '<a class="b" href="' + CFG.gumroad + '" target="_blank" rel="noopener"><span class="ic">' + LOGOS.gumroad + '</span>' +
        '<strong>' + T.gumroad + '</strong><small>' + T.gumroadHint + '</small></a>' +
      '<a class="b" href="' + CFG.liberapay + '" target="_blank" rel="noopener"><span class="ic">' + LOGOS.liberapay + '</span>' +
        '<strong>' + T.liberapay + '</strong><small>' + T.liberapayHint + '</small></a>' +
      '<a class="b" href="' + CFG.kofi + '" target="_blank" rel="noopener"><span class="ic">' + LOGOS.kofi + '</span>' +
        '<strong>' + T.kofi + '</strong><small>' + T.kofiHint + '</small></a>' +
      '</div>' +
      '<div class="card pp"><span class="ic">' + LOGOS.paypal + '</span><div class="ppb"><strong>' + T.paypal + '</strong>' +
        '<div class="mail">' + CFG.paypalEmail + '</div><div class="row"><button class="cp cpm" type="button">' + T.copyMail + '</button>' +
        '<a class="cp" href="' + CFG.paypal + '" target="_blank" rel="noopener">' + T.openPP + ' \u2197</a></div><div class="hint">' + T.paypalHint + '</div></div></div>' +
      (modal && location.hostname.replace(/^www\./, '') + location.pathname.replace(/\/$/, '') !== 'danielpro.dev/apoiar' ? '<div class="ft"><a href="' + CFG.page + '" target="_blank" rel="noopener">' + T.more + ' \u2192</a></div>' : '') +
      '</div>';
  }

  function copyText(txt, btn, label) {
    function done() { btn.textContent = T.copied; btn.classList.add('ok'); setTimeout(function () { btn.textContent = label; btn.classList.remove('ok'); }, 1800); }
    if (navigator.clipboard && window.isSecureContext) { navigator.clipboard.writeText(txt).then(done, fallback); } else fallback();
    function fallback() {
      var t = document.createElement('textarea'); t.value = txt; t.setAttribute('readonly', ''); t.style.cssText = 'position:fixed;left:-9999px;top:0';
      document.body.appendChild(t); t.select(); try { document.execCommand('copy'); done(); } catch (e) {} document.body.removeChild(t);
    }
  }

  function wire(root) {
    var cm = root.querySelector('.cpm');
    if (cm) cm.addEventListener('click', function () { copyText(CFG.paypalEmail, cm, T.copyMail); });
    var qr = root.querySelector('.qr');
    if (qr) {
      qr.addEventListener('click', function (e) { e.stopPropagation(); openLightbox(root); });
      qr.addEventListener('keydown', function (e) {
        if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); openLightbox(root); }
      });
    }
  }

  var host = null, lastFocus = null, prevOverflow = '', lightbox = null;

  function openLightbox(root) {
    if (lightbox) return;
    var lb = document.createElement('div');
    lb.className = 'lb';
    lb.innerHTML = '<button class="lb-x" type="button" aria-label="' + T.close + '">&times;</button><div class="lb-qr">' + QR + '</div>';
    lb.addEventListener('click', function (e) { if (e.target === lb) closeLightbox(); });
    lb.querySelector('.lb-x').addEventListener('click', closeLightbox);
    root.appendChild(lb);
    lightbox = lb;
    requestAnimationFrame(function () { lb.classList.add('on'); });
    setTimeout(function () { lb.querySelector('.lb-x').focus({ preventScroll: true }); }, 50);
  }

  function closeLightbox() {
    if (!lightbox) return;
    lightbox.remove();
    lightbox = null;
  }

  function open(from) {
    lang();
    lastFocus = from || document.activeElement;
    if (host) close(true);
    host = document.createElement('div');
    host.setAttribute('data-dp-apoiar-modal', '');
    var sr = host.attachShadow ? host.attachShadow({ mode: 'open' }) : host;
    sr.innerHTML = '<style>' + PANEL_CSS + '</style><div class="ov">' + panelHTML(true) + '</div>';
    var p = sr.querySelector('.p'); if (isDark()) p.classList.add('dk');
    var ov = sr.querySelector('.ov');
    ov.addEventListener('click', function (e) { if (e.target === ov) close(); });
    sr.querySelector('.x').addEventListener('click', function () { close(); });
    wire(sr);
    document.body.appendChild(host);
    prevOverflow = document.documentElement.style.overflow;
    document.documentElement.style.overflow = 'hidden';
    requestAnimationFrame(function () { ov.classList.add('on'); });
    setTimeout(function () { p.setAttribute('tabindex', '-1'); p.focus({ preventScroll: true }); }, 30);
    document.addEventListener('keydown', onKey, true);
  }
  function onKey(e) {
    if (!host) return;
    if (e.key === 'Escape') {
      e.preventDefault();
      if (lightbox) closeLightbox();
      else close();
      return;
    }
    if (e.key === 'Tab') {
      var sr = host.shadowRoot || host, f = sr.querySelectorAll('button,a[href]');
      if (!f.length) return;
      var first = f[0], last = f[f.length - 1], act = sr.activeElement || document.activeElement;
      if (e.shiftKey && act === first) { e.preventDefault(); last.focus(); }
      else if (!e.shiftKey && act === last) { e.preventDefault(); first.focus(); }
    }
  }
  function close(silent) {
    if (!host) return;
    closeLightbox();
    host.remove(); host = null;
    document.documentElement.style.overflow = prevOverflow;
    document.removeEventListener('keydown', onKey, true);
    if (!silent && lastFocus && lastFocus.focus) lastFocus.focus({ preventScroll: true });
  }

  function renderInline(el) {
    lang();
    var sr = el.attachShadow ? el.attachShadow({ mode: 'open' }) : el;
    sr.innerHTML = '<style>' + PANEL_CSS + '.p{margin:0 auto;box-shadow:0 10px 40px rgba(0,0,0,.12)}</style>' + panelHTML(false);
    var p = sr.querySelector('.p');
    function th() { p.classList.toggle('dk', isDark()); }
    th();
    new MutationObserver(th).observe(document.documentElement, { attributes: true, attributeFilter: ['class', 'data-theme'] });
    if (document.body) new MutationObserver(th).observe(document.body, { attributes: true, attributeFilter: ['class', 'data-theme'] });
    if (window.matchMedia) { var mq = matchMedia('(prefers-color-scheme: dark)'); mq.addEventListener && mq.addEventListener('change', th); }
    wire(sr);
  }

  function init() {
    var slots = document.querySelectorAll('[data-apoiar]');
    for (var i = 0; i < slots.length; i++) { if (!slots[i].__dpa) { slots[i].__dpa = 1; slots[i].appendChild(makeBtn(slots[i].getAttribute('data-apoiar'))); } }
    var inl = document.querySelectorAll('[data-apoiar-inline]');
    for (var j = 0; j < inl.length; j++) { if (!inl[j].__dpa) { inl[j].__dpa = 1; renderInline(inl[j]); } }
    if (ds.float && !slots.length && !inl.length && !document.querySelector('.dp-apoiar-float')) {
      var f = makeBtn(); f.classList.add('dp-apoiar-float'); if (ds.float === 'right') { f.style.left = 'auto'; f.style.right = '16px'; }
      document.body.appendChild(f);
    }
  }
  document.addEventListener('click', function (e) {
    var a = e.target.closest && e.target.closest('[data-apoiar-open],a[href$="#apoiar"]');
    if (a) { e.preventDefault(); open(a); }
  });
  new MutationObserver(function () {
    lang(); var ls = document.querySelectorAll('.dp-apoiar-btn .dp-l');
    for (var k = 0; k < ls.length; k++) ls[k].textContent = T.btn;
  }).observe(document.documentElement, { attributes: true, attributeFilter: ['lang'] });
  window.DPApoiar = { open: open, close: close, init: init, setVariant: function (v) { CFG.variant = v === 'paid' ? 'paid' : 'free'; } };
  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', init); else init();
  if (/#apoiar$/.test(location.hash)) setTimeout(open, 300);
})();
