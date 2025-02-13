package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/smithy-go/middleware"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/neonmei/challenge_urlshortener/domain"
	clientMock "github.com/neonmei/challenge_urlshortener/mocks/clients"
	"github.com/neonmei/challenge_urlshortener/platform/config"
	"github.com/neonmei/challenge_urlshortener/platform/repositories/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBackendErrorHandling(t *testing.T) {
	cfg := config.Load()
	ctx := context.Background()
	dynamoClient := clientMock.NewMockDynamoDbClient(t)
	repo := NewDynamoURLRepository(cfg, dynamoClient)
	validItem := domain.ShortURL{
		ID:        validId,
		Upstream:  *validURL,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	// REF: ALl operations should fail
	dynamoErr := errors.New("dynamo backend failed")
	dynamoClient.On("UpdateItem", mock.Anything, mock.Anything).Return(nil, dynamoErr)
	dynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(nil, dynamoErr)
	dynamoClient.On("PutItem", mock.Anything, mock.Anything).Return(nil, dynamoErr)

	err := repo.Delete(ctx, validId)
	assert.ErrorIs(t, err, domain.ErrUnavailableRepo)

	result, err := repo.Get(ctx, validId)
	assert.ErrorIs(t, err, domain.ErrUnavailableRepo)
	assert.Nil(t, result)

	err = repo.Save(ctx, validItem)
	assert.ErrorIs(t, err, domain.ErrUnavailableRepo)
}

func TestBackendGetNotFound(t *testing.T) {
	cfg := config.Load()
	ctx := context.Background()
	dynamoClient := clientMock.NewMockDynamoDbClient(t)
	repo := NewDynamoURLRepository(cfg, dynamoClient)

	// validItem := domain.ShortURL{
	// 	ID:        validId,
	// 	Upstream:  *validURL,
	// 	CreatedBy: validAuthor,
	// 	CreatedAt: time.Now(),
	// 	Enabled:   true,
	// }

	dynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(&awsDynamodb.GetItemOutput{
		Item:           map[string]types.AttributeValue{},
		ResultMetadata: middleware.Metadata{},
	}, nil)

	result, err := repo.Get(ctx, validId)
	assert.ErrorIs(t, err, domain.ErrURLNotFound)
	assert.Nil(t, result)
}

func TestBackendGetFound(t *testing.T) {
	cfg := config.Load()
	ctx := context.Background()
	dynamoClient := clientMock.NewMockDynamoDbClient(t)
	repo := NewDynamoURLRepository(cfg, dynamoClient)

	validItem := domain.ShortURL{
		ID:        validId,
		Upstream:  *validURL,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	itemDto := dtos.FromDomain(validItem)
	itemDynamo, err := attributevalue.MarshalMap(itemDto)
	assert.NoError(t, err)

	dynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(&awsDynamodb.GetItemOutput{
		Item:           itemDynamo,
		ResultMetadata: middleware.Metadata{},
	}, nil)

	result, err := repo.Get(ctx, validId)
	assert.NoError(t, err)
	assert.NotNil(t, itemDynamo)
	assert.Equal(t, validItem.ID, result.ID)
	assert.Equal(t, validItem.Upstream.String(), result.Upstream.String())
	assert.Equal(t, validItem.CreatedAt.Unix(), result.CreatedAt.Unix())
	assert.Equal(t, validItem.CreatedBy, result.CreatedBy)
	assert.Equal(t, validItem.Enabled, result.Enabled)
}
