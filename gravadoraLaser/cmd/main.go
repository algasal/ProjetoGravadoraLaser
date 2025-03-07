// cmd/main.go
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
		fmt.Fprintln(os.Stderr, "Erro na medição. Não foi possível classificar a peça")
		os.Exit(1)
	}

	// Classifica a peca.
	label := classificacao.Classificar(peca, maiorLado)
	fmt.Printf("classificacao: %s\n", label)
	
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
	defer conn.Close()

	// Seleciona o programa a ser executado na gravadora
	programa, err := socket.SelecionarPrograma(label)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Altera o programa da gravadora
	comando_programa := fmt.Sprintf("WX,ProgramNo=%s\r\n", programa)
	_, err = socket.EnviarRequisicao(conn, comando_programa)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Lê a resposta da alteração do programa
	resposta_alteracao := socket.LerResposta(conn)
	if resposta_alteracao == "WX,OK" {
		fmt.Printf("Definido o programa %s.\nIniciando marcação.", programa)
	} else {
		fmt.Fprintf(os.Stderr, "Erro na definição do programa: %s\nConsulte o manual", resposta_alteracao)
		os.Exit(1)
	}

	// Envia o comando para iniciar a marcação
	comando_marcacao := fmt.Sprintf("WX,StartMarking\r\n")
	_, err = socket.EnviarRequisicao(conn, comando_marcacao)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Lê a resposta da marcação
	resposta_marcacao := socket.LerResposta(conn)
	if resposta_marcacao == "WX,OK" {
		fmt.Println("Marcação concluída com sucesso.")
	} else {
		fmt.Fprintf(os.Stderr, "Erro na marcação: %s\nConsulte o manual\n", resposta_marcacao)
		os.Exit(1)
	}

}
