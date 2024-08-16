package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func checkAndInstallFluentBit() error {
	_, err := exec.LookPath("fluent-bit")
	if err != nil {
		log.Println("FluentBit not found, installing...")
		cmd := exec.Command("brew", "install", "fluent-bit")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install FluentBit: %v", err)
		}
		log.Println("FluentBit installed successfully.")
	}
	return nil
}

func configureFluentBit() error {
	config := `[INPUT]
    Name              tail
    Path              /var/log/myapp.log
    Multiline         On
    Parser_Firstline  my_parser
    DB                /var/log/myapp.db
    Tag               myapp

[PARSER]
    Name        my_parser
    Format      json
    Time_Key    timestamp
    Time_Format %Y-%m-%dT%H:%M:%S

[OUTPUT]
    Name              forward
    Match             myapp
    Host              logstash_host
    Port              5044
`
	configPath := "/etc/fluent-bit/fluent-bit.conf"
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create FluentBit config: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(config)
	if err != nil {
		return fmt.Errorf("failed to write FluentBit config: %v", err)
	}
	log.Println("FluentBit configured successfully.")
	return nil
}

func startFluentBit() error {
	cmd := exec.Command("fluent-bit", "-c", "/etc/fluent-bit/fluent-bit.conf")
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start FluentBit: %v", err)
	}
	log.Println("FluentBit started successfully.")
	return nil
}

func main() {
	err := checkAndInstallFluentBit()
	if err != nil {
		log.Fatalf("Error during FluentBit installation: %v", err)
	}

	err = configureFluentBit()
	if err != nil {
		log.Fatalf("Error during FluentBit configuration: %v", err)
	}

	err = startFluentBit()
	if err != nil {
		log.Fatalf("Error during FluentBit startup: %v", err)
	}
}
