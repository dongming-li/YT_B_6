Spine   = require('spine')
$       = Spine.$

class Account extends Spine.Model
  @configure "Account", "id", "userId", "vaultId"

  @extend Spine.Model.Ajax
  @url: "../api/account"

module.exports = Account
