// push-swap affiche la liste d'instructions la plus courte possible pour trier les nombres reçus.
//
//	./push-swap "2 1 3 6 5 8"
package main

import (
	"bufio"
	"fmt"
	"os"

	"push-swap/pile"
)

func main() {
	nombres, err := pile.Lire(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		os.Exit(1)
	}

	// Sans argument, la liste est vide : rien n'est affiché.
	sortie := bufio.NewWriter(os.Stdout)
	defer sortie.Flush()
	for _, ins := range pile.Trier(nombres) {
		fmt.Fprintln(sortie, ins)
	}
}
