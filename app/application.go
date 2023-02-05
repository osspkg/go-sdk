package app

import (
	"os"

	"github.com/deweppro/go-sdk/console"
	"github.com/deweppro/go-sdk/log"
	"github.com/deweppro/go-sdk/syscall"
)

type (
	//ENV type for enviremants (prod, dev, stage, etc)
	ENV string

	App interface {
		Logger(log log.Logger) App
		Modules(modules ...interface{}) App
		ConfigFile(filename string, configs ...interface{}) App
		Run()
	}

	_app struct {
		cfile    string
		configs  Modules
		modules  Modules
		sources  Sources
		packages *_dic
		logout   *_log
		log      log.Logger
	}
)

// New create application
func New() App {
	return &_app{
		modules:  Modules{},
		configs:  Modules{},
		packages: newDic(),
	}
}

// Logger setup logger
func (a *_app) Logger(log log.Logger) App {
	a.log = log
	return a
}

// Modules append object to modules list
func (a *_app) Modules(modules ...interface{}) App {
	for _, mod := range modules {
		switch v := mod.(type) {
		case Modules:
			a.modules = a.modules.Add(v...)
		default:
			a.modules = a.modules.Add(v)
		}
	}

	return a
}

// ConfigFile set config file path and configs models
func (a *_app) ConfigFile(filename string, configs ...interface{}) App {
	a.cfile = filename
	for _, config := range configs {
		a.configs = a.configs.Add(config)
	}

	return a
}

// Run run application
func (a *_app) Run() {
	var err error
	if len(a.cfile) == 0 {
		a.logout = newLog(&Config{
			Level:   4,
			LogFile: "/dev/stdout",
		})
		a.log = log.Default()
		a.logout.Handler(a.log)
	}
	if len(a.cfile) > 0 {
		// read config file
		a.sources = Sources(a.cfile)

		// init logger
		config := &Config{}
		if err = a.sources.Decode(config); err != nil {
			console.FatalIfErr(err, "decode config file: %s", a.cfile)
		}
		a.logout = newLog(config)
		if a.log == nil {
			a.log = log.Default()
		}
		a.logout.Handler(a.log)
		a.modules = a.modules.Add(func() log.Logger { return a.log }, ENV(config.Env))

		// decode all configs
		var configs []interface{}
		configs, err = typingRefPtr(a.configs, func(i interface{}) error {
			return a.sources.Decode(i)
		})
		if err != nil {
			a.log.WithFields(log.Fields{
				"err": err.Error(),
			}).Fatalf("decode config file")
		}
		a.modules = a.modules.Add(configs...)

		if len(config.PidFile) > 0 {
			if err = syscall.Pid(config.PidFile); err != nil {
				a.log.WithFields(log.Fields{
					"err":  err.Error(),
					"file": config.PidFile,
				}).Fatalf("create pid file")
			}
		}
	}

	a.launch()
}

func (a *_app) launch() {
	ctx := NewContext()
	result := a.steps(
		[]step{
			{
				Message: "register app dependencies",
				Call:    func() error { return a.packages.Register(a.modules...) },
			},
			{
				Message: "build app dependencies",
				Call:    func() error { return a.packages.Build() },
			},
			{
				Message: "start app dependencies",
				Call:    func() error { return a.packages.Up(ctx) },
			},
		},
		func(er bool) {
			if er {
				ctx.Close()
				return
			}
			go syscall.OnStop(ctx.Close)
			<-ctx.Done()
		},
		[]step{
			{
				Message: "stop app dependencies",
				Call:    func() error { return a.packages.Down() },
			},
		},
	)
	console.FatalIfErr(a.logout.Close(), "close log file")
	if result {
		os.Exit(1)
	}
	os.Exit(0)
}

type step struct {
	Call    func() error
	Message string
}

func (a *_app) steps(up []step, wait func(bool), down []step) bool {
	var erc int

	for _, s := range up {
		a.log.Infof(s.Message)
		if err := s.Call(); err != nil {
			a.log.WithFields(log.Fields{
				"err": err.Error(),
			}).Errorf(s.Message)
			erc++
			break
		}
	}

	wait(erc > 0)

	for _, s := range down {
		a.log.Infof(s.Message)
		if err := s.Call(); err != nil {
			a.log.WithFields(log.Fields{
				"err": err.Error(),
			}).Errorf(s.Message)
			erc++
		}
	}

	return erc > 0
}
