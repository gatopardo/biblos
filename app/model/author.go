package model

import (
        "database/sql"
	"time"
        "fmt"
        "log"


)

// *****************************************************************************
// Author
// *****************************************************************************

// Author table contains the information for each user
type Author struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
	Name      string        `db:"name" bson:"name"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
}

  type IAuthor struct {
          Id          uint32        `db:"id" bson:"id,omitempty"`
	  Name       string        `db:"name" bson:"name"`
       }
// --------------------------------------------------------
// AuthorById tenemos el author dado id
func (author * Author)AuthorById() (err error) {
        stq  :=   "SELECT id, name, created_at, updated_at FROM authors WHERE id=$1"
	err = Db.QueryRow(stq, &author.Id).Scan(&author.Id, &author.Name, &author.CreatedAt, &author.UpdatedAt)

	return  standardizeError(err)
}
// --------------------------------------------------------
// AuthorById tenemos el lenguaje dado nombre
func (author * Author)AuthorByName() (err error) {
        stq  :=   "SELECT id, name, created_at, updated_at FROM authors WHERE name=$1"
	err = Db.QueryRow(stq, &author.Name).Scan(&author.Id, &author.Name, &author.CreatedAt, &author.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------
// AuthByRe listar autores dado  reg
func AuthorByRe(reStr string) (lsAuthors []IAuthor, err error) {
//       stRe :=  "'^" + reStr+".*'"
	stq := "select a.Id, a.name as author from authors a  where a.name ~* $1  order by  a.name "
             fmt.Println(stq)
	       var rows * sql.Rows
	       rows, err = Db.Query(stq, reStr)
               if err != nil {
                    log.Println(err)
               }else{
                  defer rows.Close()
                  for rows.Next() {
			  auth := IAuthor{}
                    err = rows.Scan(&auth.Id, &auth.Name)
                     if  err != nil {
			 log.Println( err)
                          break
		     }
                    lsAuthors = append(lsAuthors, auth)
	          }
       }
	return
  }


// --------------------------------------------------------
// AuthBookById tenemos id author obtener listado libros
 func AuthBookById(Id uint32) (lsBooks []BookZ, err error) {
	stq := "select b.title, b.comment, b.year, a.name as author, e.name as editor, l.name as language from books b join authors a  on a.id = b.author_id join editors e on e.id = editor_id join languages l on l.id = b.language_id where a.id = $1 order by  a.name,b.title, e.name "
//              fmt.Println(stq)
	       var rows * sql.Rows
	       rows, err = Db.Query(stq, Id)
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

// -----------------------------------------------------
// AuthorCreate crear author
  func (author *Author)AuthorCreate()( pos int, ex error ){
         var err error
         var stmt  *sql.Stmt
	 var auth Author
         stq := "INSERT INTO authors ( name, created_at, updated_at ) VALUES ($1,$2,$3) returning id"
	 now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	     log.Println(err)
         }else{
             defer stmt.Close()
             var id uint32
	     err = stmt.QueryRow(  &author.Name,   now, now ).Scan(&id)
             if err != nil {
	           log.Println(err)
	     } else{
               author.Id = id
	       stq = "SELECT id, name FROM authors order by name"
	       var rows * sql.Rows
	       rows, err = Db.Query(stq)
               if err != nil {
                    log.Println(err)
               }else{
                  defer rows.Close()
                  for rows.Next() {
                     auth = Author{}
                     err = rows.Scan(&auth.Id, &auth.Name )
		     pos++
                     if  err != nil {
			 log.Println( err)
                          break
		     }
	             if auth.Id == id {
		         return
		     }
	          }
	      }
             }
	 }
	ex =  standardizeError(err)
	return
   }
// -----------------------------------------------------
 func  (author * Author)AuthorDeleteById()( err error){
         stqd :=  "DELETE FROM authors where id = $1"
           _, err = Db.Exec(stqd, author.Id)
         return
       }
// Delete author from databa
func (author *Author)Delete() (err error) {
	statement := "delete from authors where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(author.Id)
	return
}
// -----------------------------------------------------
// Actualizar informacion de author en la database
func (author *Author)Update() (err error) {
	statement := "update authors set name = $2 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(author.Id, author.Name)
	return
}

// Delete all authors from database
func AuthorDeleteAll() (err error) {
	statement := "delete from authors"
	_, err = Db.Exec(statement)
	return
}
// -------------------------------------------------------------
// Get number of records in authors
  func AuthorCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM authors "
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
// Get selected authors in the database and returns the list
  func SAuthors(rsearch string  ) (authors []Author, err error) {
        var author Author
       stqi :=   "SELECT id,  name, created_at, updated_at FROM authors where "
       stqm :=   " name ~* '"+ rsearch + "' "
       stqf :=   " order by  name"
       stq := stqi + stqm + stqf
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
            author = Author{}
            if err = rows.Scan(&author.Id,  &author.Name, &author.CreatedAt, &author.UpdatedAt); err != nil {
			return
		}
		authors = append(authors, author)
	}
	return
 }
// -------------------------------------------------------------
// Get limit records from offset
  func AuthorLim(lim int , offs int) (authors []Author, err error) {
        var author Author
        stq :=   "SELECT id, name, created_at, updated_at FROM authors order by name LIMIT $1 OFFSET $2"

	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
                author = Author{}
		if err = rows.Scan(&author.Id,  &author.Name, &author.CreatedAt, &author.UpdatedAt); err != nil {
			return
		}
		authors = append(authors, author)
	}
	return
 }
// -------------------------------------------------------------
// Get all authors in the database and returns the list
  func Authors() (authors []Author, err error) {
        var author Author
        stq :=   "SELECT id,  name, created_at, updated_at FROM authors order by  name"
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&author.Id, &author.Name, &author.CreatedAt, &author.UpdatedAt); err != nil {
		return
		}
		authors = append(authors, author)
	}
	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get all authors id and name in the database and returns the list
  func IAuthors() (authors []IAuthor, err error) {
        var author IAuthor
        stq :=   "SELECT id,  name FROM authors order by  name"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&author.Id, &author.Name ); err != nil {
                       return
		}
		authors = append(authors, author)
	}
	return
 }
// -------------------------------------------------------------
