# API du projet 48hAmazingAlgorithm
Ce dépôt contient le code source de l'API développée pour le projet 48hAmazingAlgorithm. L'API est construite en Go en utilisant le framework Gin et interagit avec une base de données MongoDB.​

## Table des matières

### [Prérequis](#Prérequis)

### [Installation](#installation)

### [Configuration](#configuration)

### [Utilisation](#utilisation)

### [Structure](#Structure)

### [Routes](#Routes)


### Prérequis
Avant de commencer, assurez-vous d'avoir installé les éléments suivants :

> Go (version 1.24 ou supérieure)

> MongoDB

### Installation
Cloner le dépôt :

```
git clone https://github.com/48hAmazingAlgorithm/Back.git
```

Installer les dépendances :

Le projet utilise beaucoup de modules Go. Pour installer les dépendances nécessaires, exécutez :
```
go mod tidy
```

### Configuration
Le projet utilise un fichier .env pour gérer les variables d'environnement sensibles.
Vous devrez également ajouter votre adresse IP sur mongoDB, afin de pouvoir accèder a la base de données. Pour cela, vous pouvez utiliser les identifiants googles situé dans le fichier .env.

### Utilisation
Pour démarrer le serveur, exécutez la commande suivante :
```
go run server.go
```
Le serveur devrait maintenant être en cours d'exécution à l'adresse http://localhost:8080.​

### Structure
Voici un aperçu de la structure des répertoires et fichiers principaux du projet :
```
Back/
├── routes/             # Contient les fichiers de définition des routes de l'API
│   └── ...             # Fichiers de routes
├── .env                # Variables d'environnement (non inclus dans le dépôt)
├── go.mod              # Fichier de gestion des modules Go
├── go.sum              # Sommes de contrôle des modules Go
├── readMe.md           # Documentation du projet
└── server.go           # Point d'entrée principal de l'application
```

### Routes
L'API expose plusieurs endpoints pour gérer les ressources. Pour plus d'informations, aller dans le fichier server.go pour voir toutes les routes. Voici quelques exemples :​

### Get 
```
/getClients -> récupère les clients
/getIndividus -> récupère les individus
/getPhotoVerso/clientId -> récupère la photo verso d'un client
```
#### Post
```
/postIndividu -> envoie les informations d'un individu
/postPhotoRecto -> envoie la photo recto
```

