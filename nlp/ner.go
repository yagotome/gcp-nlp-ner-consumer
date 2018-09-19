package nlp

import (
	"encoding/json"
	"log"
	"strings"

	// Imports the Google Cloud Natural Language API client package.
	language "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Entities ...
type Entities []*Entity

// Entity ...
type Entity struct {
	Name     string
	Type     string
	Salience float32
}

// GetNamedEntities returns the named entities present in text
func GetNamedEntities(text string) (Entities, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v\n", err)
		return nil, err
	}

	// Request entities to GCP
	response, err := analyzeEntities(ctx, client, text)
	if err != nil {
		log.Printf("Failed to analyze text: %v", err)
		return nil, err
	}

	// return response.Entities, nil
	entities := []*Entity{}
	for _, entity := range response.Entities {
		entities = append(entities, &Entity{
			Name:     entity.Name,
			Type:     entity.Type.String(),
			Salience: entity.Salience,
		})
	}

	return entities, nil
}

func analyzeEntities(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeEntitiesResponse, error) {
	return client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

// Match ...
func (entities Entities) Match(words []string, wrdIdx int) *Entity {
	for _, entity := range entities {
		entityWords := strings.Split(entity.Name, " ")
		if len(entityWords) == 0 {
			return nil
		}
		i, j := wrdIdx, 0
		for ; i < len(words) && j < len(entityWords); i, j = i+1, j+1 {
			if strings.ToUpper(words[i]) != strings.ToUpper(entityWords[j]) {
				break
			}
		}
		if j == len(entityWords) {
			return entity
		}
	}
	return nil
}

func (entities Entities) String() string {
	j, err := json.Marshal(entities)
	if err != nil {
		return ""
	}
	return string(j)
}
