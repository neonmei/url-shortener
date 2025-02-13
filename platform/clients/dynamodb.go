package clients

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/neonmei/challenge_urlshortener/platform/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type DynamoDbClient interface {
	Scan(ctx context.Context, params *awsDynamodb.ScanInput, optFns ...func(options *awsDynamodb.Options)) (*awsDynamodb.ScanOutput, error)
	Query(ctx context.Context, params *awsDynamodb.QueryInput, optFns ...func(options *awsDynamodb.Options)) (*awsDynamodb.QueryOutput, error)
	PutItem(ctx context.Context, params *awsDynamodb.PutItemInput, optFns ...func(options *awsDynamodb.Options)) (*awsDynamodb.PutItemOutput, error)
	UpdateItem(ctx context.Context, params *awsDynamodb.UpdateItemInput, optFns ...func(options *awsDynamodb.Options)) (*awsDynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, params *awsDynamodb.DeleteItemInput, optFns ...func(options *awsDynamodb.Options)) (*awsDynamodb.DeleteItemOutput, error)
	GetItem(ctx context.Context, params *awsDynamodb.GetItemInput, optFns ...func(*awsDynamodb.Options)) (*awsDynamodb.GetItemOutput, error)
}

func NewDynamoClient(appCfg config.AppConfig) (DynamoDbClient, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)
	return awsDynamodb.NewFromConfig(cfg), nil
}
