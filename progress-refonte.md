# Progress Refonte Casa Mia

> **Date de l'audit :** 2026-04-14
> **Branche :** `refonte-casamia`
> **Dernier commit :** `2a7b7fd` — fix: remove root health check route so frontend index.html is served

---

## Legende

| Symbole | Signification |
|---------|---------------|
| FAIT | Fonctionnel et conforme au cahier des charges |
| PARTIEL | Implemente mais incomplet ou avec des reserves |
| MANQUANT | Non implemente |
| BUG | Probleme identifie a corriger |

---

## 1. DESIGN GLOBAL

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 1.1 | Design sombre, elegant, premium | FAIT | Palette noire (`#100E0C`), accents dores (`#D4A853`), Playfair Display + Inter, systeme d'ombres premium |
| 1.2 | Responsive mobile | FAIT | Breakpoints 768px + 480px, hamburger menu, bottomsheet mobile |
| 1.3 | CSS complet | FAIT | `frontend/css/style.css` — 849 lignes, variables CSS, animations |

---

## 2. HERO BANNER

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 2.1 | Photo de fond `pizza-four-hero.jpg` | FAIT | `index.html` ligne 73 : `background-image: url('images/pizza-four-hero.jpg')` |
| 2.2 | Indicateur ouverture temps reel | FAIT | Badges dynamiques via `/api/status`, affiche "Ouvert" / "Ferme — prochaine ouverture : [date/heure]" |
| 2.3 | Bouton "Appelez-nous" | FAIT | Bouton avec icone telephone, ouvre popover desktop / bottomsheet mobile |
| 2.4 | Popover avec 2 numeros distincts | FAIT | Entraigues (SMS/WhatsApp 06 45 79 49 30) + Althen (Tel 04 90 36 16 33) |

---

## 3. SECTION ACTUALITES

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 3.1 | Bloc 1-2 actus avec photo sur homepage | FAIT | Section `newsSection` charge dynamiquement depuis `/api/news` (limite 5, publiees) |
| 3.2 | Editable depuis l'admin | FAIT | Admin panel : onglet "Actualites" avec CRUD complet (titre, contenu, image, publie) |

---

## 4. COORDONNEES

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 4.1 | Coordonnees sur la page d'accueil | FAIT | Section `locationsGrid` sur `index.html`, plus cachees dans un onglet Contact |
| 4.2 | Jours d'ouverture + horaires | FAIT | Charge depuis `/api/locations`, affiche horaires jour par jour |
| 4.3 | Telephones des 2 points de vente | FAIT | Entraigues + Althen avec numeros distincts |

### Verification des horaires (seed data backend)

| Jour | Entraigues (attendu) | Entraigues (implemente) | Althen (attendu) | Althen (implemente) |
|------|---------------------|------------------------|-------------------|---------------------|
| Lundi | Ferme | Ferme (null) | 18h-21h30 | 18h-21h30 |
| Mardi | 9h-13h / 16h-21h | 9h-13h / 16h-21h | 18h-21h30 | 18h-21h30 |
| Mercredi | 9h-13h / 16h-21h | 9h-13h / 16h-21h | Ferme | Ferme (null) |
| Jeudi | 9h-13h / 16h-21h | 9h-13h / 16h-21h | 18h-21h30 | 18h-21h30 |
| Vendredi | 9h-13h / 16h-21h30 | 9h-13h / 16h-21h30 | 18h-22h | 18h-22h |
| Samedi | 9h-13h / 16h-21h30 | 9h-13h / 16h-21h30 | 18h-22h | 18h-22h |
| Dimanche | Ferme | Ferme (null) | 18h-21h30 | 18h-21h30 |

**Resultat :** FAIT — Tous les horaires correspondent au cahier des charges.

---

## 5. CONTENU TEXTE — "FEU DE BOIS"

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 5.1 | Supprimer TOUTES les mentions "feu de bois" | BUG | 13 mentions restantes dans les fichiers racine (ancien site) |

### Detail des mentions restantes

| Fichier | Lignes | Contenu |
|---------|--------|---------|
| `/index.html` (racine) | 7, 9, 23, 105, 147 | Meta descriptions, JSON-LD, trust bar, texte carte |
| `/carte.html` (racine) | 7, 9, 20, 78, 87 | Meta descriptions, JSON-LD, hero description, trust bar |
| `/contact.html` (racine) | 20 | JSON-LD description |
| `/histoire.html` (racine) | 20 | JSON-LD description |
| `README.md` | 7 | Description du projet |

**Note :** Les fichiers dans `/frontend/` sont PROPRES (0 mention). Les fichiers en racine sont l'ancien site statique. Si ces fichiers racine ne sont plus servis (le backend sert `/frontend/`), le probleme est non-bloquant mais le nettoyage reste recommande.

---

## 6. CARTE ET MENUS

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 6.1 | Separation "La Carte" / "Traiteur" | FAIT | `menu.html` avec onglets `data-section="carte"` et `data-section="traiteur"` |
| 6.2 | Carte = Pizzas + Snacking | FAIT | 8 categories carte : Pizzas Tomate, Pizzas Creme, Calzones, Snacking, Focaccias, Desserts, Boissons, Supplements |
| 6.3 | Traiteur = formules aperitives, bouchees, verrines | FAIT | 6 categories traiteur : Coeurs d'apero, Bouchees, Brochettes, Verrines, Planches, Dolci |
| 6.4 | Photo par item | FAIT | Chaque `menu_item` a un champ `image_path`, affichage via `/api/images/` |
| 6.5 | Admin gestion complete | FAIT | CRUD categories + items depuis admin panel |

### Items de menu seed (70+ items)

| Section | Categories | Items (exemples) |
|---------|-----------|-----------------|
| Carte | Pizzas Tomate (12), Pizzas Creme (8), Calzones (4), Snacking (8), Focaccias (4), Desserts (6), Boissons (6), Supplements (4) | Margherita, Diavola, Burrata, Tartufo, Calzone Classique, Arancini, Tiramisu... |
| Traiteur | Coeurs d'apero (4), Bouchees (4), Brochettes (4), Verrines (2), Planches (5), Dolci (3) | Carpaccio-tomates, Feuilletes, Crevettes, Mousse saumon, Charcuterie XL, Cannoli... |

---

## 7. COMMANDE EN LIGNE

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 7.1 | Choix point de retrait | FAIT | Radio buttons Entraigues / Althen dans le panier |
| 7.2 | Entraigues = SMS/WhatsApp | FAIT | Boutons WhatsApp + SMS avec numero 06 45 79 49 30 |
| 7.3 | Althen = appel telephonique (04 90) | FAIT | Bouton appel avec numero 04 90 36 16 33 |
| 7.4 | Mention "pas de reponse = pas valide" | FAIT | Message affiche : "Tant qu'il n'y a pas de reponse de notre part, la commande n'est pas validee et l'horaire n'est pas confirme." |
| 7.5 | Ne pas mettre en avant la commande | FAIT | Panier flottant `display:none` par defaut, accessible uniquement pour ceux qui trouvent le bouton |

---

## 8. ESPACE ADMIN (BACKEND GO)

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 8.1 | Interface admin claire et soignee | FAIT | `admin.html` avec design sombre assorti, login, navigation par onglets |
| 8.2 | Gestion actualites (CRUD + photo) | FAIT | Ajouter, modifier, supprimer avec upload image + flag publie |
| 8.3 | Gestion cartes complete | FAIT | CRUD categories (nom, section, ordre) + CRUD items (nom, description, prix, photo, badge, note, disponibilite) |
| 8.4 | Auth JWT securisee | FAIT | Login, change password, reset password avec tokens temporaires |
| 8.5 | Secret path admin | FAIT | Routes admin derriere `/api/{ADMIN_SECRET_PATH}/`, configurable par env var |
| 8.6 | Rate limiting | FAIT | 10 req/min, lockout a 20 tentatives pour 15 min |
| 8.7 | Gestion des horaires admin | FAIT | `PUT /api/{PATH}/locations/{id}` permet de modifier horaires (JSONB) |

---

## 9. ARCHITECTURE TECHNIQUE

### 9.1 Structure des fichiers

| Attendu (spec) | Implemente | Statut |
|----------------|-----------|--------|
| `backend/main.go` | `backend/main.go` | FAIT |
| `backend/handlers/` | `backend/handlers/` (admin, auth, news, menu, locations, images, health) | FAIT |
| `backend/models/` | `backend/models/` (news, menu, location) | FAIT |
| `backend/migrations/` | Embedded dans `backend/services/database.go` | FAIT (approche differente mais fonctionnelle) |
| `backend/uploads/` | Configurable via `UPLOAD_DIR` (defaut: `./uploads`) | FAIT |
| `frontend/index.html` | `frontend/index.html` | FAIT |
| `frontend/carte.html` | `frontend/menu.html` (combine carte + traiteur en onglets) | FAIT (meilleure approche) |
| `frontend/traiteur.html` | Integre dans `frontend/menu.html` | FAIT (pas de page separee, onglet traiteur) |
| `frontend/style.css` | `frontend/css/style.css` | FAIT |
| `frontend/images/` | `frontend/images/` (54 fichiers) | FAIT |
| `photos/` | `photos/` (38 fichiers) | FAIT |
| `photos-analysis.json` | `photos-analysis.json` | FAIT |
| `Dockerfile` | `Dockerfile` (multi-stage Go + frontend) | FAIT |
| `docker-compose.yml` | `docker-compose.yml` + `docker-compose.local.yml` | FAIT |

### 9.2 Tables BDD

| Table attendue | Table implementee | Statut | Differences |
|----------------|------------------|--------|-------------|
| `news` | `news` | FAIT | `image_url` renomme en `image_path` (plus precis) |
| `menu_categories` | `menu_categories` | FAIT | Conforme |
| `menu_items` | `menu_items` | FAIT | Champs bonus : `badge`, `note`, `available` |
| `settings` | NON CREE | MANQUANT | Pas de table settings pour config admin (horaires gerees via `locations`) |
| — | `locations` | FAIT (bonus) | Gestion des points de vente + horaires JSONB |
| — | `admin_users` | FAIT (bonus) | Auth avec reset password |

### 9.3 Backend technique

| Composant | Statut | Detail |
|-----------|--------|--------|
| Go + Chi router | FAIT | Go 1.24, Chi v5.0.12 |
| PostgreSQL | FAIT | Driver lib/pq, migrations embedded |
| JWT auth | FAIT | golang-jwt v5, HMAC-SHA256, 24h expiration |
| CORS | FAIT | go-chi/cors configurable |
| Image upload | FAIT | UUID filenames, MIME validation, 5MB max |
| Rate limiting | FAIT | Custom implementation per-IP |
| Seed data | FAIT | Admin user + 2 locations + 14 categories + 70+ items |

---

## 10. PHOTOS ET ASSETS

### 10.1 Photos sources (`photos/`)

| Categorie | Attendu | Present | Statut |
|-----------|---------|---------|--------|
| Hero (four) | 4 | 4 | FAIT |
| Pizzas en boite | 14 | 14 | FAIT |
| Traiteur - Coeurs apero | 4 | 4 | FAIT |
| Traiteur - Bouchees | 4 | 4 | FAIT |
| Traiteur - Brochettes | 4 | 4 | FAIT |
| Traiteur - Verrines | 2 | 2 | FAIT |
| Traiteur - Planches | 4 | 4 | FAIT |
| Traiteur - Dolci | 2 | 2 | FAIT |
| **TOTAL** | **38** | **38** | **FAIT** |

### 10.2 Images frontend (`frontend/images/`)

54 fichiers : les 38 photos renommees + 16 images supplementaires (logo, variantes optimisees, hero background, photos additionnelles de plats).

### 10.3 Placement des photos

| Zone | Photo recommandee (spec) | Implementee | Statut |
|------|-------------------------|-------------|--------|
| Hero banner | `hero-pizza-polpettes-four.jpg` ou `hero-pizza-artichauts-four.jpg` | `pizza-four-hero.jpg` (background) + `hero-pizza-polpettes-four.jpg` (slideshow) | FAIT |
| Carte Pizzas | Photos associees par correspondance menu | Via `image_path` en BDD, chargement dynamique | FAIT |
| Carte Traiteur | Photos par categorie | Via `image_path` en BDD, chargement dynamique | FAIT |

---

## 11. PAGES DU SITE

| Page | Fichier | Statut | Contenu |
|------|---------|--------|---------|
| Accueil | `frontend/index.html` | FAIT | Hero + actus + coordonnees + galerie |
| Menu | `frontend/menu.html` | FAIT | Onglets Carte/Traiteur + panier + commande |
| Histoire | `frontend/histoire.html` | FAIT | Hero, story Tony & Cindy, valeurs, galerie, CTA |
| Admin | `frontend/admin.html` | FAIT | Login + gestion actus/menu/categories/settings |
| Contact | Supprime (integre a l'accueil) | FAIT | Coordonnees remontees sur homepage comme demande |

---

## 12. DEPLOIEMENT

| # | Exigence | Statut | Detail |
|---|----------|--------|--------|
| 12.1 | Dockerfile fonctionnel | FAIT | Multi-stage build, Go + frontend, seed images |
| 12.2 | docker-compose production | FAIT | Service unique `app`, dokploy-network externe |
| 12.3 | docker-compose local | FAIT | 3 services (frontend Nginx, backend Go, PostgreSQL 16) |
| 12.4 | Hebergement Dokploy | FAIT | Config docker-compose compatible, commit & push workflow |
| 12.5 | Variables d'environnement | FAIT | DATABASE_URL, JWT_SECRET, ADMIN_SECRET_PATH, ADMIN_USERNAME, ADMIN_PASSWORD, FRONTEND_URL, PORT |

---

## 13. SECURITE

| # | Point | Statut | Detail |
|---|-------|--------|--------|
| 13.1 | JWT Bearer token | FAIT | 24h expiration, HMAC-SHA256 |
| 13.2 | Bcrypt password hashing | FAIT | Cost factor 12 |
| 13.3 | Rate limiting auth | FAIT | 10 req/min, lockout 15 min |
| 13.4 | IP allowlisting (optionnel) | FAIT | Via env var ADMIN_ALLOWED_IPS |
| 13.5 | CORS configurable | FAIT | Via FRONTEND_URL env var |
| 13.6 | Security headers | FAIT | no-store, nosniff, DENY, no-referrer |
| 13.7 | Path traversal protection | FAIT | Images : bloque "..", "/", "\\" |
| 13.8 | MIME validation uploads | FAIT | jpeg, png, webp, gif uniquement |
| 13.9 | Pas de secrets dans le code | FAIT | Tout configurable par env var |

---

## SYNTHESE GLOBALE

### Compteur d'exigences

| Categorie | Total | FAIT | PARTIEL | MANQUANT | BUG |
|-----------|-------|------|---------|----------|-----|
| Design global | 3 | 3 | 0 | 0 | 0 |
| Hero banner | 4 | 4 | 0 | 0 | 0 |
| Actualites | 2 | 2 | 0 | 0 | 0 |
| Coordonnees | 3 | 3 | 0 | 0 | 0 |
| Contenu "feu de bois" | 1 | 0 | 0 | 0 | 1 |
| Carte et menus | 5 | 5 | 0 | 0 | 0 |
| Commande en ligne | 5 | 5 | 0 | 0 | 0 |
| Espace admin | 7 | 7 | 0 | 0 | 0 |
| Architecture technique | 8 | 7 | 0 | 1 | 0 |
| Photos et assets | 3 | 3 | 0 | 0 | 0 |
| Pages du site | 4 | 4 | 0 | 0 | 0 |
| Deploiement | 5 | 5 | 0 | 0 | 0 |
| Securite | 9 | 9 | 0 | 0 | 0 |
| **TOTAL** | **59** | **56** | **0** | **1** | **1** |

### Taux de completion : 95% (56/59)

---

## ACTIONS RESTANTES

### BUG — Priorite haute
1. **Supprimer les mentions "feu de bois"** dans les fichiers racine (`/index.html`, `/carte.html`, `/contact.html`, `/histoire.html`, `README.md`) — 13 occurrences a supprimer ou remplacer

### MANQUANT — Priorite basse
2. **Table `settings`** : prevue dans le spec mais non implementee. Les horaires sont geres via la table `locations` (JSONB), ce qui est une approche tout aussi valide. A creer uniquement si un besoin de parametres admin generiques apparait.

### RECOMMANDATIONS (hors spec, non bloquant)
3. **Nettoyage fichiers racine** : les fichiers HTML a la racine (`carte.html`, `contact.html`, `histoire.html`, `index.html`, `traiteur.html`, `style.css`) sont l'ancien site. Ils ne sont plus servis par le backend mais encombrent le repo. Envisager leur suppression.
4. **Rotation hero photos** : seul `hero-pizza-polpettes-four.jpg` est reference statiquement. Les 3 autres photos hero (artichauts, champignons-jambon, calzone) pourraient etre integrees dans un carousel.
5. **Structured error responses** : les erreurs API sont basiques (`http.Error`). Envisager un format JSON uniforme pour les erreurs.

---

## FICHIERS CLES DE LA REFONTE

```
casa-mia/
├── backend/
│   ├── main.go                          # Point d'entree, routing, middleware
│   ├── handlers/
│   │   ├── admin.go                     # CRUD news, categories, items, locations
│   │   ├── auth.go                      # Login, change password, reset
│   │   ├── news.go                      # API publique news
│   │   ├── menu.go                      # API publique menu
│   │   ├── locations.go                 # API publique locations + status
│   │   ├── images.go                    # Upload + serve images
│   │   └── health.go                    # Health check
│   ├── middleware/
│   │   └── admin_auth.go               # JWT validation, rate limiting, IP filter
│   ├── models/
│   │   ├── news.go                      # Struct News
│   │   ├── menu.go                      # Structs MenuCategory, MenuItem
│   │   └── location.go                  # Structs Location, LocationStatus
│   ├── services/
│   │   ├── database.go                  # PostgreSQL init + migrations
│   │   ├── jwt.go                       # JWT generation/validation
│   │   ├── seed.go                      # Donnees initiales (locations, menu, admin)
│   │   └── ratelimit/ratelimit.go       # Rate limiter per-IP
│   ├── go.mod / go.sum
│   └── .env.example
├── frontend/
│   ├── index.html                       # Page accueil (hero, actus, coords)
│   ├── menu.html                        # Page menu (carte + traiteur + panier)
│   ├── histoire.html                    # Page histoire
│   ├── admin.html                       # Panel admin
│   ├── css/style.css                    # 849 lignes CSS
│   ├── js/app.js                        # 582 lignes JS
│   ├── js/config.js                     # Config API
│   └── images/                          # 54 images (38 photos + extras)
├── photos/                              # 38 photos sources renommees
├── photos-analysis.json                 # Analyse detaillee des photos
├── Dockerfile                           # Multi-stage production build
├── docker-compose.yml                   # Production (Dokploy)
├── docker-compose.local.yml             # Developpement local
├── REFONTE-CASAMIA.md                   # Cahier des charges
└── progress-refonte.md                  # CE FICHIER — suivi de progression
```
