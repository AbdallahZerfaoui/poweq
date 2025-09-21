EXEC = poweq

build: 
	go build -o $(EXEC) main.go

re-build: clean build

solve: build
	./$(EXEC) solve -n 6 -m 2 -a 0 -tol 1e-6 -maxIter 100 -alg newton

clean:
	$(RM) -f $(EXEC)

re-solve: build solve

.PHONY: build solve re-solve clean