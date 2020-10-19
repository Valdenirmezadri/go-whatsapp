package whatsapp

import (
	"strings"
	"sync"
	"github.com/Valdenirmezadri/go-whatsapp/binary"
)
var rwm sync.RWMutex

type Store struct {
	contacts map[string]Contact
	chats    map[string]Chat
}

func (s Store) GetContacts() []*Contact {
	contacts := []*Contact{}
	rwm.RLock()
	defer rwm.RUnlock()
	if s.contacts != nil {
		for _, contact := range s.contacts {
			contacts = append(contacts, &contact)
		}
	}
	return contacts
}

func (s Store) GetContact(Jid string) *Contact {
	rwm.RLock()
	defer rwm.RUnlock()
	if s.contacts != nil {
		if u, ok := s.contacts[Jid]; ok {
			return &u
		}
	}
	return nil
}

func (s *Store) setContact(Jid string,c Contact) {
	rwm.Lock()
	defer rwm.Unlock()
	if s.contacts != nil {
		s.contacts[Jid] = c
	}
}

func (s Store) GetChats() []*Chat {
	chats := []*Chat{}
	rwm.RLock()
	defer rwm.RUnlock()
	if s.chats != nil {
		for _, chat := range s.chats {
			chats = append(chats, &chat)
		}
	}
	return chats
}

func (s Store) GetChat(Jid string) *Chat {
	rwm.RLock()
	defer rwm.RUnlock()
	if s.chats != nil {
		if c, ok := s.chats[Jid]; ok {
			return &c
		}
	}
	return nil
}

func (s *Store) setChat(Jid string,c Chat) {
	rwm.Lock()
	defer rwm.Unlock()
	if s.chats != nil {
		s.chats[Jid] = c
	}
}

type Contact struct {
	Jid    string
	Notify string
	Name   string
	Short  string
}

type Chat struct {
	Jid             string
	Name            string
	Unread          string
	LastMessageTime string
	IsMuted         string
	IsMarkedSpam    string
}

func newStore() *Store {
	return &Store{
		make(map[string]Contact),
		make(map[string]Chat),
	}
}

func (wac *Conn) updateContacts(contacts interface{}) {
	c, ok := contacts.([]interface{})
	if !ok {
		return
	}

	for _, contact := range c {
		contactNode, ok := contact.(binary.Node)
		if !ok {
			continue
		}

		jid := strings.Replace(contactNode.Attributes["jid"], "@c.us", "@s.whatsapp.net", 1)
		wac.Store.setContact(jid, Contact{
			jid,
			contactNode.Attributes["notify"],
			contactNode.Attributes["name"],
			contactNode.Attributes["short"],
		})
	}
}

func (wac *Conn) updateChats(chats interface{}) {
	c, ok := chats.([]interface{})
	if !ok {
		return
	}

	for _, chat := range c {
		chatNode, ok := chat.(binary.Node)
		if !ok {
			continue
		}

		jid := strings.Replace(chatNode.Attributes["jid"], "@c.us", "@s.whatsapp.net", 1)
		wac.Store.setChat(jid, Chat{
			jid,
			chatNode.Attributes["name"],
			chatNode.Attributes["count"],
			chatNode.Attributes["t"],
			chatNode.Attributes["mute"],
			chatNode.Attributes["spam"],
		})
	}
}
