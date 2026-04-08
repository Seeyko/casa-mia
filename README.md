# CasaMia Pizzeria

Site vitrine de **CasaMia**, pizzeria artisanale sicilienne située à Entraigues-sur-la-Sorgue (84320), Vaucluse.

## Aperçu

CasaMia est une pizzeria-traiteur-épicerie tenue par un chef sicilien, proposant des pizzas au feu de bois, un service traiteur italien et une épicerie fine de produits importés de Sicile.

## Pages

| Fichier | Description |
|---------|-------------|
| `index.html` | Page d'accueil — présentation, avis clients, galerie |
| `carte.html` | La carte — pizzas, entrées, desserts, boissons |
| `traiteur.html` | Service traiteur italien |
| `histoire.html` | Notre histoire — parcours du chef |
| `contact.html` | Coordonnées, horaires, formulaire de contact |

## Stack technique

- **HTML5** sémantique avec données structurées Schema.org (JSON-LD)
- **CSS3** pur (fichier unique `style.css`) — responsive, sans framework
- **Aucune dépendance** JavaScript de build ou framework frontend
- Images optimisées dans le dossier `images/`

## Lancer le site en local

Le site est entièrement statique. Il suffit d'ouvrir `index.html` dans un navigateur ou de lancer un serveur local :

```bash
# Avec Python
python3 -m http.server 8000

# Avec Node.js
npx serve .
```

Puis ouvrir [http://localhost:8000](http://localhost:8000).

## Structure du projet

```
casa-mia/
├── index.html          # Page d'accueil
├── carte.html          # Menu / carte
├── traiteur.html       # Service traiteur
├── histoire.html       # Notre histoire
├── contact.html        # Page contact
├── style.css           # Feuille de styles unique
└── images/             # Visuels du site
```

## Informations pratiques

- **Adresse** : Entraigues-sur-la-Sorgue, 84320
- **Téléphone** : 06 45 79 49 30
- **Horaires** : Mar–Dim, 11h30–14h00 & 18h00–22h00 (fermé le lundi)
