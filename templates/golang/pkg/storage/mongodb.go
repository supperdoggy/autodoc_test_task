package storage

import (
	"context"

	"github.com/Autodoc-Technology/interview-templates/template/golang/pkg/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type IMongoClient interface {
	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, max, skip int64) ([]models.User, error)
}

type mongoClient struct {
	logger   *zap.Logger
	client   *mongo.Client
	usersCol *mongo.Collection
}

func NewMongoClient(ctx context.Context, url string, l *zap.Logger) (IMongoClient, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		l.Error("error connecting to mongodb", zap.Error(err))
		return nil, err
	}

	return &mongoClient{
		client: client,
		logger: l,

		usersCol: client.Database("autotest").Collection("users"),
	}, nil
}

func (c *mongoClient) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	u.ID = uuid.New().String()
	_, err := c.usersCol.InsertOne(ctx, u)
	return u, err
}

func (c *mongoClient) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	res := c.usersCol.FindOne(ctx, bson.M{"_id": id, "is_deleted": false})
	if res.Err() != nil {
		return nil, res.Err()
	}

	if err := res.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *mongoClient) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	res := c.usersCol.FindOne(ctx, bson.M{"email": email, "is_deleted": false})
	if res.Err() != nil {
		return nil, res.Err()
	}

	if err := res.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *mongoClient) DeleteUser(ctx context.Context, id string) error {
	_, err := c.usersCol.UpdateByID(ctx, id, bson.M{"$set": bson.M{"is_deleted": true}})
	return err
}

func (c *mongoClient) ListUsers(ctx context.Context, max, skip int64) ([]models.User, error) {
	resp := []models.User{}
	opts := options.Find().SetLimit(max).SetSkip(skip)
	cur, err := c.usersCol.Find(ctx, bson.M{"is_deleted": false}, opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var u models.User
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}

		resp = append(resp, u)
	}
	return resp, nil
}
