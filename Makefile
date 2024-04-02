build: 
	go build -o dist/ .

tidy:
	go mod tidy

clean:
	rm -rf dist

uninstall:
	rm -f /usr/bin/stew
	rm -f /usr/share/zsh/site-functions/_stew

install:
	cp ./dist/gspot /usr/bin
	cp ./completions/_stew /usr/share/zsh/site-functions/_stew
