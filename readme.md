Homepedia API
==============

Ce projet est une API REST pour la gestion d'une base de données de logements.

Prérequis
----------

Pour exécuter ce projet, vous devez avoir Go 1.16 ou supérieur installé sur votre ordinateur.

Installation
------------

Pour installer les dépendances du projet, exécutez la commande suivante dans le répertoire racine du projet :
``` go mod tidy ```
Exécution en mode hot reload
-----------------------------

Pour exécuter l'application en mode hot reload, vous devez avoir l'outil `air` installé sur votre ordinateur.

Ensuite, exécutez la commande suivante dans le répertoire racine du projet :
``` air ```
Cette commande lancera le serveur de développement et recompilera et relancera l'application à chaque fois que vous modifierez vos fichiers Go.

Exécution en mode normal
------------------------

Pour exécuter l'application en mode normal, exécutez la commande suivante dans le répertoire racine du projet :
``` go run cmd/homepedia-api/main.go ```
Cette commande lancera l'application en utilisant le compilateur Go standard.
