Spine   = require('spine')
$       = Spine.$

class Currency extends Spine.Model
  @configure "Currency", "id", "name", "unitPrice", "pricePrediction"

  @extend Spine.Model.Ajax
  @url: "api/currency"

  @BTC = 0.0
  @BTCPrediction = "No prediction available ..."
  @ETH = 0.0
  @ETHPrediction = "No prediction available ..."
  @LTC = 0.0
  @LTCPrediction = "No prediction available ..."
  @USD = 1.0

  # priceUpdates handles a price update via websocket
  @priceUpdates: (msg) ->
    for key, val of msg.currencies
      @[key.toUpperCase()] = val
    @trigger('priceUpdate', {})

  @priceUpdate: (msg) =>
    @[msg.currency.name.toUpperCase()] = msg.currency.price
    @[msg.currency.name.toUpperCase() + "Prediction"] = msg.currency.priceprediction
    @trigger('priceUpdate', {})

module.exports = Currency
