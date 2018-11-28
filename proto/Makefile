build:
	prototool generate

lint:
	prototool format -l
	prototool lint

clean:
	find $(CURDIR)/../pkg -name '*.pb.go' -delete
	find $(CURDIR)/../app/src -name '*_pb*.js' -delete

test:
	prototool compile
