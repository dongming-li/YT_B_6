Spine   = require('spine')
$       = Spine.$

Currency = require('models/currency')
Balance = require('models/balance')
Transaction = require('models/transaction')
Datum = require('models/datum')
User = require('models/user')
AccountNumber = -1

class Home extends Spine.Controller
  className: 'home'

  # construct's the user's home page
  constructor: (args...) ->
    super(args)
    Balance.bind('refresh', @render)
    Currency.bind('priceUpdate', @render)
    @active @change

  # renders the user's home page
  render: =>
    @html require('views/home/index')
    @updateBalanceTableRows()
    @updateTransactionTableRows()
    @updatePricePredictions()

  # activates when the user's home page is changed
  change: (params) =>
    @[0].stack.Transactions.wsStart(localStorage.getItem('nymb-token'))
    AccountNumber = User.appUser.account
    Balance.fetch()
    Transaction.fetch()
    @render()

  # updates the rows of the balance table
  updateBalanceTableRows: () =>
    BalanceArea = @$('#BalanceArea')
    BalanceArea.html('')
    balances = Balance.all()

    for balance in balances
      currencyName = "error: currency not available"
      USDamnt = null
      switch balance.currencyId
        when 1
          currencyName = "Bitcoin"
          USDamnt = "$" + (balance.amount * Currency.BTC).toFixed(2)
        when 2
          currencyName = "Etherium"
          USDamnt = "$" + (balance.amount * Currency.ETH).toFixed(2)
        when 3
          currencyName = "Litecoin"
          USDamnt = "$" + (balance.amount * Currency.LTC).toFixed(2)
        when 4 then currencyName = "US Dollar"
        else null

      data = {
        currencyId: currencyName
        amount: balance.amount
        amountInUSD: USDamnt
      }
      BalanceArea.append require('views/home/balanceRow')({ data })

  # updates the rows of the transaction table
  updateTransactionTableRows: () =>
    TransactionArea = @$('#TransactionArea')
    TransactionArea.html('')
    transactions = Transaction.all()
    balances = Balance.all()

    for transaction in transactions
      currencyName = "error: currency not available"
      USDamnt = null
      switch transaction.currencyId
        when 1
          currencyName = "Bitcoin"
          USDamnt = "$" + (transaction.amount * Currency.BTC).toFixed(2)
        when 2
          currencyName = "Etherium"
          USDamnt = "$" + (transaction.amount * Currency.ETH).toFixed(2)
        when 3
          currencyName = "Litecoin"
          USDamnt = "$" + (transaction.amount * Currency.LTC).toFixed(2)
        when 4 then currencyName = "US Dollar"
        else null

      direction = "Outbound"
      if transaction.toId is AccountNumber
        direction = "Inbound"

      data = {
        inboundOrOutbound: direction
        currency: currencyName
        amount: 900
        valueInUSD: USDamnt
        created: transaction.created
      }
      TransactionArea.append require('views/home/transactionRow')({ data })

  # updates the price predictions
  updatePricePredictions: () =>
    CurrencyPriceArea = @$('#CurrencyPriceArea')

    data = {
      btcPrice: Currency.BTC.toFixed(2)
      btcPrediction: Currency.BTCPrediction
      ethPrice: Currency.ETH.toFixed(2)
      ethPrediction: Currency.ETHPrediction
      ltcPrice: Currency.LTC.toFixed(2)
      ltcPrediction: Currency.LTCPrediction
    }
    CurrencyPriceArea.append require('views/home/currencyPriceCards')({ data })

module.exports = Home
