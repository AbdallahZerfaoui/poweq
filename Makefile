.PHONY: build solve re-solve clean
.DEFAULT_GOAL := re-build
EXEC = poweq.exe
IP = localhost
PORT = 8080

build: 
	go build -o $(EXEC) ./cmd/poweq

re-build: clean build

solve: build
	./$(EXEC) solve -n 6 -m 2 -a 0 -tol 1e-6 -maxIter 100 -alg newton

clean:
	$(RM) -f $(EXEC)

re-solve: build solve

# ----- API Testing -----
curl-equation: 
	curl -X POST http://$(IP):$(PORT)/solve \
	  -H "Content-Type: application/json" \
	  -d '@equation.json' && echo "\n"

