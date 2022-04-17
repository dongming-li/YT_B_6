# Not Your Mother's Bank (NYMB)

Not Your Mother's Bank (NYMB) is a web app for trading/exchanging cryptocurrencies. The project
RESTfully services NYMB's frontend utilizing [Gin-Gonic](https://github.com/gin-gonic/gin) and MySQL.
The frontend uses [Spine.js](http://spine.github.io/), a light weight MVC SPA javascript web framework.

## Getting Started

```bash
go run main.go -dev
./scripts/integrationTests.sh
```

### Prerequisites

[Glide](https://github.com/Masterminds/glide)
[MySQL](https://www.mysql.com/)
[npm](https://www.npmjs.com/)

### Installing

```bash
git clone git@git.linux.iastate.edu:309Fall2017/YT_B_6_NYMB.git
cd YT_B_6_NYMB
glide install
cd frontend
npm install .
```

### Running

```bash
# backend
go run main.go -dev

# frontend
cd frontend
hem server
```

## Deployment

## Built With

* [Logrus](https://github.com/sirupsen/logrus) - Logging library
* [Gin-Gonic](https://github.com/gin-gonic/gin) - HTTP router framework
* [Spine.js](http://spine.github.io/) - Web framework

## Contributing

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/tags).

## Authors

* **Blake Roberts** - [brob](https://git.linux.iastate.edu/brob)
* **Gabriel Butruille** - [gabrielb](https://git.linux.iastate.edu/gabrielb)
* **Matthew Schaffer** - [maschaff](https://git.linux.iastate.edu/maschaff)
* **Leelabari Fulbel** - [lfulbel](https://git.linux.iastate.edu/lfulbel)

