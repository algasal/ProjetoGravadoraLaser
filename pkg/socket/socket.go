// pkg/socket/socket.go
package socket

import (
	"fmt"
	"net"
	"strings"
	"time"
	"os"
)

// Conectar estabelece uma conexão TCP com o host especificado
func Conectar(host string, porta int) (net.Conn, error) {
	endereco := fmt.Sprintf("%s:%d", host, porta)
	conn, err := net.DialTimeout("tcp", endereco, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("Erro no socket: %w", err)
	}
	fmt.Printf("Conexão estabelecida com %s\n", endereco)
	return conn, nil
}

// EnviarRequisicao envia uma requisição pelo socket com os comandos para a gravadora.
func EnviarRequisicao(conn net.Conn, classificacao string) {
	programas := map[string]string{
		"A": "0008",
		"B": "0009",
		"C": "0010",
		"D": "0011",
		"E": "0012",
	}

	programa, ok := programas[classificacao]
	if !ok {
		fmt.Fprintf(os.Stderr, "Classificação inválida: %s\n", classificacao)
		return
	}

	// Alterar o programa da gravadora
	comando_programa := fmt.Sprintf("WX,ProgramNo=%s\r\n", programa)
	if _, err := conn.Write([]byte(comando_programa)); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao definir programa: %v\n", err)
		return
	}
	
	// Ler resposta da requisição
	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta do socket: %v\n", err)
		return
	}
	
	resposta := strings.TrimSpace(string(buffer[:n]))
	if resposta == "WX,OK" {
		fmt.Printf("Definido o programa %s.\nIniciando marcação.", programa)
	} else {
		fmt.Fprintf(os.Stderr, "Erro na definição do programa: %s\nConsulte o manual\n", resposta)
	}
	
	// Enviar o comando para iniciar marcação
	comando_marcar := fmt.Sprintf("WX,StartMarking\r\n")
	if _, err := conn.Write([]byte(comando_marcar)); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar marcação: %v\n", err)
		return
	}

	// Ler a resposta. da requisição
	//buffer := make([]byte, 256)
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta do socket: %v\n", err)
		return
	}

	resposta = strings.TrimSpace(string(buffer[:n]))
	if resposta == "WX,OK" {
		fmt.Println("Gravação concluída")
	} else {
		fmt.Fprintf(os.Stderr, "Erro na gravação: %s\nConsulte o manual\n", resposta)
	}
	conn.Close()
}
