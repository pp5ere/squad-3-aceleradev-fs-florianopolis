package notificacao

import (
	"encoding/json"
	entity "squad-3-aceleradev-fs-florianopolis/entities"
	"squad-3-aceleradev-fs-florianopolis/entities/logs"
	db "squad-3-aceleradev-fs-florianopolis/interfaces/db"

	utils "squad-3-aceleradev-fs-florianopolis/utils"

	"strconv"
	"time"
)

//GetLastID get the Latest
func GetLastID() int {
	dbi, err := db.Init()
	defer dbi.Database.Close()
	if err != nil {
		logs.Errorf("GetLastNotificacaoID - DB Connection", err.Error())
		return 0
	}
	type Rsp struct {
		ID int
	}
	var Response Rsp
	squery := "SELECT id FROM NOTIFICACAO order by id desc limit 1;"
	seleciona, err := dbi.Database.Query(squery)
	defer seleciona.Close()
	for seleciona.Next() {
		if err != nil {
			logs.Errorf("GetLastNotificacaoID", err.Error())
			return 0
		}
		seleciona.Scan(&Response.ID)
	}
	return (Response.ID)
}

//GetNextID get the next notificacao id
func GetNextID() int {
	return (GetLastID() + 1)
}

//InsertNotificacao insere uma notificaçao
func InsertNotificacao(request entity.Mailrequest) error {
	dbi, erro := db.Init()
	defer dbi.Database.Close()
	var ntf entity.NotificacaoLista
	ntf.ClientesDoBanco = request.TopNames
	ntf.TopFuncionariosPublicos = request.Names
	response, erro := json.Marshal(ntf)
	if erro != nil {
		logs.Errorf("InsertNotificacao", erro.Error())
	}
	result, erro := dbi.Database.Query(`INSERT INTO NOTIFICACAO (data, lista) VALUES(?, ?)`, time.Now().Format("2006-01-02 15:04:05"), response)
	defer result.Close()

	if erro != nil {
		logs.Errorf("InsertNotificacao", erro.Error())
	}
	return erro
}

//Delete Notificacao by ID
func Delete(id int) error {
	dbi, erro := db.Init()
	defer dbi.Database.Close()
	squery := `DELETE FROM NOTIFICACAO WHERE id = ` + strconv.Itoa(id)
	result, erro := dbi.Database.Query(squery)
	defer result.Close()

	return erro
}

//GetByID notificacao by ID
func GetByID(id int) (*entity.Notificacao, error) {
	var note entity.Notificacao
	dbi, erro := db.Init()
	if erro != nil {
		logs.Errorf("get(EMAILENVIADO)", erro.Error())
	}
	defer dbi.Database.Close()
	squery := `select * from NOTIFICACAO where id = "` + strconv.Itoa(id) + `" limit 1;`
	seleciona, err := dbi.Database.Query(squery)
	defer seleciona.Close()
	if err != nil {
		return nil, err
	}
	var new []byte //changed
	for seleciona.Next() {
		seleciona.Scan(&note.ID, &note.Data, &new) //changed
		json.Unmarshal(new, &note.Lista)           //changed
	}
	return &note, nil
}

//Get notificacao by data
func Get(pData time.Time) (*entity.Notificacao, error) {
	var Data time.Time
	var note entity.Notificacao
	dbi, erro := db.Init()
	if erro != nil {
		logs.Errorf("get(EMAILENVIADO)", erro.Error())
	}
	defer dbi.Database.Close()
	if pData != Data {
		formatTime := pData.Format("2006-01-02 15:04:05")
		squery := `select * from NOTIFICACAO where data = "` + formatTime + `" limit 1;`
		seleciona, err := dbi.Database.Query(squery)
		defer seleciona.Close()
		if err != nil {
			return nil, err
		}
		for seleciona.Next() {
			seleciona.Scan(&note.ID, &note.Data, &note.Lista)
		}

	} else {
		seleciona, err := dbi.Database.Query(`select * from NOTIFICACAO order by data desc limit 1;`)
		defer seleciona.Close()
		if err != nil {
			return nil, err
		}

		for seleciona.Next() {
			var d string
			var new []byte                     //changed
			seleciona.Scan(&note.ID, &d, &new) //changed
			note.Data = utils.ConvertDateTimeSQL(d)
			json.Unmarshal(new, &note.Lista) //changed
		}
	}

	return &note, nil

}
