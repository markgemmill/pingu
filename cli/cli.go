package main

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"os"
	"pingu/pkg"
	"time"
)

type Context struct {
}

type Globals struct {
	Version VersionFlag `name:"version" help:"Print version information."`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Printf("v%s\n", vars["version"])
	app.Exit(0)
	return nil
}

type EmailOptions struct {
	Email         bool   `short:"S" name:"email" group:"email-options" help:"Flag must be present for email sending."`
	EmailHost     string `short:"H" name:"email-host" group:"email-options" help:"Domain or IP of SMTP host."`
	EmailPort     int    `short:"p" name:"email-port" default:"25" group:"email-options" help:"Port of SMTP host."`
	EmailUser     string `short:"U" name:"email-user" group:"email-options" help:"SMTP user account name."`
	EmailPassword string `short:"P" name:"email-password" group:"email-options" help:"SMTP user account password."`
	EmailFrom     string `name:"email-from" group:"email-options" help:"Email FROM address."`
	EmailTo       string `name:"email-to" group:"email-options" help:"One or more TO addresses separated by semi-colon."`
	EmailCc       string `name:"email-cc" group:"email-options" help:"One or more CC addresses separated by semi-colon."`
}

func (opt *EmailOptions) Validate() error {
	console.Trace("EmailOptions.Validate() called!")
	if opt.Email == true {
		if opt.EmailHost == "" || !pkg.ValidEmail(opt.EmailFrom, true) || !pkg.ValidEmail(opt.EmailTo, false) || !pkg.ValidEmail(opt.EmailCc, false) {
			return errors.New("invalid alert configuration")
		}
	}
	return nil
}

type UrlOptions struct {
	Url       string `arg:"" name:"url" required:"" help:"Url to check."`
	StoreName string `short:"s" name:"store-name" help:"The store file name. If not supplied name will be hash of the url."`
}

type RetryOptions struct {
	Retries        int `short:"r" name:"retries" default:"0" help:"The number of times to retry after a failed check."`
	RetryIncrement int `short:"i" name:"retry-increment" default:"1" help:"The power to raise wait seconds after each retry. Maximum of 3. Example:  f '-i=3' seconds would be 1, 8, 27, etc...."`
}

func (opt *RetryOptions) Validate() error {
	if opt.Retries > 0 && (opt.RetryIncrement < 1 || opt.RetryIncrement > 3) {
		return errors.New("retry increments must be a value between 1 and 3")
	}
	return nil
}

type CheckCmd struct {
	UrlOptions
	ExpectedStatus  int      `short:"e" name:"expect-status" group:"assertion options" default:"200" help:"The expected http status."`
	ExpectedContent string   `short:"c" name:"expect-content" group:"assertion options" help:"A regular express that must match the returned content."`
	IgnorePeriod    []string `name:"ignore-period" sep:";" help:"A time span during which calls to check will be ignored. Example: 'SAT 10:00PM - SUN 1:00AM'"`
	RetryOptions
	AlertThreshold int64 `short:"a" name:"alert-threshold" default:"0" help:"Alert will be raise after this many consecutive failures."`
	Verbose        int   `short:"v" type:"counter" help:"Verbosity can have a value of 1-3. Example: --verbose=3 or -vvv."`
	EmailOptions
}

func (cmd *CheckCmd) Validate() error {
	return nil
}

func (cmd *CheckCmd) Run(ctx *Context) error {

	console = pkg.NewConsole(cmd.Verbose)

	currentTimestamp := time.Now()

	for _, ignoreText := range cmd.IgnorePeriod {
		console.Trace("Checking: %s\n", ignoreText)
		ignore := pkg.IsIgnorePeriodActive(ignoreText, currentTimestamp)
		if ignore == true {
			console.Info("%s %s\n", pkg.Red("Ignore time period:"), pkg.Yellow(ignoreText))
			console.Info("%s\n", pkg.Red("Current period is active. Exiting..."))
			os.Exit(0)
		}
	}

	record, err := pkg.CheckCommand(cmd.Url, cmd.ExpectedStatus, cmd.ExpectedContent, cmd.StoreName, console)

	if err != nil && cmd.Retries > 0 {
		retries := 1
		for retries <= cmd.Retries {
			seconds := pkg.CalculatePauseInSeconds(retries, cmd.RetryIncrement)
			console.Debug(pkg.Green("Retry #%d in %d seconds...\n"), retries, seconds/time.Second)
			time.Sleep(seconds)
			record, err = pkg.CheckCommand(cmd.Url, cmd.ExpectedStatus, cmd.ExpectedContent, cmd.StoreName, console)
			retries += 1
		}
	}

	if err != nil && record.Count >= cmd.AlertThreshold && cmd.Email == true {
		console.Dedent()
		console.Info(pkg.Yellow("Sending Email Alert...\n"))
		pkg.SendEmailAlert(
			pkg.NewSmtpServer(
				cmd.EmailHost,
				cmd.EmailPort,
				cmd.EmailUser,
				cmd.EmailPassword,
			),
			pkg.NewAlertEmail(
				cmd.EmailFrom,
				cmd.EmailTo,
				cmd.EmailCc,
			),
			cmd.Url,
			record)
	}

	return err
}

type ReportCmd struct {
	UrlOptions
	EmailOptions
}

func (cmd *ReportCmd) Run(ctx *Context) error {
	// do check command
	store := pkg.NewStore(cmd.Url, cmd.StoreName)
	store.Read()

	message := pkg.ReportMessage{
		Store: store.Data,
	}
	message.Initialize()
	fmt.Print(message.ToText())

	if cmd.Email == false {
		return nil
	}

	pkg.SendEmailReport(
		pkg.NewSmtpServer(
			cmd.EmailHost,
			cmd.EmailPort,
			cmd.EmailUser,
			cmd.EmailPassword,
		),
		pkg.NewAlertEmail(
			cmd.EmailFrom,
			cmd.EmailTo,
			cmd.EmailCc,
		),
		&message,
	)
	return nil
}

type CLI struct {
	Globals

	Check  CheckCmd  `cmd:""`
	Report ReportCmd `cmd:""`
}
