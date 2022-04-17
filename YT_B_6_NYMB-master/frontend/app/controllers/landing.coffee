Spine   = require('spine')
$       = Spine.$

class Landing extends Spine.Controller
  className: 'landing'

  # constructs the landing page
  constructor: (args...) ->
    super(args)
    @active @change

  # renders the landing page
  render: =>
    @html require('views/landing/index')

  # activates when the landing page is changed
  change: (params) =>
    @render()

module.exports = Landing
