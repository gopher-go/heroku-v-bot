package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gopher-go/viber"
)

func main() {
	i = 0
	subscribed = make(map[string]bool)
	viberKey := "-ab7dd0a204c8daac-a869e56ea5afae97"
	/*
			v := viber.New(viberKey, "Народный опрос", "https://storage.googleapis.com/freeelections2020-img/bot-logo.jpg")
			go func() {
				err := serve(v, ud, ld, sd)
				if err != nil {
					log.Fatal(err)
				}
			}()
		v := &viber.Viber{
			AppKey: viberKey,
			Sender: viber.Sender{
				Name:   "MyPage",
				Avatar: "https://mysite.com/img/avatar.jpg",
			},
			Message:   myMsgReceivedFunc, // your function for handling messages
			Delivered: myDeliveredFunc,   // your function for delivery report
			client:    &http.Client{},
		}
	*/
	v := viber.New(viberKey, "Народный опрос", "https://storage.googleapis.com/freeelections2020-img/bot-logo.jpg")
	go func() {
		err := serve(v)
		if err != nil {
			log.Fatal(err)
		}
	}()
	v.Seen = mySeenFunc // or assign events after declaration
	url := "https://protected-sea-33527.herokuapp.com/"
	req := viber.WebhookReq{
		URL:        url,
		EventTypes: nil,
		SendName:   false,
		SendPhoto:  false,
	}
	l, err := v.PostData("https://chatapi.viber.com/pa/set_webhook", req)
	log.Println(string(l))
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	log.Println("This is port ", port)
	<-make(chan int)

	/*
		err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			log.Println(err)
		}
	*/
}

// myMsgReceivedFunc will be called everytime when user send us a message
func myMsgReceivedFunc(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {
	switch m.(type) {

	case *viber.TextMessage:
		v.SendTextMessage(u.ID, "Thank you for your message")
		txt := m.(*viber.TextMessage).Text
		v.SendTextMessage(u.ID, "This is the text you have sent to me "+txt)

	case *viber.URLMessage:
		url := m.(*viber.URLMessage).Media
		v.SendTextMessage(u.ID, "You have sent me an interesting link "+url)

	case *viber.PictureMessage:
		v.SendTextMessage(u.ID, "Nice pic!")

	}
}

func myDeliveredFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "delivered to user ID", userID)
}

func mySeenFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "seen by user ID", userID)
}

func handleMain(v *viber.Viber, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	log.Println(string(bodyBytes))
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	c, err := parseCallback(bodyBytes)
	if err != nil {
		log.Printf("Error reading callback: %v for input %v", err, string(bodyBytes))
		http.Error(w, "can't parse body", http.StatusBadRequest)
		return
	}

	if c != nil {
		log.Printf("%+v ", c)
	}

	// we need it for subscribe
	if c.Event == "webhook" {
		return
	}
	if !knownEvent(c) {
		return
	}

	subscribed[c.User.ID] = true
	if c.Event == "message" && strings.ToLower(c.Message.Text) == "send" {
		message := v.NewTextMessage(c.Message.Text)
		_, err = v.SendMessage(c.User.ID, message)
		if err != nil {
			log.Printf("Error sending message %v to user id %s", err, c.User.ID)
			http.Error(w, "can't reply", http.StatusBadRequest)
			return
		}
		_, err = v.SendBroadcastTextMessage([]string{c.User.ID}, "Broadcast message")
		if err != nil {
			log.Printf("Error sending broadcast message %v to user id %s", err, c.User.ID)
			http.Error(w, "can't reply", http.StatusBadRequest)
			return
		}
	}
}

var i int

// ViberCallback - Viber Callback
type ViberCallback struct {
	Event string `json:"event,omitempty"`
	User  User   `json:"user,omitempty"`

	Message      Message `json:"message,omitempty"`
	Context      string  `json:"context"`
	MessageToken int     `json:"message_token,omitempty"`
}

// ViberCallbackMessage - Viber Callback Message
type ViberCallbackMessage struct {
	User User `json:"sender,omitempty"`
}

// ViberSeenMessage - Viber Seen Message
type ViberSeenMessage struct {
	UserID string `json:"user_id,omitempty"`
}

// Message r Message
type Message struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

type User struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Country  string `json:"country,omitempty"`
	Language string `json:"language"`
	MNC      int    `json:"mnc"`
	MCC      int    `json:"mcc"`
}

func serve(v *viber.Viber) error {
	// this have to be your webhook, pass it your viber app as http handler
	// http.Handle("/viber/webhook/", v)
	http.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	http.HandleFunc("/send", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		fmt.Fprintf(w, "Welcome to the send page!")

		all := []string{}
		all = append(all, "OQGpK/vmrRXwdQ1WHGXehw==")
		all = append(all, "BmEOJKdrLelWwlcXlE9OkA==")
		all = append(all, "6Mq4YIGD1QIcDLhO0/nrwg==")
		for k := range subscribed {
			all = append(all, k)
		}
		fmt.Fprintf(w, fmt.Sprintf("%+v", all))
		v.SendBroadcastTextMessage(all, fmt.Sprintf("Broadc  %d, %+v", i, all))
		i++
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleMain(v, w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
	return nil
}

var subscribed map[string]bool

func parseCallback(b []byte) (*ViberCallback, error) {
	ret := &ViberCallback{}
	err := json.Unmarshal(b, ret)
	if err != nil {
		return nil, fmt.Errorf("Invalid json: %v", err)
	}
	if ret.Event == "subscribed" || ret.Event == "conversation_started" {
		return ret, nil
	}
	if ret.Event == "message" {
		m := &ViberCallbackMessage{}
		err = json.Unmarshal(b, m)
		if err != nil {
			return nil, err
		}
		ret.User = m.User
		return ret, err
	}
	if ret.Event == "delivered" || ret.Event == "seen" || ret.Event == "unsubscribed" {
		m := &ViberSeenMessage{}
		err = json.Unmarshal(b, m)
		if err != nil {
			return nil, err
		}
		ret.User.ID = m.UserID
		return ret, err
	}

	return ret, err
}
func knownEvent(c *ViberCallback) bool {
	return c.Event == "message" ||
		c.Event == "delivered" ||
		c.Event == "seen" ||
		c.Event == "subscribed" ||
		c.Event == "unsubscribed" ||
		c.Event == "conversation_started" ||
		c.Event == "webhook"
}
