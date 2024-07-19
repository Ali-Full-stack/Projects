package repo

import (
	"bookstore/internal/model"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepo struct {
	MongoClient *mongo.Client
}

func NewBookRepo(mc *mongo.Client) *BookRepo {
	return &BookRepo{MongoClient: mc}
}

func (b *BookRepo) CreateMultipleBooksToMongoDB(listbooks []model.BookInfo) (*model.Response, error) {
	bookCollection := b.MongoClient.Database(os.Getenv("mongo_db")).Collection("books")

	books :=make([]interface{}, len(listbooks))
	for index, book :=range listbooks{
		books[index] = book
	}	

	_, err :=bookCollection.InsertMany(context.Background(), books)
	if err != nil {
		return nil, fmt.Errorf("error: inserting multiple books to MongoDB: %v",err)
	}
	return &model.Response{Message: "All Books Added Succesfully ."}, nil
}

func (b *BookRepo) GetMultipleBooksFromMongoDB()([]model.BookInfo, error){
	bookCollection := b.MongoClient.Database(os.Getenv("mongo_db")).Collection("books")

	cursor, err :=bookCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error: Getting multiple books From MongoDB: %v",err)
	}
	defer cursor.Close(context.Background())

	var listbooks []model.BookInfo
	for cursor.Next(context.Background()){
		var book model.BookInfo
		if err :=cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("error: cannot decode book to json struct: %v",err)
		}
		listbooks =append(listbooks, book)
	}
	if len(listbooks) == 0 {
		return nil, fmt.Errorf("no books exist in MongoDB ")
	}
	return listbooks, nil
}

func (b *BookRepo) GetMultipleBooksByCategoryFromMongoDB(categor string)([]model.BookInfo, error){
	bookCollection := b.MongoClient.Database(os.Getenv("mongo_db")).Collection("books")

	cursor, err :=bookCollection.Find(context.Background(), bson.M{"category" : categor})
	if err != nil {
		return nil, fmt.Errorf("error: Getting multiple books From MongoDB: %v",err)
	}
	defer cursor.Close(context.Background())

	var listbooks []model.BookInfo
	for cursor.Next(context.Background()){
		var book model.BookInfo
		if err :=cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("error: cannot decode book to json struct: %v",err)
		}
		listbooks =append(listbooks, book)
	}
	if len(listbooks) == 0 {
		return nil, fmt.Errorf("no books exist in MongoDB ")
	}
	return listbooks, nil
}
func (b *BookRepo) GetMultipleBooksByAuthorFromMongoDB(authorName string)([]model.BookInfo, error){
	bookCollection := b.MongoClient.Database(os.Getenv("mongo_db")).Collection("books")

	cursor, err :=bookCollection.Find(context.Background(), bson.M{"author.fullname": authorName})
	if err != nil {
		return nil, fmt.Errorf("error: Getting multiple books From MongoDB: %v",err)
	}
	defer cursor.Close(context.Background())

	var listbooks []model.BookInfo
	for cursor.Next(context.Background()){
		var book model.BookInfo
		if err :=cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("error: cannot decode book to json struct: %v",err)
		}
		listbooks =append(listbooks, book)
	}
	if len(listbooks) == 0 {
		return nil, fmt.Errorf("no books exist in MongoDB ")
	}
	return listbooks, nil
}
func(b *BookRepo) UpdateStatusOfBookInMongoDB(bookStatus model.BookStatus)(*model.Response, error){
	bookCollection := b.MongoClient.Database(os.Getenv("mongo_db")).Collection("books")
	result, err :=bookCollection.UpdateOne(context.Background(), bson.M{"title": bookStatus.Title}, bson.M{"$set" : bson.M{"rentdetails.status" : bookStatus.Status}})
	if err != nil {
		return nil, fmt.Errorf("error: Updating  book's status  in MongoDB: %v",err)
	}

	if result.ModifiedCount ==0 {
		return nil, fmt.Errorf("no books exist in MongoDB ")
	}
	return &model.Response{Message: "Book' s status updated succesfully ."}, nil
}

func (b *BookRepo) DeleteBookByISBNFromMongoDB( isbn string)(*model.Response, error){
	bookCollection := b.MongoClient.Database(os.Getenv("mongo_db")).Collection("books")
	result, err :=bookCollection.DeleteOne(context.Background(), bson.M{"isbn" : isbn})
	if err != nil {
		return nil, fmt.Errorf("error: Deleting  book by ISBN  in MongoDB: %v",err)
	}
	if result.DeletedCount ==0 {
		return nil, fmt.Errorf("no books exist in MongoDB ")
	}
	return &model.Response{Message: "Book deleted succesfully ."}, nil
}
func (b *BookRepo) CountAuthorsBookInMongoDB( )(*[]model.AuthorTotalBook, error){
	bookCollection := b.MongoClient.Database(os.Getenv("mongo_db")).Collection("books")

	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":         "$author.fullname",
				"total_books": bson.M{"$sum": 1},
			},
		},
	}
	cursor, err :=bookCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("error: Counting author' s book  in MongoDB: %v",err)
	}
	defer cursor.Close(context.Background())
	var ListAuthors []model.AuthorTotalBook
	for cursor.Next(context.Background()) {
        var result bson.M
        if err := cursor.Decode(&result); err != nil {
            panic(err)
        }
		ListAuthors = append(ListAuthors, model.AuthorTotalBook{Author: result["author.fullname"].(string), Total_books: result["total_books"].(int32) })
    }
	return &ListAuthors, nil
}