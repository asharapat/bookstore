package helper

import (
	"bookstore/models"
	"encoding/json"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Repository struct{
	mongoSession *mgo.Session
}

var Conf = models.Config {
	"localhost:27017",
	"benfica",
	"books",
}

func NewRepository() *Repository {
	session, err := mgo.Dial(Conf.MongoDBhost)
	if err != nil {
		log.Println("Cannot make session to MongoDB! ", err)
		panic(err)
	}
	return &Repository{
		mongoSession: session,
	}
}

func (r *Repository) ClearTable() {
	r.mongoSession.DB(Conf.MongoDBname).C(Conf.MongoDBCollectionName).RemoveAll(nil)
}

func parseQuery(q string) (bson.M, error) {
	if q != "" {
		outQuery := bson.M{}
		err := bson.UnmarshalJSON([]byte(q), &outQuery)
		if err != nil {
			log.Println("Error while UnmarshalingJSON ", err)
			return nil, err
		}
		return outQuery, nil
	}
	return nil, nil
}

func (r *Repository) GetBooks(query string) ([]byte, error) {
	results := make([]models.Book, 0)
	c := r.mongoSession.DB(Conf.MongoDBname).C(Conf.MongoDBCollectionName)
	parsedQuery, err := parseQuery(query)
	if err != nil {
		log.Println("Error while parsing query! ", err)
		return nil, err
	}
	if err := c.Find(parsedQuery).All(&results); err != nil {
		log.Println("Failed to write results:", err)
	}
	output, err := json.Marshal(results)
	if err != nil {
		log.Println("Error while Marshaling data", err)
		return nil, err
	}
	return output, nil

}

func (r *Repository) CreateBook(book models.Book) (string, error) {
	book.ID = bson.NewObjectId()
	err := r.mongoSession.DB(Conf.MongoDBname).C(Conf.MongoDBCollectionName).Insert(book)
	if err != nil {
		return "", err
	}
	return book.ID.Hex(), nil
}

func (r *Repository) GetBook(id string) ([]byte, error) {
	result := make([]models.Book, 0)
	c := r.mongoSession.DB(Conf.MongoDBname).C(Conf.MongoDBCollectionName)
	if err := c.FindId(bson.ObjectIdHex(id)).All(&result); err != nil {
		log.Println("Cannot find book by ID", err)
	}
	return json.Marshal(result)
}

func (r *Repository) UpdateBook(book models.Book) error {
	err := r.mongoSession.DB(Conf.MongoDBname).C(Conf.MongoDBCollectionName).Update(
		bson.M{"_id": book.ID},
		bson.M{"$set": bson.M{"name": book.Name, "author": book.Author}})
	if err != nil {
		log.Println("Cannot update item", err)
		return err
	}
	return nil
}

func (r *Repository) DeleteBook(id string) (string, error) {
	if !bson.IsObjectIdHex(id) {
		return "404", errors.New("ID is not ObjectIdHex! ")
	}
	oid := bson.ObjectIdHex(id)
	if err := r.mongoSession.DB(Conf.MongoDBname).C(Conf.MongoDBCollectionName).RemoveId(oid); err != nil {
		log.Println("Error while removing item!", err)
		return "500", err
	}
	return "OK", nil
}

