## Quickstart
### Install
```
go install github.com/abs3ntdev/haunt@latest
```

or

```
git clone https://github.com/abs3ntdev/haunt
cd haunt
make build && sudo make install
```

#### aur

```
yay -S haunt-go-git
```

### Completions
completions will be automatically installed if you used the Makefile, if you did not you can generate completions with `haunt completion [bash/fish/powershell/zsh]`

for example: `haunt completion zsh > _haunt`

you can also source the output of the completion command directly in your .zshrc with:\
      `source <(haunt completion zsh) && compdef _haunt haunt`

## Commands List

### Init Command
This command will generate a .haunt.yaml with sane default for your current project/projects.\
If there is a main.go in the root directory it will be added along with any directories inside the relative path `cmd`

    haunt init


### Run Command

```
haunt run
```

the run command allows for specifying projects by name, all provided will be ran according to the config file:

Some examples:

    haunt run
    haunt run server api

### Add Command
Add a project, the same defaults init uses will be used for new projects unless flags are provided.

    haunt add [name] [--flags]

Possible flags are:
   
    -b, --build         Enable go build
    -f, --fmt           Enable go fmt
    -g, --generate      Enable go generate
    -h, --help          help for add
    -i, --install       Enable go install (default true)
    -p, --path string   Project base path (default "./")
    -r, --run           Enable go run (default true)
    -t, --test          Enable go test
    -v, --vet           Enable go vet

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
          - src
          - cmd/coin
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
