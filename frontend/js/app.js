(function() {
  'use strict';

  const API = window.APP_CONFIG?.API_URL || '';
  var menuItemsCache = {};

  // ===== INIT =====
  document.addEventListener('DOMContentLoaded', function() {
    document.documentElement.classList.add('js-ready');
    initNav();
    initFadeIn();
    // Fallback : 1.5s après load, dévoile tous les reveal dans la partie supérieure,
    // et 3s après (ou au premier scroll) dévoile tout pour éviter contenu caché en cas de bug IO.
    setTimeout(function() {
      document.querySelectorAll('.reveal:not(.is-visible)').forEach(function(el) {
        var r = el.getBoundingClientRect();
        if (r.top < window.innerHeight * 1.5) el.classList.add('is-visible');
      });
    }, 1500);
    setTimeout(function() {
      document.querySelectorAll('.reveal:not(.is-visible)').forEach(function(el) {
        el.classList.add('is-visible');
      });
    }, 3500);

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

  // ===== SCROLL REVEAL =====
  var revealObserver = new IntersectionObserver(function(entries) {
    entries.forEach(function(entry) {
      if (entry.isIntersecting) {
        entry.target.classList.add('is-visible');
        revealObserver.unobserve(entry.target);
      }
    });
  }, { threshold: 0.12, rootMargin: '0px 0px -60px 0px' });

  function initFadeIn() {
    document.querySelectorAll('.reveal:not(.is-visible)').forEach(function(el) {
      revealObserver.observe(el);
    });
  }

  // ===== STATUS (open/closed) =====
  function loadStatus() {
    fetch(API + '/api/status')
      .then(function(r) { return r.json(); })
      .then(function(statuses) {
        var container = document.getElementById('statusBadges');
        if (!container) return;
        container.innerHTML = statuses.map(function(s, i) {
          var cls = s.is_open ? 'status--open' : 'status--closed';
          var label = s.is_open ? 'Ouvert' : 'Fermé';
          var sep = i > 0 ? '<span class="hero__meta-sep" aria-hidden="true">·</span>' : '';
          return sep +
            '<span class="status ' + cls + '">' +
              '<span class="status__dot" aria-hidden="true"></span>' +
              '<span class="status__label">' + escapeHtml(s.name) + '</span>' +
              '<span class="status__note">— ' + label +
                (s.next_change ? ' · ' + escapeHtml(s.next_change) : '') +
              '</span>' +
            '</span>';
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
          var img = n.image_path
            ? '<figure class="news__figure"><img src="' + API + '/api/images/' + escapeHtml(n.image_path) + '" alt="" loading="lazy"></figure>'
            : '<figure class="news__figure" aria-hidden="true"></figure>';
          var hasText = n.title || n.content;
          var body = hasText ? '<div class="news__body">' +
            (n.title ? '<h3 class="news__title">' + escapeHtml(n.title) + '</h3>' : '') +
            (n.content ? '<p class="news__text">' + escapeHtml(n.content) + '</p>' : '') +
            '</div>' : '';
          return '<article class="news__item">' + img + body + '</article>';
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
        var todayIdx = (new Date().getDay() + 6) % 7;
        container.innerHTML = locations.map(function(loc) {
          return '<div class="addr reveal">' +
            '<h3 class="addr__name">' + escapeHtml(loc.name) + '</h3>' +
            '<p class="addr__street">' + escapeHtml(loc.address) + '</p>' +
            '<a href="tel:' + loc.phone.replace(/\s/g, '') + '" class="addr__phone">' + escapeHtml(loc.phone) + '</a>' +
            renderHours(loc.opening_hours, todayIdx) +
            '</div>';
        }).join('');
        initFadeIn();
      })
      .catch(function() {});
  }

  function renderHours(hoursObj, todayIdx) {
    if (typeof hoursObj === 'string') {
      try { hoursObj = JSON.parse(hoursObj); } catch(e) { return ''; }
    }
    var days = ['lundi', 'mardi', 'mercredi', 'jeudi', 'vendredi', 'samedi', 'dimanche'];
    var labels = ['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim'];
    var rows = days.map(function(day, i) {
      var dh = hoursObj[day];
      var todayCls = i === todayIdx ? ' hours__row--today' : '';
      var closed = !dh || !dh.slots || !dh.slots.length;
      var valueCls = closed ? ' hours__value--closed' : '';
      var value = closed
        ? 'Fermé'
        : dh.slots.map(function(s) { return s.open + ' – ' + s.close; }).join(' · ');
      return '<div class="hours__day' + todayCls + '">' + labels[i] + '</div>' +
             '<div class="hours__value' + valueCls + todayCls + '">' + value + '</div>';
    }).join('');
    return '<div class="hours">' + rows + '</div>';
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
          html += '<p class="menu-notice">Pizzas étalées à la main, pâte à haute hydratation et longue fermentation de 48 h, mozzarella Fior di Latte, ingrédients importés d\'Italie (Naples, Pouilles, …). Pâte disponible sans gluten sur demande.</p>';
        }

        if (section === 'traiteur') {
          html += '<p class="menu-notice"><strong>Service traiteur</strong> — sur commande, 3 jours à l\'avance, à partir de 6 personnes. Livraison possible, à demander en magasin.</p>';
        }

        data.categories.forEach(function(cat, catIdx) {
          var idxStr = String(catIdx + 1).padStart(2, '0');
          var isTraiteur = section === 'traiteur';
          html += '<section class="' + (isTraiteur ? 'traiteur-group' : 'menu-category') + '">';
          html += '<header class="' + (isTraiteur ? 'traiteur-group__head' : 'menu-category__head') + '">';
          html += '<h2 class="' + (isTraiteur ? 'traiteur-group__name' : 'menu-category__name') + '">' + escapeHtml(cat.name) + '</h2>';
          html += '<span class="' + (isTraiteur ? 'traiteur-group__count' : 'menu-category__index') + '">' + idxStr + ' — ' + cat.items.length + (cat.items.length > 1 ? ' créations' : ' création') + '</span>';
          html += '</header>';

          if (isTraiteur) {
            html += '<div class="traiteur-grid">';
            cat.items.forEach(function(item) {
              html += '<article class="traiteur-item">';
              if (item.image_path) {
                html += '<div class="traiteur-item__fig"><img src="' + API + '/api/images/' + escapeHtml(item.image_path) + '" alt="" loading="lazy"></div>';
              } else {
                html += '<div class="traiteur-item__fig traiteur-item__fig--ph" aria-hidden="true">' + itemIcon(cat.name, item.name) + '</div>';
              }
              html += '<div class="traiteur-item__row">';
              html += '<h3 class="traiteur-item__name">' + escapeHtml(item.name);
              if (item.badge === 'NEW') html += ' <span class="menu-badge">Nouveau</span>';
              if (item.badge === '★') html += ' <span class="menu-badge menu-badge--star">★</span>';
              html += '</h3>';
              if (item.price) html += '<span class="traiteur-item__price price">' + escapeHtml(item.price) + '</span>';
              html += '</div>';
              if (item.description) html += '<p class="traiteur-item__desc">' + escapeHtml(item.description) + '</p>';
              html += '</article>';
            });
            html += '</div>';
          } else {
            html += '<div class="menu-list">';
            cat.items.forEach(function(item, itemIdx) {
              var mid = 'mi-' + catIdx + '-' + itemIdx;
              menuItemsCache[mid] = item;
              html += '<article class="menu-item" data-mid="' + mid + '">';
              if (item.image_path) {
                html += '<img class="menu-item__img" src="' + API + '/api/images/' + escapeHtml(item.image_path) + '" alt="" loading="lazy">';
              } else {
                html += '<span class="menu-item__img menu-item__img--ph" aria-hidden="true">' + itemIcon(cat.name, item.name) + '</span>';
              }
              html += '<div class="menu-item__body">';
              html += '<h3 class="menu-item__name">' + escapeHtml(item.name);
              if (item.badge === 'NEW') html += ' <span class="menu-badge">Nouveau</span>';
              if (item.badge === '★') html += ' <span class="menu-badge menu-badge--star">★</span>';
              html += '</h3>';
              if (item.description) html += '<p class="menu-item__desc">' + escapeHtml(item.description) + '</p>';
              if (item.note) html += '<p class="menu-item__desc"><em>' + escapeHtml(item.note) + '</em></p>';
              html += '<button type="button" class="menu-item__more" data-mid="' + mid + '" aria-label="Voir plus de détails sur ' + escapeHtml(item.name) + '">Voir plus <span aria-hidden="true">\u2192</span></button>';
              html += '</div>';
              if (item.price) html += '<span class="menu-item__price price">' + escapeHtml(item.price) + '</span>';
              html += '</article>';
            });
            html += '</div>';
          }

          html += '</section>';
        });

        if (section === 'carte') {
          html += '<aside class="menu-notice" style="margin-top:var(--s-xl);text-align:center;font-style:normal;"><strong style="color:var(--gold);">Carte de fidélité</strong> — 10 pizzas achetées, la 11<sup>e</sup> offerte. <span style="color:var(--ink-3);">(Hors pizza du moment &amp; Tartufo)</span></aside>';
        }

        container.innerHTML = html;
        initFadeIn();
      })
      .catch(function() {});
  }

  // ===== PHONE SHEET (unifié popover + bottomsheet via CSS) =====
  function initPhonePopover() {
    var btn = document.getElementById('phoneBtn');
    var sheet = document.getElementById('phoneSheet');
    var backdrop = document.getElementById('phoneBackdrop');
    var closeBtn = document.getElementById('phoneSheetClose');
    if (!btn || !sheet || !backdrop) return;

    function open() {
      sheet.setAttribute('aria-hidden', 'false');
      backdrop.setAttribute('aria-hidden', 'false');
      document.body.style.overflow = 'hidden';
    }
    function close() {
      sheet.setAttribute('aria-hidden', 'true');
      backdrop.setAttribute('aria-hidden', 'true');
      document.body.style.overflow = '';
    }
    btn.addEventListener('click', function(e) { e.preventDefault(); open(); });
    backdrop.addEventListener('click', close);
    if (closeBtn) closeBtn.addEventListener('click', close);
    document.addEventListener('keydown', function(e) {
      if (e.key === 'Escape' && sheet.getAttribute('aria-hidden') === 'false') close();
    });
  }

  // ===== HELPERS =====
  function escapeHtml(str) {
    if (!str) return '';
    var div = document.createElement('div');
    div.appendChild(document.createTextNode(str));
    return div.innerHTML;
  }

  // Icônes de fallback pour items sans photo — détection simple par mots-clés
  var ICONS = {
    pizza:   '<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="9"/><circle cx="8.5" cy="10" r="1" fill="currentColor" stroke="none"/><circle cx="14" cy="9.5" r=".8" fill="currentColor" stroke="none"/><circle cx="12" cy="14" r="1" fill="currentColor" stroke="none"/><circle cx="15.5" cy="14" r=".7" fill="currentColor" stroke="none"/></svg>',
    glass:   '<svg viewBox="0 0 24 24"><path d="M7 3h10l-1 9c0 2-2 3.5-4 3.5S8 14 8 12l-1-9z"/><path d="M12 15.5V21"/><path d="M9 21h6"/></svg>',
    dessert: '<svg viewBox="0 0 24 24"><path d="M4 20h16"/><path d="M5 20l1.5-8h11L19 20"/><path d="M8.5 12c0-2 1.5-3.5 3.5-3.5s3.5 1.5 3.5 3.5"/><path d="M12 6v2.5"/></svg>',
    platter: '<svg viewBox="0 0 24 24"><ellipse cx="12" cy="17" rx="9" ry="1.8"/><circle cx="8" cy="13" r="2"/><circle cx="13" cy="12" r="2.5"/><circle cx="17" cy="13.5" r="1.6"/></svg>',
    skewer:  '<svg viewBox="0 0 24 24"><line x1="2" y1="12" x2="22" y2="12"/><circle cx="8" cy="12" r="2.3"/><circle cx="13" cy="12" r="2.3"/><circle cx="18" cy="12" r="2.3"/></svg>',
    bread:   '<svg viewBox="0 0 24 24"><path d="M3 12c0-3 3-5 9-5s9 2 9 5-3 5-9 5-9-2-9-5z"/><path d="M7 12c1.5.8 3 1 5 1s3.5-.2 5-1"/></svg>',
    fork:    '<svg viewBox="0 0 24 24"><path d="M8 3v7a2 2 0 0 0 4 0V3"/><path d="M10 10v11"/><path d="M16 3c-1 3-1 6 0 8v10"/></svg>'
  };
  function itemIcon(catName, itemName) {
    var t = ((catName || '') + ' ' + (itemName || '')).toLowerCase();
    if (/pizza|calzone|focaccia/.test(t)) return ICONS.pizza;
    if (/vin|bi[èe]re|boisson|limonade|soda|cola|eau|jus|spritz|campari|aperol/.test(t)) return ICONS.glass;
    if (/dolci|tiramisu|cannoli|panna|dessert|g[âa]teau|chocolat|fondant|glace|pistache/.test(t)) return ICONS.dessert;
    if (/planche|charcut|fromage|c[œoe]ur|ap[ée]ro|bouch[ée]e|feuillet[ée]|navette|verrine/.test(t)) return ICONS.platter;
    if (/brochette/.test(t)) return ICONS.skewer;
    if (/pain|bruschetta|panini|sandwich/.test(t)) return ICONS.bread;
    return ICONS.fork;
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

  // ===== ITEM DETAIL SHEET (mobile bottom sheet / desktop modal) =====
  function initItemSheet() {
    var sheet = document.getElementById('itemSheet');
    var backdrop = document.getElementById('itemSheetBackdrop');
    var closeBtn = document.getElementById('itemSheetClose');
    var nameEl = document.getElementById('itemSheetName');
    var priceEl = document.getElementById('itemSheetPrice');
    var descEl = document.getElementById('itemSheetDesc');
    var noteEl = document.getElementById('itemSheetNote');
    var mediaEl = document.getElementById('itemSheetMedia');
    var addBtn = document.getElementById('itemSheetAdd');
    if (!sheet || !backdrop) return;

    var currentItem = null;

    function open(item) {
      if (!item) return;
      currentItem = item;
      nameEl.innerHTML = escapeHtml(item.name) +
        (item.badge === 'NEW' ? ' <span class="menu-badge">Nouveau</span>' : '') +
        (item.badge === '\u2605' ? ' <span class="menu-badge menu-badge--star">\u2605</span>' : '');
      if (item.price) {
        priceEl.textContent = item.price;
        priceEl.style.display = '';
      } else {
        priceEl.style.display = 'none';
      }
      if (item.description) {
        descEl.textContent = item.description;
        descEl.style.display = '';
      } else {
        descEl.style.display = 'none';
      }
      if (item.note) {
        noteEl.innerHTML = '<em>' + escapeHtml(item.note) + '</em>';
        noteEl.style.display = '';
      } else {
        noteEl.style.display = 'none';
      }
      if (item.image_path) {
        mediaEl.innerHTML = '<img src="' + API + '/api/images/' + escapeHtml(item.image_path) + '" alt="' + escapeHtml(item.name) + '">';
        mediaEl.classList.remove('item-sheet__media--ph');
      } else {
        mediaEl.innerHTML = itemIcon('', item.name);
        mediaEl.classList.add('item-sheet__media--ph');
      }
      var price = item.price ? parseFloat(String(item.price).replace(/[^\d,\.]/g, '').replace(',', '.')) : NaN;
      addBtn.style.display = isNaN(price) ? 'none' : '';
      sheet.setAttribute('aria-hidden', 'false');
      backdrop.setAttribute('aria-hidden', 'false');
      document.body.style.overflow = 'hidden';
    }
    function close() {
      sheet.setAttribute('aria-hidden', 'true');
      backdrop.setAttribute('aria-hidden', 'true');
      document.body.style.overflow = '';
      currentItem = null;
    }

    backdrop.addEventListener('click', close);
    closeBtn.addEventListener('click', close);
    document.addEventListener('keydown', function(e) {
      if (e.key === 'Escape' && sheet.getAttribute('aria-hidden') === 'false') close();
    });
    addBtn.addEventListener('click', function() {
      if (!currentItem || !currentItem.price) return;
      var price = parseFloat(String(currentItem.price).replace(/[^\d,\.]/g, '').replace(',', '.'));
      if (isNaN(price)) return;
      if (window.addToCart) window.addToCart(currentItem.name, price);
      close();
    });

    // Capture phase so that the click on "Voir plus" does not also trigger
    // the add-to-cart listener attached on the parent .menu-item.
    document.addEventListener('click', function(e) {
      var btn = e.target.closest('.menu-item__more');
      if (!btn) return;
      e.preventDefault();
      e.stopPropagation();
      var mid = btn.getAttribute('data-mid');
      open(menuItemsCache[mid]);
    }, true);
    document.addEventListener('keydown', function(e) {
      if (e.key !== 'Enter' && e.key !== ' ') return;
      var btn = e.target.closest && e.target.closest('.menu-item__more');
      if (!btn) return;
      e.preventDefault();
      e.stopPropagation();
      var mid = btn.getAttribute('data-mid');
      open(menuItemsCache[mid]);
    }, true);
  }

  // Rend les items "carte" cliquables pour ajout panier
  function injectCartButtons() {
    document.querySelectorAll('.menu-item').forEach(function(item) {
      if (item.classList.contains('menu-item--clickable')) return;
      var nameEl = item.querySelector('.menu-item__name');
      var priceEl = item.querySelector('.menu-item__price');
      if (!nameEl || !priceEl) return;
      var name = nameEl.textContent.replace(/Nouveau/gi, '').replace(/\u2605/g, '').trim();
      var priceText = priceEl.textContent.trim();
      var price = parseFloat(priceText.replace(/[^\d,\.]/g, '').replace(',', '.'));
      if (isNaN(price)) return;

      var btn = document.createElement('span');
      btn.className = 'menu-item__add';
      btn.textContent = '+';
      btn.setAttribute('aria-hidden', 'true');
      item.appendChild(btn);

      item.classList.add('menu-item--clickable');
      item.setAttribute('role', 'button');
      item.setAttribute('tabindex', '0');
      item.setAttribute('aria-label', 'Ajouter ' + name + ' au panier');
      item.addEventListener('click', function() {
        if (window.addToCart) window.addToCart(name, price);
      });
      item.addEventListener('keydown', function(e) {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          if (window.addToCart) window.addToCart(name, price);
        }
      });
    });
  }

  // Override loadMenu pour ajouter les boutons panier après chargement
  var _origLoadMenu = loadMenu;
  function loadMenuWithCart(section) {
    _origLoadMenu(section);
    if (section === 'carte') {
      setTimeout(injectCartButtons, 80);
    }
  }

  if (document.getElementById('menuContent')) {
    loadMenu = loadMenuWithCart;
    initCart();
    initItemSheet();
  }

  // Expose for inline use
  window.CasaMia = { loadMenu: loadMenu };
})();
