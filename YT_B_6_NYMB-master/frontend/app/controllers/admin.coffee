Spine   = require('spine')
$       = Spine.$

User = require('models/user')
Vault = require('models/vault')
Transaction = require('models/transaction')
Balance = require('models/balance')

class Admin extends Spine.Controller
  className: 'home'

  elements: {
    '#admin-user-link': 'userLink'
    '#admin-vault-link': 'vaultLink'
    '#admin-transaction-link': 'transactionLink'
    '#content': 'content'
    '#search-form': 'searchForm'
  }

  events: {
    'click #admin-user-link': 'showUserTab'
    'click #admin-vault-link': 'showVaultTab'
    'click #admin-transaction-link': 'showTransactionTab'
    'click #form-cancel': 'clearForm'
    'submit #search-form': 'search'
    'click .options-btn': 'editRow'
    'click .update-btn': 'updateModel'
    'click .cancel-btn': 'cancelModel'
    'click .delete-btn': 'deleteModel'
    'click .add-funds-btn': 'addFunds'
  }

  # constructs the admin page
  constructor: ->
    super
    User.bind('refresh', (models) => @updateTable(models))
    Vault.bind('refresh', (models) => @updateTable(models))
    @active @change

  # activates when the admin page is changed
  change: (params) =>
    User.fetch()
    Vault.fetch()
    if promise = Transaction.fetch()
      promise.done (data) =>
        @transactions = data
        @updateTable(data)
    @render()

  # renders the admin page
  render: =>
    @html require('views/admin/index')
    @content.html require('views/admin/user')
    @currentModel = 'User'

  # shows the user tab
  showUserTab: (e) =>
    e.preventDefault()
    @setActiveTab('User')
    @updateTable(@getModel('User').all())

  # shows the vault tab
  showVaultTab: (e) =>
    e.preventDefault()
    @setActiveTab('Vault')
    @updateTable(@getModel('Vault').all())

  # shows the transaction tab
  showTransactionTab: (e) ->
    e.preventDefault()
    @setActiveTab('Transaction')
    @updateTable(@transactions)

  # sets the active tab
  setActiveTab: (model) ->
    for k, v of @elements
      @[v].removeClass('active')
    @[model.toLowerCase() + 'Link'].addClass('active')
    @content.html require('views/admin/' + model.toLowerCase())
    @currentModel = model

  # shearches record on the backend
  search: (e) =>
    e.preventDefault()
    field = @$('#field').val()
    filter = @$('#filter').val()
    filterNum = @filterFloat(filter)
    unless isNaN(filterNum) then filter = filterNum
    if filter
      @updateTable(@getModel().findAllByAttribute(field, filter))
    else @updateTable(@getModel().all())

  # updates the displayed table
  updateTable: (records) =>
    if @currentModel is records[0]?.constructor.name or
    not records[0] or @currentModel is 'Transaction'
      tbody = @$('tbody')
      tbody.html('')
      for row in records
        rowData = { model: @currentModel, data: row }
        tbody.append require('views/admin/row')(rowData)

  # clears an input form
  clearForm: (e) =>
    e.preventDefault()
    for input in @$('#search-form input')
      $(input).val('')

  # edits a row in the table
  editRow: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    id = row.data('id')
    row.html require('views/admin/editRow')({ data: @getModel().find(id) })

  # updates the backend model in question
  updateModel: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    id = row.data('id')
    model = @getModel().find(id)
    inputs = row.find('input, select')
    attrs = {}
    for input in inputs
      attrs[input.name] = input.value
    if attrs['role'] then attrs['role'] = parseInt(attrs['role'])
    model.updateAttributes(attrs)
    if @currentModel is 'User' then model.update()
    row.replaceWith require('views/admin/row')({ model: @currentModel, data: model })

  # cancels updates to the backend model
  cancelModel: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    id = row.data('id')
    model = @getModel().find(id)
    row.replaceWith require('views/admin/row')({ model: @currentModel, data: model })

  # deletes the backend model
  deleteModel: (e) =>
    e.preventDefault()
    row = $(e.currentTarget).parents('tr')
    id = row.data('id')
    @getModel().find(id).destroy()
    row.remove()

  # adds funds to a account
  addFunds: (e) =>
    e.preventDefault()
    id = $(e.currentTarget).parents('tr').data('id')
    modal = $('#global-modal')
    modal.html require('views/admin/addFunds')
    form = modal.find('form')
    form.submit (e) =>
      e.preventDefault()
      data = form.serializeObject()
      data['userId'] = id
      data['currencyId'] = @filterFloat(data['currencyId'])
      data['amount'] = @filterFloat(data['amount'])
      Balance.addFunds(data)
      modal.modal('hide')
    modal.modal()

  # gets the current model
  getModel: (model) =>
    unless model then model = @currentModel
    switch model
      when 'User' then return User
      when 'Vault' then return Vault
      when 'Transaction' then return Transaction
    false

  # sanitizes strings
  filterFloat: (value) ->
    if /^(\-|\+)?([0-9]+(\.[0-9]+)?|Infinity)$/.test(value)
      return Number(value)
    NaN

module.exports = Admin
