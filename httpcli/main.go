package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	var query cli.StringSlice
	var header cli.StringSlice
	var json string

	app := &cli.App{
		Name:                   "httpcli",
		Usage:                  "A simple HTTPCLI to make REST API requests",
		UsageText:              "httpcli [--query \"key=value\"] [--header \"key=value\"] <URL>",
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "query",
				Usage:       "Add query parameters to the request",
				Aliases:     []string{"q"},
				Destination: &query,
			},
			&cli.StringSliceFlag{
				Name:        "header",
				Usage:       "Add headers to the request",
				Destination: &header,
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "get",
				Usage:     "Fetches data from the specified URL. You can add query parameters and headers.",
				UsageText: "httpcli [--query \"key=value\"] [--header \"key=value\"] get <URL>",
				Action: func(cCtx *cli.Context) error {
					url := cCtx.Args().First()

					// fmt.Println("Get Method")
					// fmt.Println("URL:", url)
					// fmt.Println("Queries:", query.Value())
					// fmt.Println("Headers:", header.Value())
					httpSendRequest("GET", url, query, header, "")

					return nil
				},
			},
			{
				Name:      "post",
				Usage:     "Sends a POST request with an optional JSON payload and headers.",
				UsageText: "httpcli [--query \"key=value\"] [--header \"key=value\"] post [--json \"{'key1':'value1, 'key2':'value2'}\"] <URL>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "json",
						Usage:       "JSON payload for requests",
						Destination: &json,
					},
				},
				Action: func(cCtx *cli.Context) error {
					url := cCtx.Args().First()

					// fmt.Println("Post Method")
					// fmt.Println("URL:", url)
					// fmt.Println("Queries:", query.Value())
					// fmt.Println("Headers:", header.Value())
					// fmt.Println("JSON:", json)
					httpSendRequest("POST", url, query, header, json)

					return nil
				},
			},
			{
				Name:      "put",
				Usage:     "Updates existing data on the server using a PUT request.",
				UsageText: "httpcli [--query \"key=value\"] [--header \"key=value\"] put [--json \"{'key1':'value1, 'key2':'value2'}\"] <URL>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "json",
						Usage:       "JSON payload for requests",
						Destination: &json,
					},
				},
				Action: func(cCtx *cli.Context) error {
					url := cCtx.Args().First()

					// fmt.Println("Put Method")
					// fmt.Println("URL:", url)
					// fmt.Println("Queries:", query.Value())
					// fmt.Println("Headers:", header.Value())
					// fmt.Println("JSON:", json.Value())
					httpSendRequest("PUT", url, query, header, json)

					return nil
				},
			},
			{
				Name:      "delete",
				Usage:     "Deletes a resource at the specified URL.",
				UsageText: "httpcli [--query \"key=value\"] [--header \"key=value\"] delete <URL>",
				Action: func(cCtx *cli.Context) error {
					url := cCtx.Args().First()

					// fmt.Println("Delete Method")
					// fmt.Println("URL:", url)
					// fmt.Println("Queries:", query.Value())
					// fmt.Println("Headers:", header.Value())
					httpSendRequest("DELETE", url, query, header, "")

					return nil
				},
			},
		},

		Action: func(cCtx *cli.Context) error {
			url := cCtx.Args().First()

			// fmt.Println("URL:", url)
			// fmt.Println("Queries:", query.Value())
			// fmt.Println("Headers:", header.Value())
			httpSendRequest("GET", url, query, header, "")

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func httpSendRequest(reqType string, baseUrl string, query cli.StringSlice, header cli.StringSlice, body string) {

	// handle query
	queryParams := url.Values{}
	for _, q := range query.Value() {
		parts := strings.SplitN(q, "=", 2)
		if len(parts) == 2 {
			queryParams.Add(parts[0], parts[1])
		}
	}

	fullUrl := baseUrl
	if encodedQuery := queryParams.Encode(); encodedQuery != "" {
		fullUrl += "?" + encodedQuery
	}

	// handle body
	var bodyParser io.Reader
	if body != "" {
		body = strings.ReplaceAll(body, "'", "\"")
		bodyParser = strings.NewReader(body)
	} else {
		bodyParser = nil
	}

	// send request
	client := &http.Client{}
	req, err := http.NewRequest(reqType, fullUrl, bodyParser)
	if err != nil {
		log.Fatal("request error:", err)
	}

	// add header
	req.Header.Set("Content-Type", "application/json")
	for _, h := range header.Value() {
		parts := strings.SplitN(h, "=", 2)
		if len(parts) == 2 {
			fmt.Println("Header :", parts[0], parts[1])
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	// response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("response error:", err)
	}
	defer resp.Body.Close()

	bodyByte, _ := io.ReadAll(resp.Body)
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response", string(bodyByte))
}
