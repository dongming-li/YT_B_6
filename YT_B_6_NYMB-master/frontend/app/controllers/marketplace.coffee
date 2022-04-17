Spine   = require('spine')
$       = Spine.$

class Marketplace extends Spine.Controller
  className: 'marketplace'

  elements: {
    'input#BTC-start': 'btcStartInput'
    'input#BTC-end': 'btcEndInput'
    'input#ETH-start': 'ethStartInput'
    'input#ETH-end': 'ethEndInput'
    'input#LTC-start': 'ltcStartInput'
    'input#LTC-end': 'ltcEndInput'
  }

  events: {
    'click button#BTC-submit': 'submitDates_BTC',
    'click button#ETH-submit': 'submitDates_ETH',
    'click button#LTC-submit': 'submitDates_LTC',
    'click a.granBTC': 'choseGranularityBTC',
    'click a.granETH': 'choseGranularityETH',
    'click a.granLTC': 'choseGranularityLTC',
  }

  # constructs the marketplace page
  constructor: (args...) ->
    super(args)
    @active @change

  # activates when the marketplace page is changed
  change: (params) =>
    @render()
    @createGraphs()

  # renders the marketplace page
  render: =>
    @html require('views/marketplace/index')

  # chooses the granularity of the bitcoin graph
  choseGranularityBTC: (e) ->
    e.preventDefault()
    $('.granBTC').removeClass("chosenBTC")
    $(e.target).addClass "chosenBTC"

  # chooses the granularity of the litecoin graph
  choseGranularityLTC: (e) ->
    e.preventDefault()
    $('.granLTC').removeClass("chosenLTC")
    $(e.target).addClass 'chosenLTC'

  # chooses the granularity of the etherium graph
  choseGranularityETH: (e) ->
    e.preventDefault()
    $('.granETH').removeClass("chosenETH")
    $(e.target).addClass 'chosenETH'

  # submits the dates for the bitcoin graph
  submitDates_BTC: (e) =>
    e.preventDefault()
    @submitDates('BTC')

  # submits the dates for the etherium graph
  submitDates_ETH: (e) =>
    e.preventDefault()
    @submitDates('ETH')

  # submits the dates for the litecoin graph
  submitDates_LTC: (e) =>
    e.preventDefault()
    @submitDates('LTC')

  # submits the dates for a graph in order to render that graph
  submitDates: (currency) ->
    start = @["#{currency.toLowerCase()}StartInput"].val()
    end = @["#{currency.toLowerCase()}EndInput"].val()
    granularity = @$(".chosen#{currency}").text()
    gran = 24 * 3600
    switch granularity
      when 'Every week' then gran = 7 * 24 * 3600
      when 'Every day' then gran = 24 * 3600
      when 'Every hour' then gran = 3600
      when 'every month' then gran = 7 * 24 * 3600 * 4
    @createGraphWDates("#{currency}-USD", new Date(start), new Date(end), gran)

  # creates a graph trend line
  createLine: (responses, currency) =>
    margin = { top: 20, bottom: 20, left: 40, right: 20 }
    width = @$("##{currency}-div .card-block").width() - margin.left - margin.right
    height = $(window).height() / 3 - margin.top - margin.bottom

    @$("##{currency}-div .card-block svg").remove()
    svg = d3.select("##{currency}-div .card-block").append('svg')
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)

    g = svg.append('g').attr('transform', 'translate(' + margin.left + ',' + margin.top + ')')
    parseTime = d3.timeParse('%d-%b-%y')
    x = d3.scaleTime().rangeRound([0, width])
    y = d3.scaleLinear().rangeRound([height, 0])
    line = d3.line().x((d) -> x d.date).y((d) -> y d.avg)

    data = []
    for response in responses
      data.push { date: response.date, avg: +response.avg }

    x.domain d3.extent(data, (d) -> d.date )
    y.domain d3.extent(data, (d) -> d.avg )

    g.append('g')
      .attr('transform', 'translate(0,' + height + ')')
      .call(d3.axisBottom(x)).select('.domain').remove()

    g.append('g')
      .call(d3.axisLeft(y)).append('text')
      .attr('fill', '#000').attr('transform', 'rotate(-90)')
      .attr('y', 6).attr('dy', '0.71em')
      .attr('text-anchor', 'end').text 'Price ($)'

    g.append('path').datum(data)
      .attr('fill', 'none').attr('stroke', 'steelblue')
      .attr('stroke-linejoin', 'round').attr('stroke-linecap', 'round')
      .attr('stroke-width', 1.5).attr('d', line)

  # A graph response
  class Response
    constructor: (time, low, high) ->
      d = new Date(0)
      d.setUTCSeconds time
      @date = new Date(d.toLocaleString())
      @avg = (low + high) / 2

  # creates graphs
  createGraphs: (params) =>
    today = new Date
    beg = new Date
    beg.setMonth (beg.getMonth() - 1) % 12
    beg.setHours 0, 0, 0
    beg.setMilliseconds 0
    @createGraphWDates("BTC-USD", beg, today, 3600 * 24)
    @createGraphWDates("LTC-USD", beg, today, 3600 * 24)
    @createGraphWDates("ETH-USD", beg, today, 3600 * 24)

  # creates graphs with dates
  createGraphWDates: (currency, beg, end, granularity) =>
    responses = []
    if +beg >= +end
      swal({
        icon: 'warning'
        title: 'Oops!'
        text: 'The selected start date must come before the selected end date'
        button: 'Got it'
      })
      return
    $.ajax({
        url: "https://api.gdax.com/products/#{currency}/candles?start=#{beg.toISOString()}" +
          "&end=#{end.toISOString()}&granularity=#{granularity}"
        type: 'GET'
        dataType: 'json'
        headers: null
        success: (data) =>
            for row in data
              responses.push new Response(row[0], row[1], row[2])
            @createLine(responses, currency)
        error: (rsp) =>
          if rsp.status == 429
            @createGraphWDates(currency, beg, end, granularity)
          else
            swal({
              icon: 'warning'
              title: 'Oops!'
              text: 'The selected granularity is too fine for such a large date range.' +
                'Try decreasing one or the other!'
              button: 'Got it'
            })
    })

module.exports = Marketplace
