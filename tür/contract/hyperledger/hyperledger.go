package hyperledger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	time "time"

	"github.com/ardanlabs/kit/mapstructure"
)

type (
	HyperLedger struct {
		URL     string
		DoorID  string
		Message string
	}
	Response struct {
		TxId      string `jpath:"TxId"`
		Value     Offer  `jpath:"Value"`
		TimeStamp string `jpath:"Timestamp"`
		IsDelete  string `jpath:"IsDelete"`
	}

	QueryResult struct {
		QueryResult []Response `jpath:"queryResult"`
	}

	Offer struct {
		OfferID    int64     `jpath:"Value.offerId"`
		Free       bool      `jpath:"Value.free"`
		Price      float64   `jpath:"Value.price"`
		CheckIn    time.Time `jpath:"Value.checkIn"`
		CheckOut   time.Time `jpath:"Value.checkOut"`
		ObjectName string    `jpath:"Value.objectName"`
		OwnerName  string    `jpath:"Value.ownerName"`
		TenantPk   string    `jpath:"Value.tenantPk"`
		LandlordPk string    `jpath:"Value.landlordPk"`
	}
)

//renterID is not needed here (in contrast to the equivalent method for ethereum), since decrypting the message with the PK is proof enough for the requester's authenticity
//update: the whole method might not be needed, since getting the key from the offer that is valid in the very moment also proofs
//that the user who authenticated himself/herself successfully is allowed in right now.
func (h *HyperLedger) isAllowedAt(requestPointofTime time.Time, startTime time.Time, endTime time.Time) (allowed bool) {
	if requestPointofTime.Before(endTime) && requestPointofTime.After(startTime) {
		return true
	}
	return false
}

//By decrypting the payload with the public key that is stored in booking that is currently valid not only the user's
//authenticity is validated but also the fact that he/she is allowed in in the very moment.
func (h *HyperLedger) isAllowedIn(payload string, publicKey string) bool {
	return true
	// Decrypt function missing --> Have to talk to Max

	// key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	// ciphertext, _ := hex.DecodeString(payload)
	// nonce := []byte("64a9433eae7c")

	// block, err := aes.NewCipher(key)
	// if err != nil {
	// 	panic(err.Error())
	// 	return false
	// }

	// aesgcm, err := cipher.NewGCM(block)
	// if err != nil {
	// 	panic(err.Error())
	// 	return false
	// }

	// _, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	// if err != nil {
	// 	panic(err.Error())
	// 	return false
	// }
	// return true
}

func (h *HyperLedger) getPubKeyOfValidUser(resp []Response) string {
	currentTime := time.Now()
	for _, element := range resp {
		if element.Value.CheckIn.Before(currentTime) && element.Value.CheckOut.After(currentTime) {
			return element.Value.TenantPk
		}
	}
	return ""
}

type Payload struct {
	Args []string `json:"args"`
	Fcn  string   `json:"fcn"`
}

func (h *HyperLedger) getHistoryForOffer() []Response {
	data := Payload{Args: []string{h.DoorID}, Fcn: "getHistoryForOffer"}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
	}
	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", h.URL, body)
	if err != nil {
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var m map[string]interface{}
	json.Unmarshal(responseData, &m)
	var queryResult QueryResult
	err = mapstructure.DecodePath(m, &queryResult)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println(queryResult) // prints: map[foo:1]

	defer resp.Body.Close()
	return queryResult.QueryResult

}

func (h *HyperLedger) getBlockData() *Offer {
	fmt.Println("Starting the application...")
	response, err := http.Get(h.URL)
	var responseData []byte
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil
	} else {
		responseData, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
			return nil
		}

	}
	var offer Offer
	err = json.Unmarshal(responseData, &offer)
	if err != nil {
		fmt.Print("Unmarshalling did not work")
		log.Fatal(err)
		return nil
	}

	return &offer

}
