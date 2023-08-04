package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitorings = 3
const delay = 5

func main() {

	showIntroduction()

	for {
		showMenu()

		command := readCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("Logs")
			showLogs()
		case 0:
			fmt.Println("Saindo...")
			os.Exit(0)
		default:
			fmt.Println("Não reconhecido")
		}
	}

}

func showIntroduction() {
	fmt.Println("Olá, Usuário!")
}

func showMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
}

func readCommand() int {
	var command int

	fmt.Scan(&command)

	return command
}

func startMonitoring() {
	fmt.Println("Monitoramento")

	sites := readSitesFromArchive()

	for i := 0; i < monitorings; i++ {
		for _, site := range sites {
			fmt.Println("Testando site:", site)
			testSite(site)
		}
		time.Sleep(delay * time.Second)
	}

}

func testSite(site string) {
	response, err := http.Get(site)

	if err != nil {
		takeError(err)
	}

	if response.StatusCode == 200 {
		fmt.Println("Site:", site, "está ativo!")
		registerLog(site, true)
	} else {
		fmt.Println("Site:", site, "não está ativo. Status Code:", response.StatusCode)
		registerLog(site, false)
	}
}

func readSitesFromArchive() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")

	if err != nil {
		takeError(err)
	}

	reader := bufio.NewReader(arquivo)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()

	return sites
}

func registerLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		takeError(err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func showLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		takeError(err)
	}

	fmt.Println(string(arquivo))
}

func takeError(err error) {
	fmt.Println("ERRO ----->", err)
}
