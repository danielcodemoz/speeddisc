/*! Apoiar — framework-free support modal (free or paid wording). No remote services. */
(function () {
  'use strict';
  if (window.__dpApoiar) return;
  window.__dpApoiar = 1;

  var script = document.currentScript;
  var ds = (script && script.dataset) || {};
  var CFG = {
    variant: ds.variant === 'paid' ? 'paid' : 'free',
    theme: ds.theme || 'auto',
    accent: ds.accent || '',
    openOnLoad: ds.open === '1' || ds.open === 'true',
    page: ds.page || 'https://danielpro.dev/apoiar/',
    name: ds.name || 'Daniel Marcos Domingos Manganhe',
    paypalEmail: ds.paypalEmail || 'danielmarcos173@gmail.com',
    paypal: ds.paypal || 'https://www.paypal.com/myaccount/transfer/homepage/pay',
    kofi: ds.kofi || 'https://ko-fi.com/danielpro',
    patreon: ds.patreon || 'https://www.patreon.com/cw/DanielMarcos958',
    gumroad: ds.gumroad || 'https://danielprodev.gumroad.com/',
    liberapay: ds.liberapay || 'https://liberapay.com/danielmarcos'
  };

  var TX = {
    en: {
      btn: 'Support',
      title: 'Enjoyed my work? Support the project',
      subFree: 'Any amount helps keep these projects free and updated.',
      subPaid: 'Your support helps develop new features and keep the project evolving.',
      notePaid: 'Voluntary support — it does not replace a plan subscription.',
      mpesaHint: 'Scan the QR with the M-Pesa app',
      copied: 'Copied!',
      paypal: 'PayPal',
      paypalHint: 'PayPal → Send → this email',
      copyMail: 'Copy email',
      openPP: 'Open PayPal',
      kofi: 'Ko-fi',
      kofiHint: 'Card or PayPal',
      patreon: 'Patreon',
      patreonHint: 'Monthly or one-time',
      gumroad: 'Gumroad',
      gumroadHint: 'Digital products',
      liberapay: 'Liberapay',
      liberapayHint: 'Weekly, monthly or yearly',
      close: 'Close',
      more: 'Support page',
      enlargeQr: 'Enlarge QR'
    },
    pt: {
      btn: 'Apoiar',
      title: 'Gostou do que fiz? Apoie o projeto',
      subFree: 'Qualquer valor ajuda a manter estes projetos gratuitos e atualizados.',
      subPaid: 'O seu apoio ajuda a desenvolver novas funcionalidades e a manter o projeto em evolução.',
      notePaid: 'Apoio voluntário — não substitui a subscrição de um plano.',
      mpesaHint: 'Leia o QR com a app M-Pesa',
      copied: 'Copiado!',
      paypal: 'PayPal',
      paypalHint: 'PayPal → Enviar → este email',
      copyMail: 'Copiar email',
      openPP: 'Abrir PayPal',
      kofi: 'Ko-fi',
      kofiHint: 'Cartão ou PayPal',
      patreon: 'Patreon',
      patreonHint: 'Mensal ou único',
      gumroad: 'Gumroad',
      gumroadHint: 'Produtos digitais',
      liberapay: 'Liberapay',
      liberapayHint: 'Semanal, mensal ou anual',
      close: 'Fechar',
      more: 'Página de apoio',
      enlargeQr: 'Ampliar QR'
    }
  };

  var HEART = '<svg viewBox="0 0 24 24" width="1em" height="1em" aria-hidden="true" fill="currentColor"><path d="M12 21s-7.5-4.6-10-9.3C.4 8.4 2.2 4.5 5.9 4.1 8 3.9 9.8 5 12 7.4 14.2 5 16 3.9 18.1 4.1c3.7.4 5.5 4.3 3.9 7.6C19.5 16.4 12 21 12 21z"/></svg>';

  var root = null;
  var pending = null;
  var lastFocus = null;
  var prevOverflow = '';
  var lightboxOn = false;

  function pack() {
    return /^en/i.test(document.documentElement.lang || '') ? TX.en : TX.pt;
  }

  function scriptBase() {
    if (!script || !script.src) return '';
    return script.src.replace(/[^/]*(\?.*)?$/, '');
  }

  function isDark() {
    if (CFG.theme === 'dark') return true;
    if (CFG.theme === 'light') return false;
    var els = [document.documentElement, document.body];
    for (var i = 0; i < els.length; i++) {
      var el = els[i];
      if (!el) continue;
      var cls = ' ' + (el.className && el.className.baseVal === undefined ? el.className : '') + ' ';
      if (/ (dark|dark-mode|dark-theme|theme-dark) /.test(cls)) return true;
      var attr = (el.getAttribute('data-theme') || el.getAttribute('data-bs-theme') || el.getAttribute('data-mode') || '').toLowerCase();
      if (attr === 'dark') return true;
      if (attr === 'light' || / (light|light-mode|light-theme|theme-light) /.test(cls)) return false;
    }
    if (CFG.theme === 'default-dark') return true;
    if (CFG.theme === 'default-light') return false;
    return !!(window.matchMedia && matchMedia('(prefers-color-scheme: dark)').matches);
  }

  function ensure() {
    if (root && root.isConnected) return Promise.resolve(root);
    var existing = document.getElementById('dp-support');
    if (existing) {
      root = existing;
      return Promise.resolve(root);
    }
    if (!pending) {
      pending = fetch(scriptBase() + 'partial.html').then(function (res) {
        if (!res.ok) throw new Error('Apoiar partial failed to load (' + res.status + ')');
        return res.text();
      }).then(function (html) {
        var holder = document.createElement('div');
        holder.innerHTML = html.trim();
        var el = holder.querySelector('#dp-support');
        if (!el) throw new Error('Apoiar partial is missing #dp-support');
        document.body.appendChild(el);
        root = el;
        return root;
      }).catch(function (err) {
        pending = null;
        throw err;
      });
    }
    return pending;
  }

  function applyTheme() {
    if (CFG.accent) {
      document.documentElement.style.setProperty('--dp-support-accent', CFG.accent);
    }
    if (!root) return;
    root.setAttribute('data-mode', isDark() ? 'dark' : 'light');
    if (CFG.accent) root.style.setProperty('--dp-support-accent', CFG.accent);
  }

  function applyCopy() {
    var t = pack();
    var sub = CFG.variant === 'paid' ? t.subPaid : t.subFree;
    var note = CFG.variant === 'paid' ? t.notePaid : '';
    var buttons = document.querySelectorAll('.dp-apoiar-btn');
    for (var b = 0; b < buttons.length; b++) {
      var label = buttons[b].querySelector('.dp-l');
      if (label) label.textContent = t.btn;
      buttons[b].setAttribute('aria-label', t.btn);
    }
    if (!root) return;
    root.setAttribute('data-variant', CFG.variant);
    var nodes = root.querySelectorAll('[data-dp-text]');
    for (var i = 0; i < nodes.length; i++) {
      var key = nodes[i].getAttribute('data-dp-text');
      var value = key === 'sub' ? sub : key === 'note' ? note : t[key];
      if (value != null) nodes[i].textContent = value;
    }
    var noteEl = root.querySelector('#dp-support-note');
    if (noteEl) noteEl.hidden = CFG.variant !== 'paid';
    var panel = root.querySelector('.dp-support-panel');
    if (panel) {
      panel.setAttribute('aria-describedby', CFG.variant === 'paid' ? 'dp-support-sub dp-support-note' : 'dp-support-sub');
    }
    var arias = root.querySelectorAll('[data-dp-aria]');
    for (var a = 0; a < arias.length; a++) {
      var ariaKey = arias[a].getAttribute('data-dp-aria');
      if (t[ariaKey]) arias[a].setAttribute('aria-label', t[ariaKey]);
    }
    var nameEl = root.querySelector('[data-dp-name]');
    if (nameEl) nameEl.textContent = CFG.name;
    var mail = root.querySelector('[data-dp-email]');
    if (mail) mail.textContent = CFG.paypalEmail;
    var links = root.querySelectorAll('[data-dp-link]');
    for (var n = 0; n < links.length; n++) {
      var linkKey = links[n].getAttribute('data-dp-link');
      if (CFG[linkKey]) links[n].setAttribute('href', CFG[linkKey]);
    }
    var more = root.querySelector('.dp-support-more');
    if (more) {
      var here = location.hostname.replace(/^www\./, '') + location.pathname.replace(/\/$/, '');
      more.hidden = here === 'danielpro.dev/apoiar';
    }
  }

  function visibleFocusable(scope) {
    var list = scope.querySelectorAll('a[href], button, [tabindex]');
    var out = [];
    for (var i = 0; i < list.length; i++) {
      var el = list[i];
      if (el.disabled || el.getAttribute('aria-hidden') === 'true') continue;
      if (el.closest('[hidden]')) continue;
      if (el.tabIndex < 0) continue;
      out.push(el);
    }
    return out;
  }

  function onKey(e) {
    if (!root || root.hidden) return;
    if (e.key === 'Escape') {
      e.preventDefault();
      e.stopPropagation();
      if (lightboxOn) closeLightbox();
      else close();
      return;
    }
    if (e.key !== 'Tab') return;
    var lightbox = root.querySelector('[data-dp-lightbox]');
    var scope = lightboxOn && lightbox ? lightbox : root.querySelector('.dp-support-panel');
    if (!scope) return;
    var items = visibleFocusable(scope);
    if (!items.length) {
      e.preventDefault();
      return;
    }
    var first = items[0];
    var last = items[items.length - 1];
    var active = document.activeElement;
    var inside = scope.contains(active);
    if (e.shiftKey) {
      if (!inside || active === first || active === scope) {
        e.preventDefault();
        last.focus();
      }
    } else if (!inside || active === last) {
      e.preventDefault();
      first.focus();
    }
  }

  function openLightbox() {
    if (!root || lightboxOn) return;
    var lightbox = root.querySelector('[data-dp-lightbox]');
    var slot = root.querySelector('[data-dp-lightbox-qr]');
    var qr = root.querySelector('[data-dp-qr] svg');
    if (!lightbox || !slot || !qr) return;
    slot.replaceChildren(qr.cloneNode(true));
    lightbox.hidden = false;
    lightboxOn = true;
    requestAnimationFrame(function () { lightbox.classList.add('is-on'); });
    var closer = lightbox.querySelector('[data-dp-lightbox-close]');
    if (closer) setTimeout(function () { closer.focus({ preventScroll: true }); }, 30);
  }

  function closeLightbox() {
    if (!root || !lightboxOn) return;
    var lightbox = root.querySelector('[data-dp-lightbox]');
    if (lightbox) {
      lightbox.classList.remove('is-on');
      lightbox.hidden = true;
    }
    lightboxOn = false;
    var qr = root.querySelector('[data-dp-qr]');
    if (qr) qr.focus({ preventScroll: true });
  }

  function copyText(text, btn) {
    var t = pack();
    function done() {
      btn.textContent = t.copied;
      btn.classList.add('is-ok');
      setTimeout(function () {
        btn.textContent = t.copyMail;
        btn.classList.remove('is-ok');
      }, 1800);
    }
    function fallback() {
      var area = document.createElement('textarea');
      area.value = text;
      area.setAttribute('readonly', '');
      area.style.cssText = 'position:fixed;left:-9999px;top:0';
      document.body.appendChild(area);
      area.select();
      try { document.execCommand('copy'); done(); } catch (err) {}
      document.body.removeChild(area);
    }
    if (navigator.clipboard && window.isSecureContext) navigator.clipboard.writeText(text).then(done, fallback);
    else fallback();
  }

  function wire() {
    if (!root || root.__dpWired) return;
    root.__dpWired = 1;
    root.querySelector('[data-dp-overlay]').addEventListener('click', function (e) {
      if (e.target === e.currentTarget) close();
    });
    root.querySelector('[data-dp-close]').addEventListener('click', function () { close(); });
    root.querySelector('[data-dp-qr]').addEventListener('click', function () { openLightbox(); });
    var lightbox = root.querySelector('[data-dp-lightbox]');
    lightbox.addEventListener('click', function (e) {
      if (e.target === lightbox) closeLightbox();
    });
    lightbox.querySelector('[data-dp-lightbox-close]').addEventListener('click', function () { closeLightbox(); });
    var copyBtn = root.querySelector('[data-dp-copy]');
    copyBtn.addEventListener('click', function () { copyText(CFG.paypalEmail, copyBtn); });
  }

  function open(from) {
    lastFocus = from || document.activeElement;
    ensure().then(function () {
      applyTheme();
      applyCopy();
      wire();
      if (!root.hidden && root.classList.contains('is-open')) return;
      root.hidden = false;
      prevOverflow = document.documentElement.style.overflow;
      document.documentElement.style.overflow = 'hidden';
      document.addEventListener('keydown', onKey, true);
      requestAnimationFrame(function () { root.classList.add('is-open'); });
      var panel = root.querySelector('.dp-support-panel');
      setTimeout(function () { if (panel) panel.focus({ preventScroll: true }); }, 30);
    }).catch(function (err) {
      console.error(err);
    });
  }

  function close() {
    if (!root || root.hidden) return;
    closeLightbox();
    root.classList.remove('is-open');
    root.hidden = true;
    document.documentElement.style.overflow = prevOverflow;
    document.removeEventListener('keydown', onKey, true);
    if (lastFocus && lastFocus.focus) lastFocus.focus({ preventScroll: true });
  }

  function setVariant(variant) {
    CFG.variant = variant === 'paid' ? 'paid' : 'free';
    applyCopy();
  }

  function makeBtn(variant) {
    var t = pack();
    var btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'dp-apoiar-btn' + (variant === 'link' ? ' dp-link' : variant === 'nav' ? ' dp-nav' : '');
    btn.setAttribute('aria-label', t.btn);
    btn.setAttribute('aria-haspopup', 'dialog');
    btn.innerHTML = HEART + '<span class="dp-l">' + t.btn + '</span>';
    btn.addEventListener('click', function (e) {
      e.preventDefault();
      open(btn);
    });
    return btn;
  }

  function init() {
    var slots = document.querySelectorAll('[data-apoiar], [data-dp-support]');
    for (var i = 0; i < slots.length; i++) {
      if (slots[i].__dpa) continue;
      slots[i].__dpa = 1;
      var kind = slots[i].getAttribute('data-dp-support');
      if (kind == null) kind = slots[i].getAttribute('data-apoiar');
      slots[i].appendChild(makeBtn(kind || ''));
    }
    var inline = document.querySelectorAll('[data-apoiar-inline], [data-dp-support-inline]');
    if (!slots.length && !inline.length && ds.float && !document.querySelector('.dp-apoiar-float')) {
      var floating = makeBtn('');
      floating.classList.add('dp-apoiar-float');
      if (ds.float === 'right') {
        floating.style.left = 'auto';
        floating.style.right = '16px';
      }
      document.body.appendChild(floating);
    }
    applyTheme();
    applyCopy();
  }

  document.addEventListener('click', function (e) {
    var opener = e.target.closest && e.target.closest('[data-apoiar-open], [data-dp-support-open], a[href$="#apoiar"]');
    if (!opener || opener.classList.contains('dp-apoiar-btn')) return;
    e.preventDefault();
    open(opener);
  });

  function watch(el) {
    if (!el) return;
    new MutationObserver(function (records) {
      for (var i = 0; i < records.length; i++) {
        if (records[i].attributeName === 'lang') applyCopy();
        else applyTheme();
      }
    }).observe(el, {
      attributes: true,
      attributeFilter: ['class', 'data-theme', 'data-bs-theme', 'data-mode', 'lang']
    });
  }
  watch(document.documentElement);
  if (document.body) watch(document.body);
  else document.addEventListener('DOMContentLoaded', function () { watch(document.body); });
  if (window.matchMedia) {
    var mq = matchMedia('(prefers-color-scheme: dark)');
    if (mq.addEventListener) mq.addEventListener('change', applyTheme);
  }

  var api = { open: open, close: close, init: init, setVariant: setVariant };
  window.SupportModal = api;
  window.DPApoiar = api;

  function boot() {
    init();
    if (CFG.openOnLoad) open();
    else if (/#apoiar$/.test(location.hash)) setTimeout(open, 300);
  }
  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', boot);
  else boot();
})();
