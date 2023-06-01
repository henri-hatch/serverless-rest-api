package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	r := gin.Default()
	r.GET("/ping", PingPong)
	r.POST("/challenge", Challenge)

	ginLambda = ginadapter.New(r)
}

func PingPong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Challenge(c *gin.Context) {
	type SlackRequest struct {
		Type      string `json:"type"`
		Token     string `json:"token"`
		Challenge string `json:"challenge"`
	}

	var request SlackRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	if request.Type != "url_verification" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request type",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"body":       request.Challenge,
	})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(HandleRequest)
}
