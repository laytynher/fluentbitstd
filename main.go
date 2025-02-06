package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func verificarEInstalarFluentBit() error {
	_, err := exec.LookPath("fluent-bit")
	if err != nil {
		log.Println("FluentBit não encontrado, instalando...")
		cmd := exec.Command("brew", "install", "fluent-bit")
		erro := cmd.Run()
		if erro != nil {
			return fmt.Errorf("falha ao instalar FluentBit: %v", erro)
		}
		log.Println("FluentBit instalado com sucesso.")
	}
	return nil
}

func configurarFluentBit() error {
	configuracao := `[INPUT]
    Name              tail
    Path              /var/log/myapp.log
    Multiline         On
    Parser_Firstline  meu_parser
    DB                /var/log/myapp.db
    Tag               myapp

[PARSER]
    Name        meu_parser
    Format      json
    Time_Key    timestamp
    Time_Format %Y-%m-%dT%H:%M:%S

[OUTPUT]
    Name              forward
    Match             myapp
    Host              logstash_host
    Port              5044
`
	caminhoConfig := "/etc/fluent-bit/fluent-bit.conf"
	arquivo, err := os.Create(caminhoConfig)
	if err != nil {
		return fmt.Errorf("falha ao criar configuração do FluentBit: %v", err)
	}
	defer arquivo.Close()

	_, err = arquivo.WriteString(configuracao)
	if err != nil {
		return fmt.Errorf("falha ao escrever configuração do FluentBit: %v", err)
	}
	log.Println("FluentBit configurado com sucesso.")
	return nil
}

func iniciarFluentBit() error {
	cmd := exec.Command("fluent-bit", "-c", "/etc/fluent-bit/fluent-bit.conf")
	erro := cmd.Start()
	if erro != nil {
		return fmt.Errorf("falha ao iniciar FluentBit: %v", erro)
	}
	log.Println("FluentBit iniciado com sucesso.")
	return nil
}

func main() {
	erro := verificarEInstalarFluentBit()
	if erro != nil {
		log.Fatalf("Erro durante a instalação do FluentBit: %v", erro)
	}

	erro = configurarFluentBit()
	if erro != nil {
		log.Fatalf("Erro durante a configuração do FluentBit: %v", erro)
	}

	erro = iniciarFluentBit()
	if erro != nil {
		log.Fatalf("Erro ao iniciar o FluentBit: %v", erro)
	}
}
