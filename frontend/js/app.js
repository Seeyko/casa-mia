(function() {
  'use strict';

  const API = window.APP_CONFIG?.API_URL || '';

  // ===== INIT =====
  document.addEventListener('DOMContentLoaded', function() {
    initNav();
    initFadeIn();

    // Page-specific init
    if (document.getElementById('statusBadges')) {
      loadStatus();
      loadNews();
      loadLocations();
      setInterval(loadStatus, 5 * 60 * 1000);
    }

    if (document.getElementById('menuContent')) {
      initMenuTabs();
      // loadMenu will be overridden by loadMenuWithCart if cart system is present
      setTimeout(function() { loadMenu('carte'); }, 10);
    }

    initPhonePopover();
  });

  // ===== NAVIGATION =====
  function initNav() {
    var hamburger = document.getElementById('hamburger');
    var nav = document.getElementById('nav');
    var overlay = document.getElementById('navOverlay');
    if (!hamburger) return;

    hamburger.addEventListener('click', function() {
      var isOpen = nav.classList.toggle('open');
      hamburger.classList.toggle('active');
      overlay.classList.toggle('active');
      hamburger.setAttribute('aria-expanded', isOpen);
      document.body.style.overflow = isOpen ? 'hidden' : '';
    });

    overlay.addEventListener('click', function() {
      nav.classList.remove('open');
      hamburger.classList.remove('active');
      overlay.classList.remove('active');
      hamburger.setAttribute('aria-expanded', 'false');
      document.body.style.overflow = '';
    });

    var header = document.getElementById('header');
    window.addEventListener('scroll', function() {
      header.classList.toggle('header--scrolled', window.scrollY > 20);
    }, { passive: true });
  }

  // ===== FADE IN =====
  var fadeObserver = new IntersectionObserver(function(entries) {
    entries.forEach(function(entry) {
      if (entry.isIntersecting) {
        var el = entry.target;
        el.classList.add('visible');
        fadeObserver.unobserve(el);
        // Remove fade-in class after transition to prevent any re-trigger
        setTimeout(function() { el.classList.remove('fade-in'); }, 700);
      }
    });
  }, { threshold: 0.1, rootMargin: '0px 0px -40px 0px' });

  function initFadeIn() {
    document.querySelectorAll('.fade-in:not(.visible)').forEach(function(el) {
      fadeObserver.observe(el);
    });
  }

  // ===== STATUS (open/closed) =====
  function loadStatus() {
    fetch(API + '/api/status')
      .then(function(r) { return r.json(); })
      .then(function(statuses) {
        var container = document.getElementById('statusBadges');
        if (!container) return;
        container.innerHTML = statuses.map(function(s) {
          var cls = s.is_open ? 'status-badge--open' : 'status-badge--closed';
          var label = s.is_open ? 'Ouvert' : 'Fermé';
          return '<div class="status-badge ' + cls + '">' +
            '<span class="status-badge__dot"></span>' +
            '<span class="status-badge__label"><strong>' + escapeHtml(s.name) + '</strong> — ' + label + '</span>' +
            '<span class="status-badge__info">' + escapeHtml(s.next_change) + '</span>' +
            '</div>';
        }).join('');
      })
      .catch(function() {});
  }

  // ===== NEWS =====
  function loadNews() {
    fetch(API + '/api/news')
      .then(function(r) { return r.json(); })
      .then(function(news) {
        var container = document.getElementById('newsGrid');
        var section = document.getElementById('newsSection');
        if (!container || !news.length) {
          if (section) section.style.display = 'none';
          return;
        }
        section.style.display = '';
        container.innerHTML = news.map(function(n) {
          var img = n.image_path ? '<img class="news-card__img" src="' + API + '/api/images/' + escapeHtml(n.image_path) + '" alt="" loading="lazy">' : '';
          var date = new Date(n.created_at).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' });
          var hasText = n.title || n.content;
          var body = hasText ? '<div class="news-card__body">' +
            (n.title ? '<h3 class="news-card__title">' + escapeHtml(n.title) + '</h3>' : '') +
            (n.content ? '<p class="news-card__text">' + escapeHtml(n.content) + '</p>' : '') +
            '<p class="news-card__date">' + date + '</p>' +
            '</div>' : '';
          return '<div class="news-card fade-in">' + img + body + '</div>';
        }).join('');
        initFadeIn();
      })
      .catch(function() {});
  }

  // ===== LOCATIONS =====
  function loadLocations() {
    fetch(API + '/api/locations')
      .then(function(r) { return r.json(); })
      .then(function(locations) {
        var container = document.getElementById('locationsGrid');
        if (!container) return;
        container.innerHTML = locations.map(function(loc) {
          var hours = formatHours(loc.opening_hours);
          return '<div class="location-card fade-in">' +
            '<h3 class="location-card__name">' + escapeHtml(loc.name) + '</h3>' +
            '<div class="location-card__detail"><span class="location-card__icon"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg></span><span>' + escapeHtml(loc.address) + '</span></div>' +
            '<div class="location-card__detail"><span class="location-card__icon"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg></span><span>' + hours + '</span></div>' +
            '<a href="tel:' + loc.phone.replace(/\s/g, '') + '" class="location-card__phone"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/></svg> ' + escapeHtml(loc.phone) + '</a>' +
            '</div>';
        }).join('');
        initFadeIn();
      })
      .catch(function() {});
  }

  function formatHours(hoursObj) {
    if (typeof hoursObj === 'string') {
      try { hoursObj = JSON.parse(hoursObj); } catch(e) { return 'Horaires indisponibles'; }
    }
    var days = ['lundi', 'mardi', 'mercredi', 'jeudi', 'vendredi', 'samedi', 'dimanche'];
    var lines = [];
    days.forEach(function(day) {
      var dh = hoursObj[day];
      var cap = day.charAt(0).toUpperCase() + day.slice(1);
      if (!dh || !dh.slots || !dh.slots.length) {
        lines.push(cap + ' : Fermé');
      } else {
        var times = dh.slots.map(function(s) { return s.open + '-' + s.close; }).join(' / ');
        lines.push(cap + ' : ' + times);
      }
    });
    return lines.join('<br>');
  }

  // ===== MENU =====
  function initMenuTabs() {
    document.querySelectorAll('.menu-tab').forEach(function(tab) {
      tab.addEventListener('click', function() {
        document.querySelectorAll('.menu-tab').forEach(function(t) { t.classList.remove('active'); });
        tab.classList.add('active');
        loadMenu(tab.dataset.section);
      });
    });
  }

  function loadMenu(section) {
    fetch(API + '/api/menu?section=' + section)
      .then(function(r) { return r.json(); })
      .then(function(data) {
        var container = document.getElementById('menuContent');
        if (!container) return;

        var html = '';

        if (section === 'carte') {
          html += '<div class="menu-notice">Toutes nos pizzas sont élaborées à base de produits frais en provenance d\'Italie, avec une mozzarella Fior di Latte, une pâte maison à haute hydratation et un temps de pousse de 48h minimum.</div>';
        }

        if (section === 'traiteur') {
          html += '<div class="traiteur-note"><p><strong>Service traiteur</strong> — Sur commande, minimum 3 jours à l\'avance.<br>Minimum 6 personnes. Livraison possible (renseignez-vous en magasin).</p></div>';
        }

        data.categories.forEach(function(cat) {
          html += '<div class="menu-category">';
          html += '<h2 class="menu-category__title">' + escapeHtml(cat.name) + '</h2>';
          html += '<hr class="decorative-line" style="margin-left:0;">';
          html += '<div class="menu-grid">';

          cat.items.forEach(function(item) {
            var highlightClass = item.badge === '★' ? ' menu-item--highlight' : '';
            html += '<div class="menu-item' + highlightClass + '">';

            if (item.image_path) {
              html += '<img class="menu-item__img" src="' + API + '/api/images/' + escapeHtml(item.image_path) + '" alt="" loading="lazy">';
            }

            html += '<div class="menu-item__header">';
            html += '<h3>' + escapeHtml(item.name);
            if (item.badge === 'NEW') html += '<span class="badge-new">NEW</span>';
            if (item.badge === '★') html += '<span class="badge-star"> ★</span>';
            html += '</h3>';
            if (item.price) html += '<span class="menu-item__price">' + escapeHtml(item.price) + '</span>';
            html += '</div>';

            if (item.description) html += '<p>' + escapeHtml(item.description) + '</p>';
            if (item.note) html += '<p class="menu-item__note">' + escapeHtml(item.note) + '</p>';

            html += '</div>';
          });

          html += '</div></div>';
        });

        if (section === 'carte') {
          html += '<div class="fidelity-banner"><p class="main">Carte de fidélité digitale — 10 pizzas achetées = 11ème offerte*</p><p class="sub">*sauf pizza du moment ou tartufo</p></div>';
        }

        container.innerHTML = html;
        initFadeIn();
      })
      .catch(function() {});
  }

  // ===== PHONE POPOVER =====
  function initPhonePopover() {
    var btn = document.getElementById('phoneBtn');
    if (!btn) return;

    btn.addEventListener('click', function(e) {
      e.preventDefault();
      var isMobile = window.innerWidth <= 768;

      if (isMobile) {
        var bs = document.getElementById('phoneBottomsheet');
        var ov = document.getElementById('phoneOverlay');
        if (bs && ov) {
          bs.classList.add('active');
          ov.classList.add('active');
          document.body.style.overflow = 'hidden';
        }
      } else {
        var pop = document.getElementById('phonePopover');
        var popOv = document.getElementById('phonePopoverOverlay');
        if (pop && popOv) {
          pop.classList.add('active');
          popOv.classList.add('active');
        }
      }
    });

    // Close bottomsheet
    var closeSheet = function() {
      var bs = document.getElementById('phoneBottomsheet');
      var ov = document.getElementById('phoneOverlay');
      if (bs) bs.classList.remove('active');
      if (ov) ov.classList.remove('active');
      document.body.style.overflow = '';
    };
    var phoneOv = document.getElementById('phoneOverlay');
    if (phoneOv) phoneOv.addEventListener('click', closeSheet);

    // Close desktop popover
    var closePopover = function() {
      var pop = document.getElementById('phonePopover');
      var popOv = document.getElementById('phonePopoverOverlay');
      if (pop) pop.classList.remove('active');
      if (popOv) popOv.classList.remove('active');
    };
    var popoverOv = document.getElementById('phonePopoverOverlay');
    if (popoverOv) popoverOv.addEventListener('click', closePopover);
  }

  // ===== HELPERS =====
  function escapeHtml(str) {
    if (!str) return '';
    var div = document.createElement('div');
    div.appendChild(document.createTextNode(str));
    return div.innerHTML;
  }

  // ===== CART SYSTEM =====
  function initCart() {
    var fab = document.getElementById('cartFab');
    var badge = document.getElementById('cartBadge');
    var drawerEl = document.getElementById('cartDrawer');
    var overlayEl = document.getElementById('cartOverlay');
    var closeBtn = document.getElementById('cartClose');
    var cartEmptyEl = document.getElementById('cartEmpty');
    var cartItemsEl = document.getElementById('cartItems');
    var cartSubtotalEl = document.getElementById('cartSubtotal');
    var cartTotalEl = document.getElementById('cartTotal');
    var cartNextBtn = document.getElementById('cartNextBtn');
    var cartBackBtn = document.getElementById('cartBackBtn');
    var cartStep1 = document.getElementById('cartStep1');
    var cartStep2 = document.getElementById('cartStep2');
    var cartNameEl = document.getElementById('cartName');
    var cartPhoneEl = document.getElementById('cartPhone');
    var cartTimeEl = document.getElementById('cartTime');
    var cartAddressEl = document.getElementById('cartAddress');

    if (!fab || !drawerEl) return;

    var STORAGE_KEY = 'casamia_cart';
    var cart = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]');

    function saveCart() { localStorage.setItem(STORAGE_KEY, JSON.stringify(cart)); }
    function getTotal() { return cart.reduce(function(s, i) { return s + i.price * i.qty; }, 0); }
    function getCount() { return cart.reduce(function(s, i) { return s + i.qty; }, 0); }
    function formatPrice(n) { return n.toFixed(2).replace('.', ',') + '\u20AC'; }

    function addToCart(name, price) {
      var existing = cart.find(function(i) { return i.name === name; });
      if (existing) { existing.qty++; } else { cart.push({ name: name, price: price, qty: 1 }); }
      saveCart(); renderBadge(); renderCartItems();
      fab.classList.remove('cart-fab--bounce');
      void fab.offsetWidth;
      fab.classList.add('cart-fab--bounce');
    }
    window.addToCart = addToCart;

    function renderBadge() {
      var count = getCount();
      badge.textContent = count;
      badge.style.display = count > 0 ? 'flex' : 'none';
      fab.style.display = 'flex';
    }

    function renderCartItems() {
      var hasItems = cart.length > 0;
      cartEmptyEl.style.display = hasItems ? 'none' : 'block';
      cartSubtotalEl.style.display = hasItems ? 'block' : 'none';
      cartItemsEl.innerHTML = '';
      cart.forEach(function(item, idx) {
        var row = document.createElement('div');
        row.className = 'cart-item-row';
        row.innerHTML =
          '<div class="cart-item-row__info">' +
            '<span class="cart-item-row__name">' + escapeHtml(item.name) + '</span>' +
            '<span class="cart-item-row__line-total">' + formatPrice(item.price * item.qty) + '</span>' +
          '</div>' +
          '<div class="cart-item-row__controls">' +
            '<button class="cart-qty-btn" data-idx="' + idx + '" data-action="minus">\u2212</button>' +
            '<span class="cart-qty-val">' + item.qty + '</span>' +
            '<button class="cart-qty-btn" data-idx="' + idx + '" data-action="plus">+</button>' +
          '</div>';
        cartItemsEl.appendChild(row);
      });
      cartItemsEl.querySelectorAll('.cart-qty-btn').forEach(function(btn) {
        btn.addEventListener('click', function() {
          var i = parseInt(this.getAttribute('data-idx'));
          if (this.getAttribute('data-action') === 'minus') {
            cart[i].qty--;
            if (cart[i].qty <= 0) cart.splice(i, 1);
          } else { cart[i].qty++; }
          saveCart(); renderBadge(); renderCartItems();
        });
      });
      var totalStr = formatPrice(getTotal());
      cartTotalEl.textContent = totalStr;
      var t2 = document.getElementById('cartTotal2');
      if (t2) t2.textContent = totalStr;
      updateSendLinks();
    }

    function getSelectedLocation() {
      var checked = document.querySelector('input[name="cartLocation"]:checked');
      return checked ? checked.value : 'entraigues';
    }

    function buildOrderText() {
      var lines = ['Commande CasaMia', '---'];
      cart.forEach(function(item) {
        lines.push(item.qty + 'x ' + item.name + ' - ' + formatPrice(item.price * item.qty));
      });
      lines.push('---');
      lines.push('Total: ' + formatPrice(getTotal()));
      if (cartNameEl.value.trim()) lines.push('Nom: ' + cartNameEl.value.trim());
      if (cartPhoneEl.value.trim()) lines.push('Tel: ' + cartPhoneEl.value.trim());
      var mode = document.querySelector('input[name="cartMode"]:checked').value;
      lines.push('Mode: ' + (mode === 'livraison' ? 'Livraison' : 'Retrait en boutique'));
      if (mode === 'livraison' && cartAddressEl.value.trim()) {
        lines.push('Adresse: ' + cartAddressEl.value.trim());
      }
      var loc = getSelectedLocation();
      lines.push('Point de retrait: ' + (loc === 'althen' ? 'Althen-des-Paluds' : 'Entraigues-sur-la-Sorgue'));
      if (cartTimeEl.value) lines.push('Heure souhaitée : ' + cartTimeEl.value.replace(':', 'h'));
      return lines.join('\n');
    }

    function updateSendLinks() {
      var text = buildOrderText();
      var encoded = encodeURIComponent(text);
      var loc = getSelectedLocation();
      var whatsappBtn = document.getElementById('btnWhatsApp');
      var smsBtn = document.getElementById('btnSMS');
      var callBtn = document.getElementById('btnCall');

      if (loc === 'entraigues') {
        whatsappBtn.href = 'https://wa.me/33645794930?text=' + encoded;
        smsBtn.href = 'sms:0645794930?body=' + encoded;
        callBtn.href = 'tel:0645794930';
        whatsappBtn.style.display = '';
        smsBtn.style.display = '';
      } else {
        callBtn.href = 'tel:0490361633';
        whatsappBtn.style.display = 'none';
        smsBtn.style.display = 'none';
      }
    }

    // Location change updates send buttons
    document.querySelectorAll('input[name="cartLocation"]').forEach(function(radio) {
      radio.addEventListener('change', updateSendLinks);
    });

    // Delivery mode toggle
    document.querySelectorAll('input[name="cartMode"]').forEach(function(radio) {
      radio.addEventListener('change', function() {
        cartAddressEl.style.display = this.value === 'livraison' ? 'block' : 'none';
        updateSendLinks();
      });
    });

    [cartNameEl, cartPhoneEl, cartTimeEl, cartAddressEl].forEach(function(el) {
      if (el) el.addEventListener('input', updateSendLinks);
    });

    function openDrawer() {
      renderCartItems();
      drawerEl.classList.add('cart-drawer--open');
      overlayEl.classList.add('cart-overlay--visible');
      document.body.style.overflow = 'hidden';
    }

    function closeDrawer() {
      drawerEl.classList.remove('cart-drawer--open');
      overlayEl.classList.remove('cart-overlay--visible');
      document.body.style.overflow = '';
    }

    fab.addEventListener('click', openDrawer);
    closeBtn.addEventListener('click', closeDrawer);
    overlayEl.addEventListener('click', closeDrawer);

    cartNextBtn.addEventListener('click', function() {
      cartStep1.style.display = 'none';
      cartStep2.style.display = 'flex';
      var recap = document.getElementById('cartOrderRecap');
      if (recap && cart.length) {
        var html = '<strong>Votre commande :</strong><br>';
        cart.forEach(function(item) {
          html += item.qty + 'x ' + item.name + ' \u2014 ' + formatPrice(item.price * item.qty) + '<br>';
        });
        recap.innerHTML = html;
      }
      var t2 = document.getElementById('cartTotal2');
      if (t2) t2.textContent = cartTotalEl.textContent;
      updateSendLinks();
    });

    cartBackBtn.addEventListener('click', function() {
      cartStep2.style.display = 'none';
      cartStep1.style.display = 'flex';
    });

    // Copy button
    document.getElementById('btnCopy').addEventListener('click', function() {
      var text = buildOrderText();
      navigator.clipboard.writeText(text).then(function() {
        var btn = document.getElementById('btnCopy');
        btn.innerHTML = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg> Copie !';
        setTimeout(function() {
          btn.innerHTML = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg> Copier';
        }, 2000);
      });
    });

    renderBadge();
  }

  // Inject add-to-cart buttons into menu items
  function injectCartButtons() {
    document.querySelectorAll('.menu-item').forEach(function(item) {
      if (item.querySelector('.cart-add-btn')) return;
      var headerEl = item.querySelector('.menu-item__header');
      if (!headerEl) return;
      var h3 = headerEl.querySelector('h3');
      var priceSpan = headerEl.querySelector('.menu-item__price');
      if (!h3 || !priceSpan) return;
      var name = h3.textContent.replace(/NEW/g, '').replace(/\u2605/g, '').trim();
      var priceText = priceSpan.textContent.trim();
      var price = parseFloat(priceText.replace('\u20AC', '').replace(',', '.'));
      if (isNaN(price)) return;
      var btn = document.createElement('button');
      btn.className = 'cart-add-btn';
      btn.textContent = '+';
      btn.setAttribute('aria-label', 'Ajouter ' + name);
      btn.addEventListener('click', function(e) {
        e.stopPropagation();
        if (window.addToCart) window.addToCart(name, price);
      });
      headerEl.appendChild(btn);
    });
  }

  // Override loadMenu to inject cart buttons after menu loads
  var origLoadMenu = loadMenu;
  function loadMenuWithCart(section) {
    fetch(API + '/api/menu?section=' + section)
      .then(function(r) { return r.json(); })
      .then(function(data) {
        var container = document.getElementById('menuContent');
        if (!container) return;

        var html = '';
        if (section === 'carte') {
          html += '<div class="menu-notice">Toutes nos pizzas sont élaborées à base de produits frais en provenance d\'Italie, avec une mozzarella Fior di Latte, une pâte maison à haute hydratation et un temps de pousse de 48h minimum.</div>';
        }
        if (section === 'traiteur') {
          html += '<div class="traiteur-note"><p><strong>Service traiteur</strong> — Sur commande, minimum 3 jours à l\'avance.<br>Minimum 6 personnes. Livraison possible (renseignez-vous en magasin).</p></div>';
        }

        data.categories.forEach(function(cat) {
          html += '<div class="menu-category">';
          html += '<h2 class="menu-category__title">' + escapeHtml(cat.name) + '</h2>';
          html += '<hr class="decorative-line" style="margin-left:0;">';
          html += '<div class="menu-grid">';
          cat.items.forEach(function(item) {
            var highlightClass = item.badge === '\u2605' ? ' menu-item--highlight' : '';
            html += '<div class="menu-item' + highlightClass + '">';
            if (item.image_path) {
              html += '<img class="menu-item__img" src="' + API + '/api/images/' + escapeHtml(item.image_path) + '" alt="" loading="lazy">';
            }
            html += '<div class="menu-item__body">';
            html += '<div class="menu-item__header">';
            html += '<h3>' + escapeHtml(item.name);
            if (item.badge === 'NEW') html += '<span class="badge-new">NEW</span>';
            if (item.badge === '\u2605') html += '<span class="badge-star"> \u2605</span>';
            html += '</h3>';
            html += '<span class="menu-item__dots"></span>';
            if (item.price) html += '<span class="menu-item__price">' + escapeHtml(item.price) + '</span>';
            html += '</div>';
            if (item.description) html += '<p>' + escapeHtml(item.description) + '</p>';
            if (item.note) html += '<p class="menu-item__note">' + escapeHtml(item.note) + '</p>';
            html += '</div></div>';
          });
          html += '</div></div>';
        });

        if (section === 'carte') {
          html += '<div class="fidelity-banner"><p class="main"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 9a3 3 0 0 1 0 6v2a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-2a3 3 0 0 1 0-6V7a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2Z"/><path d="M13 5v2"/><path d="M13 17v2"/><path d="M13 11v2"/></svg> Carte de fidélité digitale — 10 pizzas achetées = 11ème offerte*</p><p class="sub">*sauf pizza du moment ou tartufo</p></div>';
        }

        container.innerHTML = html;
        initFadeIn();

        // Only inject cart buttons on "carte" section (pizzas/snacking)
        if (section === 'carte') {
          setTimeout(injectCartButtons, 50);
        }
      })
      .catch(function() {});
  }

  // Replace loadMenu if on menu page
  if (document.getElementById('menuContent')) {
    loadMenu = loadMenuWithCart;
    initCart();
  }

  // Expose for inline use
  window.CasaMia = { loadMenu: loadMenu };
})();
