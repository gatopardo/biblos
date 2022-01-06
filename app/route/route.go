package route

import (
	"net/http"
        "log"

	"github.com/gatopardo/biblos/app/controller"
	"github.com/gatopardo/biblos/app/route/middleware/acl"
	hr "github.com/gatopardo/biblos/app/route/middleware/httprouterwrapper"
	"github.com/gatopardo/biblos/app/route/middleware/logrequest"
	"github.com/gatopardo/biblos/app/route/middleware/pprofhandler"
	"github.com/gatopardo/biblos/app/model"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

  var Flogger  *log.Logger

// Load returns the routes and middleware
func Load() http.Handler {
           Flogger.Println("HTTP routes Load")
	return middleware(routes())
}

// LoadHTTPS returns the HTTP routes and middleware
func LoadHTTPS() http.Handler {
           Flogger.Println("HTTPS routes LoadHTTPS")
	return middleware(routes())
}

// LoadHTTP returns the HTTPS routes and middleware
func LoadHTTP() http.Handler {
           Flogger.Println("HTTPS routes LoadHTTP")
	return middleware(routes())

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
           Flogger.Println("HTTP redirect")
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *httprouter.Router {
	r := httprouter.New()
	// Set 404 handler
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.
		New().
		ThenFunc(controller.Static)))

	// Home page
	r.GET("/", hr.Handler(alice.
		New().
		ThenFunc(controller.IndexGET)))

	// Login
	r.GET("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginGET)))
	r.POST("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginPOST)))
        r.GET("/jlogin/:cuenta/:password", hr.Handler(alice.
                New(acl.DisallowAuth).
                ThenFunc(controller.JLoginGET)))
	r.GET("/logout", hr.Handler(alice.
		New().
		ThenFunc(controller.LogoutGET)))

// Register
	r.GET("/user/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisterGET)))
	r.POST("/user/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisterPOST)))
//          Register update
	r.GET("/user/update/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisUpGET)))
	r.POST("/user/update/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisUpPOST)))
//          List
	r.GET("/user/list/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisLisGET)))
//          Delete
	r.GET("/user/delete/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisDeleteGET)))
	r.POST("/user/delete/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisDeletePOST)))

// Language
        r.GET("/biblos/jlang/:re", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.JLangListGET)))
	r.GET("/language/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.LangGET)))
	r.POST("/language/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.LangPOST)))
	// languages update
	r.GET("/language/update/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.LangUpGET)))
	r.POST("/language/update/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.LangUpPOST)))
//          List
	r.GET("/language/list/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.LangLisGET)))
//          Delete
	r.GET("/language/delete/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.LangDeleteGET)))
	r.POST("/language/delete/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.LangDeletePOST)))

// Editor
        r.GET("/biblos/jeditor/:re", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.JEditorListGET)))
	r.GET("/editor/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditGET)))
	r.POST("/editor/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditPOST)))
	// editor update
	r.GET("/editor/update/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditUpGET)))
	r.POST("/editor/update/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditUpPOST)))
//          List
	r.GET("/editor/list/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditLisGET)))
//          Delete
	r.GET("/editor/delete/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditDeleteGET)))
	r.POST("/editor/delete/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditDeletePOST)))
	// search
	r.POST("/editor/search", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditorSearchPOST)))


// Author
        r.GET("/biblos/jauthor/:re", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.JAuthorListGET)))
	r.GET("/author/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorGET)))
	r.POST("/author/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorPOST)))
	// author update
	r.GET("/author/update/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorUpGET)))
	r.POST("/author/update/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorUpPOST)))
//          List
	r.GET("/author/list/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorLisGET)))
//          Delete
	r.GET("/author/delete/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorDeleteGET)))
	r.POST("/author/delete/:id/:pagi", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorDeletePOST)))
	// search
	r.POST("/author/search", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorSearchPOST)))

	// Biblos
        r.GET("/biblos/jbook/:re", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.JBookListGET)))
	r.GET("/biblos/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookGET)))
	r.POST("/biblos/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookPOST)))

//         biblos update
	r.GET("/biblos/update/:id/:pagi/:orig", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookUpGET)))
	r.POST("/biblos/update/:id/:pagi/:orig", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookUpPOST)))
//             search
	r.POST("/biblos/search", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookSearchPOST)))
//          List
	r.GET("/biblos/list/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookLisGET)))
//          Delete
	r.GET("/biblos/delete/:id/:pagi/:orig", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookDeleteGET)))
	r.POST("/biblos/delete/:id/:pagi/:orig", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookDeletePOST)))
//          List given language
	r.GET("/biblos/lang/list/:id/:pg/:at", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookLangGET)))
//          List given Editor
	r.GET("/biblos/edit/list/:id/:pg/:at", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookEditGET)))
//          List given Author
	r.GET("/biblos/author/list/:id/:pg/:at", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookAuthGET)))

//         pdf
	r.GET("/biblos/pdf/book", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookPdfGET)))
	r.GET("/biblos/pdf/bookauth/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookAuthPdfGET)))
	r.GET("/biblos/pdf/bookedit/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookEditPdfGET)))
	r.GET("/biblos/pdf/booklang/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BookLangPdfGET)))
	r.GET("/author/pdf/author", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AuthorPdfGET)))

	// About
	r.GET("/about", hr.Handler(alice.
		New().
		ThenFunc(controller.AboutGET)))


	// Enable Pprof
	r.GET("/debug/pprof/*pprof", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(pprofhandler.Handler)))

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, model.Store, model.Name)
	cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "token"
	csrfbanana.SingleToken = false
	h = cs

      Flogger.Println("middleware", model.Name)

	// Log every request:1
	h = logrequest.Handler(h, Flogger)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
