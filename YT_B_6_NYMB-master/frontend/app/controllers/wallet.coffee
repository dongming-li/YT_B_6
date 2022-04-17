Spine   = require('spine')
$       = Spine.$

Balance = require('models/balance')
User = require('models/user')
AccountNumber = -1

class Wallet extends Spine.Controller
  className: 'wallet'

  events: {
    'click .option-btn': 'editRow'
    'click .update-btn': 'updateBalance'
    'click .cancel-btn': 'cancelBalance'
    'click .delete-btn': 'deleteBalance'
    'click .create-btn': 'createBalance'
  }

  # constructs the wallets page
  constructor: ->
    super
    Balance.bind('refresh', (models) => @updateTableRows(models))
    @active @change

  # renders the wallets page
  render: =>
    @html require('views/wallet/index')

  # activates when the wallets page is changed
  change: (params) =>
    AccountNumber = User.appUser.account
    Balance.fetch(@render())

  # updates the balance table
  updateTableRows: (balances) =>
    content = @$('#content')
    content.html('')
    for balance in balances
        content.append require('views/wallet/row')({ data: balance })

  # edits a balance table row
  editRow: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    id = row.data('id')
    balance = Balance.find(id)
    console.log(balance)
    row.html require('views/wallet/editRow')({ data: balance })

  # updates a balance
  updateBalance: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    id = row.data('id')
    balance = Balance.find(id)
    inputs = row.find('input')
    attrs = {}
    for input in inputs
      #parses as a float because the only update that matters is amount
      attrs[input.name] = parseFloat(input.value)
    balance.updateAttributes(attrs)
    balance.update()
    Balance.fetch()
    row.replaceWith require('views/wallet/row')({ data: balance })

  # cancels updates to a balance
  cancelBalance: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    id = row.data('id')
    balance = Balance.find(id)
    row.replaceWith require('views/wallet/row')({ data: balance })

  # deletes a balance
  deleteBalance: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    id = row.data('id')
    Balance.find(id).destroy()
    row.remove()

  # creates a balance
  createBalance: (e) =>
    e.preventDefault()
    modal = $('#global-modal')
    accountIdSet = false
    if AccountNumber > 0
      modal.html require('views/wallet/addFunds')
      accountIDSet = true
    else
      modal.html require('views/wallet/addFundsToAccount')

    form = modal.find('form')
    form.submit (e) =>
      e.preventDefault()
      data = form.serializeObject()
      #Note that id: -1 needs to stay or Spine.js sets it to a tempID that's a string
      #id: -1 can't be in the database, so b.save() won't ever call update()
      b = null
      if accountIDSet
        b = new Balance({
          id: -1
          accountId: AccountNumber
          currencyId: parseInt(data['currency'])
          amount: parseFloat(data['amount'])
        })
      else
        b = new Balance({
          id: -1
          accountId: parseInt(data['accountId'])
          currencyId: parseInt(data['currency'])
          amount: parseFloat(data['amount'])
        })

      b.save()
      Balance.fetch()
      console.log(b)
      modal.modal('hide')
      @updateTableRows(Balance.all())
    modal.modal()

module.exports = Wallet
