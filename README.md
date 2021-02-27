# go-fyne-compte

Le jeu du **Compte est Bon** est un jeu bas� sur le calcul mental, issu du jeu t�l�vis� ** Des Chiffres et des Lettres**. A partir d'un tirage de 6 plaques, il faut trouver un nombre entre 101 et 999, en utilisant les 4 op�rations arithm�tiques d'addition, soustraction, multiplication ou division. Les r�sultats interm�diaires doivent �tre des entiers positifs.

Ce jeu est un *bon exercice de codage* pour apprendre un nouveau langage, ici [Go](https://golang.org/).  
L'algorithme est un *brute force*: le programme essaie toutes les combinaisons de calcul possible, entre toutes les 6 plaques du tirage initial, ce qui donne ensuite un jeu de 5 nombres, pour chacune des 4 op�rations, puis le programme it�re r�cursivement jusqu'� trouver le compte exact. La meilleure solution approch�e est gard�e � chaque �tape, pour le cas o� le compte exact ne peut pas �tre trouv�.  
Quelques optimisations permettent de r�duire le nombre de combinaisons possibles:
- comme les soustractions et les divisions doivent produire des entiers positifs, le programme classe les nombres dans l'ordre croissant
- les additions et les multiplications sont commutatives, il suffit de faire **A + B**, pas besoin de refaire **B + A**
- les soustractions qui aboutissent � un r�sultat **inf�rieur ou �gal � 0** sont inutiles
- les multiplications et les divisions par **1** sont inutiles
- les divisions avec un r�sultat non entier sont inutiles

Pour l'IHM, le framework utilis� est [fyne](https://github.com/fyne-io/fyne), dans sa version **v2** (no compatible avec les versions v1.x).  
Ce framework est suffisant pour faire une IHM minimaliste qui permet au jeu de fonctionner correctement.  
