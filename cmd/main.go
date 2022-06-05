package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	aegirdungeons "github.com/aegir-tactics/aegir-dungeons"
)

func main() {
	logger := logrus.New()
	logger.Info("server starting game")

	config, err := aegirdungeons.NewGameConfig()
	if err != nil {
		logger.Fatal("error configuring:", err)
	}

	router := gin.Default()
	go router.Run(":" + config.Port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	game, err := aegirdungeons.NewGame(logger, config)
	if err != nil {
		logger.Fatal("error creating game:", err)
	}
	defer game.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		<-sigs
		close(sigs)
		cancelFunc()
	}()

	if err := game.InitializeAndWait(ctx, (time.Minute * 2)); err != nil {
		logger.Fatal("error initializing game:", err)
	}
	go func() {
		time.Sleep(time.Second * 1)
		game.StartRound(ctx)
	}()
	game.Play(ctx)
}
