// Copyright 2016 zm@huantucorp.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
/*
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| -_- |)
                  O\  =  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
         佛祖保佑       永无BUG
*/
package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/globalways/common/logger"
	"github.com/globalways/dvip/api"
	"github.com/globalways/dvip/g"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

var (
	version = kingpin.Flag("version", "show version of dvip").Short('v').Default("false").Bool()
	cfgFile = kingpin.Flag("cfg", "config file location").Short('c').Default(filepath.Join("config", "cfg.json")).String()
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	kingpin.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	cfg := g.ParseConfig(*cfgFile)
	// init log
	logger.StartLogger(cfg.Debug, logger.UDP, cfg.Logstash, g.APP_ID, map[string]interface{}{})

	// api
	go api.Start()

	handleSignal(os.Getpid())
}

func handleSignal(pid int) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	logrus.WithField("pid", pid).Info("has registered signal notify.")

	for {
		s := <-sigs
		logrus.Infof("has received signal: %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logrus.Info("is graceful shutting down...")

			logrus.WithField("pid", pid).Info("has exited")
			os.Exit(0)
		}
	}
}
