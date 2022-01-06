package controller

import (
	"log"
	"net/http"
        "strings"
        "fmt"

	"github.com/gatopardo/biblos/app/model"
	"github.com/gatopardo/biblos/app/shared/passhash"
//	"github.com/gatopardo/biblos/app/shared/recaptcha"
	"github.com/gatopardo/biblos/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )

  // ---------------------------------------------------

// RegisterGET despliega la pagina del usuario
func RegisterGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "register/register"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
// view.Repopulate([]string{"cuenta", "password", "level"}, r.Form, v.Vars)
	v.Render(w)
 }
// ---------------------------------------------------
// ---------------------------------------------------
// RegisterPOST procesa la forma enviada con los datos
func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
        var user model.User
        var shad model.Shadow
	sess := model.Instance(r)
//        var params httprouter.Params
//	params = context.Get(r, "params").(httprouter.Params)
//        action      := params.ByName("action")
        action        := r.FormValue("action")
//  fmt.Println(action)
        if ! (strings.Compare(action,"Cancelar") == 0) {
            // Validate with required fields
           if validate, missingField := view.Validate(r, []string{"cuenta", "level"}); !validate {
               sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
               sess.Save(r, w)
               RegisterGET(w, r)
               return
	   }
           // Validate with Google reCAPTCHA
//	   if !recaptcha.Verified(r) {
//		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
//		sess.Save(r, w)
//		RegisterGET(w, r)
//		return
//	   }
           stUuid := model.CreateUUID()
           // Get form values
           rPasswd         := r.FormValue("password")
	   user.Cuenta      = r.FormValue("cuenta")
           vPasswd         := r.FormValue("password_verify")
           user.Level, _   = atoi32( r.FormValue("level"))
           user.Uuid       = stUuid
           if strings.Compare(rPasswd, vPasswd) != 0{
		log.Println(rPasswd, " * ", vPasswd)
		sess.AddFlash(view.Flash{"Claves distintas no posible", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/user/register", http.StatusFound)
		return
           }
            pass, errp := passhash.HashString(rPasswd)
	   // If password hashing failed
	   if errp != nil {
		log.Println(errp)
                sess.AddFlash(view.Flash{"Problema encriptando clave.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/user/register", http.StatusFound)
		return
	   }
           user.Password = pass
	   err := (&user).UserByCuenta()
           if err == model.ErrNoResult { // Exito:  no hay usuario creado aun 
              ex := (&user).UserCreate()
//            log.Println("Creating user")
	      if ex != nil {  // uyy como fue esto ? 
                 log.Println(ex)
                 sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                 sess.Save(r, w)
                 return
	       } else {  // todo bien
                 shad.UserId    = user.Id
                 shad.Uuid      = stUuid
                 shad.Password  = rPasswd
                 if  err = (&shad).ShadCreate() ; err != nil{
                     sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                     log.Println( err)
                     sess.Save(r, w)
                    return
                 }
                 sess.AddFlash(view.Flash{"Cuenta creada: " +user.Cuenta, view.FlashSuccess})
                 sess.Save(r, w)
	      }
           }
         }
	// Display list
	http.Redirect(w, r, "/user/list/1", http.StatusFound)
  }

// ---------------------------------------------------
// RegisUpGET despliega la pagina del usuario
func RegisUpGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        var user model.User
	// necesitamos user id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
         user.Id = id
	// Obtener usuario dado id
	 err := (&user).UserById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
      sess.AddFlash(view.Flash{"Es raro. No tenemos usuario.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, "/user/list/1", http.StatusFound)
           return
	}
// log.Println("RegisUpGet to render view ", user.Cuenta)
	// Display the view
	v := view.New(r)
	v.Name = "register/regisupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Vars["cuenta"] = user.Cuenta
        v.Vars["ulevel"]  = user.Level
//    Refill any form fields
//	view.Repopulate([]string{"cuenta", "level"}, r.Form, v.Vars)
        v.Render(w)
   }
// ---------------------------------------------------

// RegisUpPOST procesa la forma enviada con los datos
func RegisUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
         var user model.User
	// Get session
	sess := model.Instance(r)
//	 Validate with required fields
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	user.Id, _ = atoi32(params.ByName("id"))
        SPag        := params.ByName("pagi")
        path :=  fmt.Sprintf("/user/list/%s", SPag)
        action      := params.ByName("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	if validate, missingField := view.Validate(r, []string{"cuenta", "level"}); !validate {
		sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		RegisUpGET(w, r)
		return
	}
	// Get form values
	user.Cuenta     = r.FormValue("cuenta")
        user.Level, _   = atoi32( r.FormValue("level"))
             err =  user.Update()
             if err == nil{
                 sess.AddFlash(view.Flash{"Cuenta actualizada exitosamente para: " +user.Cuenta, view.FlashSuccess})
		sess.Save(r, w)
             } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
		sess.Save(r, w)
             }
       }
	http.Redirect(w, r, path, http.StatusFound)
     }
//------------------------------------------------
// RegisLisGET displays the register page
func RegisLisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        SId := params.ByName("id")
        Id,_ := atoi32(SId)

        posact = int(Id)
        offset =  posact   -  1
        TotalCount := model.UsersCount()
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay usuarios.", view.FlashError})
             sess.Save(r, w)
             return
          }else{
        offset = offset * limit

        lisUsers, err := model.UsersLim(limit, offset)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Usuarios.", view.FlashError})
           sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "register/regislis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
       lusers := make([]SUser, len(lisUsers))
        lev, ok := v.Vars["level"].(uint32)
        if ok {
           for i, usu := range lisUsers  {
              lusers[i].User  =  usu
              lusers[i].Level = lev
              lusers[i].Pagi = posact
           }

           v.Vars["LisRegis"] = lusers
           v.Vars["level"]    = lev
        }
        numberOfBtns     :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns        :=  createSliceForBtns(numberOfBtns, posact)
        v.Vars["slice"]   =  sliceBtns
        v.Vars["current"] =  posact
	v.Render(w)
      }
 }

//------------------------------------------------
// UserDeleteGET handles the note deletion
 func RegisDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
        sess := model.Instance(r)
        v    := view.New(r)
        v.Name = "register/registerdelete"
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_   := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        path :=  fmt.Sprintf("/user/list/%s", SPag)
// fmt.Println(Id,SPag)
        var user model.User
        user.Id = Id
        err := (&user).UserById()
        if err != nil {
            log.Println(err)
            fmt.Println(err)
            sess.AddFlash(view.Flash{"Error Usuario no hallado.", view.FlashError})
                http.Redirect(w, r, path, http.StatusFound)
        }else{
           v.Vars["cuenta"]  = user.Cuenta
           v.Vars["level"]    = user.Level
      }
// fmt.Println(path)
	   v.Vars["token"]  = csrfbanana.Token(w, r, sess)
           v.Render(w)

  }
//------------------------------------------------
// UserDeletePOST handles the note deletion
 func RegisDeletePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
        sess := model.Instance(r)
        var user model.User
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        action      := params.ByName("action")
        path :=  fmt.Sprintf("/user/list/%s", SPag)
        if ! (strings.Compare(action,"Cancelar") == 0) {
            user.Id = Id
            err := (&user).UserById()
if err != nil {
log.Println(err)
sess.AddFlash(view.Flash{"Error Usuario no hallado.", view.FlashError})
            }else{
	        // Get database result
// 	        err := user.UserDelete()
// 	        if err != nil {
// 		      log.Println(err)
// 		      sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
// 	        } else {
// 		        sess.AddFlash(view.Flash{"Usuario borrado!", view.FlashSuccess})
//              }
            }
         }
        sess.Save(r, w)
        http.Redirect(w, r, path, http.StatusFound)
  }
//------------------------------------------------


