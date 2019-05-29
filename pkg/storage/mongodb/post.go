package mongodb

import (
	"github.com/das-frama/website/pkg/post"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type postRepo struct {
	col *mongo.Collection
}

// NewPostRepo creates a page repos to provide access to the storage.
func NewPostRepo(db *mongo.Database) post.Repository {
	return &postRepo{db.Collection("post")}
}

// FindByID returns a post with provided id.
func (r *postRepo) FindByID(id string) (*post.Post, error) {
	post := &post.Post{}
	objectID, _ := primitive.ObjectIDFromHex(id)
	if err := r.col.FindOne(nil, bson.M{"_id": objectID}).Decode(&post); err != nil {
		return nil, err
	}

	return post, nil
}

// FindBySlug returns a post with provided id.
func (r *postRepo) FindBySlug(slug string) (*post.Post, error) {
	post := &post.Post{}
	if err := r.col.FindOne(nil, bson.M{"slug": slug}).Decode(&post); err != nil {
		return nil, err
	}

	return post, nil
}

// FindAll returns all stored posts.
func (r *postRepo) FindAll() ([]*post.Post, error) {
	cursor, err := r.col.Find(nil, bson.D{})
	if err != nil {
		return nil, err
	}

	var posts []*post.Post
	for cursor.Next(nil) {
		post := &post.Post{}
		if err := cursor.Decode(post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	// Check if the cursor encountered any errors while iterating.
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// Create creates a post in the storage.
func (r *postRepo) Insert(post *post.Post) error {
	res, err := r.col.InsertOne(nil, post)
	if err != nil {
		return err
	}
	post.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

// Update updates the post in storage.
func (r *postRepo) Update(post *post.Post) error {
	_, err := r.col.UpdateOne(nil, bson.M{"_id": post.ID}, bson.M{"$set": post})
	return err
}

// Delete deletes the post with provided id from storage.
func (r *postRepo) Delete(id string) (bool, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := r.col.DeleteOne(nil, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}

	return res.DeletedCount > 0, nil
}
