// Package pile contient les deux piles, les 11 instructions et l'algorithme de tri.
package pile

import (
	"errors"
	"strconv"
	"strings"
)

// ErrArguments : un argument n'est pas un entier, ou il y a un doublon.
var ErrArguments = errors.New("arguments invalides")

// ErrInstruction : l'instruction n'existe pas ou est mal écrite.
var ErrInstruction = errors.New("instruction inconnue")

// Piles : les piles a et b. L'élément d'indice 0 est le haut de la pile.
type Piles struct {
	A, B []int
}

// Lire transforme les arguments en nombres. Elle accepte "2 1 3" comme 2 1 3.
func Lire(args []string) ([]int, error) {
	nombres := []int{}
	vus := map[int]bool{}
	for _, arg := range args {
		for _, mot := range strings.Fields(arg) {
			n, err := strconv.Atoi(mot)
			if err != nil {
				return nil, ErrArguments // pas un entier (ou trop grand)
			}
			if vus[n] {
				return nil, ErrArguments // doublon
			}
			vus[n] = true
			nombres = append(nombres, n)
		}
	}
	return nombres, nil
}

// Nouvelles crée les piles : tous les nombres dans a, b vide.
func Nouvelles(nombres []int) *Piles {
	return &Piles{A: append([]int{}, nombres...), B: []int{}}
}

// EstTriee : a est dans l'ordre croissant (le plus petit en haut) et b est vide.
func (p *Piles) EstTriee() bool {
	return len(p.B) == 0 && croissant(p.A)
}

func croissant(s []int) bool {
	for i := 1; i < len(s); i++ {
		if s[i-1] > s[i] {
			return false
		}
	}
	return true
}

// Executer applique une instruction par son nom ("sa", "pb", "rrr"...).
func (p *Piles) Executer(nom string) error {
	switch nom {
	case "sa":
		echanger(p.A)
	case "sb":
		echanger(p.B)
	case "ss":
		echanger(p.A)
		echanger(p.B)
	case "pa":
		pousser(&p.B, &p.A)
	case "pb":
		pousser(&p.A, &p.B)
	case "ra":
		tourner(p.A)
	case "rb":
		tourner(p.B)
	case "rr":
		tourner(p.A)
		tourner(p.B)
	case "rra":
		tournerInverse(p.A)
	case "rrb":
		tournerInverse(p.B)
	case "rrr":
		tournerInverse(p.A)
		tournerInverse(p.B)
	default:
		return ErrInstruction
	}
	return nil
}

// echanger : les 2 premiers éléments changent de place (sa, sb).
func echanger(s []int) {
	if len(s) >= 2 {
		s[0], s[1] = s[1], s[0]
	}
}

// pousser : le haut de « depuis » passe en haut de « vers » (pa, pb).
func pousser(depuis, vers *[]int) {
	if len(*depuis) == 0 {
		return
	}
	haut := (*depuis)[0]
	*depuis = (*depuis)[1:]
	*vers = append(*vers, 0)
	copy((*vers)[1:], *vers)
	(*vers)[0] = haut
}

// tourner : le premier élément devient le dernier (ra, rb).
func tourner(s []int) {
	if len(s) >= 2 {
		premier := s[0]
		copy(s, s[1:])
		s[len(s)-1] = premier
	}
}

// tournerInverse : le dernier élément devient le premier (rra, rrb).
func tournerInverse(s []int) {
	if len(s) >= 2 {
		dernier := s[len(s)-1]
		copy(s[1:], s[:len(s)-1])
		s[0] = dernier
	}
}
