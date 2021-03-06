package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/schollz/kiki/src/feed"
	"github.com/schollz/kiki/src/letter"
)

func handleView(c *gin.Context) (posts []feed.Post) {
	p := feed.ShowFeedParameters{}
	p.ID = c.DefaultQuery("id", "")
	p.Hashtag = c.DefaultQuery("hashtag", "")
	p.User = c.DefaultQuery("user", "")
	p.Search = c.DefaultQuery("search", "")
	p.Latest = c.DefaultQuery("latest", "") == "1"
	posts, _ = f.ShowFeed(p)
	return
}

func handleSlash(c *gin.Context) {
	posts := handleView(c)
	showPosts(c, posts)
}

func handleHome(c *gin.Context) {
	posts := handleView(c)
	posts = f.OnlyIncludePostsFromFollowing(posts)
	showPosts(c, posts)
}

func showPosts(c *gin.Context, posts []feed.Post) {

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Posts":          posts,
		"User":           f.GetUser(),
		"Friends":        f.GetUserFriends(),
		"Connected":      f.GetConnected(),
		"Hashtags":       f.GetHashTags(),
		"RegionPublic":   RegionPublic,
		"RegionPrivate":  RegionPrivate,
		"ServerName":     strings.TrimLeft(strings.TrimLeft(ServerName, "http://"), "https://"),
		"ServerNameFull": ServerName,
	})
}

// GET /img
func handleImage(c *gin.Context) {
	id := c.Param("id")
	logger.Log.Debugf("fetching image: %s", id)
	e, err := f.GetEnvelope(id)
	if err != nil {
		logger.Log.Warn(err)
		c.Data(http.StatusInternalServerError, "text/plain", []byte(err.Error()))
		return
	}
	if f.Settings.BlockPublicPhotos {
		if !f.AmFollowing(e.Sender.Public) {
			c.Data(http.StatusInternalServerError, "text/plain", []byte("not allowing public photos"))
			return
		}
	}
	mimeType := "image/jpeg"
	if strings.Contains(e.Letter.Purpose, "png") {
		mimeType = "image/png"
	}

	imageBytes, err := base64.StdEncoding.DecodeString(e.Letter.Content)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/plain", []byte(err.Error()))
		return
	}

	c.Data(http.StatusOK, mimeType, imageBytes)
}

// POST /letter
func handleLetter(c *gin.Context) (err error) {
	// bind the payload
	var p letter.Letter
	err = c.BindJSON(&p)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(500, gin.H{"status": "error", "error": err.Error()})
		return
	}
	e, err := f.ProcessLetter(p)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(500, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "message": "added " + e.ID, "envelope": e})

	// when a new letter arrives, update everything and then sync servers
	go f.UpdateEverythingAndSync()
	return
}

// GET /ping
func handlePing(c *gin.Context) {
	fmt.Printf("%+v", c.Request)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "error": "pong"})
}

// POST /envelope
func handleEnvelope(c *gin.Context) (err error) {
	// bind the payload
	var p letter.Envelope
	err = c.BindJSON(&p)
	if err != nil {
		return
	}
	err = f.ProcessEnvelope(p)
	f.SignalUpdate()
	return
}

// GET /list?user_pub=X&signature=+
func handleList(c *gin.Context) {
	pubkey := c.DefaultQuery("user_pub", "")
	signature := c.DefaultQuery("signature", "")

	idList, err := f.GetIDs(pubkey, signature)
	personalSignature, _ := f.PersonalKey.Signature(f.RegionKey)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(500, gin.H{"status": "error", "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "found IDs", "ids": idList, "personal_key": f.PersonalKey.Public, "personal_signature": personalSignature})
	}
	return
}

// GET /download/ID
// You can always download anything you want but the envelopes are transfered so that the letter is closed up.
func handleDownload(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	e, err := f.GetEnvelope(id)
	// Close up envelope
	e.Close()
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "found envelope", "envelope": e})
	}
}

// POST /sync
func handleSync(c *gin.Context) (err error) {
	// bind the payload
	type Payload struct {
		Address string `json:"address" binding"required"`
	}
	var p Payload
	err = c.BindJSON(&p)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	if p.Address == "" {
		logger.Log.Debug("only syncing servers")
		f.SyncServers()
		return
	}

	logger.Log.Debug("syncing...")
	err = f.Sync(p.Address)
	if err != nil {
		logger.Log.Error(err)
	} else {
		go f.UpdateEverything()
	}
	return
}
