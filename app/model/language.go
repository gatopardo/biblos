package model

import (
        "database/sql"
	"time"
         "strings"
//	"fmt"
        "log"

)

// *****************************************************************************
// Language
// *****************************************************************************

// Language table contains the information for each user
type Language struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
	Name      string        `db:"name" bson:"name"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
}

       type ILang struct {
          Id          uint32        `db:"id" bson:"id,omitempty"`
          Name       string        `db:"name" bson:"name"`
        }

// --------------------------------------------------------

// LangById tenemos el lenguaje dado id
func (lang * Language)LangById() (err error) {
        stq  :=   "SELECT id, name, created_at, updated_at FROM languages WHERE id=$1"
	err = Db.QueryRow(stq, &lang.Id).Scan(&lang.Id, &lang.Name, &lang.CreatedAt, &lang.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------

// LangById tenemos el lenguaje dado nombre
func (lang * Language)LangByName() (err error) {
        stq  :=   "SELECT id, name, created_at, updated_at FROM languages WHERE name=$1"
	err = Db.QueryRow(stq, &lang.Name).Scan(&lang.Id, &lang.Name, &lang.CreatedAt, &lang.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------
// LangByRe tenemos el book dado re
func LangByRe(reStr string) (lsBooks []BookZ, err error) {
	stq := "select b.title, b.comment, b.year, a.name as author, e.name as editor, l.name as language from books b join authors a  on a.id = b.author_id join editors e on e.id = editor_id join languages l on l.id = b.language_id where l.name ~  $1 order by  l.name,b.title, e.name "
	       var rows * sql.Rows
	       rows, err = Db.Query(stq,reStr)
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
// LangCreate crear language
func (lang *Language)LangCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO languages ( name, created_at, updated_at ) VALUES ($1,$2,$3) returning id" 

	now  := time.Now()

            if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
              }
             defer stmt.Close()
             var id uint32
		err = stmt.QueryRow(  &lang.Name,   now, now ).Scan(&id)
            if err == nil {
               lang.Id = id
             }
	return standardizeError(err)
}

// -----------------------------------------------------
 func  (lang * Language)LangDeleteById()( err error){
         stqd :=  "DELETE FROM languages where id = $1"
           _, err = Db.Exec(stqd, lang.Id) 
         return
       }

// Delete language from databa
func (lang *Language)Delete() (err error) {
	statement := "delete from languages where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(lang.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de language en la database
func (lang *Language)Update() (err error) {
	statement := "update languages set name = $2 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(lang.Id, lang.Name)
	return
}

// Delete all languages from database
func LangDeleteAll() (err error) {
	statement := "delete from languages"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in languages
  func LangCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM languages "   
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
// Get limit records from offset
  func LangLim(lim int , offs int) (languages []Language, err error) {
        var lang Language
        stq :=   "SELECT id, name, created_at, updated_at FROM languages order by name LIMIT $1 OFFSET $2"   

	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
                lang = Language{} 
		if err = rows.Scan(&lang.Id,  &lang.Name, &lang.CreatedAt, &lang.UpdatedAt); err != nil {
                     (&lang).Name = strings.Trim(lang.Name, " ")
			return
		}
		languages = append(languages, lang)
	}
	return
 }
// -------------------------------------------------------------
// Get all languages in the database and returns the list
  func Langs() (languages []Language, err error) {
        var lang Language
        stq :=   "SELECT id,  name, created_at, updated_at FROM languages order by  name"   
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
           if err = rows.Scan(&lang.Id, &lang.Name, &lang.CreatedAt, &lang.UpdatedAt); err != nil {
                     (&lang).Name = strings.Trim(lang.Name, " ")
		return
		}
           languages = append(languages, lang)
	}
	return
 }
// -------------------------------------------------------------
// Get all languages id and name in  the database and returns the list
  func ILangs() (languages []ILang, err error) {
        var lang ILang
        stq :=   "SELECT id,  name FROM languages order by  name"   
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
           if err = rows.Scan(&lang.Id, &lang.Name); err != nil {
           	return
		}
           languages = append(languages, lang)
	}
	return
 }
// -------------------------------------------------------------
