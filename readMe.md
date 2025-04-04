# API du projet 48hAmazingAlgorithm
Ce dépôt contient le code source de l'API développée pour le projet 48hAmazingAlgorithm. L'API est construite en Go en utilisant le framework Gin et interagit avec une base de données MongoDB.​

## Table des matières

### [Prérequis](#prérequis)

### [Installation](#installation)

### [Configuration](#configuration)

### [Utilisation](#utilisation)

### [Structure](#Structure)

### [Routes](#routes)


### Prérequis
Avant de commencer, assurez-vous d'avoir installé les éléments suivants :

> Go (version 1.24 ou supérieure)

> MongoDB

### Installation
Cloner le dépôt :

bash
Copier
Modifier
git clone https://github.com/48hAmazingAlgorithm/Back.git
Naviguer vers le répertoire du projet :

bash
Copier
Modifier
cd Back
Installer les dépendances :

Le projet utilise les modules Go. Pour installer les dépendances nécessaires, exécutez :

bash
Copier
Modifier
go mod tidy
Configuration
Le projet utilise un fichier .env pour gérer les variables d'environnement sensibles. Un fichier .env.dist est fourni comme modèle. Renommez-le en .env et configurez les variables en conséquence :​
GitHub
+1
Postman API Platform
+1

env
Copier
Modifier
# Serveur
PORT=8080
MONGO_URL=mongodb://localhost:27017/Challenge48h

# Autres configurations...
Utilisation
Pour démarrer le serveur, exécutez la commande suivante :

bash
Copier
Modifier
go run server.go
Le serveur devrait maintenant être en cours d'exécution à l'adresse http://localhost:8080.​

Structure du projet
Voici un aperçu de la structure des répertoires et fichiers principaux du projet :

bash
Copier
Modifier
Back/
├── routes/             # Contient les fichiers de définition des routes de l'API
│   ├── individu.go     # Routes liées aux individus
│   └── ...             # Autres fichiers de routes
├── .env                # Variables d'environnement (non inclus dans le dépôt)
├── .env.dist           # Modèle pour le fichier .env
├── go.mod              # Fichier de gestion des modules Go
├── go.sum              # Sommes de contrôle des modules Go
├── individu.html       # Fichier HTML pour l'interface des individus
├── readMe.md           # Documentation du projet
└── server.go           # Point d'entrée principal de l'application
Routes de l'API
L'API expose plusieurs endpoints pour gérer les ressources. Voici quelques exemples :​
GitHub

GET /individus : Récupère la liste de tous les individus.​

POST /individus : Crée un nouvel individu.​

GET /individus/:id : Récupère les détails d'un individu spécifique par son identifiant.​

PUT /individus/:id : Met à jour les informations d'un individu existant.​

DELETE /individus/:id : Supprime un individu par son identifiant.​

Pour plus de détails sur les routes disponibles, consultez le fichier routes/individu.go.​
GitHub

Contribuer
Les contributions sont les bienvenues ! Pour contribuer :

Forkez le projet.

Créez une branche pour votre fonctionnalité (git checkout -b feature/ma-fonctionnalité).

Commitez vos modifications (git commit -m 'Ajout de ma fonctionnalité').

Poussez vers la branche (git push origin feature/ma-fonctionnalité).

Ouvrez une Pull Request.

Licence
Ce projet est sous licence MIT. Voir le fichier LICENSE pour plus de détails.​

