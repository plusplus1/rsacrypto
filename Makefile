
gen:
	protoc --go_out=plugins=grpc:. commonLib/rpcLib/*.proto

clean:
	@git clean -f -d -X
