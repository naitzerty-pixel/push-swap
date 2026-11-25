package pile

import "sort"

// limiteRecherche : jusqu'à ce nombre d'éléments, on cherche la solution la
// plus courte qui existe (parcours en largeur). Au-delà, on utilise un
// algorithme glouton (« Turk »), rapide et assez économe.
const limiteRecherche = 6

// toutes : les 11 instructions du langage push-swap.
var toutes = []string{"sa", "sb", "ss", "pa", "pb", "ra", "rb", "rr", "rra", "rrb", "rrr"}

// Trier renvoie la liste d'instructions qui trie a.
func Trier(nombres []int) []string {
	if croissant(nombres) {
		return []string{} // déjà trié : 0 instruction
	}
	r := rangs(nombres)
	if len(r) <= limiteRecherche {
		return plusCourt(r)
	}
	return optimiser(turk(r))
}

// rangs remplace chaque nombre par sa place dans l'ordre trié (0 = le plus petit).
// L'ordre est le même, mais les calculs sont plus simples : 3 87 -4 → 1 2 0.
func rangs(nombres []int) []int {
	tries := append([]int{}, nombres...)
	sort.Ints(tries)
	place := make(map[int]int, len(tries))
	for i, n := range tries {
		place[n] = i
	}
	r := make([]int, len(nombres))
	for i, n := range nombres {
		r[i] = place[n]
	}
	return r
}

// ---------- Petites listes : la solution la plus courte ----------

// cle écrit l'état des deux piles sous forme de texte, pour retenir les états déjà vus.
func cle(p *Piles) string {
	b := make([]byte, 0, len(p.A)+len(p.B)+1)
	for _, n := range p.A {
		b = append(b, byte(n))
	}
	b = append(b, 255) // séparateur entre a et b
	for _, n := range p.B {
		b = append(b, byte(n))
	}
	return string(b)
}

// depuisCle refait les deux piles à partir du texte.
func depuisCle(k string) *Piles {
	p := &Piles{A: []int{}, B: []int{}}
	dansB := false
	for i := 0; i < len(k); i++ {
		switch {
		case k[i] == 255:
			dansB = true
		case dansB:
			p.B = append(p.B, int(k[i]))
		default:
			p.A = append(p.A, int(k[i]))
		}
	}
	return p
}

// plusCourt fait un parcours en largeur : on essaie toutes les suites de
// 1 instruction, puis de 2, puis de 3... La première suite qui trie est
// forcément la plus courte. Pour 6 nombres, il n'y a que 6! × 7 = 5 040 états.
func plusCourt(r []int) []string {
	trie := make([]int, len(r))
	for i := range trie {
		trie[i] = i
	}
	depart, but := cle(Nouvelles(r)), cle(Nouvelles(trie))

	type venue struct{ avant, instruction string }
	vus := map[string]venue{depart: {}}
	file := []string{depart}

	for len(file) > 0 {
		k := file[0]
		file = file[1:]
		if k == but {
			// On remonte le chemin depuis le but jusqu'au départ.
			solution := []string{}
			for k != depart {
				solution = append(solution, vus[k].instruction)
				k = vus[k].avant
			}
			for i, j := 0, len(solution)-1; i < j; i, j = i+1, j-1 {
				solution[i], solution[j] = solution[j], solution[i]
			}
			return solution
		}
		for _, ins := range toutes {
			p := depuisCle(k)
			p.Executer(ins)
			suivant := cle(p)
			if _, deja := vus[suivant]; !deja {
				vus[suivant] = venue{k, ins}
				file = append(file, suivant)
			}
		}
	}
	return nil // impossible : on peut toujours trier
}

// ---------- Grandes listes : l'algorithme « Turk » ----------

// trieur exécute les instructions et les note dans un journal.
type trieur struct {
	p       *Piles
	journal []string
}

func (t *trieur) faire(ins string, fois int) {
	for range fois {
		t.p.Executer(ins)
		t.journal = append(t.journal, ins)
	}
}

// turk :
//  1. envoie tout dans b sauf 3 nombres, que l'on trie dans a ;
//  2. ramène les nombres de b un par un, en choisissant à chaque fois celui
//     qui demande le moins de rotations pour arriver à sa place dans a ;
//  3. tourne a pour mettre le plus petit en haut.
func turk(r []int) []string {
	t := &trieur{p: Nouvelles(r)}
	milieu := len(r) / 2
	for len(t.p.A) > 3 {
		t.faire("pb", 1)
		// Les petits nombres vont au fond de b, les grands restent en haut :
		// b est grossièrement rangé, ce qui rend l'étape 2 moins chère.
		if t.p.B[0] < milieu && len(t.p.B) > 1 {
			t.faire("rb", 1)
		}
	}
	t.trierTrois()
	for len(t.p.B) > 0 {
		t.ramenerMoinsCher()
	}
	t.minimumEnHaut()
	return t.journal
}

// trierTrois trie les 3 nombres de a en 2 instructions au plus.
func (t *trieur) trierTrois() {
	a := t.p.A
	switch {
	case a[0] > a[1] && a[0] > a[2]: // le plus grand est en haut
		t.faire("ra", 1)
	case a[1] > a[0] && a[1] > a[2]: // le plus grand est au milieu
		t.faire("rra", 1)
	}
	if t.p.A[0] > t.p.A[1] {
		t.faire("sa", 1)
	}
}

// cible : la position dans a où x doit être inséré, c'est-à-dire la place
// du plus petit nombre plus grand que x (ou du minimum s'il n'y en a pas).
func cible(a []int, x int) int {
	meilleur, min := -1, 0
	for j, v := range a {
		if v > x && (meilleur == -1 || v < a[meilleur]) {
			meilleur = j
		}
		if v < a[min] {
			min = j
		}
	}
	if meilleur == -1 {
		return min
	}
	return meilleur
}

// ramenerMoinsCher calcule, pour chaque nombre de b, combien d'instructions
// il faut pour le mettre en haut de b ET mettre sa place en haut de a.
// Il y a 4 façons de faire (haut/haut, bas/bas, haut/bas, bas/haut) :
// on garde la moins chère, puis on fait pa.
func (t *trieur) ramenerMoinsCher() {
	la, lb := len(t.p.A), len(t.p.B)
	meilleurCout, mi, mj, mode := -1, 0, 0, 0
	for i, x := range t.p.B {
		j := cible(t.p.A, x)
		couts := [4]int{
			max(i, j),       // 0 : rr puis ra ou rb
			max(lb-i, la-j), // 1 : rrr puis rra ou rrb
			i + (la - j),    // 2 : rb et rra
			(lb - i) + j,    // 3 : rrb et ra
		}
		for m, c := range couts {
			if meilleurCout == -1 || c < meilleurCout {
				meilleurCout, mi, mj, mode = c, i, j, m
			}
		}
	}

	switch mode {
	case 0:
		commun := min(mi, mj)
		t.faire("rr", commun)
		t.faire("rb", mi-commun)
		t.faire("ra", mj-commun)
	case 1:
		bi, bj := lb-mi, la-mj
		commun := min(bi, bj)
		t.faire("rrr", commun)
		t.faire("rrb", bi-commun)
		t.faire("rra", bj-commun)
	case 2:
		t.faire("rb", mi)
		t.faire("rra", la-mj)
	case 3:
		t.faire("rrb", lb-mi)
		t.faire("ra", mj)
	}
	t.faire("pa", 1)
}

// minimumEnHaut tourne a dans le sens le plus court pour mettre le plus petit en haut.
func (t *trieur) minimumEnHaut() {
	m := 0
	for j, v := range t.p.A {
		if v < t.p.A[m] {
			m = j
		}
	}
	la := len(t.p.A)
	if m <= la/2 {
		t.faire("ra", m)
	} else {
		t.faire("rra", la-m)
	}
}

// ---------- Dernière passe : supprimer les instructions inutiles ----------

// annulent : deux instructions de suite qui se compensent.
var annulent = map[[2]string]bool{
	{"ra", "rra"}: true, {"rra", "ra"}: true,
	{"rb", "rrb"}: true, {"rrb", "rb"}: true,
	{"rr", "rrr"}: true, {"rrr", "rr"}: true,
	{"sa", "sa"}: true, {"sb", "sb"}: true, {"ss", "ss"}: true,
}

// fusionnent : deux instructions de suite qui n'en font qu'une.
var fusionnent = map[[2]string]string{
	{"ra", "rb"}: "rr", {"rb", "ra"}: "rr",
	{"rra", "rrb"}: "rrr", {"rrb", "rra"}: "rrr",
	{"sa", "sb"}: "ss", {"sb", "sa"}: "ss",
}

// optimiser supprime les paires inutiles et fusionne ce qui peut l'être.
func optimiser(instructions []string) []string {
	res := []string{}
	var ajouter func(string)
	ajouter = func(ins string) {
		if n := len(res); n > 0 {
			paire := [2]string{res[n-1], ins}
			if annulent[paire] {
				res = res[:n-1]
				return
			}
			if f, ok := fusionnent[paire]; ok {
				res = res[:n-1]
				ajouter(f)
				return
			}
		}
		res = append(res, ins)
	}
	for _, ins := range instructions {
		ajouter(ins)
	}
	return res
}
