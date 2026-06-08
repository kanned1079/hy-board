package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

// graphqlSchema is assembled from the types, queries, and mutations
// defined in graphql_types.go, graphql_queries.go, and graphql_mutations.go.
var graphqlSchema graphql.Schema

func init() {
	var err error
	graphqlSchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize GraphQL Schema: %v", err))
	}
}

// GraphQLRequestBody represents the incoming JSON payload for a GraphQL request.
type GraphQLRequestBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GraphQLHandler handles all /graphql requests, delegating execution to the schema.
func GraphQLHandler(c *gin.Context) {
	var body GraphQLRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         graphqlSchema,
		RequestString:  body.Query,
		VariableValues: body.Variables,
		Context:        c,
	})

	if len(result.Errors) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"errors": result.Errors,
			"data":   result.Data,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
