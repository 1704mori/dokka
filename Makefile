.PHONY: dist
dist:
	@npm run build

.PHONY: build
	CGO_ENABLED=0 go build -o dokka

.PHONY: dev
dev:
	npm run dev
