Spine   = require('spine')
$       = Spine.$

User = require('models/user')

class Transaction extends Spine.Model
  @configure "Transaction", "id", "fromId", "toId", "currencyId",
    "amount", "created", "completed", "status"

  @fetch: ->
    if User.appUser?.role is 1 then return $.ajax({ url: '../api/transaction' })
    else return false

module.exports = Transaction
