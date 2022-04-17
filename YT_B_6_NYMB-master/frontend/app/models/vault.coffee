Spine   = require('spine')
$       = Spine.$

class Vault extends Spine.Model
  @configure "Vault", "id", "name", "owner", "ownerID", "accountID", "userID"

  @extend Spine.Model.Ajax
  @url: "../api/vault"

module.exports = Vault
