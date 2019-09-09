![Heetch](heetch.png)

## Heetch zombie test 

### Tech used
* [Docker+Docker-compose](https://www.docker.com/)
* [GOLang](https://golang.org/)
    
### Prerequisites
* [Docker](https://www.docker.com/)
* [docker-compose](https://docs.docker.com/compose/install/)
* [Make](https://ftp.gnu.org/old-gnu/Manuals/make-3.79.1/html_chapter/make_2.html)

start by cloning the repo
    
    > git clone git@github.com:heetch/jose-odg-technical-test.git

then lets
    
    > cd jose-odg-technical-test

the following shows the commands

    > make help

    usage: make <command>
    
    commands:
        install              - populate vendor/ from go.mod
        unit-test            - run unit tests
        integration-test     - run integration tests
        test                 - run unit test and integration tests
        e2e-test             - build the whole app and run and e2e test
        up                   - start the driver-location container
        down                 - stop the driver-location container
        
the lets 

    > make up

now the system should be ready to start

### Testing

to execute an end-to-end test of the whole app

    > make e2e-test

and to execute unit and/or integration tests

    > make unit-test
    > make integration-test
    > make test

### Tech notes about architecture

I started with the driver-location microservice, then the zombie-driver one and finally the gateway. For what I considered most difficult.

#### Tech notes
 
* Hexagonal architecture: The root files in the microservice is the `application layer`, the internal package is the `domain layer`, with its entities and value objects (only one aggregate here). Any other package refers to some `infrastructure layer` piece.
* Dependency injection with [google wire](https://github.com/google/wire)
* CRQS with command bus and query service. I don't use to have query bus in my applications. The command bus is [chiguirez.cromberbus](https://github.com/chiguirez/cromberbus). [Chiguirez](https://github.com/chiguirez) is a github organization that me and a few friends have, where we make a few Go snippets for developing microservices.
* To bootstrap and run the microservices I have used [chiguirez.snout](https://github.com/chiguirez/snout), a tool for bootstrapping Go applications using [spf13.viper](https://github.com/spf13/viper)
* Event dispatcher. In the `driver-location` microservice domain there are domain events and aggregate root functionality to work with. I did not add an event dispatcher for the purpose of the test, but a mock implementation is used. We are working on an [event bus](https://github.com/chiguirez/eventbus), but we don't feel it ready for production yet.

### Author

* [Jose I. Ortiz de Galisteo](https://github.com/hosseio)