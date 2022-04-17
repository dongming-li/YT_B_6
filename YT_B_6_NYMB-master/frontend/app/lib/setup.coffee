require('spine')
require('spine/lib/manager')
require('spine/lib/route')
require('spine/lib/ajax')

$('nav').on('show.bs.collapse', (e) ->
  e.preventDefault()
  stack = $('div#stack')
  $(e.target)
    .addClass('show')
    .animate({
      left: "+=600"
    }, 200)
  if $(window).width() > 575
    stack.fadeOut(200, ->
      stack
        .removeClass('col-12')
        .addClass('col-sm-9 offset-sm-3')
        .fadeIn(100)
    )
)

$('nav').on('hide.bs.collapse', (e) ->
  e.preventDefault()
  sidebar = $(e.target)
  stack = $('div#stack')
  sidebar
    .animate({
      left: "-=600"
    }, 200, ->
      sidebar.removeClass('show')
    )
  if $(window).width() > 575
    stack.fadeOut(200, ->
      stack
        .removeClass('col-sm-9 offset-sm-3')
        .addClass('col-12')
        .fadeIn(100)
    )
)

jQuery.fn.extend({
  serializeObject: ->
    data = {}
    for input in this.serializeArray()
      data[input.name] = input.value
    return data
})
