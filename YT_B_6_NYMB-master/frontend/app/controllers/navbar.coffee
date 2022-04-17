Spine   = require('spine')
$       = Spine.$

Auth = require('controllers/auth')
User = require('models/user')

class Navbar extends Spine.Controller

  elements: {
    '#navbar': 'navbar'
    '#user-display': 'userDisplay'
    '#sign-in-btn': 'signInBtn'
    '#sign-up-btn': 'signUpBtn'
    'button.navbar-toggler': 'toggleBtn'
  }

  events: {
    'click button#sign-in-btn': 'signIn'
    'click button#sign-out-btn': 'signOut'
    'click button#sign-up-btn': 'signUp'
  }

  # constructs the navbar
  constructor: (args...) ->
    super(args)
    User.bind('authorized', @authorized)
    User.bind('unauthorized', @change)
    @change()

  # activates when the navbar is changed
  change: =>
    @render(User.appUser)
    if User.appUser
      @signInBtn.hide()
      @signUpBtn.hide()
    else @userDisplay.hide()

  # renders the navbar
  render: (appuser) =>
    name = appuser?.username
    if appuser and appuser.firstname and appuser.lastname
      name = "#{appuser.firstname} #{appuser.lastname}"
    @html require('views/navbar/index')({ name: name, email: appuser?.email })

  # updates the navbar
  update: (activeController) =>
    if activeController is 'landing' then @toggleBtn.hide()
    else @toggleBtn.show()

  # determines authorization
  authorized: =>
    @change()
    if @modal
      @modal.modal('hide')
      @navigate('/home')

  # sign up process for a user
  signUp: (e) =>
    e.preventDefault()
    @modal = $('#global-modal')
    @modal.html require('views/navbar/signUp')
    form = @modal.find('form')
    formMessage = form.find('span.text-danger')
    form.submit (e) =>
      e.preventDefault()
      user = new User(form.serializeObject())
      if user.password != user.password_conf
        formMessage.html('Passwords did not match.')
      else
        Auth.signUp(user)
        .done =>
          @modal.modal('hide')
          @signInBtn
            .tooltip({ trigger: 'manual' }).tooltip('show')
        .fail (data) ->
          formMessage.html(data.responseJSON.message)
    @modal.modal()

  # sign in process for a user
  signIn: (e) ->
    e.preventDefault()
    @signInBtn.tooltip('hide')
    @modal = $('#global-modal')
    @modal.html require('views/navbar/signIn')
    form = @modal.find('form')
    form.submit (e) =>
      e.preventDefault()
      data = form.serializeObject()
      Auth.signIn(data).fail =>
        form.find('span.modal-message').html('Username or password was incorrect')
    @modal.modal()

  # sign out process for a user
  signOut: (e) ->
    e.preventDefault()
    Auth.signOut()
    @navigate('/landing')

module.exports = Navbar
