package persist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"

	"github.com/elastic/go-elasticsearch/esapi"

	"github.com/elastic/go-elasticsearch"
)

func ItemSaver() (chan engine.Item, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got item #%d: %v", itemCount, item)
			itemCount++
			err := Save(client, item)
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}
		}
	}()
	return out, nil
}

func Save(client *elasticsearch.Client, item engine.Item) (err error) {
	jsonItem, err := json.Marshal(item)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:        item.Index,
		DocumentType: "",
		DocumentID:   item.Id,
		Body:         bytes.NewReader(jsonItem),
	}
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	fmt.Println(res.String())
	return nil
}
