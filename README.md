## Quickstart
```
go install github.com/abs3ntdev/haunt@latest
```

or

```
git clone https://github.com/abs3ntdev/haunt
cd haunt
make build && sudo make install
```

## Commands List

### Run Command
From **project/projects** root execute:


```
haunt init
```

then

```
haunt start
```

haunt init will add your root directory if it contains a main.go file and will add any directory inside of cmd as projects. If you wish to add additional projects run haunt add [name] or edit the config file manually. By default projects are set to run go install and go run [project]

***start*** command supports the following custom parameters:

    --name="name"               -> Run by name on existing configuration
    --path="haunt/server"     -> Custom Path (if not specified takes the working directory name)
    --generate                  -> Enable go generate
    --fmt                       -> Enable go fmt
    --test                      -> Enable go test
    --vet                       -> Enable go vet
    --install                   -> Enable go install
    --build                     -> Enable go build
    --run                       -> Enable go run
    --server                    -> Enable the web server
    --open                      -> Open web ui in default browser
    --no-config                 -> Ignore an existing config / skip the creation of a new one

Some examples:

    haunt start
    haunt start --path="mypath"
    haunt start --name="haunt" --build
    haunt start --path="haunt" --run --no-config
    haunt start --install --test --fmt --no-config
    haunt start --path="/Users/username/go/src/github.com/oxequa/haunt-examples/coin/"

### Init Command
This command will generate a .haunt.yaml with sane default for your current project/projects.\
If there is a main.go in the root directory it will be added along with any directories inside the relative path `cmd`

    haunt init

### Add Command
Add a project to an existing config file or create a new one.

    haunt add [name]

### Remove Command
Remove a project by its name

    haunt remove [name]


## Config sample

    settings:
        legacy:
            force: true             // force polling watcher instead fsnotifiy
            interval: 100ms         // polling interval
        resources:                  // files names
            outputs: outputs.log
            logs: logs.log
            errors: errors.log
    server:
        status: false               // server status
        open: false                 // open browser at start
        host: localhost             // server host
        port: 5001                  // server port
    schema:
    - name: coin
      path: cmd/coin                // project path
      env:            // env variables available at startup
            test: test
            myvar: value
      commands:               // go commands supported
        vet:
            status: true
        fmt:
            status: true
            args:
            - -s
            - -w
        test:
            status: true
            method: gb test    // support different build tools
        generate:
            status: true
        install:
            status: true
        build:
            status: false
            method: gb build    // support differents build tool
            args:               // additional params for the command
            - -race
        run:
            status: true
      args:                     // arguments to pass at the project
      - --myarg
      watcher:
          paths:                 // watched paths are relative to directory you run haunt in
          - /
          ignore_paths:          // ignored paths
          - vendor
          extensions:                  // watched extensions
          - go
          - html
          scripts:
          - type: before
            command: echo before global
            global: true
            output: true
          - type: before
            command: echo before change
            output: true
          - type: after
            command: echo after change
            output: true
          - type: after
            command: echo after global
            global: true
            output: true
          errorOutputPattern: mypattern   //custom error pattern
