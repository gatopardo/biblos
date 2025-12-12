package controller

import (
	"log"
	"net/http"
        "strings"
         "fmt"
	 "time"
	 "strconv"
        "encoding/json"

	"github.com/gatopardo/biblos/app/model"
	"github.com/gatopardo/biblos/app/shared/view"
        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
//	"github.com/jung-kurt/gofpdf"
  )
// ---------------------------------------------------

 type twoStr struct  {
      fname, dbvalue string
  }
 type mixVal struct {
      fname   string
      dbvalue uint32
  }

// ---------------------------------------------------
// BookGET crea datos de libros
 func JBookListGET(w http.ResponseWriter, r *http.Request) {
        var err error
	var lisBooks []model.BookZ
	var params httprouter.Params
	sess := model.Instance(r)
	params      = context.Get(r, "params").(httprouter.Params)
        rebook    := params.ByName("re")
//      rebook     := strToReg(reparam)
	lisBooks,err    = model.BookByRe(rebook)
//   fmt.Println("JBookListGET ",lisBooks)
         if err != nil {
              sess.AddFlash(view.Flash{"Error en listado libros ", view.FlashError})
	      log.Println(err)
	}else{
            var js []byte
            js, err =  json.Marshal(lisBooks)
            if err == nil{
//    fmt.Println(string(js))
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
        }
    }

// ---------------------------------------------------
// BookGET crea datos de libros
 func BookGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var err error
        var ilangs []model.ILang
        ilangs, err =  model.ILangs()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay lenguajes ", view.FlashError})
         }
	       var iedits  []model.IEditor
        iedits, err =  model.IEdits()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay Editoras ", view.FlashError})
         }
        var iauthors  []model.IAuthor
        iauthors, err =  model.IAuthors()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay Autores ", view.FlashError})
         }
	v := view.New(r)
	v.Name = "book/book"
        v.Vars["LisLang"]  = ilangs
        v.Vars["LisEdit"]  = iedits
        v.Vars["LisAuthor"]  = iauthors
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        view.Repopulate([]string{"title"}, r.Form, v.Vars)
	v.Render(w)
 }
// ---------------------------------------------------
//   getBookFormatData get data from form
     func getBookFormData(r *http.Request)(book model.Book){
	 book.Title          = r.FormValue("title")
	 book.Isbn           = r.FormValue("isbn")
	 book.Comment        = r.FormValue("comment")
	 book.Year, _        = atoi32(r.FormValue("year") )
         book.LanguageId, _  = atoi32(r.FormValue("language_id"))
         book.EditorId, _    = atoi32(r.FormValue("editor_id"))
         book.AuthorId, _    = atoi32(r.FormValue("author_id"))
         return
     }
// ---------------------------------------------------
// BookPOST procesa la forma enviada con los datos : crea libro
 func BookPOST(w http.ResponseWriter, r *http.Request) {
        var book model.Book
	var pos int
	sess := model.Instance(r)
        action        := r.FormValue("action")
	pagi := 1
        if ! (strings.Compare(action,"Cancelar") == 0) {
           if validate, missingField := view.Validate(r, []string{"title"}); !validate {
              sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
              sess.Save(r, w)
              BookGET(w, r)
              return
	    }
            book = getBookFormData(r)
            err := (&book).BookByName()
            if err == model.ErrNoResult { // Exito:  no hay book creado aun 
                pos, err = (&book).BookCreate()
	        if err != nil {  // uyy como fue esto ? 
                    log.Println(err)
                    sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
	         } else {  // todo bien
                          sess.AddFlash(view.Flash{"Book creado: " +book.Title, view.FlashSuccess})
			  pagi =  pos/limit
			  if  pos % limit > 0{
				  pagi += 1
			  }
	         }
             }
        }
        sess.Save(r, w)
	sp := strconv.Itoa(pagi)
	http.Redirect(w, r, "/biblos/list/"+sp, http.StatusFound)
 }

func getOrig(orig string)( dest string ){
       dest = "biblos"
      switch orig {
          case "0": dest = "biblos"
          case "1": dest = "author"
          case "2": dest = "editor"
          case "3": dest = "language"
        }
  return
}

// ---------------------------------------------------
// BookUpGET actualizar datos libros
func BookUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var book model.Book
	var params httprouter.Params
	params       = context.Get(r, "params").(httprouter.Params)
	id,_        := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
	SOrig       := params.ByName("orig")
        orig        := getOrig(SOrig)
        path        :=  fmt.Sprintf("/%s/list/%s",orig, SPag)
        book.Id      = id
	err := (&book).BookById()
	if err != nil { // Si no existe el libro
           log.Println(err)
            sess.AddFlash(view.Flash{"Es raro. No tenemos libro.", view.FlashError})
           sess.Save(r, w)
	   http.Redirect(w, r, path, http.StatusFound)
           return
	}
        var lang model.Language
        lang.Id   = book.LanguageId
        (&lang).LangById()
        var ilangs []model.ILang
        ilangs, err =  model.ILangs()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay lenguajes ", view.FlashError})
         }
        var edit model.Editor
        edit.Id      = book.EditorId
        (&edit).EditById()
        var iedits  []model.IEditor
        iedits, err =  model.IEdits()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay Editoras ", view.FlashError})
         }
        var author model.Author
        author.Id      = book.AuthorId
        (&author).AuthorById()
        var iauthors  []model.IAuthor
        iauthors, err =  model.IAuthors()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay Autores ", view.FlashError})
         }
        v                  := view.New(r)
	v.Name              = "book/bookupdate"
        v.Vars["Lang"]      = lang
        v.Vars["LisLang"]   = ilangs
        v.Vars["Edit"]      = edit
        v.Vars["LisEdit"]   = iedits
        v.Vars["Auth"]      = author
        v.Vars["LisAuthor"] = iauthors
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
//        v.Vars["title"]     = book.Title
        v.Vars["Book"]      = book
      v.Render(w)
   }
//----------------------------------------------------
//  gettin partial updates
     func getPartialUp(r *http.Request,fieldname,  dname string, ) ( sform string ){
             stTrimf  :=  strings.Trim(r.FormValue(fieldname), " ")
             stTrimp  :=   strings.Trim(dname, " ")
             if  ( stTrimf != stTrimp) {
                  sform = fmt.Sprintf(" %s  = '%s' ",fieldname, stTrimf )
             }
        return
      }

//----------------------------------------------------
//----------------------------------------------------
//  gettin partial updates
     func getPartialUpInt(r *http.Request,fieldname string,  dname uint32) ( sform string ){
             stTrimf  :=  strings.Trim(r.FormValue(fieldname), " ")
	     i, _  := atoi32(stTrimf)
	     if i  != dname {
                 sform = fmt.Sprintf( " %s = %d ", fieldname, i)
             }
        return
      }

//----------------------------------------------------
// ---------------------------------------------------
//   getBookFormUp get data from form
     func getBookFormUp(r *http.Request, bk model.Book)( st string){
         var  sf string
         var sup  []string
	 arrTwo :=  []twoStr{{"title", bk.Title},{"isbn", bk.Isbn},{"comment", bk.Comment}}
	 arrMix :=  []mixVal{{"year", bk.Year},{"language_id",bk.LanguageId}, {"editor_id", bk.EditorId}, {"author_id", bk.AuthorId}}
//	 bib :=  getBookFormData(r )
	 for _,item := range arrTwo{
	     sf = getPartialUp(r, item.fname, item.dbvalue)
	     if len(sf) >0{
                  sup = append(sup, sf)
	     }
	 }
	 for _,item := range arrMix{
	     sf = getPartialUpInt(r, item.fname, item.dbvalue)
	     if len(sf) >0{
                  sup = append(sup, sf)
	     }
	 }
/*
*/
         if len(sup) > 0 {
             layoutISO := "2006-01-02"
             t := time.Now()
	     sl := t.Format(layoutISO)
             sf = fmt.Sprintf( " updated_at = '%s' ", sl)
             sup = append(sup, sf)
             st =  strings.Join(sup, ", ")
//	     fmt.Println(st)
         }
         return st
     }
// ---------------------------------------------------
// BookUpPOST procesa la forma enviada con los datos
func BookUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var book model.Book
	sess := model.Instance(r)
	if validate, missingField := view.Validate(r, []string{"title"}); !validate {
		sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		BookUpGET(w, r)
		return
	}
        var params httprouter.Params
	params      = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	SPag        := params.ByName("pagi")
	SOrig       := params.ByName("orig")
        orig        := getOrig(SOrig)
        path        :=  fmt.Sprintf("/%s/list/%s",orig, SPag)
        action      := r.FormValue("action")
	id, _       := atoi32(SId)
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    book.Id      = id
            (&book).BookById()
	    sini        :=  "update books set "
            sr          :=  fmt.Sprintf(" where books.id = %d ", book.Id)
            st          :=  getBookFormUp(r, book)
//	     fmt.Println(st)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             st = sini + st + sr
             err =  book.BookUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Libro actualizado para: " +book.Title, view.FlashSuccess})
             } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
	}
		sess.Save(r, w)
         }
      }
	http.Redirect(w, r, path, http.StatusFound)
     }
//------------------------------------------------
// BookDeleteGET handles the note deletion
 func BookDeleteGET(w http.ResponseWriter, r *http.Request) {
        var book model.Book
            var lang model.Language
            var editor model.Editor
            var author model.Author
        var params httprouter.Params
        sess := model.Instance(r)
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        path :=  fmt.Sprintf("/biblos/list/%s", SPag)
        book.Id = Id
        err := (&book).BookById()
        if err != nil {
                log.Println(err)
                sess.AddFlash(view.Flash{"Error Libro no hallado.", view.FlashError})
                http.Redirect(w, r, path, http.StatusFound)
        }else{
            lang.Id = book.LanguageId
            err =  (&lang).LangById()
            if err != nil {
                sess.AddFlash(view.Flash{"No hay ese lenguaje ", view.FlashError})
                lang.Name = "xxxxxx"
            }
           editor.Id = book.EditorId
           err = (&editor).EditById()
            if err != nil {
                sess.AddFlash(view.Flash{"No hay ese editor ", view.FlashError})
                editor.Name = "xxxxxx"
            }
           author.Id = book.AuthorId
           err = (&author).AuthorById()
            if err != nil {
                sess.AddFlash(view.Flash{"No hay ese author ", view.FlashError})
                author.Name = "xxxxxx"
            }
       }
           v    := view.New(r)
           v.Name = "book/bookdelete"
           v.Vars["langName"]    = lang.Name
           v.Vars["editName"]    = editor.Name
           v.Vars["authName"]    = author.Name
           v.Vars["title"]       = book.Title
           v.Vars["comment"]     = book.Comment
           v.Vars["isbn"]        = book.Isbn
           v.Vars["year"]        = fmt.Sprintf("%d",book.Year)
	   v.Vars["token"]  = csrfbanana.Token(w, r, sess)
           v.Render(w)
  }

//------------------------------------------------
// BookDeletePOST handles the note deletion
 func BookDeletePOST(w http.ResponseWriter, r *http.Request) {
        var book model.Book
        var params httprouter.Params
        sess        := model.Instance(r)
        params       = context.Get(r, "params").(httprouter.Params)
        Id,_        := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        action      := r.FormValue("action")
        path        :=  fmt.Sprintf("/biblos/list/%s", SPag)
        if ! (strings.Compare(action,"Cancelar") == 0) {
           book.Id = Id
           err := (&book).BookById()
           if err != nil {
              log.Println(err)
              sess.AddFlash(view.Flash{"Error Libro no hallado.", view.FlashError})
            }else{
               err = book.Delete()
               if err != nil {
                   log.Println(err)
	           sess.AddFlash(view.Flash{"Error no posible. No borramos.", view.FlashError})
               } else {
                      sess.AddFlash(view.Flash{"Libro borrado!", view.FlashSuccess})
               }
            }
        }
        sess.Save(r, w)
	http.Redirect(w, r, path, http.StatusFound)
  }
//------------------------------------------------
// RegisSearchPOST procesa la forma enviada con los datos
func BookSearchPOST(w http.ResponseWriter, r *http.Request) {
	sess       :=  model.Instance(r)
        path       :=  fmt.Sprintf("/biblos/list/%s", "1")
        rSearch    :=  r.FormValue("bsearch")
        if rSearch == ""{
           sess.Save(r, w)
	   http.Redirect(w, r, path, http.StatusFound)
           return
         }

        lisBooks, err := model.SBooks(rSearch)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Libros.", view.FlashError})
            sess.Save(r, w)
	    http.Redirect(w, r, path, http.StatusFound)
         }
        posact     = 0
        offset     =  posact  - 1
        TotalCount :=  len(lisBooks)
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay libros.", view.FlashError})
             sess.Save(r, w)
	     http.Redirect(w, r, path, http.StatusFound)
             return
         }
        cut_Names(lisBooks )
        offset = offset * limit
        numberOfBtns      :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns         :=  createSliceForBtns(numberOfBtns, posact)

	// Display the view
	v := view.New(r)
	v.Name             = "book/booklis"
	v.Vars["token"]    = csrfbanana.Token(w, r, sess)
        v.Vars["Search"]     =  rSearch
        v.Vars["Level"]    =  sess.Values["level"]
        v.Vars["LisBook"]  = lisBooks
        v.Vars["slice"]    =  sliceBtns
        v.Vars["current"]  =  posact
	// Refill any form fields
	v.Render(w)
     }


//------------------------------------------------
// BookLisGET displays the book page
func BookLisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess       := model.Instance(r)
        var params httprouter.Params
        params     = context.Get(r, "params").(httprouter.Params)
        SId        := params.ByName("id")
        Id,_       := atoi32(SId)
        posact     = int(Id)
        offset     =  posact  - 1
        TotalCount := model.BookCount()
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay libros.", view.FlashError})
             sess.Save(r, w)
             return
         }
        offset = offset * limit
        lisBooks, err := model.BookLim(limit, offset)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Books.", view.FlashError})
            sess.Save(r, w)
         }
        cut_Names(lisBooks )
        numberOfBtns      :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns         :=  createSliceForBtns(numberOfBtns, posact)
	v                 := view.New(r)
	v.Name             = "book/booklis"
	v.Vars["token"]    = csrfbanana.Token(w, r, sess)
        v.Vars["slice"]    =  sliceBtns
        v.Vars["current"]  =  posact
        v.Vars["Level"]    =  sess.Values["level"]
        v.Vars["LisBook"]  = lisBooks
	v.Render(w)
 }
//------------------------------------------------
// BookLangGET displays the book page given language

func BookLangGET(w http.ResponseWriter, r *http.Request) {
	sess       := model.Instance(r)
        var params httprouter.Params
        params     = context.Get(r, "params").(httprouter.Params)
        SPg        := params.ByName("pg")
        SId        := params.ByName("id")
        SAt        := params.ByName("at")
        IPg,_      := atoi32(SPg)
        Id,_       := atoi32(SId)
        posact     = int(IPg)
        offset     =  posact  - 1
        TotalCount := model.BookLangCount(Id)
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay libros.", view.FlashError})
             sess.Save(r, w)
             path :=  fmt.Sprintf("/language/list/%d", 1)
             http.Redirect(w, r, path, http.StatusFound)
          }else{
             offset = offset * limit
             lisBooks, err := model.BookLangLim(Id, limit, offset)
        if err != nil {
             log.Println(err)
	     sess.AddFlash(view.Flash{"Error Listando Books.", view.FlashError})
             sess.Save(r, w)
         }
        cut_Names(lisBooks )
	v                := view.New(r)
	v.Name            = "book/booklslang"
	v.Vars["token"]   = csrfbanana.Token(w, r, sess)
        numberOfBtns     :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns        :=  createSliceForBtns(numberOfBtns, posact)
        v.Vars["slice"]   =  sliceBtns
        v.Vars["SId"]     =  SId
        v.Vars["SAt"]     = SAt
        v.Vars["current"] =  posact
        v.Vars["Level"]   =  sess.Values["level"]
        v.Vars["LisBook"] =  lisBooks
        v.Vars["LName"]   =  lisBooks[0].Language
	v.Render(w)
      }
 }
//------------------------------------------------
// BookAuthGET displays the book page given Author
func BookAuthGET(w http.ResponseWriter, r *http.Request) {
	sess       := model.Instance(r)
        var params httprouter.Params
        params     = context.Get(r, "params").(httprouter.Params)
        SId        := params.ByName("id")
        SPg        := params.ByName("pg")
        SAt        := params.ByName("at")
        IPg,_      := atoi32(SPg)
        Id,_       := atoi32(SId)
        posact     = int(IPg)
        offset     =  posact  - 1
        TotalCount := model.BookAuthCount(Id)
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay libros.", view.FlashError})
             sess.Save(r, w)
             path := fmt.Sprintf("/editor/list/%d",1 )
           http.Redirect(w, r, path, http.StatusFound)
          }else{
        offset = offset * limit
        lisBooks, err := model.BookAuthLim(Id, limit, offset)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Books.", view.FlashError})
            sess.Save(r, w)
         }
        cut_Names(lisBooks )
	v                := view.New(r)
	v.Name            = "book/booklsauth"
	v.Vars["token"]   = csrfbanana.Token(w, r, sess)
        numberOfBtns     :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns        :=  createSliceForBtns(numberOfBtns, posact)
        v.Vars["slice"]   =  sliceBtns
        v.Vars["SId"]     =  SId
        v.Vars["Level"]   =  sess.Values["level"]
        v.Vars["SPg"]     = SPg
        v.Vars["SAt"]     = SAt
        v.Vars["current"] =  posact
        v.Vars["AName"]   = lisBooks[0].Author
        v.Vars["LisBook"] = lisBooks
	v.Render(w)
      }
 }
//------------------------------------------------

//------------------------------------------------
// BookEditGET displays the book page given Editor

func BookEditGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess       := model.Instance(r)
        var params httprouter.Params
        params     = context.Get(r, "params").(httprouter.Params)
        SPg        := params.ByName("pg")
        SId        := params.ByName("id")
        SAt        := params.ByName("at")
        IPg,_      := atoi32(SPg)
        Id,_       := atoi32(SId)
        posact     := int(IPg)
        offset     :=  posact  - 1
        offset     = offset * limit
        TotalCount := model.BookEditCount(Id)
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay libros.", view.FlashError})
             sess.Save(r, w)
             path   := fmt.Sprintf("/editor/list/%d", 1)
             http.Redirect(w, r, path, http.StatusFound)
          }else{
        lisBooks, err := model.BookEditLim(Id, limit, offset)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Books.", view.FlashError})
            sess.Save(r, w)
         }
         cut_Names(lisBooks )
	v                := view.New(r)
	v.Name            = "book/booklsedit"
	v.Vars["token"]   = csrfbanana.Token(w, r, sess)
        numberOfBtns     :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns        :=  createSliceForBtns(numberOfBtns, posact)
        v.Vars["slice"]   =  sliceBtns
        v.Vars["SId"]     =  SId
        v.Vars["Level"]   =  sess.Values["level"]
        v.Vars["SPg"]     = SPg
        v.Vars["SAt"]     = SAt
        v.Vars["current"] =  posact
        v.Vars["EName"]   = lisBooks[0].Editor
        v.Vars["LisBook"] = lisBooks
	v.Render(w)
      }
 }

//------------------------------------------------
// BookPdfGET generates pdf from model.BooksN
 func BookPdfGET(w http.ResponseWriter, r *http.Request) {

      lsBooks, _  := model.BooksN()
      pdf := getBookPdf(lsBooks)
      pdf.Output(w )
  }
//------------------------------------------------
//------------------------------------------------
// BookAuthPdfGET generates pdf from model.BooksN
 func BookAuthPdfGET(w http.ResponseWriter, r *http.Request) {
      var params httprouter.Params
      params     = context.Get(r, "params").(httprouter.Params)
      SId        := params.ByName("id")
      Id,_       := atoi32(SId)

      lsBooks, _  := model.BookAuthTot(Id )
      pdf := getBookAuthPdf(lsBooks )
      pdf.Output(w )

  }
//------------------------------------------------
//------------------------------------------------
// BookEditPdfGET generates pdf from model.BooksN
 func BookEditPdfGET(w http.ResponseWriter, r *http.Request) {
      var params httprouter.Params
      params     = context.Get(r, "params").(httprouter.Params)
      SId        := params.ByName("id")
      Id,_       := atoi32(SId)

      lsBooks, _  := model.BookEditTot(Id )
      pdf    := getBookEditPdf(lsBooks )
      pdf.Output(w )

  }
//------------------------------------------------
//------------------------------------------------
// BookLangPdfGET generates pdf from model.BooksN
 func BookLangPdfGET(w http.ResponseWriter, r *http.Request) {
      var params httprouter.Params
      params     = context.Get(r, "params").(httprouter.Params)
      SId        := params.ByName("id")
      Id,_       := atoi32(SId)

      lsBooks, _  := model.BookLangTot(Id )
      pdf    := getBookLangPdf(lsBooks )
      pdf.Output(w )

  }
//------------------------------------------------
