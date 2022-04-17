Spine   = require('spine')
$       = Spine.$

User = require('models/user')

class Profile extends Spine.Controller
  className: 'profile'

  elements: {
    '#content': 'content'
  }

  events: {
    'click #editbtn': 'edit'
    'submit .update-form': 'update'
    'click .cancel-btn': 'cancel'
  }

  # constructs the profile page
  constructor: (args...) ->
    super(args)
    @active @change

  # renders the profile page
  render: =>
    @html require('views/profile/index')
    @content.append require('views/profile/info')(User.appUser)

  # activates when the profile page is changed
  change: (params) =>
    @render()

  # opens edit dialog for a user profile
  edit: (e) =>
    e.preventDefault()
    @content.html('')
    @content.append require('views/profile/edit')(User.appUser)

  # updates a user's profile
  update: (e) =>
    e.preventDefault()
    user = User.fromForm(e.target)
    if user.password is ''
      delete user.password
    else if user.confirm_password != user.password
      @$('span.form-message').html('your passwords do not match')
      return
    delete user.confirm_password
    user.id = User.appUser.id
    user.save({ ajax: false })
    user.update()

  # cancels update for a user's profile
  cancel: (e) =>
    e.preventDefault()
    @content.html('')
    @content.append require('views/profile/info')(User.appUser)

module.exports = Profile
