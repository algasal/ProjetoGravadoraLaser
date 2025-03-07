// pkg/serialport/serialport.go
package serialport

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

// SerialPortWrapper cria um wrapper para um objeto *serial.Port
type SerialPortWrapper struct {
	Port *serial.Port
}

// AbrirPorta abre a comunicação com uma porta serial.
func AbrirPorta(nome string, baud int) (*SerialPortWrapper, error) {
	config := &serial.Config{Name: nome, Baud: baud}
	porta, err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("Falha ao abrir a porta serial: %w", err)
	}
	return &SerialPortWrapper{Port:porta}, nil
}

// LerLado solicita uma medida da porta serial
func LerLado(port *serial.Port, lado int) (float64, error) {
	serialReader := bufio.NewReader(port)
	trigger := "?"
	 _, err := port.Write([]byte(trigger))
	 if err != nil {
		return 0, fmt.Errorf("Erro ao escrever na porta serial: %w", err)
	}
	resposta, err := serialReader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("Erro ao ler da porta serial: %w", err)
	}
	trimmed := strings.TrimSpace(resposta)
	valor, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, fmt.Errorf("Erro na conversão da medida do lado %d: %w", lado, err)
	}
	return valor, nil
}
