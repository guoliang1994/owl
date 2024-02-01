package owl

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"os/exec"
	"owl/log"
	"runtime"
)

/*
control-group（默认值）：当前控制组里面的所有子进程，都会被杀掉
process：只杀主进程
mixed：主进程将收到 SIGTERM 信号，子进程收到 SIGKILL 信号
*/
type KillMode string

const (
	KillMain KillMode = "process"
	KillAll  KillMode = "control-group"
)

// 从 service 复制过来，增加了  KillMode,Type
const systemdScript = `[Unit]
Description={{.Description}}
ConditionFileIsExecutable={{.Path|cmdEscape}}
{{range $i, $dep := .Dependencies}} 
{{$dep}} {{end}}

[Service]
StartLimitInterval=5
StartLimitBurst=10
Type=simple
KillMode=%s
ExecStart={{.Path|cmdEscape}}{{range .Arguments}} {{.|cmd}}{{end}}
{{if .ChRoot}}RootDirectory={{.ChRoot|cmd}}{{end}}
{{if .WorkingDirectory}}WorkingDirectory={{.WorkingDirectory|cmdEscape}}{{end}}
{{if .UserName}}User={{.UserName}}{{end}}
{{if .ReloadSignal}}ExecReload=/bin/kill -{{.ReloadSignal}} "$MAINPID"{{end}}
{{if .PIDFile}}PIDFile={{.PIDFile|cmd}}{{end}}
{{if and .LogOutput .HasOutputFileSupport -}}
StandardOutput=file:{{.LogDirectory}}/{{.Name}}.out
StandardError=file:{{.LogDirectory}}/{{.Name}}.err
{{- end}}
{{if gt .LimitNOFILE -1 }}LimitNOFILE={{.LimitNOFILE}}{{end}}
{{if .Restart}}Restart={{.Restart}}{{end}}
{{if .SuccessExitStatus}}SuccessExitStatus={{.SuccessExitStatus}}{{end}}
RestartSec=120
EnvironmentFile=-/etc/sysconfig/{{.Name}}

{{range $k, $v := .EnvVars -}}
Environment={{$k}}={{$v}}
{{end -}}

[Install]
WantedBy=multi-user.target
`

const (
	Start     = "start"
	Stop      = "stop"
	Restart   = "restart"
	Install   = "install"
	Uninstall = "uninstall"
)

// App 应用
type App struct {
	apps        []*App // 多应用
	binName     string // 可执行文件名称，英文
	name        string // 应用名称
	description string // 应用描述

	svc       service.Service // 程序注册的系统服务
	startFunc func()          // 程序启动执行的方法
	Console   *cobra.Command  // 命令行调用程序
}

func (i *App) Start(s service.Service) error {
	go i.startFunc()
	return nil
}
func (i *App) Stop(s service.Service) error {
	return nil
}

// AddApp 增加app
func (i *App) AddApp(apps ...*App) *App {
	i.apps = append(i.apps, apps...)
	return i
}

// Run 运行 应用
func (i *App) Run() {

	// 如果多个应用同时安装
	if len(i.apps) > 0 {
		for _, app := range i.apps {
			i.Console.AddCommand(app.Console)
		}
	}

	i.Console.Execute()
}

// install 安装应用
func (i *App) install() error {
	_ = i.svc.Stop()
	_ = i.svc.Uninstall()
	err := i.svc.Install()
	if err == nil {
		log.PrintLnBlue(i.name + "安装成功")
	}
	return err
}

func (i *App) uninstall() error {
	_ = i.svc.Stop() // 不处理停止错误，因为有可能没启动
	err := i.svc.Uninstall()
	if err == nil {
		log.PrintLnBlue(i.name + "卸载成功")
	}
	return err
}

func (i *App) run() error {
	err := i.svc.Run()
	return err
}

func (i *App) control(command string) error {
	var err error
	if service.Platform() == "unix-systemv" {
		terminal := exec.Command("/etc/init.d/"+i.binName, command)
		err = terminal.Run()
	} else {
		err = service.Control(i.svc, command)
	}

	if err == nil {
		log.PrintLnBlue(i.name + command + "成功")
	}
	return err
}

func NewApp(binName, name, description string, killMode KillMode, startFunc func()) *App {
	cmd := &SystemCtl{
		Command: &cobra.Command{
			Use:   binName,
			Short: binName,
			Long:  description,
		},
	}
	app := &App{
		binName:     binName,
		name:        name,
		description: description,
		startFunc:   startFunc,
		Console:     cmd.Command,
	}

	if startFunc != nil {

		cmd.regCtlCmd(app) // 注册基本的命令行命令

		options := make(service.KeyValue)
		systemdScript := fmt.Sprintf(systemdScript, killMode)
		options["SystemdScript"] = systemdScript
		svcConfig := &service.Config{
			Name:        name,
			DisplayName: name,
			Description: description,
			Option:      options,
		}
		if runtime.GOOS != "windows" {
			svcConfig.Dependencies = []string{
				"Requires=network.target",
				"After=network-online.target syslog.target"}
			svcConfig.UserName = "root"
		}
		if binName != "" {
			svcConfig.Arguments = append(svcConfig.Arguments, binName) // 系统服务调用
		}
		svcConfig.Arguments = append(svcConfig.Arguments, "run") // 系统服务启动执行的方法
		svc, err := service.New(app, svcConfig)
		if err != nil {
			return nil
		}

		app.svc = svc
	}

	return app
}

type SystemCtl struct {
	*cobra.Command
}

// AddCmd 添加命令行
func (i *SystemCtl) AddCmd(commands ...*cobra.Command) {
	i.AddCommand(commands...)
}

func (i *SystemCtl) regCtlCmd(app *App) {
	i.AddCommand(&cobra.Command{
		Use:   Install,
		Short: "安装" + app.name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.install()
		},
	})
	i.AddCommand(&cobra.Command{
		Use:   Uninstall,
		Short: "卸载" + app.name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.uninstall()
		},
	})
	i.AddCommand(&cobra.Command{
		Use:   Start,
		Short: "启动" + app.name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.control(Start)
		},
	})
	i.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: "停止" + app.name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.control(Stop)
		},
	})
	i.AddCommand(&cobra.Command{
		Use:    "run",
		Short:  "前台运行" + app.name,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.run()
		},
	})
	i.AddCommand(&cobra.Command{
		Use:   "restart",
		Short: "重启" + app.name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.control(Restart)
		},
	})
}
