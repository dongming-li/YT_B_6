Spine   = require('spine')
$       = Spine.$

class Main extends Spine.Stack

  controllers: {
    Landing: require('controllers/landing')
    Home: require('controllers/home')
    Profile: require('controllers/profile')
    Marketplace: require('controllers/marketplace')
    Transactions: require('controllers/transactions')
    Favorites: require('controllers/favorites')
    Vaults: require('controllers/vaults')
    Wallet: require('controllers/wallet')
    Admin: require('controllers/admin')
  }

  default: 'Landing'

module.exports = Main
