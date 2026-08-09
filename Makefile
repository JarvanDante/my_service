APP := my_service
BIN := bin
MYDOCKER := /Users/wangdante/D/mydocker

.PHONY: all front backend manage cron dev tidy clean docker-up docker-down docker-logs docker-rebuild

all: front backend manage cron   ## 构建全部四个二进制

front:                          ## 前台 API
	go build -o $(BIN)/frontapi   ./app/frontapi
backend:                        ## 后台 API
	go build -o $(BIN)/backendapi ./app/backendapi
manage:                         ## 总后台 API（站点侧 manage 门面；控制面见 my_manage_service）
	go build -o $(BIN)/manageapi  ./app/manageapi
cron:                           ## 定时任务
	go build -o $(BIN)/cron       ./app/cron

dev:                            ## 本地一体化运行(单进程全挂)
	gf run main.go

tidy:
	go mod tidy

clean:
	rm -rf $(BIN)

# ---- Docker：一体化入口 :8000（mydocker）----
docker-up:
	cd $(MYDOCKER) && docker compose up -d --build my_service
	@echo "探活: curl -sS http://127.0.0.1:8000/health"

docker-down:
	cd $(MYDOCKER) && docker compose stop my_service

docker-logs:
	docker logs -f my_service

docker-rebuild:
	cd $(MYDOCKER) && docker compose up -d --build --force-recreate my_service
