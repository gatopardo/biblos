package controller

import (
        "github.com/gatopardo/biblos/app/model"
        "fmt"
	"log"
	"time"
        "strconv"
        "strings"
	"github.com/jung-kurt/gofpdf"
  )

// ---------------------------------------------------
      const(
              limit       = 18
              margenlat   = 3
            )

// ---------------------------------------------------

      var (
            TotalCount  int
            offset      int
            posact      int
           )

// ---------------------------------------------------
     type SUser struct {
          User model.User
          Level  uint32
          Pagi   int
       }


       type SLang struct {
          Lang model.Language
          Level  uint32
          Pagi   int
       }


       type SAuthor struct {
          Author model.Author
          Level  uint32
          Pagi   int
       }


       type SEdit struct {
          Edit model.Editor
          Level  uint32
          Pagi   int
       }

       type SBook struct {
          Book model.Book
          Level  uint32
       }

       type TBook struct {
          Book model.BookN
          Level  uint32
          Pagi   int
       }


//------------------------------------------------
func substr(s string,pos,length int) string{
    runes:=[]rune(s)
    l := pos+length
    if l > len(runes) {
        l = len(runes)
    }
    return string(runes[pos:l])
}

//------------------------------------------------
 func cut_Names( books []model.BookN )  {
       tam := len(books)
       for i := 0; i < tam; i++ {
          books[i].Language = substr(books[i].Language, 0, 15)
          books[i].Editor   = substr(books[i].Editor,0,15)
          books[i].Author   = substr(books[i].Author,0,15)
          books[i].Title    = substr(books[i].Title,0,38)
          books[i].Isbn     = substr(books[i].Isbn,0,10)
          books[i].Comment  = substr(books[i].Comment,0,15)
        }
    }

// ---------------------------------------------------
  func atoi32( str string) (nr uint32,err error){
        i, errn := strconv.Atoi(str)
        nr  = uint32(i)
        err =  errn
        return
    }
// ---------------------------------------------------

func getNumberOfButtonsForPagination(TotalCount int, limit int) int {
    num := (int)(TotalCount / limit)
    if (TotalCount%limit > 0) {
        num++
    }
    return num
}

func createSliceForBtns(number int, posact int) []int {
    var sliceOfBtn []int
    lffin := margenlat
    rtini := number   -  margenlat  + 1
    inilf := posact   -  margenlat
    finrt := posact   +  margenlat
    if inilf < 1 {
       inilf = 1
      }
    if finrt > number  {
       finrt =  number
      }
    if lffin  > inilf  {
       lffin  = inilf - 1
    }
    if rtini  < finrt  {
        rtini = finrt  + 1
    }
    for i := 1; i <= lffin; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    for i := inilf; i <= finrt; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    for i := rtini; i <= number; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    return sliceOfBtn
 }

  func roundU(val float64) int {
      if val > 0 { return int(val+1.0) }
      return int(val)
  }

  func  leftStr(pdf * gofpdf.Fpdf, st string, w float64) string {
             lt  := pdf.GetStringWidth(st)
             dif :=  lt  - w
             if dif > 0 {
             for dif > -2 {
                st   = st[:len(st) - 4]
                lt   = pdf.GetStringWidth(st)
                dif  =  lt  - w
            }
                runes :=  []rune(st)
                st    =  string(runes)
            }
                return st
   }

//------------------------------------------------

  func getFields( book model.BookN, pdf * gofpdf.Fpdf, lw []float64) ( lin []string ) {
       str :=  strings.Trim(book.Language, " ")
       str  = leftStr(pdf, str, lw[0])
       lin = append(lin, str)
       str =  strings.Trim(book.Author, " ")
       str  = leftStr(pdf, str, lw[1])
       lin = append(lin, str)
       str =  strings.Trim(book.Title, " ")
       str  = leftStr(pdf, str, lw[2])
       lin = append(lin, str)
       stYear := fmt.Sprint(book.Year)
       lin = append(lin, stYear)
       str =  strings.Trim(book.Editor, " ")
       str  = leftStr(pdf, str, lw[4])
       lin = append(lin, str)
       return lin
  }

 func addData(pdf *gofpdf.Fpdf, lh string, lhead[]string,  lw []float64, books []model.BookN) *gofpdf.Fpdf {
        pdf.SetTitle(lh,true)
	align := []string{"L", "L", "L", "R",  "L"}
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        fill := false
	// 	Data
        nr := 26
        nrlin := 28
        ht   := 6.0
                     tr := pdf.UnicodeTranslatorFromDescriptor("")
	for i, c := range books {
                pdf.SetFont("Times", "", 16)
                pdf.SetTextColor(0, 0, 0)
                pdf.SetDrawColor(128, 0, 0)
		if i % 2 == 0 {
                  pdf.SetFillColor(255, 255, 255)
                }else{
                  pdf.SetFillColor(197, 236, 235)
		}
		fields := getFields(c, pdf, lw)
		for j, field := range fields {
//                   str  := leftStr(pdf, field, lw[j])
                     pdf.CellFormat(lw[j], ht, tr(field), "1", 0, align[j], fill, 0, "")
                }
                pdf.Ln(-1)
                fill = ! fill
		if (i == nr) || ( (i > nr  ) && ( (i - nr)  % nrlin == 0)) {
                  pdf.CellFormat(80.0, 0, "", "TR", 0, "C", false, 0, "")
		  addFooterBook(pdf)
                    pdf.AddPage()
		}
                if (i - nr) % nrlin == 0{
                   header(pdf, lhead, lw)
               }
	  }
          pdf.CellFormat(255.0, 0, "", "TR", 0, "", false, 0, "")
   return pdf
 }

func header(pdf *gofpdf.Fpdf, hdr []string, whdr []float64 ) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 16)
	pdf.SetFillColor(120, 140, 120)
        pdf.SetDrawColor(239, 234, 228)
	for i, str := range hdr {
		pdf.CellFormat(whdr[i],      5, str, "1", 0, "L", true, 0, "")
	}
	pdf.Ln(-1)
	return pdf
}

 func addFooterBook(pdf * gofpdf.Fpdf) *gofpdf.Fpdf  {
         pdf.SetX(-15)
         pdf.SetFont("Arial", "I", 8)
	 stPage := fmt.Sprintf("Pag %d/{nb}", pdf.PageNo())
         pdf.CellFormat(0, 10, stPage, "", 0, "C", false, 0, "")
         pdf.Image("static/favicons/ipi_7a.png", 100, 236, 25, 0, false, "", 0, "")
         pdf.Line(10,275,200,275)
	return pdf
  }

 func newReport(stit string) *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetFillColor(240, 80, 80)
        pdf.SetDrawColor(128, 0, 0)
	pdf.SetFont("Times", "B", 20)
	pdf.CellFormat(80, 10, stit,"",0,"L",false, 0, "")
	pdf.SetFont("Times", "B", 9)
	 sttime := time.Now().Format("2006/01/02")
	pdf.CellFormat(160, 10, sttime,"",0,"R", false, 0, "")
	pdf.Ln(8)
	return pdf
}

 func getBookPdf(books []model.BookN ) (pdf * gofpdf.Fpdf)  {
        SName  := "Libros "
        stit   := "Listado de " + SName
	pdf    = newReport(stit)
        if pdf.Err() { log.Fatalf("Fallo 1 creando PDF : %s\n", pdf.Error()) }

        stheader := []string  {"Lang", "Author", "Title", "Year", "Editor"}
	szwidth := []float64 { 28.0,    60.0,     94.0,    14.0,     60.0    }
        pdf.SetHeaderFunc(func(){     header(pdf, stheader, szwidth) } )
        if pdf.Err() { log.Fatalf("Fallo 2 header PDF : %s\n", pdf.Error()) }

	pdf.SetFooterFunc(func() {   addFooterBook(pdf )} )
        if pdf.Err() { log.Fatalf("Fallo 4 Footer PDF : %s\n", pdf.Error()) } 

        pdf = addData(pdf,stit,stheader, szwidth,  books )
        if pdf.Err() { log.Fatalf("Fallo 3 data PDF : %s\n", pdf.Error()) }

         return pdf
    }

    func getAuthor( pdf *gofpdf.Fpdf,  lw []float64, authors []model.Author) *gofpdf.Fpdf {

	pdf.SetFillColor(224, 235, 255)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.AddPage()

        w := []float64{175 }
	wSum := 175.0
        fill := false
        ht   := 6.0
        nfields := len(w)
        tr := pdf.UnicodeTranslatorFromDescriptor("")
        write := func(str string,w float64) {
                   str  = leftStr(pdf, str, w)
                   pdf.CellFormat(w, ht, tr(str), "LR", 0, "", fill, 0, "")
                }
	// 	Data
	for i, c := range authors {

                fields := []string{ c.Name }
                for j := 0; j < nfields; j++{
                     write(fields[j], w[j])
                }
                pdf.Ln(-1)
                fill = ! fill
                if (i + 1) % 35 == 0{
                   pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")
                    pdf.AddPage()
               }
          }
          pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")
        return pdf
   }

 func getAuthorPdf(authors []model.Author ) (pdf * gofpdf.Fpdf)  {

        pdf = gofpdf.New("P", "mm", "Letter", "")
        SName  := "Autores"
        pdf.SetHeaderFunc(func() {
           pdf.SetFillColor(255, 0, 0)
           pdf.SetTextColor(255, 255, 255)
           pdf.SetDrawColor(128, 0, 0)

           stit := "Listado de " + SName
           pdf.CellFormat(195, 7, stit, "1", 0, "C", true, 0, "")
           pdf.Ln(-1)

           header := []string{ "Author" }

            w := []float64{25 }
	    pdf.SetFillColor(224, 235, 255)
            pdf.SetTextColor(0, 0, 0)
             for j, str := range header {
                pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
              }
             pdf.Ln(-1)
       })

        if pdf.Err() { log.Fatalf("Fallo 1 header  : %s\n", pdf.Error())
        fmt.Println("Fallo 1 header  ", pdf.Error()) 
      }

	pdf.SetFooterFunc(func() {
              pdf.SetX(-15)
              pdf.SetFont("Arial", "I", 8)
              pdf.CellFormat(0, 10, fmt.Sprintf("Pag %d/{nb}", pdf.PageNo()),
                  "", 0, "C", false, 0, "")
//              pdf.Image("static/favicons/ipi_7a.png", 100, 236, 25, 0, false, "", 0, "")
              pdf.Line(10,275,200,275)
    })

        if pdf.Err() { log.Fatalf("Fallo 2 footer  : %s\n", pdf.Error())
            fmt.Println("Fallo 2 footer  ", pdf.Error()) 
          }

        pdf.SetTitle(SName,true)
        lw := []float64{175 }
        pdf = getAuthor( pdf,  lw, authors )
        if pdf.Err() { log.Fatalf("Fallo 3 data  : %s\n", pdf.Error())
           fmt.Println("Fallo 3 data  ", pdf.Error()) 
       }

         return pdf
	}


    func getBookAuthor( pdf *gofpdf.Fpdf, lw []float64, books []model.BookN) *gofpdf.Fpdf {

	pdf.SetFillColor(224, 235, 255)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.AddPage()

	wSum := 0.0
	for _, v := range lw {
		wSum += v
	}
        fill := false
        ht   := 6.0
        nfields := 2
        tr := pdf.UnicodeTranslatorFromDescriptor("")
        write := func(str string,w float64) {
                   str  = leftStr(pdf, str, w)
                   pdf.CellFormat(w, ht, tr(str), "LR", 0, "", fill, 0, "")
                }
	// 	Data
	for i, c := range books {
                fields := []string{ c.Title, c.Comment }
                for j := 0; j < nfields; j++{
                     write(fields[j], lw[j])
                }
                pdf.Ln(-1)
                fill = ! fill
                if (i + 1) % 16 == 0{
                   pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")
                    pdf.AddPage()
               }
	  }
                   pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")
              return pdf
    }


 func getBookAuthPdf(books []model.BookN ) (pdf * gofpdf.Fpdf )  {
	pdf = gofpdf.New("P", "mm", "A6", "")
        SName := books[0].Author
        pdf.SetHeaderFunc(func() {
           pdf.SetFillColor(255, 0, 0)
           pdf.SetTextColor(255, 255, 255)
           pdf.SetDrawColor(128, 0, 0)

           stit := strings.Trim(" Libros de " + SName, " ")
           pdf.CellFormat(90, 7, stit, "1", 0, "C", true, 0, "")
	   pdf.Ln(-1)

	   header := []string{ "Title", "Tema" }
	   w := []float64{ 65,25 }
	   wSum := 0.0
	   for _, v := range w {
		  wSum += v
	   }
	   pdf.SetFillColor(224, 255, 235)
           pdf.SetTextColor(20, 20, 20)
           for j, str := range header {
                pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
           }
           pdf.Ln(-1)
    })

        if pdf.Err() { log.Fatalf("Fallo 1 header  : %s\n", pdf.Error())
        fmt.Println("Fallo 1 header  ", pdf.Error()) 
      }

        pdf.SetFooterFunc(func() {
//            pdf.SetY(-15)
              pdf.SetX(-15)
              pdf.SetFont("Arial", "I", 8)
              pdf.CellFormat(0, 10, fmt.Sprintf("Pag %d/{nb}", pdf.PageNo()),
                  "", 0, "C", false, 0, "")
//              pdf.Image("static/favicons/ipi_7a.png", 100, 236, 25, 0, false, "", 0, "")
              pdf.Line(10,275,200,275)
              pdf.Line(5,140,98,140)
    })

        if pdf.Err() { log.Fatalf("Fallo 2 footer  : %s\n", pdf.Error())
        fmt.Println("Fallo 2 footer  ", pdf.Error()) 
      }

        pdf.SetTitle(SName,true)
	lw := []float64{ 65,25 }
        pdf = getBookAuthor( pdf,  lw, books )
        if pdf.Err() { log.Fatalf("Fallo 3 data  : %s\n", pdf.Error())
           fmt.Println("Fallo 3 data  ", pdf.Error()) 
       }


         return pdf
	}

    func getBookEdit( pdf *gofpdf.Fpdf, lw []float64, books []model.BookN) *gofpdf.Fpdf {

	pdf.SetFillColor(224, 235, 255)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.AddPage()

	wSum := 0.0
	for _, v := range lw {
		wSum += v
	}
        fill := false

        ht   := 6.0
        nfields := len(lw)
        tr := pdf.UnicodeTranslatorFromDescriptor("")
        write := func(str string,w float64) {
                   str  = leftStr(pdf, str, w)
                   pdf.CellFormat(w, ht, tr(str), "LR", 0, "", fill, 0, "")
                }
	// 	Data
	for i, c := range books {
                fields := []string{ c.Title, c.Comment, c.Author }
                for j := 0; j < nfields; j++{
                     write(fields[j], lw[j])
                }
                pdf.Ln(-1)
                fill = ! fill
                if (i + 1) % 27 == 0{
                   pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")
                    pdf.AddPage()
               }
         }
         pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")
         return pdf
    }

 func getBookEditPdf(books []model.BookN ) (pdf * gofpdf.Fpdf)  {

        pdf = gofpdf.New("P", "mm", "A5", "")
        SName := books[0].Editor
        pdf.SetHeaderFunc(func() {
           pdf.SetFillColor(255, 0, 0)
           pdf.SetTextColor(255, 255, 255)
           pdf.SetDrawColor(128, 0, 0)

           stit := strings.Trim("Editor " + SName, " ")
           pdf.CellFormat(135, 7, stit, "1", 0, "C", true, 0, "")
	   pdf.Ln(-1)

	   header := []string{ "Title", "Tema", "Author" }
	   w := []float64{ 65,25,45 }
	   wSum := 0.0
	   for _, v := range w {
		  wSum += v
	   }
	   pdf.SetFillColor(224, 255, 235)
           pdf.SetTextColor(20, 20, 20)
           for j, str := range header {
                pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
           }
           pdf.Ln(-1)
    })

        if pdf.Err() { log.Fatalf("Fallo 1 header  : %s\n", pdf.Error())
        fmt.Println("Fallo 1 header  ", pdf.Error()) 
      }

        pdf.SetFooterFunc(func() {
//            pdf.SetY(-15)
              pdf.SetX(-15)
              pdf.SetFont("Arial", "I", 8)
              pdf.CellFormat(0, 10, fmt.Sprintf("Pag %d/{nb}", pdf.PageNo()),
                  "", 0, "C", false, 0, "")
              pdf.Image("static/favicons/ipi_7a.png", 100, 236, 25, 0, false, "", 0, "")
              pdf.Line(10,275,200,275)
              pdf.Line(5,200,140,200)
    })

        if pdf.Err() { log.Fatalf("Fallo 2 footer  : %s\n", pdf.Error())
        fmt.Println("Fallo 2 footer  ", pdf.Error()) 
      }

        pdf.SetTitle(SName,true)
	lw := []float64{ 65,25,45 }
        pdf = getBookEdit( pdf,  lw, books )
        if pdf.Err() { log.Fatalf("Fallo 3 data  : %s\n", pdf.Error())
           fmt.Println("Fallo 3 data  ", pdf.Error()) 
       }

         return pdf
	}

    func getBookLang( pdf *gofpdf.Fpdf, lw []float64, books []model.BookN) *gofpdf.Fpdf {
	pdf.SetFillColor(224, 235, 255)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.AddPage()

	wSum := 0.0
	for _, v := range lw {
		wSum += v
	}
        fill := false

        ht   := 6.0
        nfields := len(lw)
        tr := pdf.UnicodeTranslatorFromDescriptor("")
        write := func(str string,w float64) {
                   str  = leftStr(pdf, str, w)
                   pdf.CellFormat(w, ht, tr(str), "LR", 0, "", fill, 0, "")
                }
	// 	Data
	for i, c := range books {
                fields := []string{ c.Title, c.Comment, c.Author }
                for j := 0; j < nfields; j++{
                     write(fields[j], lw[j])
                }
                pdf.Ln(-1)
                fill = ! fill
                if (i + 1) % 27 == 0{
                   pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")
                    pdf.AddPage()
               }
         }
         pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")

         return pdf
	}

 func getBookLangPdf(books []model.BookN ) (pdf * gofpdf.Fpdf)  {

	pdf = gofpdf.New("P", "mm", "A5", "")
        SName := books[0].Language
        pdf.SetHeaderFunc(func() {
           pdf.SetFillColor(255, 0, 0)
           pdf.SetTextColor(255, 255, 255)
           pdf.SetDrawColor(128, 0, 0)

           stit := strings.Trim("Language " + SName, " ")
           pdf.CellFormat(135, 7, stit, "1", 0, "C", true, 0, "")
	   pdf.Ln(-1)

	   header := []string{ "Title", "Tema", "Author" }
	   w := []float64{ 65,25,45 }
	   wSum := 0.0
	   for _, v := range w {
		  wSum += v
	   }
	   pdf.SetFillColor(224, 255, 235)
           pdf.SetTextColor(20, 20, 20)
           for j, str := range header {
                pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
           }
           pdf.Ln(-1)
    })

        if pdf.Err() { log.Fatalf("Fallo 1 header  : %s\n", pdf.Error())
        fmt.Println("Fallo 1 header  ", pdf.Error()) 
      }

        pdf.SetFooterFunc(func() {
//            pdf.SetY(-15)
              pdf.SetX(-15)
              pdf.SetFont("Arial", "I", 8)
              pdf.CellFormat(0, 10, fmt.Sprintf("Pag %d/{nb}", pdf.PageNo()),
                  "", 0, "C", false, 0, "")
//              pdf.Image("static/favicons/ipi_7a.png", 100, 236, 25, 0, false, "", 0, "")
              pdf.Line(10,275,200,275)
              pdf.Line(5,200,140,200)
    })

        if pdf.Err() { log.Fatalf("Fallo 2 footer  : %s\n", pdf.Error())
        fmt.Println("Fallo 2 footer  ", pdf.Error()) 
      }

        pdf.SetTitle(SName,true)
	lw := []float64{ 65,25,45 }
        pdf = getBookLang( pdf,  lw, books )
        if pdf.Err() { log.Fatalf("Fallo 3 data  : %s\n", pdf.Error())
           fmt.Println("Fallo 3 data  ", pdf.Error()) 
       }

         return pdf
	}


