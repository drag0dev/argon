module github.com/drag0dev/argon/src/cdk

replace common => ../common

require (
	common v0.0.0-00010101000000-000000000000
	github.com/aws/aws-cdk-go/awscdk/v2 v2.147.2
	github.com/aws/aws-sdk-go-v2 v1.27.1
	github.com/aws/constructs-go/constructs/v10 v10.3.0
	github.com/aws/jsii-runtime-go v1.99.0
)

require (
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/aws/aws-lambda-go v1.47.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.14.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.32.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.20.9 // indirect
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/cdklabs/awscdk-asset-awscli-go/awscliv1/v2 v2.2.202 // indirect
	github.com/cdklabs/awscdk-asset-kubectl-go/kubectlv20/v2 v2.1.2 // indirect
	github.com/cdklabs/awscdk-asset-node-proxy-agent-go/nodeproxyagentv6/v2 v2.0.3 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
)

go 1.22.4
