package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/das-frama/website/pkg/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDB returns a new MongoDB storage based on config.
func NewDB(config *app.Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Get client.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DBURI))
	if err != nil {
		return nil, err
	}

	// Check if mongodb server is available.
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.New("mongodb is not reachable")
	}

	return client.Database(config.DBName), nil
}
