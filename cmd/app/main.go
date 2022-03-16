package main

import (
	"awesomeProjectRucenter/internal/config"
	"awesomeProjectRucenter/internal/db"
	"awesomeProjectRucenter/internal/handler"
	"awesomeProjectRucenter/internal/server"
	"awesomeProjectRucenter/internal/service"
	"awesomeProjectRucenter/internal/vmParser"
	"awesomeProjectRucenter/pkg/erx"
	"awesomeProjectRucenter/pkg/tools"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := config.GetConfig(config.GetConfigInstance())
	cfg := config.GetConfigInstance()
	log := tools.GetLogrusInstance(cfg.Logger.Path)
	if err != nil {
		log.Error(erx.New(err))
	}

	function := flag.String("f", "default", "Specify one of the commands: initVmDb, reinitVmDb, printAsyncMap -d, printSyncMap -d, startServer")
	domain := flag.String("d", "195.24.64.34", "Specify DOMAIN of API to be parsed (to Database Folder or to Folder with logs)")
	flag.Parse()
	dbConn, err := db.NewSqliteDB("./data/vm.db")
	if err != nil {
		log.Error(err)
		return
	}
	defer dbConn.Close()

	repository := db.NewRepository(dbConn)
	serv := service.NewService(repository)
	h := handler.NewHandler(serv)

	switch *function {
	case "initVmDb":
		err = serv.InitDbWithData()
		if err != nil {
			log.Error(err)
			return
		}
	case "reinitVmDb":
		err = serv.InitDbWithData()
		if err != nil {
			log.Error(err)
			return
		}
	case "printAsyncMap":
		m, err := vmParser.PrintAsyncMap(*domain)
		if err != nil {
			log.Error(err)
		}
		fmt.Println(m)
	case "printSyncMap":
		m, err := vmParser.PrintSyncMap(*domain)
		if err != nil {
			log.Error(err)
		}
		vmParser.PrintSortedMap(m)

	case "startServer":
		srv := new(server.Server)
		port := fmt.Sprintf("%v", cfg.Server.Port)
		go func() {
			if err := srv.Run(port, h.InitRoutes()); err != nil {
				log.Error(erx.New(err))
			}
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit

		log.Info("Monitorlogs Shutting Down")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error(erx.New(err))
		}

		if err := dbConn.Close(); err != nil {
			log.Error(erx.New(err))
		}

	default:
		println("Expected flag (-f)")
		flag.PrintDefaults()
	}

}

func removeDuplicates(nums []int) int {
	k := 0
	for i, v := range nums {
		if v != nums[k] && k != i {
			k++
			nums[k] = nums[i]
		}
	}
	//nums = nums[:k]
	return k + 1
}
