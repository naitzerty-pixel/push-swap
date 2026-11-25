# push-swap

Trier une liste d'entiers avec deux piles (`a` et `b`) et 11 instructions, en utilisant le moins d'instructions possible.

Le projet contient deux programmes en Go (bibliothèque standard uniquement) :

- **push-swap** calcule et affiche la liste d'instructions qui trie les nombres ;
- **checker** lit des instructions, les exécute et affiche `OK` si `a` est triée et `b` vide, sinon `KO`.

## Compiler et lancer

```bash
go build -o push-swap ./cmd/push-swap
go build -o checker ./cmd/checker
```

(Sous Windows, les fichiers s'appellent `push-swap.exe` et `checker.exe`.)

```console
$ ./push-swap "2 1 3 6 5 8"
pb
pb
ra
sa
rrr
pa
pa
$ ARG="4 67 3 87 23"; ./push-swap "$ARG" | wc -l
6
$ ARG="4 67 3 87 23"; ./push-swap "$ARG" | ./checker "$ARG"
OK
$ ./push-swap "0 one 2 3"
Error
```

Les tests :

```bash
go test ./...
go test -v ./pile    # affiche aussi le nombre d'instructions obtenu
```

## Organisation

```
push-swap/
├── go.mod
├── cmd/
│   ├── push-swap/main.go   lit les arguments, affiche les instructions
│   └── checker/main.go     lit les instructions, affiche OK / KO
└── pile/
    ├── pile.go             les piles, la lecture des arguments, les 11 instructions
    ├── tri.go              l'algorithme de tri
    └── pile_test.go        les tests
```

Le code commun aux deux programmes est dans le paquet `pile` : `checker` et `push-swap` exécutent les instructions avec **le même code**, ce qui évite qu'ils ne soient pas d'accord.

## Les piles

Une pile est un `[]int` : l'élément d'indice `0` est le **haut** de la pile.

| Instruction | Ce que fait le code |
|---|---|
| `sa`, `sb`, `ss` | `echanger` : `s[0], s[1] = s[1], s[0]` |
| `pa`, `pb` | `pousser` : retire le haut d'une pile et l'insère en haut de l'autre |
| `ra`, `rb`, `rr` | `tourner` : le premier passe en dernier |
| `rra`, `rrb`, `rrr` | `tournerInverse` : le dernier passe en premier |

Une instruction sur une pile vide ou trop petite ne fait rien (pas de plantage).

## L'algorithme

### 1. Les rangs

On remplace chaque nombre par sa place dans l'ordre trié : `4 67 3 87 23` devient `1 3 0 4 2`. L'ordre ne change pas, mais les calculs sont plus simples (le but est toujours `0 1 2 3 ...`).

### 2. Jusqu'à 6 nombres : la solution la plus courte (parcours en largeur)

On essaie toutes les suites de 1 instruction, puis de 2, puis de 3… (un **BFS**). La première suite qui trie est forcément **la plus courte qui existe**. On retient les états déjà vus pour ne jamais repasser deux fois par le même.

C'est possible parce que le nombre d'états est petit : pour 6 nombres, il y a 6! façons de les ordonner × 7 façons de les répartir entre `a` et `b`, soit **5 040 états**.

### 3. Plus de 6 nombres : l'algorithme « Turk » (glouton)

1. On envoie tout dans `b` sauf 3 nombres. Au passage, les petits nombres (sous la moitié) sont envoyés au fond de `b` avec `rb` : `b` est grossièrement rangé.
2. On trie les 3 nombres restants dans `a` (2 instructions au plus).
3. On ramène les nombres de `b` **un par un**. Pour chaque nombre de `b`, on calcule le coût :
   - pour le mettre en haut de `b` ;
   - et pour mettre en haut de `a` **sa place** (le plus petit nombre de `a` plus grand que lui).
   
   Il y a 4 façons de faire : tout vers le haut (`rr`), tout vers le bas (`rrr`), ou un dans chaque sens. On garde le nombre et la façon **les moins chers**, puis `pa`.
4. On tourne `a` dans le sens le plus court pour mettre le plus petit en haut.

### 4. Dernière passe : supprimer l'inutile

`optimiser` relit la liste : `ra` suivi de `rra` s'annulent, `ra` suivi de `rb` devient `rr`, etc.

## Résultats

| Nombres | Instructions |
|---|---|
| 3 | 2 au pire |
| 5 | 8 au pire (toutes les permutations testées) |
| 6 | 10 au pire (toutes les permutations testées) |
| 100 | environ 580 en moyenne |
| 500 | environ 4 530 en moyenne |

## Les erreurs

`Error` est affiché sur la **sortie d'erreur** quand :

- un argument n'est pas un entier (`one`, `1.5`, un nombre trop grand) ;
- un nombre apparaît deux fois ;
- pour `checker`, une instruction n'existe pas ou est mal écrite (`rra ` avec une espace, `RA`).

Sans argument, les deux programmes n'affichent rien. Le `checker` ignore les lignes vides (par exemple la dernière ligne de `echo -e "sa\n"`).

## Les tests (`pile/pile_test.go`)

- lecture des arguments, y compris les erreurs ;
- l'exemple du sujet, étape par étape ;
- **toutes** les permutations de 1 à 6 nombres sont triées ;
- 20 listes au hasard de 7, 100 et 500 nombres sont triées, sous un nombre maximal d'instructions ;
- l'exemple `4 67 3 87 23` en 6 instructions ;
- la passe d'optimisation.
