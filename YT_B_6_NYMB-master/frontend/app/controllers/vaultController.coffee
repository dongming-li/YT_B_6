Spine   = require('spine')
$       = Spine.$

User = require('models/user')
Vault = require('models/vault')
Permissions = require('models/permissions')

class VaultController extends Spine.Controller
    className: 'vaultController'

    events: {
        'click .add-vault-btn': 'addVault'
        'click .delete-vault' : 'deleteVault'
        'click .options-btn': 'editVault'
    }

    # constructs teh vaultController
    constructor: ->
        super
        @routes
            "/vaultController/:id": (params) ->
                console.log(params.id)
                showPermissions(params.id)
        
        @active @change

    # renders the vaultController page
    render: =>
        @html require('views/vaultController/index')
        
        # @showVaults(showPermissions)
    
    # shows the vault permissions
    showPermissions: (vaultID)=>
        
        
        vault = Vault.find(vaultID)
        console.log(vault)
        # vaultData = { data: vault } 

        # modal = $('#global-modal')
        # modal.addClass("modal fade bd-example-modal-lg")
        # modal.html require('views/vaults/permission')(vaultData)
        # form = modal.find('form')

        # tbody = modal.find('tbody')
        # tbody.html('')
        
        # $.ajax
        #     url: 'http://localhost:9294/api/vault/'+ vaultID+'/permission'
        #     type: 'GET'
        #     # data: JSON.stringify(vault)
        #     dataType: 'json'
        #     success: (data) ->
        #         console.log(data)
        #         length = data.length
        #         i = 0
        #         console.log(length)
        #         while i < length
        #             rowData = { data: data[i] }
        #             console.log(data[i])
        #             tbody.append require('views/vaults/permissionRow')(rowData)
        #             i++    
                
        #         # return
        #     error: (data) ->
        #         # console.log(data)
        #         return
        # modal.modal()

    # deletes a vault permission
    deletePermission: (e) =>
        e.preventDefault()
        tr = $(e.currentTarget).parents('td').parents('tr')
        vaultID = tr[0].children[0].textContent
        vault = Vault.find(vaultID)
        console.log(vault)
        vaultData = { data: vault }   
        modal = $('#global-modal')
        modal.html require('views/vaults/deletePermission')(vaultData)
        form = modal.find('form')
        form.submit (e) =>
            e.preventDefault()
            Permissions.destroy(PermissionID)
            modal.modal('hide')
            @change()
        modal.modal()
    
    # activates when the vault page is changed
    change: (params) =>
        User.fetch() 
        Vault.fetch()
        Permissions.fetch()
        # @render()


module.exports = VaultController
