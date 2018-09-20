package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/euskadi31/go-tokenizer"
	"github.com/yagotome/gcp-nlp-ner-consumer/csv"
	"github.com/yagotome/gcp-nlp-ner-consumer/nlp"
)

func main() {
	log.Println("Starting processing...")
	groupedSentences := getGroupedSentences()

	groupedSentences = map[string][]string{"Roubos em Geral": []string{"Roubaram meu carro. Na Barra da Tijuca, em frente ao Barra Shopping! Incidente ocorreu hoje de manhã, 10:30."}}

	for group, sentences := range groupedSentences {
		createTsv(group, sentences)
	}
	log.Println("Finished processing")
}

func createTsv(group string, sentences []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(sentences))
	mutex := sync.Mutex{}

	f, err := os.Create(fmt.Sprintf("output/%v.tsv", group))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, text := range sentences {
		go func(text string) {
			defer wg.Done()
			entities, err := nlp.GetNamedEntities(text)
			if err != nil {
				log.Println(err)
			}
			// cleanText := stringutil.RemovePunctuation(text)
			// phrases := stringutil.SplitByPunctuation(cleanText)
			cleanText := text
			phrases := []string{cleanText}

			t := tokenizer.New()
			t.KeepSeparator()
			r := regexp.MustCompile(`^[\s\t\r\n]*$`)
			mutex.Lock()
			for _, phrase := range phrases {
				words := t.Tokenize(phrase)
				// words := strings.Split(phrase, " ")
				for i := 0; i < len(words); i++ {
					if r.MatchString(words[i]) {
						continue
					}
					if entity := entities.Match(words, i); entity != nil {
						f.WriteString(entity.Name + "\t" + entity.Type + "\n")
						entityWords := strings.Split(entity.Name, " ")
						i += len(entityWords) - 1
					} else {
						f.WriteString(words[i] + "\t" + "O" + "\n")
					}
				}
				f.WriteString("\n")
			}
			mutex.Unlock()
		}(text)
		time.Sleep(time.Millisecond)
	}

	wg.Wait()
}

func getGroupedSentences() map[string][]string {
	subjectFilter := []string{"Roubos em Geral", "Roubo carga/ veículo", "Tráfico de drogas", `Tráfico de Drogas \ Armas`, "Homicidio", "Homicídios", "Armas"}
	groupedSentences, err := csv.ExtractColumnGroupedBy("app_dd.csv", "relato", "assunto", subjectFilter)
	if err != nil {
		log.Fatal(err)
	}
	groupedSentences["Homicídios"] = append(groupedSentences["Homicídios"], groupedSentences["Homicidio"]...)
	delete(groupedSentences, "Homicidio")
	delete(groupedSentences, "")
	return groupedSentences
}
