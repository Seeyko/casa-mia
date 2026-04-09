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
      loadMenu('carte');
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
  function initFadeIn() {
    var els = document.querySelectorAll('.fade-in');
    var observer = new IntersectionObserver(function(entries) {
      entries.forEach(function(entry) {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
          observer.unobserve(entry.target);
        }
      });
    }, { threshold: 0.1, rootMargin: '0px 0px -40px 0px' });
    els.forEach(function(el) { observer.observe(el); });
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
          var label = s.is_open ? 'Ouvert' : 'Ferme';
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
          return '<div class="news-card fade-in">' +
            img +
            '<div class="news-card__body">' +
            '<h3 class="news-card__title">' + escapeHtml(n.title) + '</h3>' +
            '<p class="news-card__text">' + escapeHtml(n.content) + '</p>' +
            '<p class="news-card__date">' + date + '</p>' +
            '</div></div>';
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
            '<div class="location-card__detail"><span class="location-card__icon">&#9906;</span><span>' + escapeHtml(loc.address) + '</span></div>' +
            '<div class="location-card__detail"><span class="location-card__icon">&#128337;</span><span>' + hours + '</span></div>' +
            '<a href="tel:' + loc.phone.replace(/\s/g, '') + '" class="location-card__phone">' + escapeHtml(loc.phone) + '</a>' +
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
        lines.push(cap + ' : Ferme');
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
          html += '<div class="menu-notice">Toutes nos pizzas sont elaborees a base de produits frais en provenance d\'Italie, avec une mozzarella Fior di Latte, une pate maison a haute hydratation et un temps de pousse de 48h minimum.</div>';
        }

        if (section === 'traiteur') {
          html += '<div class="traiteur-note"><p><strong>Service traiteur</strong> — Sur commande, minimum 3 jours a l\'avance.<br>Minimum 6 personnes. Livraison possible (renseignez-vous en magasin).</p></div>';
        }

        data.categories.forEach(function(cat) {
          html += '<div class="menu-category fade-in">';
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
          html += '<div class="fidelity-banner"><p class="main">Carte de fidelite digitale — 10 pizzas achetees = 11eme offerte*</p><p class="sub">*sauf pizza du moment ou tartufo</p></div>';
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

    var isMobile = window.innerWidth <= 768;

    if (isMobile) {
      btn.addEventListener('click', function(e) {
        e.preventDefault();
        var bs = document.getElementById('phoneBottomsheet');
        var ov = document.getElementById('phoneOverlay');
        if (bs && ov) {
          bs.classList.add('active');
          ov.classList.add('active');
          document.body.style.overflow = 'hidden';
        }
      });

      var closeSheet = function() {
        var bs = document.getElementById('phoneBottomsheet');
        var ov = document.getElementById('phoneOverlay');
        if (bs && ov) {
          bs.classList.remove('active');
          ov.classList.remove('active');
          document.body.style.overflow = '';
        }
      };

      var ov = document.getElementById('phoneOverlay');
      if (ov) ov.addEventListener('click', closeSheet);
    } else {
      btn.addEventListener('click', function(e) {
        e.preventDefault();
        var pop = document.getElementById('phonePopover');
        if (pop) pop.classList.toggle('active');
      });

      document.addEventListener('click', function(e) {
        var pop = document.getElementById('phonePopover');
        if (pop && !e.target.closest('.phone-popover-wrapper')) {
          pop.classList.remove('active');
        }
      });
    }
  }

  // ===== HELPERS =====
  function escapeHtml(str) {
    if (!str) return '';
    var div = document.createElement('div');
    div.appendChild(document.createTextNode(str));
    return div.innerHTML;
  }

  // Expose for inline use
  window.CasaMia = { loadMenu: loadMenu };
})();
