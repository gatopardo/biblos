package controller

import (
	"log"
	"net/http"
        "strings"
        "fmt"
//     gr "github.com/mikeshimura/goreport"
//	"strconv"   
	"encoding/json"

	"github.com/gatopardo/biblos/app/model"
	"github.com/gatopardo/biblos/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
  // ---------------------------------------------------
// JLangGET crea datos de libros
 func JLangListGET(w http.ResponseWriter, r *http.Request) {
        var err error
	var lisBooks []model.BookZ
	var params httprouter.Params
	sess := model.Instance(r)
	params      = context.Get(r, "params").(httprouter.Params)
	reauth      := params.ByName("re")
	lisBooks,err    = model.LangByRe(reauth)
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

// LangGET despliega la pagina del usuario
func LangGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "language/language"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
//      Refill any form fields
// view.Repopulate([]string{"name"}, r.Form, v.Vars)
	v.Render(w)
 }
// ---------------------------------------------------
// ---------------------------------------------------
// LangPOST procesa la forma enviada con los datos
func LangPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
        var lang model.Language
	sess := model.Instance(r)
        action        := r.FormValue("action")
//  fmt.Println(action)
        if ! (strings.Compare(action,"Cancelar") == 0) {

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"name"}); !validate {
          sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
          sess.Save(r, w)
          LangGET(w, r)
          return
	}
	// Get form values
	lang.Name      = r.FormValue("name")

	 err := (&lang).LangByName()
         if err == model.ErrNoResult { // Exito:  no hay usuario creado aun 
            ex := (&lang).LangCreate()
	    if ex != nil {  // uyy como fue esto ? 
              log.Println(ex)
              sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
              sess.Save(r, w)
                   return
	   } else {  // todo bien
           sess.AddFlash(view.Flash{"Lenguaje creado: " +lang.Name, view.FlashSuccess})
              sess.Save(r, w)
	   }
        }
      }
// log.Println("LangPOST b")
	http.Redirect(w, r, "/language/list/1", http.StatusFound)
 }

// ---------------------------------------------------
// LangUpGET despliega la pagina del usuario
func LangUpGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        var lang model.Language
	// necesitamos language id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
         lang.Id = id
	// Obtener usuario dado id
	 err := (&lang).LangById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No tenemos lenguaje.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, "/language/list/1", http.StatusFound)
           return
	}
	// Display the view
	v := view.New(r)
	v.Name = "language/langupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Vars["name"] = lang.Name
//    Refill any form fields
//	view.Repopulate([]string{"name"}, r.Form, v.Vars)
        v.Render(w)
   }
// ---------------------------------------------------

// LangUpPOST procesa la forma enviada con los datos
func LangUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
         var lang model.Language
	// Get session
	sess := model.Instance(r)
//        Validate with required fields
        if validate, missingField := view.Validate(r, []string{"name"}); !validate {
           sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
            sess.Save(r, w)
            LangUpGET(w, r)
            return
         }
	// Get form values
	lang.Name = r.FormValue("name")

        var params httprouter.Params
	params       = context.Get(r, "params").(httprouter.Params)
	lang.Id, _   = atoi32(params.ByName("id"))
        SPag        := params.ByName("pagi")
        path        :=  fmt.Sprintf("/language/list/%s", SPag)
       action      := params.ByName("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {

 // Vamos a actualizar 
             err =  lang.Update()
             if err == nil{
                 sess.AddFlash(view.Flash{"Lenguaje actualizado exitosamente para: " +lang.Name, view.FlashSuccess})
		sess.Save(r, w)
             } else {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
		sess.Save(r, w)
            }
	}
	// Display the page
	http.Redirect(w, r, path, http.StatusFound)

//	LangUpGET(w, r)
     }
//------------------------------------------------
// LangLisGET displays the language page

func LangLisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        SId := params.ByName("id")
        Id,_ := atoi32(SId)

        posact = int(Id)
        offset = posact  - 1
        TotalCount := model.LangCount()
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay lenguajes.", view.FlashError})
             sess.Save(r, w)
             return
          }else{
        offset = offset * limit

        lisLangs, err := model.LangLim(limit, offset)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Lenguajes.", view.FlashError})
            sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "language/langlis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
       langs := make([]SLang, len(lisLangs))
        lev, ok := v.Vars["level"].(uint32)
        if ok {
           for i, lang := range lisLangs  {
              langs[i].Lang  =  lang
//              langs[i].Level = lev
//              langs[i].Pagi = posact
           }

           v.Vars["LisLang"] = langs
        }

        numberOfBtns     :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns        :=  createSliceForBtns(numberOfBtns, posact)
        v.Vars["slice"]   =  sliceBtns
        v.Vars["current"] =  posact
        v.Vars["Level"]    = lev

	v.Render(w)
      }
 }
//------------------------------------------------
//  LangReportGET handles report
/*
 func LangReporGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        TotalCount := model.LangCount()
        if TotalCount == 0 {
	     sess.AddFlash(view.Flash{"No hay lenguajes.", view.FlashError})
             sess.Save(r, w)
             return 
          }
        lisLangs, err := model.Langs() 
        if err != nil {
  	    log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Lenguajes.", view.FlashError})
            sess.Save(r, w)
         }
        r := gr.CreateGoReport()
	font1 := gr.FontMap{
		FontName: "IPAex",
		FileName: "ttf//ipaexg.ttf",
	}
	fonts := []*gr.FontMap{&font1}
	r.SetFonts(fonts) 
        d := new(S1Detail)
	r.RegisterBand(gr.Band(*d), gr.Detail)
	h := new(S1Header)
	r.RegisterBand(gr.Band(*h), gr.PageHeader)

        r.Records = gr.ReadTextFile("sales1.txt", 7)
	//fmt.Printf("Records %v \n", r.Records)
	r.SetPage("A4", "mm", "L")
	r.Execute("simple1.pdf")

  }
type S1Detail struct {
}

func (h S1Detail) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h S1Detail) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	report.Font("IPAexG", 12, "")
	y := 2.0
	report.Cell(15, y, cols[0])
}

type S1Header struct {
}

func (h S1Header) GetHeight(report gr.GoReport) float64 {
	return 30
}
func (h S1Header) Execute(report gr.GoReport) {
	report.Font("IPAexG", 14, "")
	report.Cell(50, 15, "Lenguajes")
}
*/
//------------------------------------------------

// LangDeleteGET handles the note deletion
 func LangDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
        sess := model.Instance(r)
        v    := view.New(r)
        v.Name ="language/languagedelete"
        var lang model.Language
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        path :=  fmt.Sprintf("/language/list/%s", SPag)
//    fmt.Println(Id,SPag)
        lang.Id = Id
        err := (&lang).LangById()
        if err != nil{
            log.Println(err)
            sess.AddFlash(view.Flash{"Error Language no hallado.", view.FlashError})
            http.Redirect(w, r, path, http.StatusFound)
        }else{
            canti := model.BookLangCount(Id)
            v.Vars["name"]  = lang.Name
            v.Vars["canti"]  = fmt.Sprintf("%d",canti)
        }
// fmt.Println(path)
        v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Render(w)
  }

//------------------------------------------------

// LangDeletePOST handles the note deletion
 func LangDeletePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
        sess := model.Instance(r)
        var lang model.Language
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        action      := params.ByName("action")
         path :=  fmt.Sprintf("/language/list/%s", SPag)
//    fmt.Println(Id,SPag)
        if ! (strings.Compare(action,"Cancelar") == 0) {
            lang.Id = Id
// Get database result
//    	    err := lang.Delete()
// 	    if err != nil {
// 		log.Println(err)
// //		sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
// 		sess.Save(r, w)
// 	     } else {
// 		sess.AddFlash(view.Flash{"Author borrado!", view.FlashSuccess})
// 		sess.Save(r, w)
// 	     }
          }
           sess.Save(r, w)
//       fmt.Println(path)
           http.Redirect(w, r, path, http.StatusFound)
        }

//------------------------------------------------



