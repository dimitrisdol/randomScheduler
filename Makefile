GO = go
BIN = kube-scheduler
VERSION := 0.0.1

DOCKER = docker
WHO ?= jmmisd
IMG_BASENAME = randomschedulertest01
LOCAL_IMG := $(WHO)/$(IMG_BASENAME)
IMG_TAG ?= $(VERSION)
LOCAL_REGISTRY = localhost:5000


.PHONY: all
all: image push-local-reg

.PHONY: image
image:
	$(DOCKER) build --no-cache -f ./Dockerfile --build-arg BIN=$(BIN) \
		--build-arg VERSION=$(VERSION) -t $(LOCAL_IMG):$(IMG_TAG) .
	$(DOCKER) image prune --force --filter label=stage=builder

.PHONY: push-local-reg
push-local-reg:
	$(DOCKER) tag $(LOCAL_IMG):$(IMG_TAG) \
		$(LOCAL_REGISTRY)/$(IMG_BASENAME):$(IMG_TAG)
	$(DOCKER) push $(LOCAL_REGISTRY)/$(IMG_BASENAME):$(IMG_TAG)

.PHONY: check build
check:
	$(GO) build -v ./...
build:
	$(GO) build -v -o $(BIN) \
		-ldflags '-X k8s.io/component-base/version.gitVersion=v$(VERSION) -w' \
		./cmd/kube-scheduler/main.go


.PHONY: clean-image clean-local clean
clean: clean-local clean-image
clean-local:
	-$(RM) -v $(BIN)
clean-image:
	-$(DOCKER) rmi $(LOCAL_IMG):$(IMG_TAG)
