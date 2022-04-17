Spine   = require('spine')
$       = Spine.$

class Permissions extends Spine.Model
  @configure "Permissions", "id", "userID", "vaultID", "requestTransaction", "approveTransaction", "addUser","removeUser", "addFunds", "removeFunds", "userName" 

  @extend Spine.Model.Ajax
  @url: "../api/permissions"

module.exports = Permissions
