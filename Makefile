build:
	go build -o build/git-lab

clean:
	rm -rf build/

rebuild:
	rm -rf build/ && go build -o build/git-lab

