require('lib/setup')

Spine = require('spine')

Auth = require('controllers/auth')
Navbar = require('controllers/navbar')
Sidebar = require('controllers/sidebar')
Main = require('controllers/main')

User = require('models/user')

class App extends Spine.Controller
  constructor: (agrs...) ->
    super(args)

    @navbar = new Navbar({ el: 'nav#navbar' })
    @sidebar = new Sidebar({ el: 'nav#sidebar' })
    @main = new Main({ el: 'div#stack' })
    @append @navbar, @sidebar, @main

    User.one('initial_auth', (params) =>
      @sidebar.change()
      Spine.Route.setup()
      unless params.authorized then @navigate('/landing')
    )

    @routes {
      '/landing': (params) ->
        @doRoute('landing')
      '/home': (params) ->
        @doRoute('home')
      '/profile': (params) =>
        @doRoute('profile')
      '/marketplace': (params) =>
        @doRoute('marketplace')
      '/transactions': (params) =>
        @doRoute('transactions')
      '/wallet': (params) =>
        @doRoute('wallet')
      '/vaults': (params) =>
        @doRoute('vaults')
      '/favorites': (params) =>
        @doRoute('favorites')
      '/admin': (params) =>
        if User.appUser?.role is 1 then @doRoute('admin')
    }

    @auth = new Auth()

  doRoute: (controller) =>
    if controller isnt 'landing' and User.appUser is undefined
      @navigate('/landing')
      return
    @log('navigating to ' + controller)
    @sidebar.update(controller)
    @navbar.update(controller)
    if controller is 'landing'
      @main.el.removeClass('col-sm-9 offset-sm-3').addClass('col-12')
    else
      @main.el.removeClass('col-12').addClass('col-sm-9 offset-sm-3')
    @main[controller.charAt(0).toUpperCase() + controller.slice(1)].active()

module.exports = App
