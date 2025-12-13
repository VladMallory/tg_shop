IMAGE := telegram-go-bot
TAG := dev
NAME := parfume-bot

# ---------- BUILD ----------

.PHONY: build
build:
	@docker build -q -t $(IMAGE):$(TAG) .

.PHONY: build-amd64
build-amd64:
	@docker buildx build -q --platform linux/amd64 -t $(IMAGE):amd64 --load .

# ---------- RUN ----------

.PHONY: run
run:
	docker -f $(NAME) 2>/dev/null || true
	docker run -it \
		--platform linux/amd64 \
		--env-file .env \
		-v $(PWD):/app \
		-w /app \
		--name $(NAME) \
		$(IMAGE):amd64

.PHONY: run-amd64
run-amd64: build-amd64
	@docker rm -f $(NAME) 2>/dev/null || true
	@docker run -d \
		--env-file .env \
		-v $(PWD):/app \
		-w /app \
		--name $(NAME) \
		$(IMAGE):$(TAG) > /dev/null
	@docker logs -f $(NAME)
	
# ---------- LOGS ----------
.PHONY: logs
logs:
	docker logs -f $(NAME)

# ---------- Combo ----------

# При make по умолчанию будем билдить и сразу запускать контейнер с логами
.PHONY: build-and-run
build-and-run: build run

# ---------- DEFAULT ----------
.DEFAULT_GOAL := build-and-run