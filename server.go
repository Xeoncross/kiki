package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/schollz/kiki/src/feed"
	"github.com/schollz/kiki/src/logging"
	"github.com/schollz/kiki/src/web"
)

var (
	// Port defines what port the carrier should listen on
	PublicPort  = "8003"
	PrivatePort = "8004"

	NoSync bool
	// Location defines where to open up the kiki database
	Location = "."
	f        feed.Feed
	log      = logging.Log
)

func MiddleWareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request
		log.Debug(fmt.Sprintf("%v %v %v", c.Request.RemoteAddr, c.Request.Method, c.Request.URL))
		// Add base headers
		AddCORS(c)
		// Run next function
		c.Next()
	}
}

func minus(a, b int) string {
	return strconv.FormatInt(int64(a)-int64(b), 10)
}

// Run will start the server listening
func Run() (err error) {
	// Startup feed
	log.Debug("opening feed")
	f, err = feed.Open(Location)
	fmt.Println(err)
	if err != nil {
		log.Debug(err)
		log.Debug("feed not found, creating new one")
		var err2 error
		f, err2 = feed.New(Location)
		if err2 != nil {
			return err2
		}
	}
	log.Debug("opened feed")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(f feed.Feed) {
		<-c
		f.Cleanup()
		os.Exit(1)
	}(f)
	defer f.Cleanup()

	if !NoSync {
		go f.DoSyncing()
	}

	// Startup server
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(MiddleWareHandler(), gin.Recovery()) // Standardize logs
	r.HTMLRender = loadTemplates("index.tmpl", "client.html")
	r.HEAD("/", func(c *gin.Context) { // handler for the uptime robot
		c.String(http.StatusOK, "OK")
	})
	r.GET("/", func(c *gin.Context) {
		p := feed.ShowFeedParameters{}
		p.ID = c.DefaultQuery("id", "")
		p.Channel = c.DefaultQuery("channel", "")
		p.User = c.DefaultQuery("user", "")
		p.Search = c.DefaultQuery("search", "")
		p.Latest = c.DefaultQuery("latest", "") == "1"

		posts, _ := f.ShowFeed(p)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Posts":   posts,
			"User":    f.GetUser(),
			"Friends": f.GetUserFriends(),
		})
	})

	r.GET("/client", func(c *gin.Context) {
		p := feed.ShowFeedParameters{}
		p.ID = c.DefaultQuery("id", "")
		p.Channel = c.DefaultQuery("channel", "")
		p.User = c.DefaultQuery("user", "")
		p.Search = c.DefaultQuery("search", "")
		p.Latest = c.DefaultQuery("latest", "") == "1"
		posts, _ := f.ShowFeed(p)
		c.HTML(http.StatusOK, "client.html", gin.H{
			"Posts": posts,
			"User":  f.GetUser(),
		})
	})

	r.GET("/api/v1/posts.json", func(c *gin.Context) {
		p := feed.ShowFeedParameters{}
		p.ID = c.DefaultQuery("id", "")
		p.Channel = c.DefaultQuery("channel", "")
		p.User = c.DefaultQuery("user", "")
		p.Search = c.DefaultQuery("search", "")
		p.Latest = c.DefaultQuery("latest", "") == "1"
		posts, err := f.ShowFeed3(p)
		if err != nil {
			respondWithJSON(c, "", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data": gin.H{
				"posts": posts,
			},
		})
	})

	r.GET("/api/v1/user.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data": gin.H{
				"user": f.GetUser(),
			},
		})
	})

	r.GET("/api/v1/fiends.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data": gin.H{
				"friends": f.GetUserFriends(),
			},
		})
	})

	r.GET("/static/:file", func(c *gin.Context) {
		file := c.Param("file")
		filename := "static/" + file
		data, err := Asset(filename)
		if err != nil {
			c.String(http.StatusInternalServerError, "uh-oh")
		} else {
			c.Data(http.StatusOK, web.ContentType(filename), data)
		}
	})
	r.GET("/ping", handlePing)
	r.GET("/img/:id", handleImage)
	r.POST("/letter", handlerLetter)       // post to put in letter (local only)
	r.OPTIONS("/letter", handlePing)       // post to put in letter (local only)
	r.POST("/sync", handlerSync)           // tell server to sync with another server (local only)
	r.GET("/list", handleList)             // GET list of all envelope IDs
	r.POST("/envelope", handlerEnvelope)   // post to put into database (public)
	r.GET("/download/:id", handleDownload) // download a specific envelope
	r.GET("/test", func(c *gin.Context) {
		message := ""
		f.TestStuff()
		c.JSON(http.StatusOK, gin.H{"success": err == nil, "message": message})
	})

	// PUBLIC FACING ROUTES
	publicRouter := gin.New()
	publicRouter.Use(MiddleWareHandler(), gin.Recovery())
	publicRouter.GET("/ping", handlePing)             // PING a kiki server to see if it is available
	publicRouter.GET("/list", handleList)             // GET list of all envelope IDs
	publicRouter.POST("/envelope", handlerEnvelope)   // post to put into database (public)
	publicRouter.GET("/download/:id", handleDownload) // download a specific envelope
	go (func() {
		log.Infof("Running public router on 0.0.0.0:%s", PublicPort)
		err = publicRouter.Run(":" + PublicPort)
		if err != nil {
			log.Error(err)
			panic(err)
		}
	})()

	// private routes bind to localhost
	log.Infof("Running private router on localhost:%s", PrivatePort)
	err = r.Run("localhost:" + PrivatePort)
	if nil != err {
		log.Error(err)
		return
	}
	return
}

func respondWithJSON(c *gin.Context, message string, err error) {
	if nil != err {
		log.Error(fmt.Sprintf("%v %v %v [%v]", c.Request.RemoteAddr, c.Request.Method, c.Request.URL, 500))
		log.Warn(err)
		c.JSON(500, gin.H{"status": "error", "error": err.Error()})
		return
	}
	log.Info(fmt.Sprintf("%v %v %v [%v]", c.Request.RemoteAddr, c.Request.Method, c.Request.URL, 200))
	log.Debug(message)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": message})
}

func handlerLetter(c *gin.Context) {
	respondWithJSON(c, "letter added", handleLetter(c))
}

func handlerEnvelope(c *gin.Context) {
	respondWithJSON(c, "envelope added", handleEnvelope(c))
}

func handlerSync(c *gin.Context) {
	respondWithJSON(c, "synced", handleSync(c))
}

func readFormFile(file *multipart.FileHeader) (data []byte, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, src)
	data = buf.Bytes()
	return
}

func AddCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
}

// func ValidateLocalAddress(c *gin.Context) (valid bool) {
// 	clientIP, errIP := web.GetClientIPHelper(c.Request)
// 	if errIP != nil {
// 		return
// 	}
// 	log.Debugf("Got IP adddress: '%s'", clientIP)
// 	return clientIP == "127.0.0.1" || clientIP == "::1"
// }

func loadTemplates(list ...string) multitemplate.Render {
	r := multitemplate.New()
	funcMap := template.FuncMap{
		"minus": minus,
	}
	for _, x := range list {
		templateString, err := Asset("templates/" + x)
		if err != nil {
			panic(err)
		}
		tmplMessage, err := template.New(x).Funcs(funcMap).Parse(string(templateString))
		if err != nil {
			panic(err)
		}
		r.Add(x, tmplMessage)
	}
	return r
}
