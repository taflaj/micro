// access: models/models.go

package models

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/taflaj/micro/modules/messaging"
)

// func unauthorized(w http.ResponseWriter, address string) {
// 	err := fmt.Errorf("Host at %v is not authorized to change access levels", address)
// 	messaging.FailNotAuthorized(w, err)
// }

func (db *DB) changeAccess(w http.ResponseWriter, msg *messaging.Message, reset bool) error {
	data := *msg.Decompose()
	// logger.GetLogger().Printf("DEBUG data=%#v", data)
	var err error
	// source, _ := messaging.IPtoInt(msg.IP)
	source := msg.IP
	level, _ := db.GetAccess(source, messaging.Access)
	// logger.GetLogger().Printf("DEBUG level=%#v", level)
	if level.CanWrite {
		service := msg.Command[2]
		if service == "" {
			err = fmt.Errorf("Invalid service '%v'", service)
			messaging.Fail(w, http.StatusBadRequest, err)
		} else {
			ip, err := messaging.IPtoInt(data["ip"])
			if err != nil {
				messaging.Fail(w, http.StatusBadRequest, err)
			} else if reset {
				db.resetAccess(ip, service)
				messaging.Ok(w)
			} else {
				addressID, err := db.addAddress(ip, data["owner"], data["remarks"])
				// logger.GetLogger().Printf("DEBUG addressID=%v, err=%v", addressID, err)
				if err == nil {
					serviceID, err := db.addService(service)
					if err == nil {
						access := data["access"]
						if len(access) != 2 {
							messaging.FailBadRequest(w, errors.New("Missing argument"))
						} else {
							canRead := access[0] == 'r'
							canWrite := access[0] == 'w' || access[1] == 'w'
							_, err := db.setAccess(addressID, serviceID, canRead, canWrite)
							if err == nil {
								// fmt.Fprintf(w, "address: %v, service: %v, id: %v", addressID, serviceID, id)
								messaging.Ok(w)
							} else {
								messaging.FailInternal(w, err)
							}
						}
					} else {
						messaging.FailInternal(w, err)
					}
				} else {
					messaging.FailInternal(w, err)
				}
			}
		}
	} else {
		err = fmt.Errorf("Host at %v is not authorized to change access levels", messaging.IPtoString(msg.IP))
		messaging.Fail(w, http.StatusUnauthorized, err)
		// unauthorized(w, msg.IP)
	}
	// logger.GetLogger().Printf("DEBUG err=%v", err)
	return err
}

// SetAccess sets an access level for someone
func (db *DB) SetAccess(w http.ResponseWriter, msg *messaging.Message) error {
	return db.changeAccess(w, msg, false)
}

// ResetAccess restores access level to the default
func (db *DB) ResetAccess(w http.ResponseWriter, msg *messaging.Message) error {
	return db.changeAccess(w, msg, true)
}
