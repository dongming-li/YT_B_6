Spine   = require('spine')
$       = Spine.$

User = require('models/user')

class Balance extends Spine.Model
  @configure "Balance",
  "id", "accountId", "currencyId", "amount"

  @extend Spine.Model.Ajax
  @url: "../api/balance"

  # addFunds adds funds via ajax request
  @addFunds: (data) ->
    $.ajax({
      url: '../api/account/addFunds'
      type: 'POST'
      contentType: 'application/json'
      data: JSON.stringify(data)
    })

module.exports = Balance
