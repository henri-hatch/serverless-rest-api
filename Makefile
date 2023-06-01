build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main main.go
deploy_prod: build
	serverless deploy --stage prod --aws-profile henrihatch