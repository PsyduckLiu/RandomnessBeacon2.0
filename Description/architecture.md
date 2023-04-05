# Folder Backbone

This folder contains all the basical components used in Randomness Beacon, such as class group, timed commitment, grpc communication, and signature. The folder itself does not contribute to the execution of RB and can be considered a technical prototype.

# Folder BulletinBoard

This project implements a bulletin board based on a simple HTTP server. All files on the server can be accessed and downloaded by users(generators, collectors, and other third-party users). It guarantees that all users can get the same value at nearly the same time.

```[]
BulletinBoard
│   main.go 
│
└───config
│   initConfig.go
│   readConfig.go
│   remoteConfig.go
│   writeConfig.go
│   
└───crypto/binaryquadraticform
│   binaryquadratic.go
│   guide.go
│
└───download
│
└───RBC
│   newLeaderpb
│   newOutputpb
│   outputpb
│   pkMsgpb
│   proposalpb
│
└───result
│
└───signature
│   generateSig.go
│   verifySig.go
│
└───util
│   util.go
```

# Folder Collector

This project implements a collector who receives TC sets from some generators and runs a protocol to make an agreement on a certian random number.

```[]
Collector
│   main.go 
│
└───config
│   readConfig.go
│   remoteConfig.go
│   writeConfig.go
│   
└───crypto/binaryquadraticform
│   binaryquadratic.go
│   guide.go
│
└───download
│
└───RBC
│   newLeaderpb
│   newOutputpb
│   outputpb
│   pkMsgpb
│   proposalpb
│   submitpb
│   tcMsgpb
│
└───result
│
└───signature
│   generateSig.go
│   verifySig.go
│
└───util
│   util.go
│
└───watch
│   watchOutput.go
```

# Folder Generator

This project implements a generator who generates TC sets and sends them to a corresponding collector.

```[]
Collector
│   main.go 
│
└───config
│   readConfig.go
│   remoteConfig.go
│   
└───crypto/binaryquadraticform
│   binaryquadratic.go
│   guide.go
│
└───download
│
└───result
│
└───tcMsgpb
│
└───util
│   util.go
│
└───watch
│   watchOutput.go
```
