Spine   = require('spine')
$       = Spine.$

User = require('models/user')

class Sidebar extends Spine.Controller

  elements: {
    '#landing-link': 'landingLink'
    '#home-link': 'homeLink'
    '#profile-link': 'profileLink'
    '#marketplace-link': 'marketplaceLink'
    '#transactions-link': 'transactionsLink'
    '#favorites-link': 'favoritesLink'
    '#vaults-link': 'vaultsLink'
    '#wallet-link': 'walletLink'
    '#admin-link': 'adminLink'
  }

  # constrctor for the sidebar
  constructor: () ->
    super
    @change()
    @el.hide()

  # activates when the sidebar is changed
  change: =>
    @render()

  # renders the sidebar page
  render: =>
    @html require('views/sidebar/index')(User.appUser)

  # updates the sidebar by setting the link which corresponds to the activeController to active
  # this is called any time the user switches controllers
  update: (activeController) =>
    if activeController is 'landing' then @el.hide()
    else @el.show()
    for k, v of @elements
      @[v].removeClass('active')
    @[activeController + 'Link'].addClass('active')

module.exports = Sidebar
