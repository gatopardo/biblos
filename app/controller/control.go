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
   func strToReg(s string) ( sreg string){
	   var strinic = "@#:!&/"
	   var strepl  = "[]{}^|"
	   sreg = s
	   for i, _ := range strinic {
		   c1:= strinic[i:i+1]
		   c2:= strepl[i:i+1]
	        sreg = strings.ReplaceAll(sreg, c1,c2)
	   }
       return
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

 func addData(pdf *gofpdf.Fpdf, lh string, lhead[]string,  lw []float64,fill bool, books []model.BookN) *gofpdf.Fpdf {
	align := []string{"L", "L", "L", "R",  "L"}
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.SetFillColor(197, 236, 235)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
//	pdf.SetFont("Arial", "", 12)
        pdf.SetFont("Times", "", 11)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
	pdf.AddPage()
        ht   := 6.0
        tr := pdf.UnicodeTranslatorFromDescriptor("")

	for _, c := range books {
		fields := getFields(c, pdf, lw)
		for j, field := range fields {
                     pdf.CellFormat(lw[j], ht, tr(field), "1", 0, align[j], fill, 0, "")
                }
                pdf.Ln(-1)
                fill = ! fill
	  }
          pdf.CellFormat(255.0, 0, "", "TR", 0, "", false, 0, "")
   return pdf
 }

func header(pdf *gofpdf.Fpdf,stit string,  hdr []string, whdr []float64, fill bool ) *gofpdf.Fpdf {
        fill = false
	pdf.SetFillColor(213, 234, 235)
        pdf.SetDrawColor(128, 0, 0)
	pdf.SetFont("Times", "B", 16)
	pdf.CellFormat(80, 10, stit,"",0,"C",true, 0, "")
	pdf.SetFont("Times", "B", 8)
        sttime := time.Now().Format("2006/01/02")
	pdf.CellFormat(110, 10, sttime,"",0,"R", true, 0, "")
	pdf.Ln(8)

	pdf.SetFillColor(120, 140, 120)
        pdf.SetDrawColor(239, 234, 228)
	pdf.SetFont("Times", "B", 12)
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
//         pdf.Image("static/favicons/ipi_7a.png", 100, 236, 25, 0, false, "", 0, "")
//         pdf.Line(10,275,200,275)
	return pdf
  }

 func getBookPdf(books []model.BookN ) (pdf * gofpdf.Fpdf)  {
	pdf = gofpdf.New("P", "mm", "Letter", "")
        fill := false
        if pdf.Err() { log.Fatalf("Fallo 1 creando PDF : %s\n", pdf.Error()) }

	pdf.SetTitle("Libros",true)
        stit   := "Listado de Libros"
        stheader := []string  {"Lang", "Author", "Title", "Year", "Editor"}
	szwidth := []float64 { 18.0,    55.0,     80.0,    10.0,     60.0    }

        pdf.SetHeaderFunc(func(){ header(pdf, stit, stheader, szwidth, fill) } )
        if pdf.Err() { log.Fatalf("Fallo 2 header PDF : %s\n", pdf.Error()) }

	pdf.SetFooterFunc(func() {   addFooterBook(pdf )} )
        if pdf.Err() { log.Fatalf("Fallo 4 Footer PDF : %s\n", pdf.Error()) } 

        pdf = addData(pdf,stit,stheader, szwidth,fill,  books )
        if pdf.Err() { log.Fatalf("Fallo 3 data PDF : %s\n", pdf.Error()) }
         return pdf
    }

    func getAuthor( pdf *gofpdf.Fpdf,  lw []float64, authors []model.Author) *gofpdf.Fpdf {
//	pdf.SetFillColor(224, 235, 255)
        pdf.SetFillColor(197, 236, 235)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.AddPage()

	wSum := 175.0
        fill := false
        ht   := 6.0
        nfields := len(lw)
        tr := pdf.UnicodeTranslatorFromDescriptor("")
        write := func(str string,w float64) {
                   str  = leftStr(pdf, str, w)
                   pdf.CellFormat(w, ht, tr(str), "LR", 0, "", fill, 0, "")
                }
	// 	Data
	for _, c := range authors {
                fields := []string{ c.Name }
                for j := 0; j < nfields; j++{
                     write(fields[j], lw[j])
                }
                pdf.Ln(-1)
                fill = ! fill
          }
          pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")

        return pdf
   }

 func getAuthorPdf(authors []model.Author ) (pdf * gofpdf.Fpdf)  {
        pdf = gofpdf.New("P", "mm", "Letter", "")
	SName := "Autores"
        stit := "Listado de Autores"
        lw := []float64{175 }
        pdf.SetHeaderFunc(func() {
           header := []string{ "Author" }
           w := []float64{25 }

	pdf.SetFillColor(247, 216, 151)
        pdf.SetDrawColor(128, 0, 0)
	pdf.SetFont("Times", "B", 16)
	pdf.CellFormat(80, 10, stit,"",0,"C",true, 0, "")
	pdf.SetFont("Times", "B", 8)
        sttime := time.Now().Format("2006/01/02")
	pdf.CellFormat(110, 10, sttime,"",0,"R", true, 0, "")
	pdf.Ln(8)

	pdf.SetFillColor(120, 140, 120)
        pdf.SetDrawColor(239, 234, 228)
	pdf.SetFont("Times", "B", 12)
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
//              pdf.Line(10,275,200,275)
    })

        if pdf.Err() { log.Fatalf("Fallo 2 footer  : %s\n", pdf.Error())
            fmt.Println("Fallo 2 footer  ", pdf.Error()) 
          }

        pdf.SetTitle(SName,true)
        pdf = getAuthor( pdf,  lw, authors )
        if pdf.Err() { log.Fatalf("Fallo 3 data  : %s\n", pdf.Error())
           fmt.Println("Fallo 3 data  ", pdf.Error()) 
        }

         return pdf
	}

  func getAuthFields( book model.BookN, pdf * gofpdf.Fpdf, lw []float64) ( lin []string ) {
       str :=  strings.Trim(book.Title, " ")
       str  = leftStr(pdf, str, lw[0])
       lin = append(lin, str)
       str =  strings.Trim(book.Comment, " ")
       str  = leftStr(pdf, str, lw[1])
       lin = append(lin, str)
       stYear := fmt.Sprint(book.Year)
       lin = append(lin, stYear)
       str =  strings.Trim(book.Editor, " ")
       str  = leftStr(pdf, str, lw[3])
       lin = append(lin, str)
       return lin
  }

    func getBookAuthor( pdf *gofpdf.Fpdf, lw []float64, fill bool, books []model.BookN) *gofpdf.Fpdf {
	align := []string{"L", "L", "R",  "L"}
        pdf.SetFillColor(197, 236, 235)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
        pdf.SetFont("Times", "", 11)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.AddPage()

	wSum := 0.0
	for _, v := range lw {
		wSum += v
	}
        ht   := 6.0
        tr := pdf.UnicodeTranslatorFromDescriptor("")
	// 	Data

	for _, c := range books {
                fields := getAuthFields( c, pdf, lw )
		for j, field := range fields {
                     pdf.CellFormat(lw[j], ht, tr(field), "1", 0, align[j], fill, 0, "")
                }
                pdf.Ln(-1)
                fill = ! fill
	  }
           pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")
             return pdf
    }


 func getBookAuthPdf(books []model.BookN ) (pdf * gofpdf.Fpdf )  {
	pdf = gofpdf.New("P", "mm", "Letter", "")
        SName := books[0].Author
	pdf.SetTitle(SName,true)
        stit := strings.Trim(" Libros de " + SName, " ")
        header := []string{ "Title", "Tema", "Year", "Editor" }
        lw := []float64{ 95,45, 14,64  }
            fill := false

        pdf.SetHeaderFunc(func() {
	   wSum := 0.0
	   for _, v := range lw { wSum += v }
            fill = false

	 pdf.SetFillColor(213, 234, 235)
         pdf.SetDrawColor(128, 0, 0)
         pdf.SetFont("Times", "B", 16)
         pdf.CellFormat(60, 10, stit,"",0,"C",true, 0, "")
         pdf.SetFont("Times", "B", 8)
         sttime := time.Now().Format("2006/01/02")
         pdf.CellFormat(80, 10, sttime,"",0,"R", true, 0, "")
         pdf.Ln(8)

         pdf.SetFillColor(120, 140, 120)
         pdf.SetDrawColor(239, 234, 228)
         pdf.SetFont("Times", "B", 12)

           for j, str := range header {
                pdf.CellFormat(lw[j], 7, str, "1", 0, "C", true, 0, "")
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
//            pdf.Image("static/favicons/ipi_7a.png", 100, 236, 25, 0, false, "", 0, "")
//              pdf.Line(10,275,200,275)
    })

        if pdf.Err() { log.Fatalf("Fallo 2 footer  : %s\n", pdf.Error())
        fmt.Println("Fallo 2 footer  ", pdf.Error())
      }

        pdf.SetTitle(SName,true)
        pdf = getBookAuthor( pdf,  lw,fill, books )
        if pdf.Err() { log.Fatalf("Fallo 3 data  : %s\n", pdf.Error())
           fmt.Println("Fallo 3 data  ", pdf.Error())
       }

         return pdf
	}

  func getEditFields( book model.BookN, pdf * gofpdf.Fpdf, lw []float64) ( lin []string ) {
       str :=  strings.Trim(book.Title, " ")
       str  = leftStr(pdf, str, lw[0])
       lin = append(lin, str)
       str =  strings.Trim(book.Comment, " ")
       str  = leftStr(pdf, str, lw[1])
       lin = append(lin, str)
       stYear := fmt.Sprint(book.Year)
       lin = append(lin, stYear)
       str =  strings.Trim(book.Author, " ")
       str  = leftStr(pdf, str, lw[3])
       lin = append(lin, str)
       return lin
  }

    func getBookEdit( pdf *gofpdf.Fpdf, lw []float64,fill bool, books []model.BookN) *gofpdf.Fpdf {

	align := []string{"L", "L", "R",  "L"}
        pdf.SetFillColor(197, 236, 235)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
        pdf.SetFont("Times", "", 11)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.AddPage()

	wSum := 0.0
	for _, v := range lw {
		wSum += v
	}
        ht   := 6.0
        tr := pdf.UnicodeTranslatorFromDescriptor("")
	// 	Data
	for _, c := range books {
                fields := getEditFields( c, pdf, lw )
		for j, field := range fields {
                     pdf.CellFormat(lw[j], ht, tr(field), "1", 0, align[j], fill, 0, "")
                }
                pdf.Ln(-1)
                fill = ! fill
         }
         pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")

         return pdf
    }

 func getBookEditPdf(books []model.BookN ) (pdf * gofpdf.Fpdf)  {
        pdf = gofpdf.New("P", "mm", "Letter", "")
        SName := books[0].Editor
	pdf.SetTitle(SName,true)
        stit := strings.Trim("Editor " + SName, " ")
        header := []string{ "Title", "Tema", "Year", "Author" }
        lw := []float64{ 95,45, 14,64  }
            fill := false
        pdf.SetHeaderFunc(func() {
	   wSum := 0.0
	   for _, v := range lw {
		  wSum += v
	   }
            fill = false

	 pdf.SetFillColor(213, 234, 235)
         pdf.SetDrawColor(128, 0, 0)
         pdf.SetFont("Times", "B", 16)
         pdf.CellFormat(60, 10, stit,"",0,"C",true, 0, "")
         pdf.SetFont("Times", "B", 8)
         sttime := time.Now().Format("2006/01/02")
         pdf.CellFormat(80, 10, sttime,"",0,"R", true, 0, "")
         pdf.Ln(8)

         pdf.SetFillColor(120, 140, 120)
         pdf.SetDrawColor(239, 234, 228)
         pdf.SetFont("Times", "B", 12)

           for j, str := range header {
                pdf.CellFormat(lw[j], 7, str, "1", 0, "C", true, 0, "")
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
//              pdf.Line(5,200,140,200)
    })

        if pdf.Err() { log.Fatalf("Fallo 2 footer  : %s\n", pdf.Error())
            fmt.Println("Fallo 2 footer  ", pdf.Error()) 
         }

        pdf.SetTitle(SName,true)
        pdf = getBookEdit( pdf,  lw,fill, books )
        if pdf.Err() { log.Fatalf("Fallo 3 data  : %s\n", pdf.Error())
            fmt.Println("Fallo 3 data  ", pdf.Error()) 
        }
         return pdf
	}

  func getLangFields( book model.BookN, pdf * gofpdf.Fpdf, lw []float64) ( lin []string ) {
       str  :=  strings.Trim(book.Author, " ")
       str  = leftStr(pdf, str, lw[0])
       lin  = append(lin, str)
       str  =  strings.Trim(book.Title, " ")
       str  = leftStr(pdf, str, lw[1])
       lin  = append(lin, str)
       str  =  strings.Trim(book.Comment, " ")
       str  = leftStr(pdf, str, lw[2])
       lin  = append(lin, str)
       stYear := fmt.Sprint(book.Year)
       lin = append(lin, stYear)
       str =  strings.Trim(book.Editor, " ")
       str  = leftStr(pdf, str, lw[4])
       lin = append(lin, str)
       return lin
  }


    func getBookLang( pdf *gofpdf.Fpdf, lw []float64, fill bool, books []model.BookN) *gofpdf.Fpdf {

	align := []string{"L", "L", "L", "R",  "L"}
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
        pdf.SetFillColor(197, 236, 235)
	pdf.SetTextColor(0, 0, 0)
        pdf.SetDrawColor(128, 0, 0)
//	pdf.SetFont("Arial", "", 12)
        pdf.SetFont("Times", "", 11)
	pdf.SetLineWidth(.3)
        pdf.AliasNbPages("")
	pdf.AddPage()

	wSum := 0.0
	for _, v := range lw { wSum += v }
        ht   := 6.0
        tr := pdf.UnicodeTranslatorFromDescriptor("")

	// 	Data
	for _, c := range books {
                fields := getLangFields( c, pdf, lw )
		for j, field := range fields {
                     pdf.CellFormat(lw[j], ht, tr(field), "1", 0, align[j], fill, 0, "")
                }
                pdf.Ln(-1)
                fill = ! fill
         }
         pdf.CellFormat(wSum, 0, "", "TR", 0, "", false, 0, "")

         return pdf
	}

 func getBookLangPdf(books []model.BookN ) (pdf * gofpdf.Fpdf)  {
	pdf = gofpdf.New("P", "mm", "Letter", "")
        SName := books[0].Language
        fill := false
	pdf.SetTitle(SName,true)
        stit := strings.Trim("Language " + SName, " ")
        header := []string  { "Author", "Title", "Tema", "Year", "Editor"}
	lw := []float64 {       45.0,     80.0,    35.0,   10.0,     50.0    }
        pdf.SetHeaderFunc(func() {
            fill = false
            pdf.SetFillColor(213, 234, 235)
            pdf.SetDrawColor(128, 0, 0)
            pdf.SetFont("Times", "B", 16)
            pdf.CellFormat(80, 10, stit,"",0,"C",true, 0, "")
            pdf.SetFont("Times", "B", 8)
            sttime := time.Now().Format("2006/01/02")
            pdf.CellFormat(110, 10, sttime,"",0,"R", true, 0, "")
            pdf.Ln(8)

            pdf.SetFillColor(120, 140, 120)
            pdf.SetDrawColor(239, 234, 228)
            pdf.SetFont("Times", "B", 12)
            for i, str := range header {
	         pdf.CellFormat(lw[i],  5, str, "1", 0, "L", true, 0, "")
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
//              pdf.Line(10,275,200,275)
//              pdf.Line(5,200,140,200)
    })

        if pdf.Err() { log.Fatalf("Fallo 2 footer  : %s\n", pdf.Error())
        fmt.Println("Fallo 2 footer  ", pdf.Error())
      }

        pdf = getBookLang( pdf,  lw, fill, books )
        if pdf.Err() { log.Fatalf("Fallo 3 data  : %s\n", pdf.Error())
           fmt.Println("Fallo 3 data  ", pdf.Error())
       }
         return pdf
	}


