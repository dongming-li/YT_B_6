#--------- Hem setup options

defaultHem =
	baseAppRoute: "/NYMB"
	tests:
		runner: "browser"
	proxy:
		"/api":
			"host": "localhost"
			"port": 8080

#--------- main configuration setup

config =

	# main hem configuration
	hem: defaultHem

	# appliation configuration

	application:
		defaults: "spine"
		root: "./"
		js:
			libs: [
				'lib/jquery.js',
				'lib/jade_runtime.js',
				'lib/tether.js',
				'lib/bootstrap.js',
				'lib/d3.min.js',
				'lib/sockjs.min.js'
			]
			modules: [
				"spine",
				"spine/lib/ajax",
				"spine/lib/route",
				"spine/lib/manager"
			]
		test:
			after: "require('lib/setup')"

#--------- export the configuration map for hem

module.exports.config = config

#--------- customize hem

module.exports.customize = (hem) ->
	# provide hook to customize the hem instance,
	# called after config is parsed/processed.
	return
