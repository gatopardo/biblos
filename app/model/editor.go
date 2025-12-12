package model

import (
        "database/sql"
	"time"
//	"fmt"
        "log"
)

// *****************************************************************************
// Editor
// *****************************************************************************

// Editor table contains the information for each user
type Editor struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
	Name      string        `db:"name" bson:"name"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
}

  type IEditor struct {
          Id          uint32        `db:"id" bson:"id,omitempty"`
          Name       string        `db:"name" bson:"name"`
       }

// --------------------------------------------------------

// EditById tenemos el editor dado id
func (edit * Editor)EditById() (err error) {
        stq  :=   "SELECT id, name, created_at, updated_at FROM editors WHERE id=$1"
	err = Db.QueryRow(stq, &edit.Id).Scan(&edit.Id, &edit.Name, &edit.CreatedAt, &edit.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------

// EditById tenemos el editor dado nombre
func (edit * Editor)EditByName() (err error) {
        stq  :=   "SELECT id, name, created_at, updated_at FROM editors WHERE name=$1"
	err = Db.QueryRow(stq, &edit.Name).Scan(&edit.Id, &edit.Name, &edit.CreatedAt, &edit.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------
// EditByRe listar editoras dado reg
func EditByRe(reStr string) (lsEditors []IEditor, err error) {
//      stRe :=  "'^" + reStr+".*'"
	stq := "select e.id, e.name as editor from editors e  where e.name ~*  $1 order by  e.name "
//	       fmt.Println(stq)
	       var rows * sql.Rows
	       rows, err = Db.Query(stq, reStr)
               if err != nil {
                    log.Println(err)
               }else{
                  defer rows.Close()
                  for rows.Next() {
			  edit := IEditor{}
                    err = rows.Scan(&edit.Id, &edit.Name )
                     if  err != nil {
			 log.Println( err)
                          break
		     }
                    lsEditors = append(lsEditors, edit)
	          }
       }
	return
  }
// --------------------------------------------------------
// EditBookByRe tenemos id author obtener listado libros
 func EditBookById(Id uint32) (lsBooks []BookZ, err error) {
	stq := "select b.title, b.comment, b.year, a.name as author, e.name as editor, l.name as language from books b join authors a  on a.id = b.author_id join editors e on e.id = editor_id join languages l on l.id = b.language_id where e.id = $1 order by  a.name,b.title, e.name "
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
// EditCreate crear editor
func (edit *Editor)EditCreate()(pos int, ex error) {
         var err error
         var stmt  *sql.Stmt
	 var ed Editor
         stq := "INSERT INTO editors ( name, created_at, updated_at ) VALUES ($1,$2,$3) returning id"
	 now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	          log.Println(err)
          }else{
             defer stmt.Close()
             var id uint32
             err = stmt.QueryRow(  &edit.Name,   now, now ).Scan(&id)
             if err != nil {
                  log.Println(err)
	     }else{
               edit.Id = id
	       stq = "SELECT id, name FROM editors order by name"
	       var rows *sql.Rows
	       rows, err = Db.Query(stq)
               if err != nil {
                    log.Println(err)
               }else{
		  defer rows.Close()
                  for rows.Next() {
                     ed = Editor{}
                     err = rows.Scan(&ed.Id, &ed.Name )
		     pos++
                     if  err != nil {
			 log.Println( err)
                          break
		     }
	             if ed.Id == id {
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
 func  (edit * Editor)EditDeleteById()( err error){
         stqd :=  "DELETE FROM editors where id = $1"
           _, err = Db.Exec(stqd, edit.Id)
         return
       }

// Delete editor from databa
func (edit *Editor)Delete() (err error) {
	statement := "delete from editors where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(edit.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de editor en la database
func (edit *Editor)Update() (err error) {
	statement := "update editors set name = $2 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(edit.Id, edit.Name)
	return
}

// Delete all editors from database
func EditDeleteAll() (err error) {
	statement := "delete from editors"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in editors
  func EditCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM editors "
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
// Get selected editors in the database and returns the list
  func SEditors(rsearch string  ) (editors []Editor, err error) {
        var editor Editor
       stqi :=   "SELECT id,  name, created_at, updated_at FROM editors where "
       stqm :=   " name ~* '"+ rsearch + "' "
       stqf :=   " order by  name"
       stq := stqi + stqm + stqf
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
            editor = Editor{}
            if err = rows.Scan(&editor.Id,  &editor.Name, &editor.CreatedAt, &editor.UpdatedAt); err != nil {
			return
		}
		editors = append(editors, editor)
	}
	return
 }
// -------------------------------------------------------------
// Get limit records from offset


// -------------------------------------------------------------
// Get limit records from offset
  func EditLim(lim int , offs int) (editors []Editor, err error) {
        var edit Editor
        stq :=   "SELECT id, name, created_at, updated_at FROM editors order by name LIMIT $1 OFFSET $2"

	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
                edit = Editor{}
		if err = rows.Scan(&edit.Id,  &edit.Name, &edit.CreatedAt, &edit.UpdatedAt); err != nil {
			return
		}
		editors = append(editors, edit)
	}
	return
 }
// -------------------------------------------------------------
// Get all editors in the database and returns the list
  func Edits() (editors []Editor, err error) {
        var edit Editor
        stq :=   "SELECT id,  name, created_at, updated_at FROM editors order by  name"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&edit.Id, &edit.Name, &edit.CreatedAt, &edit.UpdatedAt); err != nil {
		return
		}
		editors = append(editors, edit)
	}
	return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
// Get all editors id and name in the database and returns the list
  func IEdits() (editors []IEditor, err error) {
        var edit IEditor
        stq :=   "SELECT id,  name FROM editors order by  name"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&edit.Id, &edit.Name); err != nil {

		}
		editors = append(editors, edit)
	}
	return
 }
// -------------------------------------------------------------
