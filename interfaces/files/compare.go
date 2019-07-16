package main

import (
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	entity "squad-3-aceleradev-fs-florianopolis/entities"
	"squad-3-aceleradev-fs-florianopolis/entities/logs"
	"squad-3-aceleradev-fs-florianopolis/interfaces/crud/funcpublico"
	db "squad-3-aceleradev-fs-florianopolis/interfaces/db"
	"strconv"
	"strings"
)

func openFileCSV() error {
	workPath, err := getFileName()
	if err != nil {
		logs.Errorf("openFileCSV", err.Error())
		return err
	}
	fileName := getLastFiles(workPath.Directory, 1, ".txt")
	fullPath := workPath.Directory + fileName[0]
	csvfile, err := os.Open(fullPath)
	if err != nil {
		logs.Errorf("openFileCSV", err.Error())
		return err
	}
	defer csvfile.Close()
	logs.Info("openFileCSV", "Opening file: "+fullPath)
	reader := csv.NewReader(csvfile)
	reader.Comma = ';'
	logs.Info("openFileCSV", "Reading file...")
	rawdata, err := reader.ReadAll()
	if err != nil {
		logs.Errorf("openFileCSV", err.Error())
		return err
	}
	dbi, _ := db.Init() //CONEXÃO COM DB INICIADA AQUI
	defer dbi.Database.Close()
	err = insertIntoPessoa(rawdata, dbi)
	if err != nil {
		logs.Errorf("openFileCSV", err.Error())
		return err
	}

	return err
}

//Check if person already exists in DB (by name)
func checkPersonInDB(name string, dbi *db.MySQLDatabase) (bool, int) {
	Pessoa := new(entity.FuncPublico)
	alreadyInDB := false
	Pessoa, _ = funcpublico.GetByName(name, dbi)
	if Pessoa.Nome == name {
		alreadyInDB = true
	}
	return alreadyInDB, Pessoa.ID
}

func insertIntoPessoa(rawdata [][]string, dbi *db.MySQLDatabase) error {
	if len(rawdata) > 0 {
		Pessoa := new(entity.FuncPublico)
		for i, column := range rawdata {
			if i > 0 {
				Remuneracaodomes, err := strconv.ParseFloat(changeComma(column[3]), 64)
				if err != nil {
					logs.Errorf("insertIntoPessoa", err.Error())
					return err
				}
				Redutorsalarial, err := strconv.ParseFloat(changeComma(column[8]), 64)
				if err != nil {
					logs.Errorf("insertIntoPessoa", err.Error())
					return err
				}
				Totalliquido, err := strconv.ParseFloat(changeComma(column[9]), 64)
				if err != nil {
					logs.Errorf("insertIntoPessoa", err.Error())
					return err
				}
				if Totalliquido > 20000 {
					/*type FuncPublico struct {
						ID               int
						Nome             string
						Cargo            string
						Orgao            string
						Remuneracaodomes float64
						RedutorSalarial  float64
						TotalLiquido     float64
						Updated          bool
						ClientedoBanco   bool
					}*/
					Pessoa.Nome = column[0]
					Pessoa.Cargo = column[1]
					Pessoa.Orgao = column[2]
					Pessoa.Remuneracaodomes = Remuneracaodomes
					Pessoa.RedutorSalarial = Redutorsalarial
					Pessoa.TotalLiquido = Totalliquido

					//Funcao para procurar pelo nome (Pessoa.Nome)
					alreadyExists, existingID := checkPersonInDB(Pessoa.Nome, dbi)

					/*REVER: PROCURA PELO NOME COMPLETO DA PESSOA
					CASO JÁ EXISTA NO BANCO, FAZ O UPDATE DOS DADOS NA TABELA FUNCPUBLICO
					CASO NÃO EXISTA, INSERE NOVA TUPLA NA TABELA FUNCPUBLICO
					OUTRA QUESTÃO É SE É CLIENTE DO BANCO (true ou false na variável ClientedoBanco,
					a depender do CSV importado com os nomes)*/

					//Verifica se nome já é cliente
					//Caso cliente, update dos dados e update = true
					if alreadyExists { //verifica se é cliente do banco na api e pega o ID da pessoa
						Pessoa.ID = existingID
						Pessoa.Updated = true
						Pessoa.ClientedoBanco = true //pega o valor do json do cliente
						//jsonData, err := json.Marshal(Pessoa)
						if err != nil {
							logs.Errorf("insertIntoPessoa", err.Error())
							return err
						}
						//log.Println(string(jsonData))
						//atualiza no banco
						erro := funcpublico.Update(Pessoa, dbi)
						if erro != nil {
							logs.Errorf("insertIntoPessoa", erro.Error())
						}

						//Caso não cliente, insere os dados e seta update = true
					} else {
						Pessoa.ClientedoBanco = false
						Pessoa.TotalLiquido = 0
						//Insere no banco

						erro := funcpublico.Insert(Pessoa, dbi)
						if erro != nil {
							logs.Errorf("insertIntoPessoa", erro.Error())
						}
						Pessoa.Updated = true
					}

					/*	PONTOS:
						- Caso salário líquido > 20k,
						- Verifica se nome já é cliente
						- Caso cliente, update dos dados e update = true
						- Caso não cliente, insere os dados e seta update = true
						- Caso update = false, set totalliquido = 0
						- Ao final, setar novamente update = false em todos os clientes da tabela*/
				}

			}

		}
	} else {
		err := errors.New("the csv file is empty")
		logs.Errorf("insertIntoPessoa", err.Error())
		return err
	}

	return nil
}

func changeComma(pFloatNumber string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(pFloatNumber, ",", ";"), ".", ":"), ";", "."), ":", ",")
}

func getHashFromFile(filePath string) (string, error) {
	var stringHashSHA1 string
	file, err := os.Open(filePath)
	if err != nil {
		logs.Errorf("getHashFromFile", err.Error())
		return stringHashSHA1, err
	}
	defer file.Close()
	hashSHA1 := sha1.New()
	_, err = io.Copy(hashSHA1, file)
	if err != nil {
		logs.Errorf("getHashFromFile", err.Error())
		return stringHashSHA1, err
	}
	stringHashSHA1 = hex.EncodeToString(hashSHA1.Sum(nil))
	return stringHashSHA1, err
}

func getLastFiles(pathDir string, countFiles int, extension string) []string {
	var filesName []string
	files, err := ioutil.ReadDir(pathDir)
	if err != nil {
		logs.Errorf("getHashFromFile", err.Error())
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() > files[j].ModTime().Unix()
	})
	if countFiles > 0 && extension != "" {
		for i, file := range files {
			if file.Mode().IsRegular() {
				if filepath.Ext(file.Name()) == extension {
					if len(filesName) < countFiles {
						filesName = append(filesName, files[i].Name())
					} else {
						break
					}
				}
			}
		}
	}

	return filesName
}
