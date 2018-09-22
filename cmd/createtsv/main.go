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
	"github.com/yagotome/gcp-nlp-ner-consumer/utils/stringutil"
)

var (
	progress        = 0
	sentencesAmount = 0
)

func main() {
	log.Println("Starting processing...")
	gss := getGroupedSentences()

	groupedSentences := map[string][]string{}
	for k, v := range gss {
		// groupedSentences[k] = v[:2]
		groupedSentences[k] = v
		sentencesAmount += len(groupedSentences[k])
	}

	for group, sentences := range groupedSentences {
		createTsv(group, sentences)
	}
	log.Println("Finished processing")
}

func createTsv(group string, sentences []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(sentences))
	mutex := sync.Mutex{}

	f, err := os.Create(fmt.Sprintf("output/%v.tsv", stringutil.EscapeGroup(group)))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	ef, err := os.Create(fmt.Sprintf("error/%v.tsv", stringutil.EscapeGroup(group)))
	if err != nil {
		log.Fatal(err)
	}
	defer ef.Close()
	ef2, err := os.Create(fmt.Sprintf("error/%v.csv", stringutil.EscapeGroup(group)))
	if err != nil {
		log.Fatal(err)
	}
	ef2.WriteString("assunto,relato\n")
	defer ef2.Close()

	for _, text := range sentences {
		go func(text string) {
			defer wg.Done()
			entities, err := nlp.GetNamedEntities(text)
			// err = fmt.Errorf("error fake")
			if err != nil {
				log.Println("Error on sentence:", text+"\n", err)
				ef2.WriteString(fmt.Sprintf("%v,%v\n", group, `"`+strings.Replace(text, `"`, `""`, -1)+`"`))
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
					if err != nil {
						ef.WriteString(words[i] + "\t" + "O" + "\n")
					} else if entity := entities.Match(words, i); entity != nil {
						f.WriteString(entity.Name + "\t" + entity.Type + "\n")
						entityWords := strings.Split(entity.Name, " ")
						i += len(entityWords) - 1
					} else {
						f.WriteString(words[i] + "\t" + "O" + "\n")
					}
				}
				if err != nil {
					ef.WriteString("\n")
				} else {
					f.WriteString("\n")
				}
			}
			progress++
			log.Printf("%.2f%% - %d de %d", 100*float64(progress)/float64(sentencesAmount), progress, sentencesAmount)
			mutex.Unlock()
		}(text)
		time.Sleep(150 * time.Millisecond)
	}
	wg.Wait()
}

func getGroupedSentences() map[string][]string {
	subjectFilter := []string{"Roubos em Geral", "Roubo carga/ veículo", "Tráfico de drogas", `Tráfico de Drogas \ Armas`, "Homicidio", "Homicídios", "Armas"}
	groupedSentences, err := csv.ExtractColumnGroupedBy("app_dd.csv", "relato", "assunto", subjectFilter)
	// groupedSentences, err := csv.ExtractColumnGroupedBy("error/Tráfico de drogas.csv", "relato", "assunto", subjectFilter)
	if err != nil {
		log.Fatal(err)
	}
	groupedSentences["Homicídios"] = append(groupedSentences["Homicídios"], groupedSentences["Homicidio"]...)
	delete(groupedSentences, "Homicidio")
	delete(groupedSentences, "")
	return groupedSentences
}
