package transaction

import (
	"encoding/json"
	"sync"
	"time"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/models"

	log "github.com/sirupsen/logrus"
)

const (
	getNewTransactions = "SELECT * FROM Transactions WHERE Status='new'"
	createTransaction  = "INSERT INTO Transactions(FromID, ToID, CurrencyID, Amount, Created) VALUES(?, ?, ?, ?, ?)"
	updateTransaction  = "UPDATE Transactions SET Completed=?, Status=? WHERE ID=?"
	readBalance        = "SELECT * FROM Balances WHERE AccountID=? AND CurrencyID=?"
	updateBalance      = "UPDATE Balances SET Amount=? WHERE AccountID=? AND CurrencyID=?"
)

type queue struct {
	transactions map[int]*models.Transaction
	mux          sync.RWMutex
}

func newQueue() *queue {
	q := &queue{transactions: make(map[int]*models.Transaction)}

	rows, err := database.Query(getNewTransactions)
	if err != nil {
		log.Warningf("queue: %s", err.Error())
	}

	for rows.Next() {
		transaction := &models.Transaction{}
		if err := models.ScanTransaction(transaction, rows); err != nil {
			log.Warningf("queue: %s", err.Error())
		}
		q.transactions[transaction.ID] = transaction
	}

	return q
}

// add takes a transaction, inserts it into the database,
// appends it to the transaction queue, and notifies both the from and to of the transaction
func (q *queue) add(t *models.Transaction) {
	// insert transaction into database
	result, err := database.Exec(createTransaction, t.FromID, t.ToID, t.CurrencyID, t.Amount, t.Created)
	if err != nil {
		log.Warningf("queue: %s", err.Error())
		return
	}

	// grab generated id and set the transaction's id
	id, err := result.LastInsertId()
	if err != nil {
		log.Warningf("queue: %s", err.Error())
		return
	}
	t.ID = int(id)

	// lock the queue and append the transaction
	q.mux.Lock()
	if _, ok := q.transactions[t.ID]; !ok {
		q.transactions[t.ID] = t
	}
	q.mux.Unlock()

	err = sendTransactionUpdate(t)
	if err != nil {
		log.Warningf("queue: %s", err.Error())
	}
}

func (q *queue) approve(t *models.Transaction) {
	q.mux.RLock()
	if t, ok := q.transactions[t.ID]; ok {
		t.Status = models.Approved
		_, err := database.Exec(updateTransaction, time.Now(), string(t.Status), t.ID)
		if err != nil {
			log.Warningf("queue: %s", err.Error())
			return
		}
		fromBalance := getBalance(t.FromID, t.CurrencyID)
		_, err = database.Exec(updateBalance, fromBalance.Amount-t.Amount, t.FromID, t.CurrencyID)
		if err != nil {
			log.Warningf("queue: %s", err.Error())
			return
		}
		toBalance := getBalance(t.ToID, t.CurrencyID)
		_, err = database.Exec(updateBalance, toBalance.Amount+t.Amount, t.ToID, t.CurrencyID)
		if err != nil {
			log.Warningf("queue: %s", err.Error())
			return
		}
		delete(q.transactions, t.ID)
		sendTransactionUpdate(t)
	}
	q.mux.RUnlock()
}

func (q *queue) deny(t *models.Transaction) {
	q.mux.RLock()
	if t, ok := q.transactions[t.ID]; ok {
		t.Status = models.Denied
		_, err := database.Exec(updateTransaction, time.Now(), string(t.Status), t.ID)
		if err != nil {
			log.Warningf("queue: %s", err.Error())
			return
		}
		delete(q.transactions, t.ID)
		sendTransactionUpdate(t)
	}
	q.mux.RUnlock()
}

func getBalance(accountID, currencyID int) *models.Balance {
	balance := &models.Balance{}
	rows, err := database.Query(readBalance, accountID, currencyID)
	defer rows.Close()
	if err != nil {
		log.Warningf("queue: %s", err.Error())
		return nil
	}

	if ok := rows.Next(); !ok {
		log.Warningf("queue: could not get balance for account(%d)", accountID)
		return nil
	}
	if err = rows.Scan(&balance.ID, &balance.AccountID, &balance.CurrencyID, &balance.Amount); err != nil {
		log.Warningf("queue: %s", err.Error())
		return nil
	}

	return balance
}

func (q *queue) sendNewTransactions(session *Session) error {
	q.mux.RLock()
	defer q.mux.RUnlock()
	for _, t := range q.transactions {
		if t.FromID == session.accountID || t.ToID == session.accountID {
			// make a datum with the transaction
			data := &WSDatum{
				Type:        WSTransactionType,
				Transaction: t,
			}
			// marshal and send the datum
			msg, err := json.Marshal(data)
			if err != nil {
				return err
			}
			if err = session.sock.Send(string(msg)); err != nil {
				return err
			}
		}
	}
	return nil
}

func sendTransactionUpdate(t *models.Transaction) error {
	// make a datum with the transaction
	data := &WSDatum{
		Type:        WSTransactionType,
		Transaction: t,
	}
	// marshal the datum
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	sessionsMux.RLock()
	defer sessionsMux.RUnlock()
	// send the transaction to its 'from' account
	if sesh, ok := sessions[t.FromID]; ok {
		err = sesh.sock.Send(string(msg))
	}
	// send the transaction to its 'to' account
	if sesh, ok := sessions[t.ToID]; ok {
		err = sesh.sock.Send(string(msg))
	}
	return err
}
