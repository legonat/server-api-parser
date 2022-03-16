# server-api-parser

Demo project for API parsing. Project includes server with API and parser. Start server as one instance and use another to parse data 

Application supports flags:
* initVmDb - initialize database with data from text files
* reinitVmDb - initializes database again even if it exists
* printAsyncMap -d - print map with concurrent write to one map
* printSyncMap -d - print map with cycled write to one map
* startServer - run server, domain anr port can be changed in app.conf
