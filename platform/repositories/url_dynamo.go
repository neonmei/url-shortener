package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/neonmei/challenge_urlshortener/domain/validators"
	"github.com/neonmei/challenge_urlshortener/platform/repositories/dtos"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/neonmei/challenge_urlshortener/platform/clients"
	"github.com/neonmei/challenge_urlshortener/platform/config"
)

type dynaURLRepo struct {
	tableName    string
	client       clients.DynamoDbClient
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (d *dynaURLRepo) Save(ctx context.Context, shortUrl domain.ShortURL) error {
	if err := validators.ValidateShortURL(shortUrl); err != nil {
		return err
	}

	urlItem := dtos.FromDomain(shortUrl)
	item, err := attributevalue.MarshalMap(urlItem)
	if err != nil {
		return errors.Join(errors.New("cannot serialize urlItem"), err)
	}

	newCtx, cancelFunc := context.WithTimeout(ctx, d.writeTimeout)
	defer cancelFunc()

	_, err = d.client.PutItem(newCtx, &awsDynamodb.PutItemInput{
		TableName:           &d.tableName,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(url_id)"),
	})
	if err != nil {
		return errors.Join(domain.ErrUnavailableRepo, err)
	}

	return nil
}

func (d *dynaURLRepo) Get(ctx context.Context, urlID string) (*domain.ShortURL, error) {
	newCtx, cancelFunc := context.WithTimeout(ctx, d.readTimeout)
	defer cancelFunc()

	itemResult, err := d.client.GetItem(newCtx, &awsDynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]types.AttributeValue{
			"url_id": &types.AttributeValueMemberS{Value: urlID},
		},
	})
	if err != nil {
		return nil, errors.Join(domain.ErrUnavailableRepo, err)
	}

	if len(itemResult.Item) == 0 {
		return nil, domain.ErrURLNotFound
	}

	itemModel := dtos.URLItem{}
	if err = attributevalue.UnmarshalMap(itemResult.Item, &itemModel); err != nil {
		return nil, errors.Join(domain.ErrRepoSchema, err)
	}

	shortUrl, err := itemModel.Domain()
	if err != nil {
		return nil, errors.Join(domain.ErrRepoSchema, err)
	}

	return shortUrl, nil
}

func (d *dynaURLRepo) Delete(ctx context.Context, urlID string) error {
	newCtx, cancelFunc := context.WithTimeout(ctx, d.writeTimeout)
	defer cancelFunc()

	_, err := d.client.UpdateItem(newCtx, &awsDynamodb.UpdateItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"url_id": &types.AttributeValueMemberS{Value: urlID},
		},
		UpdateExpression: aws.String("SET enabled = :enabled"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":enabled": &types.AttributeValueMemberBOOL{Value: false},
		},
	})
	if err != nil {
		return errors.Join(domain.ErrUnavailableRepo, err)
	}
	return nil
}

func NewDynamoURLRepository(cfg config.AppConfig, client clients.DynamoDbClient) domain.URLRepository {
	return &dynaURLRepo{
		tableName:    cfg.Dynamo.TableName,
		client:       client,
		readTimeout:  cfg.Dynamo.ReadTimeout,
		writeTimeout: cfg.Dynamo.WriteTimeout,
	}
}
