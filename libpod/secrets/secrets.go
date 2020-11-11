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

func secretCreate(name string, data []byte, driver string) (string, error) {
	exist, err := checkSecretExist(name)
	if err != nil {
		return "", err
	}
	if exist {
		fmt.Println("secret name in use")
		return "", errors.New("secret name in use")
	}

	secr := new(Secret)
	secr.Name = name
	secr.ID = stringid.GenerateNonCryptoID()
	secr.Driver = driver

	err = storeSecret(*secr)
	if err != nil {
		return "", err
	}
	err = storeSecretData(secr.ID, data)
	if err != nil {
		return "", err
	}
	return secr.ID, nil
}

func secretRm(nameOrID string) (string, error) {
	exist, err := checkSecretExist(nameOrID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("secret is not made")
	}
	id, err := deleteSecret(nameOrID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func secretInspect(nameOrID string) (*InspectSecret, error) {
	secret, err := getSecret(nameOrID)
	if err != nil {
		return nil, err
	}
	inspect := new(InspectSecret)
	inspect.ID = secret.ID
	inspect.Spec.Name = secret.Name
	return inspect, nil
}

func secretLs() ([]Secret, error) {
	secrets, err := getAllSecrets()
	if err != nil {
		return nil, err
	}
	var ls []Secret
	for _, v := range secrets {
		ls = append(ls, v)

	}
	return ls, nil

}

//secret db operations////////////////////////////////////////////////////////////////////////////////////////////////////////

func getDB() (*DB, error) {
	_, err := os.Stat(secretsFile) // For read access.
	if err != nil {
		if os.IsNotExist(err) {
			secretDBInit()
		} else {
			return nil, err
		}
	}
	file, err := os.Open(secretsFile)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("herre")
		return nil, err
	}
	db := new(DB)
	if err = json.Unmarshal(byteValue, db); err != nil {
		return nil, err
	}
	return db, nil

}

func getNameAndID(nameOrID string) (name string, id string, err error) {
	db, err := getDB()
	if err != nil {
		return "", "", err
	}
	idMap := db.NameToID
	for secretName, id := range idMap {
		if nameOrID == secretName || nameOrID == id {
			return secretName, id, nil
		}
	}
	return "", "", errors.Wrapf(define.ErrNoSuchSecret, "No secret with name or id %s", nameOrID)
}

func checkSecretExist(nameOrID string) (bool, error) {
	_, _, err := getNameAndID(nameOrID)
	if err != nil {
		if errors.Cause(err) == define.ErrNoSuchSecret {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getAllSecrets() (map[string]Secret, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	return db.Secrets, nil
}

func getSecret(nameOrID string) (*Secret, error) {
	_, id, err := getNameAndID(nameOrID)
	if err != nil {
		return nil, err
	}
	allSecrets, err := getAllSecrets()
	if err != nil {
		return nil, err
	}
	secret := allSecrets[id]
	return &secret, nil
}

func storeSecret(entry Secret) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	db.Secrets[entry.ID] = entry
	db.NameToID[entry.Name] = entry.ID

	marshalled, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(secretsFile, marshalled, 0644)
	if err != nil {
		return err
	}

	return nil

}

func deleteSecret(nameOrID string) (string, error) {
	name, id, err := getNameAndID(nameOrID)
	if err != nil {
		return "", err
	}
	db, err := getDB()
	if err != nil {
		return "", err
	}
	delete(db.Secrets, id)
	delete(db.NameToID, name)
	marshalled, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(secretsFile, marshalled, 0644)
	if err != nil {
		return "", err
	}
	return id, nil

}

func secretDBInit() error {
	f, err := os.Create(secretsFile)
	if err != nil {
		return err
	}
	f.WriteString("{\"secrets\":{},\"nametoid\":{}}")
	if err != nil {
		return err
	}
	f.Close()
	if err != nil {
		return err
	}
	return nil
}

//functions that will be interfaces///////////////////////////////////////////////////////////////////////////////////////

type secretsData struct {
	secretData map[string][]byte
}

func getAllSecretData() (*secretsData, error) {
	_, err := os.Stat(secretsDataFile) // For read access.
	if err != nil {
		if os.IsNotExist(err) {
			createSecretDataFile()
		} else {
			return nil, err
		}
	}

	file, err := os.Open(secretsDataFile)

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
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

}

func createSecretDataFile() error {
	f, err := os.Create(secretsDataFile)
	if err != nil {
		return nil
	}
	_, err = f.WriteString("{}")
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func storeSecretData(id string, data []byte) error {
	allSecrets, err := getAllSecretData()
	if err != nil {
		return err
	}
	allSecrets.secretData[id] = data
	marshalled, err := json.MarshalIndent(allSecrets.secretData, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./secretstore.json", marshalled, 0644)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

}

func deleteSecretData(id string) error {
	allSecrets, err := getAllSecretData()
	if err != nil {
		return err
	}
	delete(allSecrets.secretData, id)
	marshalled, err := json.MarshalIndent(allSecrets.secretData, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./secretstore.json", marshalled, 0644)
	if err != nil {
		return err
	}
	return nil
}

type Driver interface {
	getAllSecretData() (*secretsData, error)
	lookupSecretData(ID string) ([]byte, error)
	createSecretDataFile() error
	storeSecretData(id string, data []byte)
	deleteSecretData(id string)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {

	// id1, _ := secretCreate("secret1", []byte("secret1"), "file")
	// secretCreate("secret2", []byte("secret2"), "file")
	// secretRm(id1)

	fmt.Println(secretLs())

}
