Spine   = require('spine')
$       = Spine.$

class Favorite extends Spine.Model
  @configure "Favorite",
  "id", "userId", "accountId", "favoriteName",
  "username", "vaultname"

  @extend Spine.Model.Ajax
  @url: "../api/favorite"

module.exports = Favorite