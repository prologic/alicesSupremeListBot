package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/prologic/bitcask"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	groups map[int64]bool
	admins map[int]bool
)

const (
	listDirectory string = "db/lists/"
	adminLocation string = "db/admins"
	groupLocation string = "db/groups"
)

func init() {
	admins = make(map[int]bool)
	groups = make(map[int64]bool)

	admindb, _ := bitcask.Open(adminLocation)
	defer admindb.Close()
	err := admindb.Fold(func(key string) error {
		ikey, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return err
		}
		admins[int(ikey)] = true
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to load admins: %v", err)
	}

	groupdb, _ := bitcask.Open(groupLocation)
	defer groupdb.Close()
	err = groupdb.Fold(func(key string) error {
		ikey, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return err
		}
		groups[ikey] = true
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to load admins: %v", err)
	}
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token: os.Getenv("TELEGRAMTOKEN"),
		// You can also set custom API URL. If field is empty it equals to "https://api.telegram.org"
		// URL:    "http://195.129.111.17:8012",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/admin", func(m *tb.Message) {
		if admins[m.Sender.ID] {
			b.Send(m.Chat, handleAdmin(m))
		} else {
			fmt.Printf("Unauthorized User: %d", m.Sender.ID)
			b.Send(m.Chat, "Unauthorized User")
		}
	})

	b.Handle("/list", func(m *tb.Message) {
		if groups[m.Chat.ID] {
			b.Send(m.Chat, handleList(m.Payload))
		} else {
			b.Send(m.Chat, "Unauthorized Group")
		}
	})

	b.Handle("/oof", func(m *tb.Message) {
		b.Send(m.Chat, "oof")
	})

	b.Start()
}

func handleAdmin(m *tb.Message) string {
	args := strings.SplitN(m.Payload, " ", 2)
	if len(args) != 2 || strings.ToLower(args[0]) != "add" {
		return "Invalid Command"
	}
	switch {
	case strings.ToLower(args[1]) == "group":
		return addGroup(m)
	case strings.ToLower(args[1]) == "admin":
		return addAdmin(m)
	default:
		return "Invalid Command"
	}
}

func addGroup(m *tb.Message) string {
	db, err := bitcask.Open(groupLocation)
	if err != nil {
		fmt.Printf("db error opening groupdb: %v", err)
		return "Failed to add group to db."
	}
	defer db.Close()

	db.Put(strconv.FormatInt(m.Chat.ID, 10), []byte("1"))
	groups[m.Chat.ID] = true
	return "Group added."
}

func addAdmin(m *tb.Message) string {
	db, err := bitcask.Open(adminLocation)
	if err != nil {
		fmt.Printf("db error opening groupdb: %v", err)
		return "Failed to add group to db."
	}
	defer db.Close()
	for _, entities := range m.Entities {
		if entities.User != nil {

			db.Put(strconv.FormatInt(int64(entities.User.ID), 10), []byte("1"))
			admins[entities.User.ID] = true
		}
	}
	return "User(s) added."
}

func handleList(payload string) string {
	args := strings.SplitN(payload, " ", 2)
	switch {
	case len(args) == 1 && args[0] != "":
		return printList(args[0])
	case len(args) > 1:
		return addToList(args[0], args[1])
	default:
		return "Invalid Command"
	}
}

func printList(list string) string {
	if list == "" {
		return "List needs a name."
	}
	db, err := bitcask.Open(listDirectory + list)
	if err != nil {
		fmt.Printf("db error opening list at %s: %v", listDirectory+list, err)
		return "Failed to open List."
	}
	defer db.Close()
	items := make([]string, 0)
	err = db.Fold(func(key string) error {
		items = append(items, key)
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to load List: %v", err)
	}
	sort.Strings(items)
	listString := ""
	for _, item := range items {
		listString += "• "
		entry, err := db.Get(item)
		if err != nil {
			fmt.Printf("Getting item failed: %v", err)
			return "Failed to load List items."
		}
		listString += string(entry)
		listString += "\n"
	}
	if listString == "" {
		return "<-- List is Empty -->"
	}
	return listString
}

func addToList(list, message string) string {
	db, err := bitcask.Open(listDirectory + list)
	if err != nil {
		fmt.Printf("db error opening list at %s: %v", listDirectory+list, err)
		return "Failed to add to list."
	}
	defer db.Close()
	db.Put(strconv.FormatInt(time.Now().Unix(), 10), []byte(message))
	return "Item added."
}
