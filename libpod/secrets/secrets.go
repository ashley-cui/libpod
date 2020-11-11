package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/containers/podman/v2/libpod/define"
	"github.com/containers/storage/pkg/stringid"
	"github.com/pkg/errors"
)

// temporary becasue i dont remember what dan said LOL
var secretsFile = "./secrets.json"         //this stores secretID and metadata
var secretsDataFile = "./secretstore.json" //this stores the actual secret data,, file driver i guess

type DB struct { //this is basically database design
	Secrets  map[string]Secret `json:"secrets"`  // secrets bucket , literally {"secrets": [secret 1, secret 2, etc]}
	NameToID map[string]string `json:"nametoid"` // bucket of map of names to id's, literally {"NameToID": {"name":"ID" "name":"ID"}}
}

type Secret struct {
	Name     string            `json:"name"`
	ID       string            `json:"id"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Driver   string            `json:"driver"`
}

type InspectSecret struct {
	ID        string
	CreatedAt string //(?)
	UpdatedAt string
	Spec      struct {
		Name string
	}
}

func secreteCreate(name string, data []byte, driver string) (string, error) {
	exist, err := checkSecretExist(name)
	if err != nil {
		return "", err
	}
	if exist {
		return "", errors.New("secret name in use")
	}
	//for each in getsecrets()
	// check if name is there, if not,,,,,, name a ok to go
	secr := new(Secret)
	secr.Name = name
	secr.ID = stringid.GenerateNonCryptoID()
	secr.Driver = driver
	storeSecret(*secr)
	storeSecretData(secr.ID, data)
	return secr.ID, nil
}

func secretRm(nameOrID string) (string, error) {
	exist, err := checkSecretExist(nameOrIS)

	if !exist {
		return "", errors.New("secret is not made")
	}
	id, _ := deleteSecret(nameOrID)
	deleteSecretData((id))
	return id, nil
}

//secret db operations////////////////////////////////////

func getDB() (*DB, error) {
	file, err := os.Open(secretsFile) // For read access.
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file does not exist")
			return nil, nil
		}
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(file)
	db := new(DB)
	json.Unmarshal(byteValue, db)
	return db, nil

}

func lookupID(nameOrID string) (string, string, error) {
	db, _ := getDB()
	idMap := db.NameToID
	for secretName, id := range idMap {
		if nameOrID == secretName || nameOrID == id {
			return secretName, id, nil
		}
	}
	return "", "", errors.Wrapf(define.ErrNoSuchSecret, "No secret with name or id %s", nameOrID)
}

func checkSecretExist(nameOrID string) (bool, error) {
	_, _, err := lookupID(nameOrID)
	if err != nil {
		if errors.Cause(err) == define.ErrNoSuchSecret {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getAllSecrets() (map[string]Secret, error) {
	db, _ := getDB()
	return db.Secrets, nil

}
func storeSecret(entry Secret) error {
	db, _ := getDB()
	db.Secrets[entry.ID] = entry
	db.NameToID[entry.Name] = entry.ID
	marshalled, _ := json.MarshalIndent(db, "", "  ")
	_ = ioutil.WriteFile("./secretoutput.json", marshalled, 0644)
	fmt.Println(db)
	return nil

}
func deleteSecret(nameOrID string) (string, error) {
	name, id, _ := lookupID(nameOrID)
	db, _ := getDB()
	delete(db.Secrets, id)
	delete(db.NameToID, name)
	marshalled, _ := json.MarshalIndent(db, "", "  ")
	_ = ioutil.WriteFile("./secretoutput.json", marshalled, 0644)
	return id, nil

}

// func createSecretsFile(){
// 	file, _ := os.Create(secretsFile){

// 	}
// }

//interfaces///////////////////////////////////////////////////////////////////////////////////////
type secretsData struct {
	secretData map[string][]byte
}

// func write(id string, data []byte) error { //this should be an interface function
// 	//if encryption is needed later on, then the key would be generated in this function and stored with id
// 	//open file
// 	//write entry[id]=data
// 	//close file
// }

func getAllSecretData() (*secretsData, error) {
	file, err := os.Open(secretsDataFile) // For read access.
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file does not exist")
			return nil, nil
		}
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(file)
	file.Close()
	secrets := new(secretsData)
	json.Unmarshal([]byte(byteValue), &secrets.secretData)
	return secrets, nil

}

func lookupSecretData(ID string) ([]byte, error) {
	allSecrets, err := getAllSecretData()
	if err != nil {
		return nil, err
	}
	return allSecrets.secretData[ID], nil

	// if

}

func createSecretDataFile() {
	f, _ := os.Create(secretsDataFile)
	f.WriteString("{}")
	f.Close()

}

func storeSecretData(id string, data []byte) {
	allSecrets, _ := getAllSecretData()
	allSecrets.secretData[id] = data
	marshalled, _ := json.MarshalIndent(allSecrets.secretData, "", "  ")
	_ = ioutil.WriteFile("./dataoutput.json", marshalled, 0644)

}

func deleteSecretData(id string) {
	allSecrets, _ := getAllSecretData()
	delete(allSecrets.secretData, id)
	marshalled, _ := json.MarshalIndent(allSecrets.secretData, "", "  ")
	_ = ioutil.WriteFile("./dataoutput.json", marshalled, 0644)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// func Hello() {
// 	fmt.Println(storeSecret())
// }

func main() {
	// Hello()
	// fmt.Println(getDB())
	// fmt.Println(createSecret("secret2", []byte("hkjdfg"), "file"))
	// data, _ := lookupSecretData("id1")
	// fmt.Println(string(data))
	// deleteSecretData("aaa")
	deleteSecret("secret1")

}
