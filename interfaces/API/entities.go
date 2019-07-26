package api

import (
	"crypto/rsa"
	"time"
	entity "squad-3-aceleradev-fs-florianopolis/entities"
	mail "squad-3-aceleradev-fs-florianopolis/services/MailSender/src"
	"github.com/gorilla/mux"
)

const (
	//Success Code for response 
	Success = 1
	//Empty Code for response
	Empty 	= 2
	//Error Code for response
	Error   = 3
 )
 

//Result pattern response for all requests
type Result struct {
	Result    string   `json:"Result"`
	Code	  int	   `json:"Code,omitempty"`	
	Token     string   `json:"token,omitempty"`
	Warn      *Warn     `json:"Warn,omitempty"`
	Warns     *WarnList `json:"WarnList,omitempty"`
	Mail      *MailType `json:"Mail,omitempty"`
	DataResum *Resum    `json:"DataResum,omitempty"`
	Usermails *[]mail.Target `json:"UsermailList,omitempty"`
}

//App the struct for the app
type App struct {
	Router    *mux.Router
	Database  string
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

type passT struct {
	Subject string `json:"Subject"`
	Message string `json:"Message"`
	Target MailType `json:"Target"`
}


// Credentials struct
type credentials struct {
	Password string `json:"password"`
	Usermail string `json:"usermail"`
}

type tokenSt struct {
	Token string `json:"token,omitempty"`
}

//Warn References to a Warn
/*
type Warn struct {
	ID int `json:"id"`
	Date time.Time `json:"Date"`
	Data []byte `json:"Data"`
	SendedMail Mail
}
*/

//WarnList References to a list of Warns
type WarnList struct {
	Warns *[]Warn    `json:"Warns,omitempty"`
	Date  *time.Time `json:"Date,omitempty"`
}

//MailType References to a Mail
type MailType struct {
	ID   int    `json:"ID,omitempty"`
	Name string `json:"Name,omitempty"`
	Mail string `json:"Mail,omitempty"`
}

//Resum references to a resum of data
type Resum struct {
	Name string    `json:"Name,omitempty"`
	Date time.Time `json:"Date,omitempty"`
	Data []Data    `json:"DATA,omitempty"`
}

//Data references to a Data
type Data struct {
	Name     string   `json:"Name,omitempty"`
	LineName string   `json:"LineName,omitempty"`
	ColName  string   `json:"ColName,omitempty"`
	Lines    []string `json:"Lines,omitempty"`
	Cols     []string `json:"Cols,omitempty"`
}

//ListaClientes define json struct
type ListaClientes struct {
	Nome string `json:"nome,omitempty"`
}

//DataEmailUsuario define json struct
type DataEmailUsuario struct{
	Data  time.Time   `json:"data"`
}

type Warn struct{
	ID    int         `json:"id"`
	Data  time.Time   `json:"data"`
	Lista entity.FuncPublico `json:"pessoa"`
	EmailsEnviados entity.EmailEnviado `json:"emails"`
}