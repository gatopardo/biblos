package model

import (
        "database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"time"
	"strings"
	"crypto/rand"
        "log"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
	Uuid        string        `db:"uuid" bson:"uuid,omitempty"`
	Cuenta      string        `db:"cuenta" bson:"cuenta"`
	Password    string        `db:"password" bson:"password"`
	Level       uint32        `db:"level" bson:"level"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
}

// User table contains text password
type Shadow struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
	UserId      uint32        `db:"userid" bson:"id,omitempty"`
	Uuid        string        `db:"uuid" bson:"uuid,omitempty"`
	Password    string        `db:"password" bson:"password"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
}
   type Jperson struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
	Cuenta      string        `db:"cuenta" bson:"cuenta"`
	Uuid        string        `db:"uuid" bson:"uuid,omitempty"`
	Nivel       uint32        `db:"nivel" bson:"nivel"`
	Email       string        `db:"email" bson:"email"`
   }


// --------------------------------------------------------
// Crear una nueva sesion para un usuario
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid,  user_id, created_at, updated_at) values ($1, $2, $3, $4) returning id, uuid,  user_id, created_at, updated_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// usar QueryRow para retornar una fila y buscar el id para struct Session 
	err = stmt.QueryRow(CreateUUID(),  user.Id, time.Now()).Scan(&session.Id, &session.Uuid,  &session.UserId, &session.CreatedAt)
	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid,  user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
		Scan(&session.Id, &session.Uuid,  &session.UserId, &session.CreatedAt)
	return
}

// --------------------------------------------------------
// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, user_id, created_at FROM sessions WHERE uuid = '$1'", session.Uuid). Scan(&session.Id, &session.Uuid,  &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// --------------------------------------------------------
// Delete session from database
func (session *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return
}

// --------------------------------------------------------
// Obtener usuario desde la la session
func (session *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, cuenta,password level  created_at FROM users WHERE id = $1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Cuenta, &user.Password, &user.Level, &user.CreatedAt)
	return
}

// Delete all sessions from database
func SessionDeleteAll() (err error) {
	statement := "delete from sessions"
	_, err = Db.Exec(statement)
	return
}

// -----------------------------------------

// UserById tenemos el usuario dado id
func (user * User)UserById() (err error) {
//        stq  :=   "SELECT id, uuid, cuenta, password,level FROM users WHERE cuenta=? LIMIT 1"
        stq  :=   "SELECT id, uuid, cuenta, password,level, created_at FROM users WHERE id=$1"
	switch ReadConfig().Type {
	case TypeMySQL:
		err = Db.QueryRow(stq, &user.Id).Scan(&user.Id, &user.Uuid, &user.Cuenta, &user.Password, &user.Level, &user.CreatedAt)

	case TypePostgreSQL:
	err = Db.QueryRow(stq, &user.Id).Scan(&user.Id, &user.Uuid,  &user.Cuenta, &user.Password, &user.Level, &user.CreatedAt)

	default:
	err = ErrCode
	}
	return  standardizeError(err)
}

// --------------------------------------------------------

// UserByCuenta gets user information from cuenta
func (user *User)UserByCuenta() ( error) {
	var err error
//        stq  :=   "SELECT id, uuid, cuenta, password,level FROM users WHERE cuenta=? LIMIT 1"
        stq  :=   "SELECT id, uuid, cuenta, password,level, created_at FROM users WHERE cuenta=$1"
//	result := User{
	switch ReadConfig().Type {
	case TypeMySQL:
//		err = SQL.Get(&result, stq, cuenta)
		err = Db.QueryRow(stq, &user.Cuenta).Scan(&user.Id, &user.Uuid, &user.Cuenta, &user.Password, &user.Level, &user.CreatedAt)

	case TypePostgreSQL:
		err = Db.QueryRow(stq, &user.Cuenta).Scan(&user.Id, &user.Uuid, &user.Cuenta, &user.Password, &user.Level, &user.CreatedAt)

	default:
		err = ErrCode
	}
	return   standardizeError(err)
}

// --------------------------------------------------------

// JpersByCuenta gets user information from cuenta
func (jpers *Jperson)JPersByCuenta()(pass string,  ex error) {
/*     var err error
        stq  :=   "SELECT  u.id,u.cuenta, u.password,u.uuid, u.level FROM users u WHERE u.cuenta=$1"
//	fmt.Println(stq)
         err = Db.QueryRow(stq,&jpers.Cuenta).Scan(&jpers.Id, &jpers.Cuenta, &pass, &jpers.Uuid, &jpers.Nivel )
     if err != nil {
	 jpers.Cuenta   = strings.Trim(jpers.Cuenta, " ")
	 jpers.Uuid     = strings.Trim(jpers.Uuid, " ")
//	 jpers.Email    = strings.Trim(jpers.Email, " ")
    }
     ex             = standardizeError(err)
    return

=======
*/
        stq  :=   "SELECT  u.id,u.cuenta, u.password,u.uuid, u.level FROM users u WHERE u.cuenta = $1"
	row := Db.QueryRow(stq,jpers.Cuenta)
	err := row.Scan(&jpers.Id, &jpers.Cuenta, &pass, &jpers.Uuid, &jpers.Nivel )
	ex     = standardizeError(err)
	switch  err {
       case sql.ErrNoRows:
	   fmt.Println("JpersBy Cuenta ex", ex)
        case nil:
	 jpers.Cuenta   = strings.Trim(jpers.Cuenta, " ")
	 jpers.Uuid     = strings.Trim(jpers.Uuid, " ")
//         fmt.Println("JpersByCuenta ", jpers.Cuenta, jpers.Uuid)
	 //	 jpers.Email    = strings.Trim(jpers.Email, " ")
         default:
              fmt.Println("JpersByCuenta panic", ex)
//               panic(err)
       }
	return
}

// -----------------------------------------------------
// UserCreate crear usuario
func (u *User)UserCreate() error {
         var err error
         var stmt  *sql.Stmt
	 stq := "INSERT INTO users (uuid, cuenta, password, level, created_at, updated_at ) VALUES ($1,$2,$3,$4, $5, $6) returning id"

	now  := time.Now()
//        stuu :=  createUUID()

            if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
              }
             defer stmt.Close()
             var id uint32
		err = stmt.QueryRow( &u.Uuid, &u.Cuenta, &u.Password, &u.Level, now, now ).Scan(&id)
            if err == nil {
               u.Id = id
             }
	return standardizeError(err)
}

// -----------------------------------------------------
// ShadCreate creates shadow
func (sd *Shadow)ShadCreate() error {
	var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO shadows (user_id, uuid,  password, created_at, updated_at ) VALUES ($1,$2,$3,$4, $5) returning id"
	now  := time.Now()
            if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
              }
             defer stmt.Close()
             var id uint32
             err = stmt.QueryRow(&sd.UserId, &sd.Uuid, &sd.Password,now, now ).Scan(&id)
            if err == nil {
               sd.Id = id
             }
	return standardizeError(err)
}

// -----------------------------------------------------
 func  (user * User)UserDeleteById()( err error){
         stqd :=  "DELETE FROM users where id = $1"
           _, err = Db.Exec(stqd, user.Id)
         return
       }

// Delete user from databa
func (user *User) UserDelete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de usuario en la database
func (user *User)Update() (err error) {
	statement := "update users set cuenta = $2, password = $3, level = $4 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Cuenta, user.Password, user.Level)
	return
}

// Delete all users from database
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in users
  func UsersCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM users "
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
  func UsersLim(lim int , offs int) (users []User, err error) {
        stq :=   "SELECT id, uuid, cuenta, password, level, created_at FROM users order by level, cuenta LIMIT $1 OFFSET $2"
	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Cuenta, &user.Password, &user.Level, &user.CreatedAt); err != nil {

			return
		}
		users = append(users, user)
	}
	return
 }
// -------------------------------------------------------------
// Get all users in the database and returns the list
  func Users() (users []User, err error) {
        stq :=   "SELECT id, uuid, cuenta, password, level, created_at FROM users order by level, cuenta"
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Cuenta, &user.Password, &user.Level, &user.CreatedAt); err != nil {

			return
		}
		users = append(users, user)
	}
	return
 }
// -------------------------------------------------------------
 func CreateUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Println("No se genera UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122 
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
 }




