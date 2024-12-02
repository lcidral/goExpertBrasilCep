package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

// apis de consulta - prod
const (
	urlBrasil = "https://brasilapi.com.br/api/cep/v1/%s"
	urlViaCep = "http://viacep.com.br/ws/%s/json/"

	BRASIL_API = "BrasilAPI"
	VIA_CEP    = "ViaCep"
)

// struct principal pra retorno padronizado
type Endereco struct {
	Cep, Estado, Cidade, Rua, Bairro string
}

// brasil api
type BrasilAPIResp struct {
	Cep    string `json:"cep"`
	State  string `json:"state"`
	City   string `json:"city"`
	Street string `json:"street"`
	Area   string `json:"neighborhood"` // aka bairro
}

// via cep
type ViaCEPResp struct {
	Cep    string `json:"cep"`
	Rua    string `json:"logradouro"`
	Bairro string `json:"bairro"`
	Cidade string `json:"localidade"`
	Uf     string `json:"uf"`
}

type apiResult struct {
	endereco Endereco
	api      string
	err      error
}

// converters para normalizar a resposta
func (resp *BrasilAPIResp) toEndereco() Endereco {
	return Endereco{
		Cep:    resp.Cep,
		Estado: resp.State,
		Cidade: resp.City,
		Rua:    resp.Street,
		Bairro: resp.Area,
	}
}

func (resp *ViaCEPResp) toEndereco() Endereco {
	return Endereco{
		Cep:    resp.Cep,
		Estado: resp.Uf,
		Cidade: resp.Cidade,
		Rua:    resp.Rua,
		Bairro: resp.Bairro,
	}
}

// tenta brasil api primeiro
func buscaBrasilApi(cep string, result chan<- apiResult) {
	url := fmt.Sprintf(urlBrasil, cep)
	resp, err := http.Get(url)
	if err != nil {
		result <- apiResult{Endereco{}, BRASIL_API, err}
		return
	}
	defer resp.Body.Close()

	// resposta
	buf, _ := io.ReadAll(resp.Body)

	// parse
	var dados BrasilAPIResp
	if err = json.Unmarshal(buf, &dados); err != nil {
		result <- apiResult{Endereco{}, BRASIL_API, err}
		return
	}

	result <- apiResult{dados.toEndereco(), BRASIL_API, nil}
}

// tenta via cep em paralelo
func buscaViaCep(cep string, result chan<- apiResult) {
	url := fmt.Sprintf(urlViaCep, cep)
	resp, err := http.Get(url)
	if err != nil {
		result <- apiResult{Endereco{}, VIA_CEP, err}
		result <- apiResult{Endereco{}, VIA_CEP, err}
		return
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)
	var dados ViaCEPResp

	if err = json.Unmarshal(buf, &dados); err != nil {
		result <- apiResult{Endereco{}, VIA_CEP, err}
		return
	}

	result <- apiResult{dados.toEndereco(), VIA_CEP, nil}
}

func buscaCep(cep string) {
	// canais pros resultados
	results := make(chan apiResult, 2) // guarda os 2 resultados
	done := make(chan bool)            // sync

	// chama as duas APIs
	go buscaBrasilApi(cep, results)
	go buscaViaCep(cep, results)

	// pega o mais rapido
	select {
	case r := <-results:
		if r.err != nil {
			fmt.Printf("erro %v\n", r.err)
			return
		}

		// formata saida
		fmt.Printf("\nAPI......: %s\n", r.api)
		fmt.Printf("CEP......: %s\n", r.endereco.Cep)
		fmt.Printf("Estado...: %s\n", r.endereco.Estado)
		fmt.Printf("Cidade...: %s\n", r.endereco.Cidade)
		fmt.Printf("Rua......: %s\n", r.endereco.Rua)
		fmt.Printf("Bairro...: %s\n\n", r.endereco.Bairro)

		// limpa segunda resposta
		go func() {
			<-results // descarta
			done <- true
		}()

	case <-time.After(1 * time.Second): // timeout
		fmt.Println("demorou demais :(")
	}

	// espera cleanup
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond): // da um tempinho
	}
}

func main() {
	// flags
	cepFlag := flag.String("cep", "", "cep para buscar")
	flag.Parse()

	if *cepFlag == "" {
		fmt.Println("Informe um cep, use -cep")
		return
	}

	buscaCep(*cepFlag)
}
