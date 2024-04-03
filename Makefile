binaryName = stew
build: 
	go build -o dist/ .

tidy:
	go mod tidy

clean:
	rm -rf dist

uninstall:
	rm -f /usr/bin/${binaryName}
	rm -f /usr/share/zsh/site-functions/_${binaryName}

install:
	cp ./dist/${binaryName} /usr/bin
	cp ./completions/_${binaryName} /usr/share/zsh/site-functions
