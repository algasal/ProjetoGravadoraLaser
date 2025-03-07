// pkg/medida/medida.go
package medida

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gravadoraLaser/pkg/serialport"
)

// SelecionarPeca pede ao usuário que selecione a peça a ser medida
func SelecionarPeca() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(`Selecione a peça a ser medida:
    1- E1531
    2- E1536
> `)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Erro ao ler a entrada:", err)
			continue
		}
		input = strings.TrimSpace(input)
		switch input {
		case "1":
			return "E1531"
		case "2":
			return "E1536"
		default:
			fmt.Fprintln(os.Stderr, "Peça não encontrada, digite novamente.")
		}
	}
}

// MedirPeca mede os três lados e retorna a maior medida
func MedirPeca(wrapper *serialport.SerialPortWrapper) (float64, error) {
	var lados []float64
	for i := 1; i <= 3; i++ {
		fmt.Printf("Pressione '?' para ler o lado %d\n> ", i)
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		medida, err := serialport.LerLado(wrapper.Port, i)
		if err != nil {
			return 0, err
		}
		lados = append(lados, medida)
		fmt.Printf("Lado %d: %f\n\n", i, medida)
	}

	max := lados[0]
	for _, lado := range lados {
		if lado > max {
			max = lado
		}
	}
	fmt.Printf("Maior lado: %f\n", max)
	return max, nil
}
