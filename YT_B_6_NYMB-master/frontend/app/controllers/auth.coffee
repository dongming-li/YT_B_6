Spine   = require('spine')
$       = Spine.$

User = require('models/user')

class Auth extends Spine.Controller

    constructor: (args...) ->
        super(args)
        @token = localStorage.getItem('nymb-token')
        @tokenExpiration = localStorage.getItem('nymb-token-exp')
        if @token and (+new Date(@tokenExpiration) > +Date.now())
            $.ajaxSetup({ headers: { Authorization: 'Bearer ' + @token } })
            Auth.signInWithStoredToken(@token, @tokenExpiration)
        else User.trigger('initial_auth', { authorized: false })

    # signUp is called by the navbar controller after a client "submits" the sign up form
    @signUp: (user) ->
        delete user.cid
        delete user.password_conf
        $.ajax({
            url: '../api/user'
            type: 'POST'
            contentType: 'application/json'
            data: JSON.stringify(user)
        })

    # signIn is called by the navbar controller after a client "submits" the sign in form
    @signIn: (data) =>
        $.ajax({
            url: '../api/login'
            type: 'POST'
            contentType: 'application/json'
            data: JSON.stringify(data)
        }).done (response) =>
            @storeToken(response.token, response.expire)
            @fetchAppUser()
        .fail => @signOut()

    # signInWithStoredToken should only be called by this constructor
    @signInWithStoredToken: (@token, @tokenExpiration) =>
        $.ajax({ url: '../api/auth/token' })
        .done (response) =>
            @storeToken(response.token, response.expire)
            @fetchAppUser().complete -> User.trigger('initial_auth', { authorized: true })
        .fail =>
            @signOut()
            User.trigger('initial_auth', { authorized: false })

    # fetchAppUser should only be called by this controller
    @fetchAppUser: =>
        if !@token or (+new Date(@tokenExpiration) < +Date.now())
            @signOut()
            return false
        $.ajax({ url: '../api/auth/user' })
        .done (response) ->
            User.appUser = response
            User.trigger('authorized')
        .fail => @signOut()

    # signOut is called by the navbar controller when a user clicks "sign out"
    @signOut: =>
        localStorage.removeItem('nymb-token')
        localStorage.removeItem('nymb-token-exp')
        delete @token
        delete @tokenExpiration
        delete User.appUser
        User.trigger('unauthorized')

    # storeToken should only be called by this controller
    @storeToken: (@token, @expire) ->
        localStorage.setItem('nymb-token', @token)
        localStorage.setItem('nymb-token-exp', @expire)
        $.ajaxSetup({ headers: { Authorization: 'Bearer ' + @token } })

module.exports = Auth
