class Datum

class Token extends Datum

  constructor: (@token) ->
    unless typeof @token is "string" then throw new Error("Datum.Token takes a token string")
    @type = 'token'

class Message extends Datum

  # constructor creates a Messge
  constructor: (to, frm, message) ->
    unless typeof to is "string"
      throw new Error("Datum.Message 'to' takes a string")
    unless typeof frm is "string"
      throw new Error("Datum.Message 'from' takes a string")
    unless typeof message is "string"
      throw new Error("Datum.Message 'message' takes a string")
    @type = 'message'
    @message = {
      'to': to
      'from': frm
      'message': message
    }

class Transaction extends Datum

  constructor: (@transaction) ->
    @type = 'transaction'

Datum.Token = Token
Datum.Message = Message
Datum.Transaction = Transaction
module.exports = Datum
