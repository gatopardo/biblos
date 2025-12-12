package controller

import (
	"log"
	"net/http"
        "strconv"
          "strings"
          "fmt"
        "encoding/json"

	"github.com/gatopardo/biblos/app/model"
	"github.com/gatopardo/biblos/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
// ---------------------------------------------------
// JAuthorListGET listar autores
 func JAuthListGET(w http.ResponseWriter, r *http.Request) {
        var err error
	var lisAuthors []model.IAuthor
	var params httprouter.Params
	sess := model.Instance(r)
	params      = context.Get(r, "params").(httprouter.Params)
        reauth    := params.ByName("re")
//      reauth     :=  strToReg(reparam)
//   fmt.Println("JAuthorListGET ", reauth)
	lisAuthors,err    = model.AuthorByRe(reauth)
   fmt.Println("JAuthorListGET ",lisAuthors)
         if err != nil {
              sess.AddFlash(view.Flash{"Error en listado autores ", view.FlashError})
	      log.Println(err)
	}else{
            var js []byte
            js, err =  json.Marshal(lisAuthors)
            if err == nil{
//    fmt.Println(string(js))
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
        }
    }

// ---------------------------------------------------
// JAuthBookListGET listar libros dado author
 func JAuthBookListGET(w http.ResponseWriter, r *http.Request) {
        var err error
	var lisBooks []model.BookZ
	var params httprouter.Params
	sess := model.Instance(r)
	params      = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))

	lisBooks,err    = model.AuthBookById(Id)
//       fmt.Println("JAuthBookListGET ",lisBooks)
         if err != nil {
              sess.AddFlash(view.Flash{"Error listado libros ", view.FlashError})
	      log.Println(err)
	}else{
            var js []byte
            js, err =  json.Marshal(lisBooks)
            if err == nil{
//         fmt.Println(string(js))
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
        }
    }

// ---------------------------------------------------
// AuthorGET despliega la pagina del usuario
func AuthorGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "author/author"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
 }
// ---------------------------------------------------
// ---------------------------------------------------
// AuthorPOST procesa la forma enviada con los datos
   func AuthorPOST(w http.ResponseWriter, r *http.Request) {
        var author model.Author
	var pos int
	sess := model.Instance(r)
        action        := r.FormValue("action")
	pagi := 1
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    if validate, missingField := view.Validate(r, []string{"name"}); !validate {
               sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
               sess.Save(r, w)
               AuthorGET(w, r)
               return
	}
	author.Name      = r.FormValue("name")
        err := (&author).AuthorByName()
        if err == model.ErrNoResult { // Exito:  no hay author creado aun 
            pos, err = (&author).AuthorCreate()
	    if err != nil {  // uyy como fue esto ? 
              log.Println(err)
              sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
              sess.Save(r, w)
	    } else {  // todo bien
               sess.AddFlash(view.Flash{"Author creado: " +author.Name, view.FlashSuccess})
               pagi = pos /limit
	       if pos % limit > 0{
		       pagi += 1
	       }
	   }
        }
       }
        sess.Save(r, w)
	sp := strconv.Itoa(pagi)
	http.Redirect(w, r, "/author/list/" + sp, http.StatusFound)
 }

// ---------------------------------------------------
// AuthorUpGET despliega la pagina del usuario
func AuthorUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var author model.Author
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
        author.Id = id
	err := (&author).AuthorById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No tenemos author.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, "/author/list/1", http.StatusFound)
           return
	}
	// Display the view
	v := view.New(r)
	v.Name = "author/authorupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Vars["name"] = author.Name
//    Refill any form fields
//	view.Repopulate([]string{"name"}, r.Form, v.Vars)
        v.Render(w)
   }
// ---------------------------------------------------

// AuthorUpPOST procesa la forma enviada con los datos
func AuthorUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var author model.Author
	sess := model.Instance(r)
	if validate, missingField := view.Validate(r, []string{"name"}); !validate {
		sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		AuthorUpGET(w, r)
		return
	}
        var params httprouter.Params
	params         = context.Get(r, "params").(httprouter.Params)
	author.Id, _   = atoi32(params.ByName("id"))
        SPag          := params.ByName("pagi")
        path          :=  fmt.Sprintf("/author/list/%s", SPag)
	author.Name    = r.FormValue("name")
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err =  author.Update()
             if err == nil{
                 sess.AddFlash(view.Flash{"Autor actualizado exitosamente para: " +author.Name, view.FlashSuccess})
             } else {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
	     }
		sess.Save(r, w)
	}
	http.Redirect(w, r, path, http.StatusFound)
     }
//------------------------------------------------
// AuthorLisGET displays the author page
func AuthorLisGET(w http.ResponseWriter, r *http.Request) {
        var params httprouter.Params
	sess       := model.Instance(r)
        params      = context.Get(r, "params").(httprouter.Params)
        SId        := params.ByName("id")
        Id,_       := atoi32(SId)
        posact      = int(Id)
        offset      =  posact  - 1
        TotalCount := model.AuthorCount()
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay lenguajes.", view.FlashError})
             sess.Save(r, w)
             return
          }else{
        offset = offset * limit
        lisAuthors, err := model.AuthorLim(limit, offset)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Authors.", view.FlashError})
            sess.Save(r, w)
         }
	v                  := view.New(r)
	v.Name              = "author/authorlis"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        numberOfBtns       :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns          :=  createSliceForBtns(numberOfBtns, posact)
        v.Vars["slice"]     =  sliceBtns
        v.Vars["current"]   =  posact
        v.Vars["LisAuthor"] =  lisAuthors
        v.Vars["Level"]     =  sess.Values["level"]
	v.Render(w)
      }
 }

//------------------------------------------------
// AuthorDeleteGET handles the note deletion
 func AuthorDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        v    := view.New(r)
        v.Name ="author/authordelete"
        var author model.Author
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        path :=  fmt.Sprintf("/author/list/%s", SPag)
        author.Id = Id
        err := (&author).AuthorById()
        if err != nil{
                log.Println(err)
		sess.AddFlash(view.Flash{"Error Autor no hallado.", view.FlashError})
                http.Redirect(w, r, path, http.StatusFound)
          }else{
            canti := model.BookAuthCount(Id)
            v.Vars["name"]  = author.Name
            v.Vars["canti"]  = fmt.Sprintf("%d",canti)
       }
        v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Render(w)
  }

//------------------------------------------------
// AuthorDeletePOST handles the note deletion
 func AuthorDeletePOST(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var author model.Author
        var params httprouter.Params
        params      = context.Get(r, "params").(httprouter.Params)
        Id,_       := atoi32(params.ByName("id"))
	SPag       := params.ByName("pagi")
        action     := r.FormValue("action")
         path      :=  fmt.Sprintf("/author/list/%s", SPag)
        if ! (strings.Compare(action,"Cancelar") == 0) {
            author.Id = Id
            err := author.Delete()
            if err != nil {
                log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
            } else {
                     sess.AddFlash(view.Flash{"Author borrado!", view.FlashSuccess})
             }
          }
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
        }
//------------------------------------------------
// AuthorSearchPOST procesa la forma enviada con los datos
func AuthorSearchPOST(w http.ResponseWriter, r *http.Request) {
	sess       :=  model.Instance(r)
        path       :=  fmt.Sprintf("/author/list/%s", "1")
        rSearch    :=  r.FormValue("bsearch")
        if rSearch == ""{
           sess.Save(r, w)
	   http.Redirect(w, r, path, http.StatusFound)
           return
         }

        lisAuthors, err := model.SAuthors(rSearch)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Autores.", view.FlashError})
            sess.Save(r, w)
	    http.Redirect(w, r, path, http.StatusFound)
         }
        posact     = 0
        offset     =  posact  - 1
        TotalCount :=  len(lisAuthors)
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay autores.", view.FlashError})
             sess.Save(r, w)
	     http.Redirect(w, r, path, http.StatusFound)
             return
         }
	// Display the view
	v := view.New(r)
	v.Name             = "author/authorlis"
	v.Vars["token"]    = csrfbanana.Token(w, r, sess)
        v.Vars["Level"]    =  sess.Values["level"]
        v.Vars["Search"]     = rSearch
        v.Vars["LisAuthor"]  = lisAuthors
        v.Vars["current"]  =  posact
	v.Render(w)
     }

//------------------------------------------------
//  AuthorPdfGet generates pdf from model.Author
func AuthorPdfGET(w http.ResponseWriter, r *http.Request) {
      lsAuths, _  := model.Authors( )
      pdf := getAuthorPdf(lsAuths )
      pdf.Output(w )

  }

