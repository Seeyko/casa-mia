# Refonte Casa Mia - Document de référence

## 1. Contexte

Casa Mia est une pizzeria avec deux points de vente :
- **Entraigues-sur-la-Sorgue** : 51 Rue Laurent Bertrand, 84320 — Tel: 06 45 79 49 30
- **Althen-des-Paluds** : 254 Avenue Ernest Perrin, 84210 — Tel: 04 90 36 16 33

Le site actuel est un site statique HTML/CSS hébergé sur le repo `Seeyko/casa-mia`. La refonte vise à passer à une vraie application avec backend Go et panel admin.

**Stack technique cible :** Go + Chi router + PostgreSQL + JWT auth, basé sur le repo `Seeyko/Personnal-Website` (site perso tomandrieu.com qui a déjà auth, admin panel, etc.).

**Hébergement :** Dokploy (déjà configuré). Ne pas déployer directement — commit & push uniquement, validation en local avant déploiement.

---

## 2. Cahier des charges (besoins client)

### 2.1 Design global
- Passer à un design **plus sombre, élégant, premium**

### 2.2 Hero banner
- Changer la photo de fond par celle du hero de la page Contact actuelle (`pizza-four-hero.jpg`)
- Ajouter un **indicateur d'ouverture en temps réel** : "Ouvert" ou "Fermé — prochaine ouverture : [date/heure]"
- Ajouter un bouton **"Appelez-nous"** qui ouvre un popover (bottomsheet sur mobile) avec les deux numéros distincts : Entraigues et Althen

### 2.3 Section actualités (sous le hero)
- Bloc pour 1 à 2 actualités avec photo
- Editable depuis l'admin

### 2.4 Coordonnées
- Remonter les coordonnées des deux pizzerias (jours d'ouverture, horaires, téléphones) **directement sur la page d'accueil**
- Ne plus les cacher dans un onglet "Contact"

**Horaires :**
| | Entraigues | Althen-des-Paluds |
|---|---|---|
| Lundi | Fermé | 18h-21h30 |
| Mardi | 9h-13h / 16h-21h | 18h-21h30 |
| Mercredi | 9h-13h / 16h-21h | Fermé |
| Jeudi | 9h-13h / 16h-21h | 18h-21h30 |
| Vendredi | 9h-13h / 16h-21h30 | 18h-22h |
| Samedi | 9h-13h / 16h-21h30 | 18h-22h |
| Dimanche | Fermé | 18h-21h30 |

### 2.5 Contenu texte
- **Supprimer TOUTES les mentions "feu de bois" / "cuisson au feu de bois"** partout sur le site. C'est factuellement faux.

### 2.6 Carte et menus
- Séparer en deux sections/onglets : **"La Carte"** (Pizzas + Snacking) et **"Traiteur"** (formules apéritives, bouchées, verrines)
- Chaque item de chaque carte doit pouvoir avoir une photo

### 2.7 Commande en ligne
- Ajouter un choix de point de retrait :
  - **Entraigues** : SMS/WhatsApp
  - **Althen** : appel téléphonique uniquement (numéro fixe 04 90)
- Afficher la mention : *"Tant qu'il n'y a pas de réponse de notre part, la commande n'est pas validée et l'horaire n'est pas confirmé."*
- **Ne pas mettre en avant la commande** (pas de CTA visible, accessible uniquement pour ceux qui la trouvent)

### 2.8 Espace admin (backend Go)
- Interface admin claire, soignée, agréable à utiliser
- **Gestion des actualités** : ajouter, modifier, supprimer (avec photo)
- **Gestion complète des cartes** (pizzas, snacking, traiteur) : ajouter, modifier, supprimer des items avec prix, description et photo par item

---

## 3. Inventaire photos

38 photos analysées et catégorisées (voir `photos-analysis.json` pour le détail complet). Disponibles dans le dossier `photos/`.

### 3.1 Photos HERO (pizza/calzone devant four à bois en flammes) — 4 photos
Photos prioritaires pour le hero banner et les visuels forts du site.

| Fichier | Sujet | Nom dans le repo |
|---|---|---|
| IMG_9079 | Pizza artichauts/légumes devant four flammes | `hero-pizza-artichauts-four.jpg` |
| IMG_9081 | Pizza champignons/jambon devant four flammes | `hero-pizza-champignons-jambon-four.jpg` |
| IMG_9088 | Pizza polpettes/saucisse devant four flammes | `hero-pizza-polpettes-four.jpg` |
| IMG_9095 | Calzone doré dans four flammes | `hero-calzone-four.jpg` |

### 3.2 Pizzas en boîte (galerie carte Pizzas) — 14 photos

| Fichier | Pizza | Correspondance menu | Nom dans le repo |
|---|---|---|---|
| IMG_9089 | Crème jambon cru / pesto | Pizza crème jambon cru pesto | `pizza-creme-jambon-cru-pesto.jpg` |
| IMG_9093 | Quattro Stagioni | Pizza 4 saisons | `pizza-quattro-stagioni.jpg` |
| IMG_9096 | Diavola (salami piquant) | Pizza Diavola | `pizza-diavola-new.jpg` |
| IMG_9097 | Crème jambon cru / pesto (variante) | Pizza crème jambon cru pesto | `pizza-creme-jambon-cru-pesto-2.jpg` |
| IMG_9098 | Truffe noire (crème) | Pizza Tartufo | `pizza-tartufo.jpg` |
| IMG_9099 | Boscaiola bacon/oignons (crème) | Pizza Boscaiola | `pizza-boscaiola-bacon-oignons.jpg` |
| IMG_9100 | Bresaola / mozza / roquette / balsamique | Pizza Bresaola balsamique | `pizza-bresaola.jpg` |
| IMG_9101 | Bresaola simple sur roquette | Pizza Bresaola | `pizza-bresaola-2.jpg` |
| IMG_9102 | Thon / olives / poivrons | Pizza Tonno / Siciliana | `pizza-tonno.jpg` |
| IMG_9103 | Boscaiola (variante) | Pizza Boscaiola | `pizza-boscaiola-2.jpg` |
| IMG_9104 | Truffe noire (variante) | Pizza Tartufo | `pizza-tartufo-2.jpg` |
| IMG_9105 | Bresaola / balsamique (variante) | Pizza Bresaola balsamique | `pizza-manzo-balsamique.jpg` |
| IMG_9106 | Bresaola simple (variante) | Pizza Bresaola | `pizza-manzo-balsamique-2.jpg` |
| IMG_9112 | Burrata / tomates colorées / roquette | Pizza Burrata / Estiva | `pizza-burrata.jpg` |

### 3.3 Traiteur — 18 photos

#### Coeurs d'apéro (4)
| Nom dans le repo | Sujet |
|---|---|
| `coeur-apero-carpaccio-tomates.jpg` | Planche charcuterie / carpaccio |
| `coeur-apero-saumon-gravelax.jpg` | Saumon Gravelax sur roquette |
| `coeur-apero-saumon-gravelax-2.jpg` | Saumon Gravelax grand plateau |
| `coeur-apero-legumes-croquants.jpg` | Légumes croquants + sauces + bruschetta |

#### Bouchées (4)
| Nom dans le repo | Sujet |
|---|---|
| `bouchees-feuilletes-assortis.jpg` | Feuilletés salés (palmiers, roulés pesto, tourtes) |
| `bouchees-melon-jambon.jpg` | Melon / jambon cru sur boule mozza |
| `bouchees-mini-burgers-focaccia.jpg` | Mini burgers polpettes + focaccia |
| `bouchees-navettes-saumon.jpg` | Navettes saumon / tarama pavot |

#### Brochettes (4)
| Nom dans le repo | Sujet |
|---|---|
| `brochettes-tomate-mozza-speck.jpg` | Tomate cerise / mozza / jambon cru / pesto |
| `brochettes-crevettes.jpg` | Crevettes marinées basilic |
| `brochettes-bresaola-balsamique.jpg` | Bresaola / roquette / parmesan / balsamique |
| `brochettes-assortiment.jpg` | Vue d'ensemble brochettes variées |

#### Verrines (2)
| Nom dans le repo | Sujet |
|---|---|
| `verrine-mousse-saumon.jpg` | Mousse de saumon fumé |
| `verrine-tartare-boeuf.jpg` | Tartare boeuf à l'italienne |

#### Planches (4)
| Nom dans le repo | Sujet |
|---|---|
| `planche-charcuterie-multi.jpg` | Plateaux individuels charcuterie |
| `planche-charcuterie-fromages.jpg` | Charcuterie + fromages (provolone, pecorino) |
| `planche-charcuterie-premium.jpg` | Planche XL avec rose de coppa |
| `planche-camembert-charcuterie.jpg` | Camembert rôti au four + charcuterie |

#### Dolci (2)
| Nom dans le repo | Sujet |
|---|---|
| `dolci-mini-cannoli.jpg` | Mini cannoli siciliens |
| `dolci-assortiment.jpg` | Assortiment tiramisu café, choco pistache, panna cotta |

### 3.4 Note sur les doublons/variantes
- IMG_9081 envoyé 2x (lot 1 + lot 3) — 1 seule copie gardée
- IMG_9103-9106 sont des variantes proches de IMG_9098-9101 (angles/versions différents) — toutes conservées pour le choix final

---

## 4. Placement recommandé des photos

### Page d'accueil
- **Hero banner** : `hero-pizza-polpettes-four.jpg` ou `hero-pizza-artichauts-four.jpg` (les plus impactantes)
- **Section actus** : photos au choix selon l'actu du moment

### Carte — Pizzas
Chaque item du menu pizza devrait être associé à sa photo correspondante (colonne "Correspondance menu" dans le tableau 3.2).

### Carte — Traiteur
- **Coeurs d'apéro** : les 4 photos de la section 3.3
- **Bouchées** : les 4 photos bouchées
- **Brochettes** : les 4 photos brochettes
- **Verrines** : les 2 photos verrines
- **Planches** : les 4 photos planches
- **Dolci** : les 2 photos dolci

### Galerie / visuels d'ambiance
Les 4 photos HERO (devant le four) peuvent être réutilisées en background ou en galerie pour donner l'ambiance "pizzeria authentique".

---

## 5. Architecture technique cible

```
casa-mia/
├── backend/           # Go + Chi + PostgreSQL + JWT
│   ├── main.go
│   ├── handlers/      # admin, auth, news, menu, upload
│   ├── models/        # News, MenuCategory, MenuItem
│   ├── migrations/
│   └── uploads/       # images uploadées
├── frontend/          # HTML/CSS/JS statique servi par Go
│   ├── index.html     # accueil (hero, actus, coordonnées)
│   ├── carte.html     # La Carte (pizzas + snacking)
│   ├── traiteur.html  # Traiteur
│   ├── style.css
│   └── images/
├── photos/            # photos originales analysées
├── photos-analysis.json
├── docker-compose.yml
├── Dockerfile
└── REFONTE-CASAMIA.md # ce document
```

**Backend basé sur** `Seeyko/Personnal-Website` : réutiliser auth JWT, admin panel pattern (secret path + API key), structure Chi router.

**Tables BDD :**
- `news` : id, title, content, image_url, published_at, created_at
- `menu_categories` : id, name, section (carte|traiteur), sort_order
- `menu_items` : id, category_id, name, description, price, image_url, sort_order
- `settings` : id, key, value (pour horaires, etc.)

---

## 6. Fichiers dans ce repo

| Fichier/Dossier | Description |
|---|---|
| `photos/` | 38 photos analysées et renommées (prêtes à l'emploi) |
| `photos-analysis.json` | Analyse détaillée de chaque photo (catégorie, sujet, description, correspondance menu, usage recommandé) |
| `REFONTE-CASAMIA.md` | Ce document |
| `*.html` / `style.css` / `images/` | Site actuel (à refaire) |
