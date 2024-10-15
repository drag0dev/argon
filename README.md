# Argon
Argon is an AWS cloud-native movie and TV show streaming service.

## Features
- **Streaming** - watch any available movie or TV show
- **Rating** - rate and review content
- **Content Management** - upload, edit, and delete movies and TV shows
- **Personalized feed** - recommending movies and TV shows to the users based on their interests
- **Notification** - users get notified of all new content that interests them

## Prerequisites 
- Golang
- AWS CLI
- AWS CDK

## Deployment
Before deploying the project, all necessary lambdas and additional libraries have to be compiled and zipped
```
cd ./src/build-system
go run main.go ffmpeg
go run main.go build
```

After all the previous commands have been executed, the project is ready for deployment
```
cdk deploy --app 'ACCOUNT_ID=your_account_id go run main.go'
```