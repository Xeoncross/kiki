package feed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/schollz/kiki/src/database"
	"github.com/schollz/kiki/src/keypair"
	"github.com/schollz/kiki/src/letter"
	"github.com/schollz/kiki/src/logging"
	"github.com/schollz/kiki/src/purpose"
	"github.com/schollz/kiki/src/utils"
)

// New generates a new feed based on the location to find the identity file, the database, and the settings
func New(location ...string) (f Feed, err error) {
	locationToSaveData := "."
	if len(location) > 0 {
		locationToSaveData = location[0]
	}

	locationToSaveData, err = filepath.Abs(locationToSaveData)
	if err != nil {
		return
	}
	f = Feed{
		storagePath: locationToSaveData,
		Settings:    GenerateSettings(),
	}

	// initialize
	err = f.init()
	return
}

// Open will load a feed from the specified location
func Open(locationToFeed string) (f Feed, err error) {
	bFeed, err := ioutil.ReadFile(path.Join(locationToFeed, "feed.json"))
	if err != nil {
		return
	}
	err = json.Unmarshal(bFeed, &f)
	if err != nil {
		return
	}

	// initialize
	err = f.init()
	return
}

// init initializes the kiki instance
func (f *Feed) init() (err error) {
	if f.RegionKey.Public == "" {
		// define region key
		f.RegionKey, err = keypair.FromPair("rbcDfDMIe8qXq4QPtIUtuEylDvlGynx56QgeHUZUZBk=",
			"GQf6ZbBbnVGhiHZ_IqRv0AlfqQh1iofmSyFOcp1ti8Q=") // define region key
		if err != nil {
			return
		}
	}

	f.log = logging.Log
	f.log.Info(f.storagePath)
	f.db = database.Setup(f.storagePath)

	// Setup identity file
	f.log.Debug("setting up personalKey")
	identityFile := path.Join(f.storagePath, "identity.json")
	if _, err := os.Stat(identityFile); os.IsNotExist(err) {
		var err2 error
		// generate a new personal key
		f.PersonalKey = keypair.New()

		// add the friends key
		err2 = f.AddFriendsKey()
		if err2 != nil {
			return err2
		}

		// block the region public key from being used as a sender, ever
		err2 = f.ProcessLetter(letter.Letter{
			To:      []string{"public"},
			Purpose: purpose.AssignBlock,
			Content: f.RegionKey.Public,
		})
		if err2 != nil {
			err2 = errors.Wrap(err2, "setup")
			return err2
		}
	}

	// overwrite the feed file
	feedBytes, err := json.MarshalIndent(f, "", " ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join(f.storagePath, "feed.json"), feedBytes, 0644)
	return
}

// ProcessLetter will determine where to put the letter
func (f Feed) ProcessLetter(l letter.Letter) (err error) {
	if !purpose.Valid(l.Purpose) {
		err = errors.New("invalid purpose")
		return
	}

	// rewrite the letter.To array so that it contains
	// public keys that are valid
	newTo := []string{}
	for _, to := range l.To {
		switch to {
		case "public":
			newTo = append(newTo, f.RegionKey.Public)
		case "self":
			// automatically done when adding any letter
			// this just put here for pedantic reasons
		case "friends":
			friendsKeyPairs, err2 := f.db.GetKeysFromSender(f.PersonalKey.Public)
			if err2 != nil {
				return err2
			}
			for _, friendsKeyPair := range friendsKeyPairs {
				newTo = append(newTo, friendsKeyPair.Public)
			}
		default:
			_, err2 := keypair.FromPublic(to)
			if err2 != nil {
				f.log.Infof("Not a valid public key: '%s'", to)
			} else {
				newTo = append(newTo, to)
			}
		}
	}
	l.To = newTo

	// seal the letter
	if f.PersonalKey == f.RegionKey {
		err = errors.New("cannot post with region key")
		return
	}
	e, err := l.Seal(f.PersonalKey, f.RegionKey)
	if err != nil {
		return
	}

	err = f.db.AddEnvelope(e)
	if err != nil {
		return
	}

	err = f.UnsealLetters()
	return
}

// UnsealLetters will go through unopened envelopes and open them and then add them to the f.db. Also go through and purge bad letters (invalidated letters)
func (f Feed) UnsealLetters() (err error) {
	lettersToPurge := []string{}
	envelopes, err := f.db.GetAllEnvelopes(false)
	if err != nil {
		return err
	}
	keysToTry := []keypair.KeyPair{f.PersonalKey, f.RegionKey}
	for _, envelope := range envelopes {
		if err := envelope.Validate(f.RegionKey); err != nil {
			// add to purge
			lettersToPurge = append(lettersToPurge, envelope.ID)
		}
		ue, err := envelope.Unseal(keysToTry, f.RegionKey)
		if err != nil {
			f.log.Debug(err)
			continue
		}
		f.log.Debug(ue.Letter)
		err = f.db.AddEnvelope(ue)
		if err != nil {
			f.log.Debug(err)
			continue
		}
	}

	// purge invalid letters
	if len(lettersToPurge) > 0 {
		err = f.db.RemoveLetters(lettersToPurge)
	}
	return
}

func (f Feed) ShowFeed() (err error) {
	envelopes, err := f.db.GetAllEnvelopes(true)
	if err != nil {
		return
	}
	f.log.Debugf("Found %d envelopes", len(envelopes))
	for _, e := range envelopes {
		if e.Letter.Purpose != purpose.ShareText {
			continue
		}
		senderName, err2 := f.db.GetName(e.Sender.Public)
		if err2 != nil {
			f.log.Warn(err2)
			senderName = e.Sender.Public
		}
		fmt.Printf("%s (%s) [%s]:\n%s\n\n", senderName, e.Sender.Public, utils.TimeAgo(e.Timestamp), e.Letter.Content)

	}
	return
}

// AddFriendsKey will generate a new friends key and post it to the feed
func (f Feed) AddFriendsKey() (err error) {
	// generate a key for friends
	myfriends := keypair.New()
	if err != nil {
		err = errors.Wrap(err, "AddFriendsKey")
		return
	}
	myfriendsByte, err := json.Marshal(myfriends)

	// share the friends key with yourself
	err = f.ProcessLetter(letter.Letter{
		To:      []string{"self"},
		Purpose: purpose.ShareKey,
		Content: string(myfriendsByte),
	})
	if err != nil {
		err = errors.Wrap(err, "AddFriendsKey")
		return
	}

	// block the friends public key from being used as a sender, ever
	err = f.ProcessLetter(letter.Letter{
		To:      []string{"public"},
		Purpose: purpose.AssignBlock,
		Content: myfriends.Public,
	})
	if err != nil {
		err = errors.Wrap(err, "AddFriendsKey")
		return
	}

	return
}

// GetEnvelope will return an envelope with the given ID
func (f Feed) GetEnvelope(id string) (e letter.Envelope, err error) {
	return f.db.GetEnvelopeFromID(id)
}
