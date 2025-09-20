EXEC = poweq

build: 
	go build -o $(EXEC) main.go

run: build
	./$(EXEC) -n 6 -m 2 -a 0 -tol 1e-6 -maxIter 100 -alg newton

clean:
	$(RM) -f $(EXEC)

re: build run

.PHONY: build run clean