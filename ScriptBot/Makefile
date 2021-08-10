build:
	mkdir -p Output
	go build -o Output/ScriptBot cmd/main.go
	cp -T cmd/.envexample Output/.envexample
	cp -T deploy.sh Output/deploy.sh
	cp ServiceCheck Output/ServiceCheck
clean:
	rm -rf Output 2>&1 > /dev/null
