##@ Protobufs

.PHONY: protos
protos: $(tools/buf) ## Compile all protobufs
	$(tools/buf) generate

.PHONY: proto-extension
proto-extension: $(tools/buf) ## Compile internal/extension/proto/extension.proto
	$(tools/buf) generate internal/extension
