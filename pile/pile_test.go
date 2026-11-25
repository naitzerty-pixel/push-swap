package pile

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// verifier exécute les instructions et dit si elles trient bien les nombres.
func verifier(nombres []int, instructions []string) bool {
	p := Nouvelles(nombres)
	for _, ins := range instructions {
		if p.Executer(ins) != nil {
			return false
		}
	}
	return p.EstTriee()
}

func TestLire(t *testing.T) {
	if n, err := Lire([]string{"2 1 3"}); err != nil || !slices.Equal(n, []int{2, 1, 3}) {
		t.Errorf(`"2 1 3" : obtenu %v, %v`, n, err)
	}
	if n, err := Lire([]string{"2", "-1", "3"}); err != nil || !slices.Equal(n, []int{2, -1, 3}) {
		t.Errorf(`2 -1 3 : obtenu %v, %v`, n, err)
	}
	if n, err := Lire(nil); err != nil || len(n) != 0 {
		t.Errorf("sans argument : obtenu %v, %v", n, err)
	}
	for _, mauvais := range []string{"0 one 2 3", "1 2 2", "1.5 2", "99999999999999999999", "1 +-2"} {
		if _, err := Lire([]string{mauvais}); err == nil {
			t.Errorf("%q aurait dû donner une erreur", mauvais)
		}
	}
}

// L'exemple du sujet, étape par étape.
func TestInstructionsExempleDuSujet(t *testing.T) {
	p := Nouvelles([]int{2, 1, 3, 6, 8, 5})
	etapes := []struct {
		instructions []string
		a, b         []int
	}{
		{[]string{"sa"}, []int{1, 2, 3, 6, 8, 5}, []int{}},
		{[]string{"pb", "pb", "pb"}, []int{6, 8, 5}, []int{3, 2, 1}},
		{[]string{"rb"}, []int{6, 8, 5}, []int{2, 1, 3}},
		{[]string{"rrr"}, []int{5, 6, 8}, []int{3, 2, 1}},
		{[]string{"pa", "pa", "pa"}, []int{1, 2, 3, 5, 6, 8}, []int{}},
	}
	for _, e := range etapes {
		for _, ins := range e.instructions {
			if err := p.Executer(ins); err != nil {
				t.Fatal(err)
			}
		}
		if !slices.Equal(p.A, e.a) || !slices.Equal(p.B, e.b) {
			t.Fatalf("après %v : a=%v b=%v, attendu a=%v b=%v", e.instructions, p.A, p.B, e.a, e.b)
		}
	}
	if !p.EstTriee() {
		t.Error("la pile devrait être triée")
	}
}

func TestInstructionInconnue(t *testing.T) {
	p := Nouvelles([]int{1, 2})
	for _, mauvaise := range []string{"rra ", "RA", "pc", ""} {
		if p.Executer(mauvaise) == nil {
			t.Errorf("%q aurait dû être refusée", mauvaise)
		}
	}
}

// Les instructions sur une pile vide ou trop petite ne plantent pas.
func TestPilesVides(t *testing.T) {
	p := Nouvelles([]int{})
	for _, ins := range toutes {
		p.Executer(ins)
	}
	p = Nouvelles([]int{7})
	for _, ins := range toutes {
		p.Executer(ins)
	}
}

func TestDejaTrie(t *testing.T) {
	if n := len(Trier([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})); n != 0 {
		t.Errorf("liste déjà triée : %d instructions, 0 attendue", n)
	}
}

// L'exemple du sujet : 6 instructions pour "4 67 3 87 23".
func TestExempleUsage(t *testing.T) {
	nombres := []int{4, 67, 3, 87, 23}
	ins := Trier(nombres)
	if !verifier(nombres, ins) || len(ins) > 6 {
		t.Errorf("%v : %d instructions %v", nombres, len(ins), ins)
	}
}

// Toutes les permutations de 1 à 6 nombres sont triées.
func TestToutesLesPetitesListes(t *testing.T) {
	for n := 1; n <= 6; n++ {
		pire := 0
		permutations(n, func(p []int) {
			ins := Trier(p)
			if !verifier(p, ins) {
				t.Fatalf("%v n'est pas trié par %v", p, ins)
			}
			pire = max(pire, len(ins))
		})
		t.Logf("%d nombres : %d instructions au pire", n, pire)
	}
}

// Les grandes listes (7, 100, 500 nombres) sont triées, avec un nombre
// d'instructions raisonnable.
func TestGrandesListes(t *testing.T) {
	cas := []struct{ taille, limite int }{
		{7, 30},
		{100, 700},
		{500, 5500},
	}
	r := rand.New(rand.NewPCG(1, 2))
	for _, c := range cas {
		total, pire := 0, 0
		for essai := 0; essai < 20; essai++ {
			nombres := r.Perm(c.taille)
			for i := range nombres {
				nombres[i] = nombres[i]*3 - 1000 // pas seulement 0..n-1
			}
			ins := Trier(nombres)
			if !verifier(nombres, ins) {
				t.Fatalf("%d nombres : la liste n'est pas triée", c.taille)
			}
			total += len(ins)
			pire = max(pire, len(ins))
		}
		t.Logf("%d nombres : %d instructions en moyenne, %d au pire", c.taille, total/20, pire)
		if pire > c.limite {
			t.Errorf("%d nombres : %d instructions au pire, limite %d", c.taille, pire, c.limite)
		}
	}
}

func TestOptimiser(t *testing.T) {
	avant := []string{"ra", "rb", "pb", "ra", "rra", "rra", "rrb", "sa", "sa", "pa"}
	apres := []string{"rr", "pb", "rrr", "pa"}
	if res := optimiser(avant); !slices.Equal(res, apres) {
		t.Errorf("obtenu %v, attendu %v", res, apres)
	}
}

// permutations appelle f avec chaque ordre possible de 0..n-1.
func permutations(n int, f func([]int)) {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	var generer func(k int)
	generer = func(k int) {
		if k == n {
			f(append([]int{}, p...))
			return
		}
		for i := k; i < n; i++ {
			p[k], p[i] = p[i], p[k]
			generer(k + 1)
			p[k], p[i] = p[i], p[k]
		}
	}
	generer(0)
}
