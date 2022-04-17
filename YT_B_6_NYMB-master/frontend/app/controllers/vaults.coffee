Spine   = require('spine')
$       = Spine.$

User = require('models/user')
Vault = require('models/vault')
Permissions = require('models/permissions')
Balance = require('models/balance')
Currency = require('models/currency')
Account = require('models/account')
maxVaultID = 0

class Vaults extends Spine.Controller
    className: 'vaults'

    # key value pairings of events to functions
    events: {
        'click .add-user-btn': 'addUser'
        'click .showFunds-btn': 'showFunds'
        'click .add-vault-btn': 'addVault'
        'click .delete-vault': 'deleteVault'
        'click .moreInfo': 'showMore'
        'click .delete-user': 'deleteUser'
        'click .backToAllVaults': 'backToAllVaults'
    }

    
    constructor: (args...) ->
        super(args)
        Vault.bind('refresh', (models) => @showVaults(models))
        @active @change

    # creates initial view of the vaults page
    render: =>
        @html require('views/vaults/index')
        @showVaults(Vault.all())
    # goes back to the initial view of the vaults page but also hides the button again.
    backToAllVaults: (e) =>
        e.preventDefault()
        $(".backToAllVaults").css("visibility", "hidden")
        @render()

    # @param {event}, the click of the delete user buttoj
    # removes a user from a vault
    deleteUser: (e) =>
        e.preventDefault()
        Permissions.fetch()
        tr = $(e.currentTarget).parents('td').parents('tr')
        permissionID = tr[0].children[0].textContent
        perm = Permissions.find(permissionID)
        Permissions.destroy(permissionID)
        @rendor()
    # @param {event}, the click of the edit permissions button
    # shows permission information for a vault on a per user basis
    showMore: (e) =>
        e.preventDefault()

        tr = $(e.currentTarget).parents('td').parents('tr')
        vaultID = tr[0].children[0].textContent
        vault = Vault.find(vaultID)

        block = @$('#content.card-block')
        block.html(' ')
        block.append require('views/vaults/permission')({ data: vault })
        $(".backToAllVaults").css("visibility", "visible")
        tbody = block.find('tbody')
        tbody.html(' ')

        $.ajax({
            url: 'http://localhost:9294/api/vault/' + vault.id + '/allpermissions'
            type: 'GET'
            dataType: 'json'
            success: (data) ->
                mPerms = null
                for perm in data
                    if perm.userID == User.appUser.id
                        mPerms = perm
                for perm in data
                    if perm.userID != 1
                        data = { data: perm, myPerms: mPerms }
                        tbody.append require('views/vaults/permissionRow')(data)
                Permissions.fetch()
        })

    # @param {event}, the click of the show funds buttons
    # shows the funds of a chosen vault
    showFunds: (e) ->
        e.preventDefault()
        tr = $(e.currentTarget).parents('td').parents('tr')
        vaultID = tr[0].children[0].textContent
        vault = Vault.find(vaultID)
        vaultData = { data: vault }
        modal = $('#global-modal')
        modal.html require('views/vaults/showFunds')(vaultData)
        modal.modal()

        tbody = modal.find('tbody')
        tbody.html(' ')
        userBalances = Balance.all()
        for balance in userBalances
            if balance.accountId == vault.account
                rowData = { data: balance }
                curId = balance.currencyId
                cur = Currency.find(curId)
                data = { name: cur.name, amount: balance.amount }
                tbody.append require('views/vaults/showFundsRow')(data)
    # @param {event}, the click of the delete vault button
    # deletes the vault and removes all permission for the vault. Also refreshes the page
    deleteVault: (e) =>
        e.preventDefault()
        tr = $(e.currentTarget).parents('td').parents('tr')
        vaultID = tr[0].children[0].textContent
        vault = Vault.find(vaultID)

        vaultData = { data: vault }
        modal = $('#global-modal')
        modal.html require('views/vaults/deleteVault')(vaultData)
        form = modal.find('form')
        form.submit (e) =>
            e.preventDefault()
            Vault.destroy(vaultID)
            modal.modal('hide')
            @change()
        modal.modal()

    # a function that refreshes the page
    change: (params) =>
        User.fetch()
        Vault.fetch()
        Permissions.fetch()
        Balance.fetch()
        Currency.fetch()
        Account.fetch()
        @render()
    # @param {event}, the click of the add user to a vault button
    # adds a user to a vault, gives them basic permissions by default
    addUser: (e) =>
        e.preventDefault()
        tr = $(e.currentTarget).parents('td').parents('tr')
        id = tr[0].children[0].textContent
        vault = Vault.find(id)
        modal = $('#global-modal')
        modal.html require('views/vaults/addUser')({ data: vault })
        form = modal.find('form')

        form.submit (e) =>
            e.preventDefault()
            data = form.serializeObject()
            name = User.find(data["userid"]).username
            perm = {
                id: -1
                userID: parseInt(data["userid"])
                vaultID: parseInt(id)
                requestTransaction: true
                approveTransaction: false
                addUser: false
                removeUser: false
                addFunds: true
                removeFunds: false
                userName: name
            }
            permission = new Permissions(perm)
            permission.bind('save', () ->
                alert(name + " has been given permission to your vault"))
            permission.save()
            modal.modal('hide')
            @change()
        modal.modal()

    # @param {event}, the click of the add user to add vault button
    # creates a new vault with the user as the owner.
    addVault: (e) =>
        e.preventDefault()
        modal = $('#global-modal')
        modal.html require('views/vaults/addVault')
        form = modal.find('form')
        ownerName = User.find(User.appUser.id).username

        form.submit (e) =>
            e.preventDefault()
            data = form.serializeObject()
            vaultData = {
                id: -1
                ownerID: User.appUser.id
                owner: ownerName
                name: data["name"]
            }
            vault = new Vault(vaultData)
            vault.save()
            modal.modal('hide')
            @change()
        modal.modal()

    # @param tbody, the table body object
    # @param vault, the vault that is to be displayed in the row
    # adds a row to the table that will only allow certain elements be seen for the editting of based on permissions
    getPermAndShowRow: (tbody, vault) ->
        $.ajax({
            url: 'http://localhost:9294/api/vault/' + vault.id + '/permission'
            type: 'GET'
            dataType: 'json'
            success: (data) ->
                rowData = { vInfo: vault, permission: data }
                tbody.append require('views/vaults/row')(rowData)
        })
    # @param vaults, all of the vaults that the user has permissions to
    # shows the vaults for a user
    showVaults: (vaults) =>
        tbody = @$('tbody')
        tbody.html(' ')
        for vault in vaults
            if vault.userID == User.appUser.id
                @getPermAndShowRow(tbody, vault)


module.exports = Vaults
