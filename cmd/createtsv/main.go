package main

import (
	"fmt"
	"log"

	"github.com/yagotome/gcp-nlp-ner-consumer/csv"
)

func main() {
	log.Println("Starting processing...")
	subjectFilter := []string{"Roubos em Geral", "Roubo carga/ veículo", "Tráfico de drogas", `Tráfico de Drogas \ Armas`, "Homicidio", "Homicídios", "Armas"}
	sentencesGrouped, err := csv.ExtractColumnGroupedBy("app_dd.csv", "relato", "assunto", subjectFilter)
	if err != nil {
		log.Fatal(err)
	}
	sentencesGrouped["Homicídios"] = append(sentencesGrouped["Homicídios"], sentencesGrouped["Homicidio"]...)
	delete(sentencesGrouped, "Homicidio")
	delete(sentencesGrouped, "")
	fmt.Println(sentencesGrouped)
}
