Spine   = require('spine')
$       = Spine.$

class User extends Spine.Model
  @configure "User",
    "id", "username", "email", "firstname",
    "lastname", "role", "account", "password"

  @extend Spine.Model.Ajax
  @url: "../api/user"

module.exports = User
