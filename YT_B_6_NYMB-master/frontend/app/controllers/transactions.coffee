Spine   = require('spine')
$       = Spine.$

Currency = require('models/currency')
Account = require('models/account')
Datum = require('models/datum')
Transaction = require('models/transaction')
User = require('models/user')
Vault = require('models/vault')

class Transactions extends Spine.Controller
  className: 'transactions'

  elements: {
    'input#message': 'messageInput'
    '#btc-price': 'btcPrice'
    '#eth-price': 'ethPrice'
    '#ltc-price': 'ltcPrice'
  }

  events: {
    'click #message-form': 'sendMessage'
    'submit #transaction-form': 'sendTransaction'
    'click .approve-btn': 'approveTransaction'
    'click .deny-btn': 'denyTransaction'
    'click .cancel-btn': 'denyTransaction'
  }

  # constructor for the transactions page
  constructor: (args...) ->
    super(args)
    Account.bind('refresh', @render)
    Currency.bind('priceUpdate', @change)
    Transaction.bind('create update', @updateTransactionsTable)
    @active @change

  # renders the transactions page
  render: =>
    accountId = 0
    if a = Account.findByAttribute('userId', User.appUser?.id)
      accountId = a.id
    transactions = Transaction.findAllByAttribute('fromId', accountId)
    transactions.push.apply(transactions, Transaction.findAllByAttribute('toId', accountId))
    data = {
      accountId: accountId
      accounts: []
      currencies: [
        Currency.BTC.toFixed(2)
        Currency.ETH.toFixed(2)
        Currency.LTC.toFixed(2)
      ]
      transactions: transactions
    }
    accounts = Account.all()
    for account in accounts
      if account.userId isnt 0
        if user = User.find(account.userId)
          data.accounts.push({
            id: account.id
            name: user.username
          })
      else
        if vault = Vault.find(account.vaultId)
          data.accounts.push({
            id: account.id
            name: vault.name
          })
    @html require('views/transactions/index')(data)

  # activates when the transactions page is changed
  change: (params) =>
    @wsStart(localStorage.getItem('nymb-token'))
    unless params
      User.fetch()
      Vault.fetch()
      Account.fetch()
      @render()
    @btcPrice.html(Currency.BTC.toFixed(2))
    @ethPrice.html(Currency.ETH.toFixed(2))
    @ltcPrice.html(Currency.LTC.toFixed(2))

  # updates the transactions table
  updateTransactionsTable: =>
    account = Account.findByAttribute('userId', User.appUser.id)
    if id = account?.id
      transactions = Transaction.findAllByAttribute('fromId', id)
      transactions.push.apply(transactions, Transaction.findAllByAttribute('toId', id))
      tbody = @$('tbody')
      tbody.html('')
      for transaction in transactions
        if ['approved', 'denied'].indexOf(transaction.status) is -1
          tbody.append require('views/transactions/row')({ data: transaction, accountId: id })

  # starts the websocket connection
  wsStart: (token) ->
    if @sock?.readyState != 1
      @sock = new SockJS(window.location.origin + '/api/ws')
      @sock.onopen = @wsOpen(token)
      @sock.onmessage = @wsMessage
      @sock.onclose = @wsClose

  # opens the websocket connection
  wsOpen: (token) ->
    (e) =>
      @log('ws connection opened')
      @sock.send(JSON.stringify(new Datum.Token(token)))

  # creates a websocket message
  wsMessage: (e) =>
    msg = JSON.parse(e.data)
    if msg.type is "currencies"
      Currency.priceUpdates(msg)
    else if msg.type is "currency"
      Currency.priceUpdate(msg)
    else if msg.type is "transaction"
      if Transaction.exists(msg.transaction.id)
        Transaction.find(msg.transaction.id).updateAttributes(msg.transaction)
      else Transaction.create(msg.transaction)
    else @log(msg)

  # closes the websocket
  wsClose: (e) =>
    @log('ws connection closed')

  # sends a websocket message
  sendMessage: (e) ->
    e.preventDefault()
    if (msg = @messageInput.val()) isnt ""
      datum = new Datum.Message("someone", "this guy", msg)
      @log(JSON.stringify(datum))
      @sock.send(JSON.stringify(datum))

  # sends a transaction
  sendTransaction: (e) =>
    e.preventDefault()
    form = $(e.target)
    transaction = {
      'id': 0
      'fromId': Account.findByAttribute('userId', User.appUser.id).id
      'toId': @filterFloat(form.find("select[name='toId']").val())
      'currencyId': @filterFloat(form.find("select[name='currencyId']").val())
      'amount': @filterFloat(form.find("input[name='amount']").val())
      'created': new Date().toISOString()
      'status': 'new'
    }
    unless transaction.fromId is transaction.toId
      datum = new Datum.Transaction(transaction)
      @sock.send(JSON.stringify(datum))

  # approves a transaction
  approveTransaction: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    @sendTransactionUpdate(row, 'approve')

  # denies a transaction
  denyTransaction: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    @sendTransactionUpdate(row, 'deny')

  # sends a transaction update
  sendTransactionUpdate: (row, status) =>
    id = row.data('id')
    transaction = Transaction.find(id)
    transaction.status = status
    datum = new Datum.Transaction(transaction)
    @sock.send(JSON.stringify(datum))
    row.remove()

  # filters floats
  filterFloat: (value) ->
    if /^(\-|\+)?([0-9]+(\.[0-9]+)?|Infinity)$/.test(value)
      return Number(value)
    NaN

module.exports = Transactions
