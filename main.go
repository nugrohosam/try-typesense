package main

import (
	"fmt"
	"log"

	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func main() {
	app := fiber.New()

	client := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("Hu52dwsas2AdxdE"))

	schema := &api.CollectionSchema{
		Name: "person",
		Fields: []api.Field{
			{
				Name: "name",
				Type: "string",
			},
			{
				Name: "age",
				Type: "int32",
			},
			{
				Name: "country",
				Type: "string",
			},
		},
	}

	client.Collections().Create(schema)

	app.Get("/", func(c *fiber.Ctx) error {
		paramQ := c.Query("q", "Adam")
		quereyBy := "name"
		sortBy := pointer.String("age:desc")

		fmt.Println(paramQ)
		searchParameters := &api.SearchCollectionParams{
			Q:       paramQ,
			QueryBy: quereyBy,
			SortBy:  sortBy,
		}

		results, err := client.Collection("person").Documents().Search(searchParameters)
		if err != nil {
			fmt.Println(err)
			return c.SendString("Error search!")
		}

		fmt.Println(*results.Found)
		fmt.Println(*results.OutOf)
		fmt.Println(*results.Page)
		fmt.Println(*results.SearchTimeMs)

		defer func() {
			hits := *results.Hits
			for _, val := range hits {
				fmt.Println(*val.Document)
			}
		}()

		return c.SendString("Success get ")
	})

	app.Post("/", func(c *fiber.Ctx) error {

		documentFaker := struct {
			ID      string `faker:"uuid_digit"`
			Name    string `faker:"name"`
			Age     int    `faker:"oneof: 15, 27, 61"`
			Country string `faker:"oneof: USA, ID, AUS, GER, NED"`
		}{}

		err := faker.FakeData(&documentFaker)

		document := struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Age     int    `json:"age"`
			Country string `json:"country"`
		}{
			ID:      documentFaker.ID,
			Name:    documentFaker.Name,
			Age:     documentFaker.Age,
			Country: documentFaker.Country,
		}

		if err != nil {
			fmt.Println(err)
			return c.SendString("Error create!")
		}

		fmt.Println(document)
		if _, err := client.Collection("person").Documents().Create(document); err != nil {
			fmt.Println(err)
			return c.SendString("Error create!")
		}

		return c.SendString("Success create!")
	})

	log.Fatal(app.Listen(":3000"))
}
