package graph

import (
	"fmt"
	"time"

	"kate/services/graphql/internal/store"

	"github.com/graphql-go/graphql"
)

var taskType *graphql.Object
var storeInstance *store.Store

func InitSchema(s *store.Store) graphql.Schema {
	storeInstance = s

	taskType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Task",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.ID},
			"title":       &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"done":        &graphql.Field{Type: graphql.Boolean},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"tasks": &graphql.Field{
				Type: graphql.NewList(taskType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					tasks := storeInstance.GetAll()
					result := make([]interface{}, len(tasks))
					for i, t := range tasks {
						result[i] = t
					}
					return result, nil
				},
			},
			"task": &graphql.Field{
				Type: taskType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					task, ok := storeInstance.GetByID(id)
					if !ok {
						return nil, nil
					}
					return task, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createTask": &graphql.Field{
				Type: taskType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewInputObject(graphql.InputObjectConfig{
							Name: "CreateTaskInput",
							Fields: graphql.InputObjectConfigFieldMap{
								"title": &graphql.InputObjectFieldConfig{
									Type: graphql.NewNonNull(graphql.String),
								},
								"description": &graphql.InputObjectFieldConfig{
									Type: graphql.String,
								},
							},
						}),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input, _ := p.Args["input"].(map[string]interface{})
					title := input["title"].(string)
					var description *string
					if desc, ok := input["description"].(string); ok {
						description = &desc
					}
					id := fmt.Sprintf("t_%d", time.Now().UnixNano())
					task := storeInstance.Create(id, title, description)
					return task, nil
				},
			},
			"updateTask": &graphql.Field{
				Type: taskType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewInputObject(graphql.InputObjectConfig{
							Name: "UpdateTaskInput",
							Fields: graphql.InputObjectConfigFieldMap{
								"title":       &graphql.InputObjectFieldConfig{Type: graphql.String},
								"description": &graphql.InputObjectFieldConfig{Type: graphql.String},
								"done":        &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
							},
						}),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					input, _ := p.Args["input"].(map[string]interface{})

					var title *string
					if t, ok := input["title"].(string); ok {
						title = &t
					}
					var description *string
					if d, ok := input["description"].(string); ok {
						description = &d
					}
					var done *bool
					if d, ok := input["done"].(bool); ok {
						done = &d
					}

					task, ok := storeInstance.Update(id, title, description, done)
					if !ok {
						return nil, fmt.Errorf("task not found")
					}
					return task, nil
				},
			},
			"deleteTask": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					return storeInstance.Delete(id), nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		panic(err)
	}
	return schema
}
