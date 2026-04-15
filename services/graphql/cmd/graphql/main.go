package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"kate/services/graphql/graph"
	"kate/services/graphql/internal/store"

	"github.com/graphql-go/graphql"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	s := store.NewStore()
	schema := graph.InitSchema(s)

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  params.Query,
			VariableValues: params.Variables,
			OperationName:  params.OperationName,
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// GraphQL Playground
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
		<html>
		<head>
			<title>GraphQL Playground</title>
			<style>
				body { margin: 0; padding: 0; font-family: monospace; }
				#playground { height: 100vh; width: 100vw; }
			</style>
		</head>
		<body>
			<div id="playground"></div>
			<script src="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
			<script>
				GraphQLPlayground.init(document.getElementById('playground'), {
					endpoint: '/graphql'
				})
			</script>
		</body>
		</html>`
		w.Write([]byte(html))
	})

	log.Printf("GraphQL server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
