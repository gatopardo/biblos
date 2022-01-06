package controller

import (
	"log"
	"net/http"
        "strings"
	"strconv"
        "fmt"
	"encoding/json"

	"github.com/gatopardo/biblos/app/model"
	"github.com/gatopardo/biblos/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)
  // ---------------------------------------------------
// JEditGET crea datos de libros
 func JEditorListGET(w http.ResponseWriter, r *http.Request) {
        var err error
	var lisBooks []model.BookZ
	var params httprouter.Params
	sess := model.Instance(r)
	params      = context.Get(r, "params").(httprouter.Params)
	reedit      := params.ByName("re")
	lisBooks,err    = model.EditByRe(reedit)
         if err != nil {
              sess.AddFlash(view.Flash{"Error listado libros ", view.FlashError})
	      log.Println(err)
	}else{
            var js []byte
            js, err =  json.Marshal(lisBooks)
            if err == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
        }
    }
// ---------------------------------------------------
// view.Repopulate([]string{"name"}, r.Form, v.Vars)

// EditGET despliega la pagina del usuario
func EditGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
	v := view.New(r)
	v.Name = "editor/editor"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
 }
// ---------------------------------------------------
// ---------------------------------------------------
// EditPOST procesa la forma enviada con los datos
func EditPOST(w http.ResponseWriter, r *http.Request) {
        var edit model.Editor
	var pos int
	sess      := model.Instance(r)
        action    := r.FormValue("action")
	pagi := 1
        if ! (strings.Compare(action,"Cancelar") == 0) {
           if validate, missingField := view.Validate(r, []string{"name"}); !validate {
               sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
               sess.Save(r, w)
               EditGET(w, r)
               return
           }
           edit.Name      = r.FormValue("name")
	   err := (&edit).EditByName()
           if err == model.ErrNoResult { // Exito:  no hay editor creado aun 
               pos, err = (&edit).EditCreate()
	       if err != nil {  // uyy como fue esto ? 
                   log.Println(err)
                   sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
	       } else {  // todo bien
                   sess.AddFlash(view.Flash{"Editor creado: " +edit.Name, view.FlashSuccess})
		   pagi = pos / limit
		   if pos % limit > 0{
			    pagi += 1
		   }
	   }
        }
      }
       sess.Save(r, w)
       sp := strconv.Itoa(pagi)
      http.Redirect(w, r, "/editor/list/"+ sp, http.StatusFound)
 }

// ---------------------------------------------------
// EditUpGET despliega la pagina del usuario
func EditUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var edit model.Editor
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
        edit.Id = id
	err := (&edit).EditById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No tenemos editor.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, "/editor/list/1", http.StatusFound)
           return
	}
	v               := view.New(r)
	v.Name           = "editor/editupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Vars["name"]   = edit.Name
//	view.Repopulate([]string{"name"}, r.Form, v.Vars)
        v.Render(w)
   }
// ---------------------------------------------------

// EditUpPOST procesa la forma enviada con los datos
func EditUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
         var edit model.Editor
	sess := model.Instance(r)
	if validate, missingField := view.Validate(r, []string{"name"}); !validate {
		sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		EditUpGET(w, r)
		return
	}
	edit.Name = r.FormValue("name")
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	edit.Id, _ = atoi32(params.ByName("id"))
        SPag        := params.ByName("pagi")
             err =  edit.Update()
             if err == nil{
                 sess.AddFlash(view.Flash{"Lenguaje actualizado exitosamente para: " +edit.Name, view.FlashSuccess})
             } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
	}
	sess.Save(r, w)
        path :=  fmt.Sprintf("/editor/list/%s", SPag)
	http.Redirect(w, r, path, http.StatusFound)
//	EditUpGET(w, r)
     }
//------------------------------------------------
// EditLisGET displays the editor page
func EditLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        SId := params.ByName("id")
        Id,_ := atoi32(SId)
        posact = int(Id)
        offset = int(Id)  - 1
        TotalCount := model.EditCount()
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay lenguajes.", view.FlashError})
             sess.Save(r, w)
             return
          }else{
        offset = offset * limit
        lisEdits, err := model.EditLim(limit, offset)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Editores.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "editor/editlis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
        numberOfBtns        :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns           :=  createSliceForBtns(numberOfBtns, posact) 
        v.Vars["slice"]      =  sliceBtns
        v.Vars["current"]    =  posact
        v.Vars["Level"]      =  sess.Values["level"]
        v.Vars["LisEdit"]    = lisEdits
	v.Render(w)
      }
 }

// EditDeleteGET handles the note deletion
 func EditDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        v    := view.New(r)
        v.Name ="editor/editordelete"
        var edit model.Editor
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
        Id,_        := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        path        :=  fmt.Sprintf("/editor/list/%s", SPag)
        edit.Id      = Id
        err         := (&edit).EditById()
        if err != nil{
            log.Println(err)
            sess.AddFlash(view.Flash{"Error Editor no hallado.", view.FlashError})
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
            canti := model.BookEditCount(Id)
            v.Vars["name"]  = edit.Name
            v.Vars["canti"]  = fmt.Sprintf("%d",canti)
        v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Render(w)
  }

//------------------------------------------------
// EditorDeletePOST handles the note deletion
 func EditDeletePOST(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var editor model.Editor
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
        Id,_        := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        action      := r.FormValue("action")
        path        :=  fmt.Sprintf("/editor/list/%s", SPag)
        if ! (strings.Compare(action,"Cancelar") == 0) {
            editor.Id = Id
            err      := editor.Delete()
            if err   != nil {
                log.Println(err)
                sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
            } else {
                sess.AddFlash(view.Flash{"Editor borrado!", view.FlashSuccess})
            }
         }
         sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
        }

//------------------------------------------------
// EditorSearchPOST procesa la forma enviada con los datos
func EditorSearchPOST(w http.ResponseWriter, r *http.Request) {
	sess       :=  model.Instance(r)
        path       :=  fmt.Sprintf("/editor/list/%s", "1")
        rSearch    :=  r.FormValue("bsearch")
        if rSearch == ""{
           sess.Save(r, w)
	   http.Redirect(w, r, path, http.StatusFound)
           return
         }

        lisEditors, err := model.SEditors(rSearch)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Editores.", view.FlashError})
            sess.Save(r, w)
	    http.Redirect(w, r, path, http.StatusFound)
         }
        posact     = 0
        offset     =  posact  - 1
        TotalCount :=  len(lisEditors)
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay editores.", view.FlashError})
             sess.Save(r, w)
	     http.Redirect(w, r, path, http.StatusFound)
             return
         }
	// Display the view
	v := view.New(r)
	v.Name             = "editor/editlis"
	v.Vars["token"]    = csrfbanana.Token(w, r, sess)
        v.Vars["Level"]    =  sess.Values["level"]
        v.Vars["Search"]   = rSearch
        v.Vars["LisEdit"]  = lisEditors
        v.Vars["current"]  =  posact
	v.Render(w)
     }
//------------------------------------------------

