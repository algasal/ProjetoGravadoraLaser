// cmd/gravadoraLaser/main.go
package main

import (
	"fmt"
	"os"
	"bufio"
	
	"gravadoraLaser/pkg/classificacao"
	"gravadoraLaser/pkg/medida"
	"gravadoraLaser/pkg/serialport"
	"gravadoraLaser/pkg/socket"
)

func main() {
	// Abre a comunicação com a porta serial
	serialWrapper, err := serialport.AbrirPorta("COM7", 9600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro na comunicação serial: %v\n", err)
		os.Exit(1)
	}
	defer serialWrapper.Port.Close()

	// Seleciona a peça a ser medida
	peca := medida.SelecionarPeca()
	fmt.Printf("\nFazendo a medição da peça %s\n\n", peca)

	// Mede a peca.
	maiorLado, err := medida.MedirPeca(serialWrapper)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro na medição. Não foi possível classeificar a peça")
		os.Exit(1)
	}

	// Classifica a peca.
	classe := classificacao.Classificar(peca, maiorLado)
	fmt.Printf("classificacao: %s\n", classe)
	
	//Aguardar o posicionamento da peça na gravadora
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nAbra o Finder no Marking Builder para posicionar a peça na gravadora.")
	fmt.Printf("Aperte enter quando estiver pronto para iniciar a gravação.")
	_, err = reader.ReadString('\n')
	
	// Conectar na gravadora pelo socket TCP
	conn, err := socket.Conectar("192.168.1.3", 50002)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Envia a instrução para a gravadora pelo socket
	socket.EnviarRequisicao(conn, classe)
}
