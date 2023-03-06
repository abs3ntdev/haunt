pkgname := haunt
build: ${pkgname}

${pkgname}: $(shell find . -name '*.go')
	mkdir -p bin
	go build -o bin/${pkgname} .

decode:
	go run hack/unpack/main.go

pack-site:
	go-bindata -pkg haunt -o src/haunt/bindata.go -fs assets/...

completions:
	mkdir -p completions
	./bin/${pkgname} completion zsh > completions/_${pkgname}
	./bin/${pkgname} completion bash > completions/${pkgname}
	./bin/${pkgname} completion fish > completions/${pkgname}.fish

run:
	go run main.go

tidy:
	go mod tidy

clean:
	rm -rf bin
	rm -rf completions
	rm -rf assets

uninstall:
	rm -f /usr/bin/${pkgname}
	rm -f /usr/share/zsh/site-functions/_${pkgname}
	rm -f /usr/share/bash-completion/completions/${pkgname}
	rm -f /usr/share/fish/vendor_completions.d/${pkgname}.fish

install:
	cp bin/${pkgname} /usr/bin
	bin/${pkgname} completion zsh > /usr/share/zsh/site-functions/_${pkgname}
	bin/${pkgname} completion bash > /usr/share/bash-completion/completions/${pkgname}
	bin/${pkgname} completion fish > /usr/share/fish/vendor_completions.d/${pkgname}.fish
