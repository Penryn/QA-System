package asynq

import (
	"QA-System/internal/handler/queue"
	"strconv"

	"github.com/hibiken/asynq"
)

var (
	Client *asynq.Client
	Srv *asynq.Server
)


func AsynqInit() {
	cfg := NewConfig()
	port:=strconv.Itoa(cfg.port)
	Client = asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.host+":"+port, 
		DB: (cfg.db)+1,
		Username: cfg.user, 
		Password: cfg.password,
	})

	Srv = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: cfg.host+":"+port, 
			DB: (cfg.db)+1,
			Username: cfg.user, 
			Password: cfg.password,
		},
        asynq.Config{
            Concurrency: 10, //并发数
			RetryDelayFunc: asynq.DefaultRetryDelayFunc,//重试延迟
        },
	)
	mux := asynq.NewServeMux()
    mux.HandleFunc(queue.TypeSubmitSurvey, queue.HandleSubmitSurveyTask)

    if err := Srv.Run(mux); err != nil {
        panic(err)
    }
}