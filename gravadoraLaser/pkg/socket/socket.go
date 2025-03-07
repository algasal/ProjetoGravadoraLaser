// pkg/socket/socket.go
package socket

import (
	"fmt"
	"net"
	"strings"
	"time"
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

// SelecionarPrograma define qual programa da gravadora será utilizado para a gravação
func SelecionarPrograma(classificacao string) (string, error) {
	programas := map[string]string{
		"A": "0008",
		"B": "0009",
		"C": "0010",
		"D": "0011",
		"E": "0012",
	}

	programa, ok := programas[classificacao]
	if !ok {
		return "", fmt.Errorf("Classificação inválida: %s\n", classificacao)
	}
	
	return programa, nil
}


// EnviarRequisicao envia um comando para a gravadora
func EnviarRequisicao(conn net.Conn, requisicao string) (int, error) {
	n, err := conn.Write([]byte(requisicao))
	if err != nil {
		return 0, fmt.Errorf("Erro ao definir programa: %v\n", err)
	}

	return n, nil
}

// LerResposta lê a resposta da gravadora para o comando executado
func LerResposta(conn net.Conn) string {
	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Sprintf("Erro ao ler a resposta do socket: %v\n", err)
	}

	resposta := strings.TrimSpace(string(buffer[:n]))
	return resposta
}
