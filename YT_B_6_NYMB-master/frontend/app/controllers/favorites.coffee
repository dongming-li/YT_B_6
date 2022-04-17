Spine   = require('spine')
$       = Spine.$

User = require('models/user')
Favorite = require('models/favorite')

class Favorites extends Spine.Controller
    className: 'favorites'

    elements: {
        '#user-favorites-link': 'userLink'
        '#vault-favorites-link': 'vaultLink'
    }

    events: {
        'click #user-favorites-link': 'showUserTab'
        'click #vault-favorites-link': 'showVaultTab'
        'click .update-btn': 'updateFavorite'
        'click .cancel-btn': 'cancelUpdate'
        'click .delete-btn': 'deleteFavorite'
        'click .create-btn': 'addFavorite'
        'click .options-btn': 'editRow'
    }

    # constructs the favorites page
    constructor: ->
        super
        Favorite.bind('refresh', (models) => @updateTable(models))
        @active @change

    # renders the favorites page
    render: =>
        @html require('views/favorites/index')
        @currentTab = 'User'

    # activates when the favorites page is changed
    change: (params) =>
        Favorite.fetch()
        @render()

    # shows the favorite users
    showUserTab: (e) =>
        e.preventDefault()
        @setActiveTab('User')

    # shows the favorite vaults
    showVaultTab: (e) =>
        e.preventDefault()
        @setActiveTab('Vault')

    # sets the active tab
    setActiveTab: (tab) ->
        for k, v of @elements
            @[v].removeClass('active')
        @[tab.toLowerCase() + 'Link'].addClass('active')
        @currentTab = tab
        @updateTable(Favorite.all())

    # redraws the favorites table
    updateTable: (favorites) =>
        tbody = @$('#rows')
        tbody.html('')
        for row in favorites
            if row.username.Valid && @currentTab == 'User' || row.vaultname.Valid && @currentTab == 'Vault'
                tbody.append require('views/favorites/row')( { tab: @currentTab, data: row } )
    
    # adds a favorite
    addFavorite: (e) ->
        e.preventDefault()
        modal = $('#global-modal')
        modal.html require('views/favorites/addFavorite')
        form = modal.find('form')
        form.submit (e) =>
            e.preventDefault()
            data = form.serializeObject()

            fav = new Favorite({
                id: -1
                userId: User.appUser.id
                accountId: parseInt(data['accountId'])
                favoriteName: data['name']
            })

            fav.create({idx: 0})
            Favorite.fetch()
            modal.modal('hide')
            @updateTable(Favorite.all())
        modal.modal()

    # deletes a favorite
    deleteFavorite: (e) ->
        e.preventDefault()
        row = $(e.currentTarget).parents('tr')
        id = row.data('id')
        Favorite.find(id).destroy()
        row.remove()

    # readies a favorites row for editing
    editRow: (e) ->
        e.preventDefault()
        row = $(e.currentTarget).parents('tr')
        id = row.data('id')
        row.replaceWith require('views/favorites/editRow')( { tab: @currentTab, data: Favorite.find(id) } )

    # cancels an update to a favorite
    cancelUpdate: (e) ->
        e.preventDefault()
        row = $(e.currentTarget).parents('tr')
        id = row.data('id')
        fav = Favorite.find(id)
        row.replaceWith require('views/favorites/row')({ tab: @currentTab, data: Favorite.find(id)})
    
    # updates a favorite
    updateFavorite: (e) ->
        e.preventDefault()
        row = $(e.currentTarget).parents('tr')
        id = row.data('id')
        fav = Favorite.find(id)
        input = row.find('input')[0]
        console.log(input.value)
        fav.updateAttributes( { favoriteName: input.value } )
        fav.update()
        Favorite.fetch()
        row.replaceWith require('views/favorites/row')({tab: @currentTab, data: fav})
module.exports = Favorites
    