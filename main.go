package main

import (
	"fmt"
	"log"

	"os"

	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

type structFaker struct {
	ID        string `faker:"uuid_digit"`
	Name      string `faker:"name"`
	Age       int    `faker:"oneof: 15, 27, 61"`
	Country   string `faker:"oneof: USA, ID, AUS, GER, NED"`
	Bio       string `faker:"paragraph"`
	Religion  string `faker:"oneof: MOESLEM, CRIST, HINDU"`
	Birthdate string `faker:"date"`
}

type structDocument struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Country   string `json:"country"`
	Bio       string `json:"bio"`
	Religion  string `json:"religion"`
	Birthdate string `json:"birthdate"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hostTypesense := os.Getenv("HOST_TYPESENSE")
	apiKeyTypesense := os.Getenv("API_KEY_TYPESENSE")
	app := fiber.New()

	client := typesense.NewClient(
		typesense.WithServer(hostTypesense),
		typesense.WithAPIKey(apiKeyTypesense))

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
				Name: "bio",
				Type: "string",
			},
			{
				Name: "religion",
				Type: "string",
			},
			{
				Name: "birthdate",
				Type: "string",
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

		defer func() {
			hits := *results.Hits
			for _, val := range hits {
				fmt.Println(*val.Document)
			}
		}()

		return c.SendString("Success get ")
	})

	app.Post("/", func(c *fiber.Ctx) error {

		documentFaker := structFaker{}
		err := faker.FakeData(&documentFaker)

		document := structDocument{
			ID:        documentFaker.ID,
			Name:      documentFaker.Name,
			Age:       documentFaker.Age,
			Country:   documentFaker.Country,
			Bio:       documentFaker.Bio,
			Religion:  documentFaker.Religion,
			Birthdate: documentFaker.Birthdate,
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
