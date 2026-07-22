APP := my_service
BIN := bin

.PHONY: all front backend manage cron dev tidy clean

all: front backend manage cron   ## 构建全部四个二进制

front:                          ## 前台 API
	go build -o $(BIN)/frontapi   ./app/frontapi
backend:                        ## 后台 API
	go build -o $(BIN)/backendapi ./app/backendapi
manage:                         ## 总后台 API
	go build -o $(BIN)/manageapi  ./app/manageapi
cron:                           ## 定时任务
	go build -o $(BIN)/cron       ./app/cron

dev:                            ## 本地一体化运行(单进程全挂)
	gf run main.go

tidy:
	go mod tidy

clean:
	rm -rf $(BIN)
