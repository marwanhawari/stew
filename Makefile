build: 
	go build -o dist/ .

tidy:
	go mod tidy

clean:
	rm -rf dist

uninstall:
	rm -f ~/.local/bin/stew
	rm -f /usr/share/zsh/site-functions/_stew

install:
	cp ./dist/stew ~/.local/bin/stew
	cp ./completions/_stew /usr/share/zsh/site-functions/_stew
