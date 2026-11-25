// checker lit des instructions sur l'entrée standard, les exécute sur les
// nombres reçus en argument, puis affiche OK si a est triée et b vide, sinon KO.
//
//	./push-swap "3 2 1 0" | ./checker "3 2 1 0"
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"push-swap/pile"
)

// erreur affiche « Error » sur la sortie d'erreur et arrête le programme.
func erreur() {
	fmt.Fprintln(os.Stderr, "Error")
	os.Exit(1)
}

func main() {
	nombres, err := pile.Lire(os.Args[1:])
	if err != nil {
		erreur()
	}
	if len(nombres) == 0 {
		return // pas d'argument : rien à afficher
	}

	piles := pile.Nouvelles(nombres)
	lecteur := bufio.NewScanner(os.Stdin)
	for lecteur.Scan() {
		// \r : fin de ligne Windows. Une ligne vide (souvent la dernière,
		// avec echo -e "...\n") ne contient aucune instruction.
		ligne := strings.TrimSuffix(lecteur.Text(), "\r")
		if ligne == "" {
			continue
		}
		if piles.Executer(ligne) != nil {
			erreur() // instruction inconnue ou mal écrite (ex. "rra " ou "RA")
		}
	}
	if lecteur.Err() != nil {
		erreur()
	}

	if piles.EstTriee() {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
