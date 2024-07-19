package handler

import (
	"bookstore/internal/model"
	"bookstore/internal/repo"
	"encoding/json"
	"log"
	"net/http"
)

type BookHandler struct {
	BookRepo *repo.BookRepo
}

func NewBookHandler(br *repo.BookRepo) *BookHandler {
	return &BookHandler{BookRepo: br}
}

// @Router  				/books [post]
// @Summary 			Create Multiple Books
// @Description 		This method addes Multiple Books to Database
// @Security 				BearerAuth
// @Tags					 BOOK
// @accept					json
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param 					password    header    string    true    "Password"
// @Param 					body    body    []model.BookInfo    true  "Books"
// @Success					201 	{object}   model.Response		"All Books Added Succesfully ."
// @Failure					 400 {object} error "Invalid Request: incorrect book information !!"
// @Failure					 403 {object} error "Unauthorized access"
func (b *BookHandler) CreateMultipleBooks(w http.ResponseWriter, r *http.Request) {
	var listbooks []model.BookInfo
	if err := json.NewDecoder(r.Body).Decode(&listbooks); err != nil {
		http.Error(w, "Invalid Request: incorrect book information !!", http.StatusBadRequest)
		return
	}
	resp, err := b.BookRepo.CreateMultipleBooksToMongoDB(listbooks)
	if err != nil {
		log.Println(err)
		http.Error(w, "Request Denied : Unable to Add Books !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)

}
////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/books [get]
// @Summary 			Gets Multiple Books
// @Description 		This method Gets  Multiple Books From Database
// @Security 				BearerAuth
// @Tags					 BOOK
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Success					200 	{object}   []model.BookInfo		
// @Failure					 500 {object} error "Request Denied : Unable to Get Books !"
// @Failure					 403 {object} error "Unauthorized access"
func (b *BookHandler) GetMultipleBooks(w http.ResponseWriter, r *http.Request) {
	listbooks, err := b.BookRepo.GetMultipleBooksFromMongoDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Request Denied : Unable to Get Books !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(listbooks)
}
////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/books/category/{category} [get]
// @Summary 			Gets Multiple Books
// @Description 		This method Gets  Multiple Books By category From Database
// @Security 				BearerAuth
// @Tags					 BOOK
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param				    category path  	string 		 true  "Book's category"
// @Success					200 	{object}   []model.BookInfo		
// @Failure					 500 {object} error "Request Denied : Unable to Get Books !"
// @Failure					 403 {object} error "Unauthorized access"
func (b *BookHandler) GetMultipleBooksByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	listbooks, err := b.BookRepo.GetMultipleBooksByCategoryFromMongoDB(category)
	if err != nil {
		log.Println(err)
		http.Error(w, "Request Denied : Unable to Get Books !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(listbooks)
}

////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/books/author/{author} [get]
// @Summary 			Gets Multiple Books
// @Description 		This method Gets  Multiple Books By AuthorsFrom Database
// @Security 				BearerAuth
// @Tags					 BOOK
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param				    author  path  	  string 	 true  "Book's  author name"
// @Success					200 	{object}   []model.BookInfo		
// @Failure					 500 {object} error "Request Denied : Unable to Get Books !"
// @Failure					 403 {object} error "Unauthorized access"
func (b *BookHandler) GetMultipleBooksByAuthor(w http.ResponseWriter, r *http.Request) {
	authorName :=r.PathValue("author")
	listbooks, err := b.BookRepo.GetMultipleBooksByAuthorFromMongoDB(authorName)
	if err != nil {
		log.Println(err)
		http.Error(w, "Request Denied : Unable to Get Books !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(listbooks)
}
////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/books/status  [put]
// @Summary 			Updates   Book's status
// @Description 		This method updates   Book's status  in Database
// @Security 				BearerAuth
// @Tags					 BOOK
// @accept					json
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param 					password    header    string    true    "Password"
// @Param 					body    body    model.BookStatus    true  "Book's Status"
// @Success					202 	{object}   model.Response		"Book' s status updated succesfully ."
// @Failure					 400 {object} error "Invalid Request: incorrect book information !!"
// @Failure					 403 {object} error "Unauthorized access"
func (b *BookHandler) UpdateStatusOfBook(w http.ResponseWriter, r *http.Request) {
	var bookStatus model.BookStatus
	if err := json.NewDecoder(r.Body).Decode(&bookStatus); err != nil {
		http.Error(w, "Invalid Request: incorrect book information !!", http.StatusBadRequest)
		return
	}
	resp, err :=b.BookRepo.UpdateStatusOfBookInMongoDB(bookStatus)
	if err != nil {
		log.Println(err)
		http.Error(w, "Request Denied : Unable to Update  Book 's Status !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(resp)
}
////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/books/{isbn}  [delete]
// @Summary 			Deletes   Book 
// @Description 		This method Deletes   Book By ISBN  From Database
// @Security 				BearerAuth
// @Tags					 BOOK
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param 					password    header    string    true    "Password"
// @Param 					isbn    path    	string       true  "Book ISBN"
// @Success					201 	{object}   model.Response		"Book deleted succesfully ."
// @Failure					 500 {object} error "Request Denied : Unable to Delete  Book !"
// @Failure					 403 {object} error "Unauthorized access"
func (b *BookHandler) DeleteBookByISBN(w http.ResponseWriter, r *http.Request) {
	isbn :=r.PathValue("isbn")
	resp, err :=b.BookRepo.DeleteBookByISBNFromMongoDB(isbn)
	if err != nil {
		log.Println(err)
		http.Error(w, "Request Denied : Unable to Delete  Book !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)
}

////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/books/count  [get]
// @Summary 			Counts Autors   Book 
// @Description 		This method Counts  author 's  Book By ISBN  From Database
// @Security 				BearerAuth
// @Tags					 BOOK
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Success					201 	{object}   model.Response		"Book deleted succesfully ."
// @Failure					 500 {object} error "Request Denied : Unable to Delete  Book !"
// @Failure					 403 {object} error "Unauthorized access"
func (b *BookHandler) CountAuthorsBook(w http.ResponseWriter, r *http.Request) {
	ListAuthors, err :=b.BookRepo.CountAuthorsBookInMongoDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Request Denied : Unable to count authors total  book !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(ListAuthors)
}
