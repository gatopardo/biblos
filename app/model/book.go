package model

import (
        "database/sql"
	"time"
        "strings"
	"fmt"
        "log"
        )

// *****************************************************************************
// Book
// *****************************************************************************

// Book table contains the information for each book
type Book struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
        LanguageId  uint32        `db:"languageid" bson:"languageid,omitempty"`
        EditorId  uint32          `db:"editorid" bson:"editorid,omitempty"`
        AuthorId  uint32          `db:"authorid" bson:"authorid,omitempty"`
	Title       string        `db:"title" bson:"title"`
	Isbn        string        `db:"isbn" bson:"isbn"`
	Comment     string        `db:"comment" bson:"comment"`
        Year        uint32        `db:"year" bson:"year"`
	CreatedAt   time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at" bson:"updated_at"`
}

// --------------------------------------------------------

// Book table contains the information for each book and language, editor and author
type BookN struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
        Language    string        `db:"language" bson:"language,omitempty"`
        Editor      string        `db:"editor" bson:"editor,omitempty"`
        Author      string        `db:"author" bson:"author,omitempty"`
	Title       string        `db:"title" bson:"title"`
	Isbn        string        `db:"isbn" bson:"isbn"`
	Comment     string        `db:"comment" bson:"comment"`
        Year        uint32        `db:"year" bson:"year"`
	CreatedAt   time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at" bson:"updated_at"`
}

type BookZ struct {
	Title       string        `db:"title" bson:"title"`
	Comment     string        `db:"comment" bson:"comment"`
        Year        uint32        `db:"year" bson:"year"`
        Author      string        `db:"author" bson:"author,omitempty"`
        Editor      string        `db:"editor" bson:"editor,omitempty"`
        Language    string        `db:"language" bson:"language,omitempty"`
}


// --------------------------------------------------------
func (book * Book) NoBlankb(){
            book.Title =  strings.Trim(book.Title, " ")
            book.Isbn =  strings.Trim(book.Isbn, " ")
            book.Comment =  strings.Trim(book.Comment, " ")
       }
// --------------------------------------------------------
// BookById tenemos el book dado id
func (book * Book)BookById() (err error) {
        stq  :=   "SELECT language_id, editor_id, author_id, title, isbn, comment, year, created_at, updated_at FROM books WHERE id=$1"
	err = Db.QueryRow(stq, book.Id).Scan( &book.LanguageId, &book.EditorId, &book.AuthorId, &book.Title, &book.Isbn, &book.Comment,&book.Year ,&book.CreatedAt, &book.UpdatedAt)

	return  standardizeError(err)
}
// --------------------------------------------------------
// BookByRe tenemos el book dado re
func BookByRe(reStr string) (lsBooks []BookZ, err error) {
	stq := "select b.title, b.comment, b.year, a.name as author, e.name as editor, l.name as language from books b join authors a  on a.id = b.author_id join editors e on e.id = editor_id join languages l on l.id = b.language_id where b.title ~ $1 order by b.title, a.name, e.name"
	       var rows * sql.Rows
	       rows, err = Db.Query(stq, reStr)
               if err != nil {
                    log.Println(err)
               }else{
                  defer rows.Close()
                  for rows.Next() {
			  bib := BookZ{}
                     err = rows.Scan(&bib.Title, &bib.Comment, &bib.Year, &bib.Author, &bib.Editor, &bib.Language)
                     if  err != nil {
			 log.Println( err)
                          break
		     }
                    lsBooks = append(lsBooks, bib)
	          }
       }
	return
  }

// --------------------------------------------------------

// BookByName tenemos el book dado title
func (book * Book)BookByName() (err error) {
        stq  := "SELECT language_id, editor_id, author_id, title, isbn, comment, created_at, updated_at  FROM books WHERE title=$1"
	err = Db.QueryRow(stq, &book.Title).Scan(&book.Id, &book.LanguageId, &book.EditorId, &book.AuthorId, &book.Title, &book.Isbn, &book.Comment, &book.CreatedAt, &book.UpdatedAt)

	return  standardizeError(err)
}

// -----------------------------------------------------
// BookCreate crear book
func (book *Book)BookCreate()( pos int, ex error) {
         var err error
         var id uint32
	 var bib Book
         stq := "INSERT INTO books ( language_id, editor_id, author_id, title, isbn, comment,year, created_at, updated_at ) VALUES ($1,$2,$3, $4, $5, $6, $7, $8, $9) returning id"
	 now  := time.Now()
         err = Db.QueryRow( stq, book.LanguageId, book.EditorId, book.AuthorId, book.Title, book.Isbn, book.Comment, book.Year, now, now ).Scan(&id)
            if err != nil {
                 log.Println(err)
	    }else{
	book.Id = id
               stq =  "SELECT b.id,  b.title, b.isbn, b.comment, b.year FROM books b order by title "
	       var rows * sql.Rows
	       rows, err = Db.Query(stq)
               if err != nil {
                    log.Println(err)
               }else{
                  defer rows.Close()
                  for rows.Next() {
                     bib = Book{}
                     err = rows.Scan(&bib.Id, &bib.Title, &bib.Isbn, &bib.Comment, &bib.Year)
		     pos++
                     if  err != nil {
			 log.Println( err)
                          break
		     }
	             if bib.Id == id {
			     return
		     }
	          }
	      }

            }
	    ex = standardizeError(err)
	return
}
// -----------------------------------------------------
// BookUpdate update book
func (book *Book)BookUpdate(stq string) (err error) {
             _, err = Db.Exec(stq )
	    return standardizeError(err)
}

// -----------------------------------------------------
// -----------------------------------------------------
 func  (book * Book)BookDeleteById()( err error){
         stqd :=  "DELETE FROM books where id = $1"
           _, err = Db.Exec(stqd, book.Id)
         return
       }

// Delete language from databa
func (book *Book)Delete() (err error) {
	statement := "delete from books where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(book.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de language en la database
func (book *Book)Update() (err error) {
	statement := "update books set language_id = $2, editor_id = $3, author_id =$4,  title = $5, isbn = $6, comment = $7 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.Id, book.LanguageId, book.EditorId, book.AuthorId, book.Title, book.Isbn, book.Comment )
	return
}

// -----------------------------------------------------
// Delete all books from database
func BookDeleteAll() (err error) {
	statement := "delete from books"
	_, err = Db.Exec(statement)
	return
}

// -----------------------------------------------------

// -------------------------------------------------------------
// Get number of records in books
  func BookCount( ) ( count int) {
        stq :=  "SELECT count(*) FROM books "
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
            err = rows.Scan(&count)
	    if err != nil {
	         return
	    }
        }
	return
 }
// -------------------------------------------------------------
// Get number of records in books with all restrictions
  func BookCountAll( ) ( count int) {
        stq :=  "SELECT COUNT(*) FROM books b, languages l, editors e, authors a where b.language_id = l.id and b.editor_id = e.id and b.author_id = a.id "
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
	    err = rows.Scan(&count)
	    if err != nil {
	         return
	    }
        }
	return
 }
// -------------------------------------------------------------
// Get number of records in books for a given author
  func BookAuthCount(Id  uint32 ) ( count int) {

        stq :=  "SELECT COUNT(*) FROM books b, languages l, editors e, authors a where b.language_id = l.id and b.editor_id = e.id and b.author_id = a.id  and a.id = $1"
        count = 0
	rows, err := Db.Query(stq, Id)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
	    err = rows.Scan(&count)
	    if err != nil {
	         return
	    }
        }
	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get number of records in books for a given editor
  func BookEditCount(Id  uint32 ) ( count int) {

        stq :=  "SELECT COUNT(*) FROM books b, languages l, editors e, authors a where b.language_id = l.id and b.editor_id = e.id and b.author_id = a.id  and e.id = $1"

	rows, err := Db.Query(stq, Id)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
            err = rows.Scan(&count)
	    if err != nil {
	         return
	    }
        }
	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get number of records in books for a given language
  func BookLangCount(Id  uint32 ) ( count int) {

        stq :=  "SELECT COUNT(*) FROM books b, languages l, editors e, authors a where b.language_id = l.id and b.editor_id = e.id and b.author_id = a.id  and l.id = $1"

	rows, err := Db.Query(stq, Id)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
            err = rows.Scan(&count)
	    if err != nil {
	         return
	    }
        }
	return
 }
// -------------------------------------------------------------
// Get selected books in the database and returns the list
  func SBooks(rsearch string  ) (books []BookN, err error) {
        var book BookN
       var stqi, stqf, stq, stq1, stq2, stq3   string
        rsearch    =   strings.Trim(rsearch, " ")
        nCount     := strings.Count(rsearch, "@")
        arSt := strings.Split(rsearch, "@")
        if arSt[0] != ""{
            stq1  = " and  b.title ~* '"+arSt[0]+ "' "
        }
        if nCount >= 1 && arSt[1] != "" {
            stq2  = " and  a.name ~* '"+arSt[1]+ "' "
        }
        if nCount >= 2 && arSt[2] != "" {
            stq3  = " and  b.comment ~* '" +arSt[2]+ "' "
        }
        stqi =  "SELECT b.id, l.name as language, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b join languages l on ( b.language_id = l.id ) join editors e on ( b.editor_id = e.id ) join authors a on ( b.author_id = a.id ) "

        stqf =    " order by a.name, b.title"
        stq =  stqi  + stq1 + stq2 + stq3 + stqf

	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
            book = BookN{}
            err = rows.Scan(&book.Id, &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt)
            if  err != nil {
                             fmt.Println( err)
                             log.Println( err)
			return
		}
		books = append(books, book)
	}
	return
 }
// -----------------------------------------------------
// -------------------------------------------------------------
// Get limit records from offset
  func BookLim(lim int , offs int) (books []BookN, err error) {
        var book BookN
        stq :=  "SELECT b.id, l.name as language, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b left outer join languages l on ( b.language_id = l.id ) left outer join editors e on ( b.editor_id = e.id ) left outer join authors a on ( b.author_id = a.id )  order by b.title  LIMIT $1 OFFSET $2"
	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
            book = BookN{}
            err = rows.Scan(&book.Id, &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt)
            if  err != nil {
                             fmt.Println( err)
                             log.Println( err)
//			return
		}
		books = append(books, book)
	}
	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get limit records from offset
  func BookAuthLim(Id uint32,lim int , offs int) (books []BookN, err error) {
        var book BookN

        stq :=  "SELECT b.id, l.name as language, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b, languages l,  editors e,  authors a where b.language_id = l.id and b.editor_id = e.id and  b.author_id = a.id  and a.id = $3  order by title  LIMIT $1 OFFSET $2"

	rows, err := Db.Query(stq, lim, offs, Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
             book = BookN{}

             if err = rows.Scan(&book.Id, &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt);  err != nil {

			return
		}
             books = append(books, book)
	}

	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get all records given author
  func BookAuthTot(Id uint32) (books []BookN, err error) {
        var book BookN

        stq :=  "SELECT b.id, l.name as language, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b, languages l,  editors e,  authors a where b.language_id = l.id and b.editor_id = e.id and  b.author_id = a.id  and a.id = $1  order by title  "   

	rows, err := Db.Query(stq,  Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
             book = BookN{}

             if err = rows.Scan(&book.Id, &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt);  err != nil {

			return
		}
               books = append(books, book)
	}

	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get all records given editor
  func BookEditTot(Id uint32) (books []BookN, err error) {
        var book BookN

        stq :=  "SELECT b.id, l.name as language, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b, languages l,  editors e,  authors a where b.language_id = l.id and b.editor_id = e.id and  b.author_id = a.id  and e.id = $1  order by author, title "


	rows, err := Db.Query(stq,  Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
             book = BookN{}

             if err = rows.Scan(&book.Id, &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt);  err != nil {

			return
		}
             books = append(books, book)
	}

	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get limit records from offset
  func BookEditLim(Id uint32,lim int , offs int) (books []BookN, err error) {
        var book BookN

        stq :=  "SELECT b.id, l.name as language, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b, languages l,  editors e,  authors a where b.language_id = l.id and b.editor_id = e.id and  b.author_id = a.id  and e.id = $3  order by title  LIMIT $1 OFFSET $2"

	rows, err := Db.Query(stq, lim, offs, Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
             book = BookN{}

             if err = rows.Scan(&book.Id, &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt);  err != nil {

			return
		}
                  books = append(books, book)
	}

	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get limit records from offset
  func BookLangLim(Id uint32,lim int , offs int) (books []BookN, err error) {
        var book BookN

        stq :=  "SELECT b.id, l.name as language, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b, languages l,  editors e,  authors a where b.language_id = l.id and b.editor_id = e.id and  b.author_id = a.id  and l.id = $3  order by title  LIMIT $1 OFFSET $2"

	rows, err := Db.Query(stq, lim, offs, Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
             book = BookN{}

             if err = rows.Scan(&book.Id, &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt);  err != nil {

			return
		}
               books = append(books, book)
	}

	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get total records for a language
  func BookLangTot(Id uint32) (books []BookN, err error) {
        var book BookN

        stq :=  "SELECT b.id, l.name as language, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b, languages l,  editors e,  authors a where b.language_id = l.id and b.editor_id = e.id and  b.author_id = a.id  and l.id = $1  order by author, title "

	rows, err := Db.Query(stq, Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
             book = BookN{}

             if err = rows.Scan(&book.Id, &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt);  err != nil {

			return
		}
                books = append(books, book)
	}

	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get all books in the database and returns the list
  func BooksN() (books []BookN, err error) {
     var book BookN
        stq :=  "SELECT l.name as lang, e.name as editor, a.name as author, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b, languages l, editors e, authors a where b.language_id = l.id and b.editor_id = e.id and b.author_id = a.id   order by title"  
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
           if err = rows.Scan( &book.Language, &book.Editor, &book.Author, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt);  err != nil {
		return
		}
		books = append(books, book)
	}
	return
 }
// -------------------------------------------------------------
// Get all books in the database and returns the list
  func Books() (books []Book, err error) {
     var book Book
        stq :=  "SELECT b.id, b.language_id, b.editor_id, b.author_id, b.title, b.isbn, b.comment, b.year, b.created_at, b.updated_at FROM books b order by title"  
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
           if err = rows.Scan(&book.Id, &book.LanguageId, &book.EditorId, &book.AuthorId, &book.Title, &book.Isbn, &book.Comment, &book.Year, &book.CreatedAt, &book.UpdatedAt);  err != nil {
		return
             }
		books = append(books, book)
	}
	return
 }
// -------------------------------------------------------------
