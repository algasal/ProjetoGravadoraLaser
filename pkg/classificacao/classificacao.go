// pkg/classificacao/classificacao.go
package classificacao

// Categoria define um valor mínimo e máximo para as dimensões da peça, e uma label para a categoria.
type Categoria struct {
	Min, Max  float64
	Label     string
}

// Dimensões define os limites para as medidas de cada peça
var Dimensoes = map[string][]Categoria{
	"E1531": {
		{16.270, 16.280, "E"},
		{16.280, 16.290, "D"},
		{16.290, 16.300, "C"},
		{16.300, 16.310, "B"},
	},
	"E1536": {
		{16.270, 16.280, "A"},
		{16.280, 16.290, "B"},
		{16.290, 16.300, "C"},
		{16.300, 16.310, "D"},
	},
}

// Classificar retorna a Categoria com base na peça escolhida e no maior lado medido
func Classificar(peca string, maiorLado float64) string {
	for _, cat := range Dimensoes[peca] {
		if maiorLado >= cat.Min && maiorLado < cat.Max {
			return cat.Label
		}
	}
	return "Classificação inválida, dimensão fora dos limites"
}
