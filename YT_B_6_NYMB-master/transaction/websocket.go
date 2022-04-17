package transaction

import (
	"encoding/json"
	"errors"
	"sync"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/models"

	ajwt "github.com/appleboy/gin-jwt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/igm/sockjs-go/sockjs"
	log "github.com/sirupsen/logrus"
)

const (
	readAccountIDFromUserEmail = "SELECT a.ID FROM Users AS u JOIN Accounts AS a ON u.ID=a.UserID WHERE u.Email=?"
)

var (
	auth        *ajwt.GinJWTMiddleware
	sessions    map[int]*Session
	sessionsMux sync.RWMutex
)

// Session defines the properties of a client's websocket connection
type Session struct {
	sock      sockjs.Session
	accountID int
	channel   chan []byte
}

// SetAuthMiddleware configures the transaction websocket authentication middleware
func SetAuthMiddleware(mw *ajwt.GinJWTMiddleware) {
	auth = mw
}

// WebSocketHandler returns a handler used to read and write to
// a websocket associated with every authenticated client connection
func WebSocketHandler(sock sockjs.Session) {
	log.Debug("websocket: new sockjs session established")
	var closedSession = make(chan struct{})

	session := &Session{
		sock:      sock,
		accountID: 0,
		channel:   make(chan []byte),
	}

	go webSocketWriter(session, closedSession)
	webSocketReader(session, closedSession)

	close(closedSession)
	log.Debug("websocket: sockjs session closed")
}

// webSocketWriter sends service updates to the client
// updates include:
//     closed session
//     currency price updates
//     transaction updates
// TODO: add a chanel for client to client messaging notifications
func webSocketWriter(session *Session, closedSession chan struct{}) {
	for {
		select {
		case <-closedSession:
			return
		case data, ok := <-session.channel:
			if ok {
				if err := session.sock.Send(string(data)); err != nil {
					log.Warningf("websocket: %s", err.Error())
				}
			}
		}
	}
}

func webSocketReader(session *Session, closedSession chan struct{}) {
	msgCount := 0
	for {
		// wait for client message
		msg, err := session.sock.Recv()
		if err != nil {
			log.Warningf("websocket: %s", err.Error())
			break
		}

		// parse the request as a WSDatum
		req := &WSDatum{}
		err = json.Unmarshal([]byte(msg), req)
		if err != nil {
			log.Warningf("websocket: %s", err.Error())
			break
		}

		// process the WSDatum. If this fails, break from the
		// for loop which closes the session
		msgCount++
		if ok := process(req, session, msgCount); !ok {
			log.Warningf("websocket: couldn't process %s for account(%d)", req.Type, session.accountID)
			break
		}
	}
}

// process takes a WSDatum from the client and processes the request and returns false if the request fails
func process(datum *WSDatum, session *Session, msgCount int) bool {
	reqOk := false
	switch datum.Type {

	case WSTokenType:
		if msgCount != 1 {
			log.Warningf("websocket: user sent token request more than once")
			break
		}
		reqOk = processToken(datum, session)
		break

	case WSTransactionType:
		if msgCount == 1 {
			log.Warningf("websocket: user didn't send token request")
			break
		}
		reqOk = processTransaction(datum, session)
		break

	case WSMessageType:
		if msgCount == 1 {
			log.Warningf("websocket: user didn't send token request")
			break
		}
		reqOk = processMessage(datum, session)
		break

	default:
		log.Warningf("websocket: didn't understand WSDatum(%d)", datum.Type)
		break
	}

	return reqOk
}

// processToken processes a token request and returns false if the token isn't valid
func processToken(datum *WSDatum, session *Session) bool {
	if datum.Token == nil {
		log.Warning("websocket: token request sent without token")
		return false
	}

	email := ""
	// parse the request's token and validate
	token, err := jwt.Parse(*datum.Token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(auth.SigningAlgorithm) != t.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		if id, ok := t.Claims.(jwt.MapClaims)["id"].(string); ok {
			email = id
		}

		if err := t.Claims.Valid(); err != nil {
			return nil, errors.New("invalid token")
		}
		return auth.Key, nil
	})

	if err != nil {
		log.Warningf("websocket: %s", err.Error())
		return false
	}

	if email == "" {
		log.Debug("websocket: couldn't get id from claims")
		return false
	}

	if !token.Valid {
		log.Warningf("websocket: user sent invalid token")
		return false
	}

	if ok := sendInitialWSCurrenciesDatum(session); !ok {
		log.Warningf("websocket: couldn't send initial currencies")
		return false
	}

	rows, err := database.Query(readAccountIDFromUserEmail, email)
	if err != nil {
		log.Warningf("websocket: %s", err.Error())
		return false
	}

	accountID := 0
	if rows.Next() {
		rows.Scan(&accountID)
	}

	session.accountID = accountID
	if err = q.sendNewTransactions(session); err != nil {
		log.Warningf("websocket: %s", err.Error())
		return false
	}

	sessionsMux.Lock()
	sessions[accountID] = session
	sessionsMux.Unlock()

	return true
}

// processTransaction processes a transaction request
func processTransaction(datum *WSDatum, session *Session) bool {
	if datum.Transaction == nil {
		log.Warningf("websocket: transaction request sent without transaction")
		return false
	}

	switch datum.Transaction.Status {
	case models.Approve:
		q.approve(datum.Transaction)
		break
	case models.Approved:
		break
	case models.Deny:
		q.deny(datum.Transaction)
		break
	case models.Denied:
		break
	case models.New:
		q.add(datum.Transaction)
		break
	default:
		log.Warningf("websocket: transaction request sent without valid transaction status")
		return false
	}

	return true
}

// processMessage processes a message request
func processMessage(datum *WSDatum, session *Session) bool {
	if datum.Message == nil {
		log.Warningf("websocket: message request sent without message")
		return false
	}

	log.Debugf("websocket: message from %s to %s -> %s",
		datum.Message.From, datum.Message.To, datum.Message.Message)

	resp, err := json.Marshal(datum)
	if err != nil {
		log.Warningf("websocket: %s", err.Error())
	}

	session.sock.Send(string(resp))
	return true
}

// sendInitialWSCurrenciesDatum takes a session and sends
// a list of all currencies to the session's client
func sendInitialWSCurrenciesDatum(session *Session) bool {
	currs := map[string]float64{}
	for name, price := range getCurrencies() {
		currs[name] = price
	}
	data := &WSDatum{
		Type:       WSCurrenciesType,
		Currencies: currs,
	}
	msg, err := json.Marshal(data)
	if err != nil {
		log.Warningf("websocket: %s", err.Error())
		return false
	}
	if err = session.sock.Send(string(msg)); err != nil {
		log.Warningf("websocket(processToken): %s", err.Error())
		return false
	}

	return true
}
