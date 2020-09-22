#### Documentation
- https://confluince.fujie.agency/project/swift

#### Tests
- https://github.com/jaminas/swift.tests
Codeception

#### CI/CD
- USE swarm cluster deploy with project "Swift"
- Or download docker container "swift.node" - https://github.com/jaminas/fujie-docker

#### Metrics
- https://graphana.fujie.agency/

#### Dependencies
- go.mod


## Manual
#### INSTALATION TO DEVELOPMENT
- Copy files to `GOPATH/src` directory

#### BUILD
- Run command on Linux: `go build main.go`
- MACOS: `env GOOS=linux GOARCH=amd64 GOARM=7 go build main.go`

#### RUN
- Run command `main` 

#### DEPLOY
- Build it
- Copy executable file `main` to server and rename it
- Copy `conf` directory to server and configure it
- Make systemd (or other) script to start executable file

