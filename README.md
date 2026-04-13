# CasaMia Pizzeria

Site vitrine de **CasaMia**, pizzeria artisanale sicilienne avec deux points de vente :
- **Entraigues-sur-la-Sorgue** : 51 Rue Laurent Bertrand, 84320 — Tel: 06 45 79 49 30
- **Althen-des-Paluds** : 254 Avenue Ernest Perrin, 84210 — Tel: 04 90 36 16 33

## Stack technique

- **Backend** : Go 1.24 + Chi router + PostgreSQL + JWT auth
- **Frontend** : HTML/CSS/JS statique servi par le backend Go
- **Déploiement** : Docker (Dokploy)

## Lancer en local

```bash
docker-compose -f docker-compose.local.yml up --build
```

Puis ouvrir [http://localhost:3000](http://localhost:3000).

## Structure du projet

```
casa-mia/
├── backend/           # API Go (handlers, models, services, middleware)
├── frontend/          # Pages HTML + CSS + JS + images
├── photos/            # 38 photos sources renommées
├── Dockerfile         # Build production multi-stage
├── docker-compose.yml # Production (Dokploy)
└── docker-compose.local.yml  # Développement local
```

## Documentation

Voir `REFONTE-CASAMIA.md` pour le cahier des charges complet et `progress-refonte.md` pour le suivi de progression.
