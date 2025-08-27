package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vrld/absicht/internal"
)

var emailPath string
var goEditNoTimeToWaste bool

func init() {
	rootCmd.Flags().StringVarP(&emailPath, "file", "f", "-", "Read initial email from this path; `-' means stdin")
	rootCmd.Flags().BoolVarP(&goEditNoTimeToWaste, "edit", "e", false, "Start editing the email right away")
	rootCmd.Flags().StringP("sendmail", "s", "msmtp -t --read-envelope-from", "Command to send mail; mail will be piped to stdin")
}

var rootCmd = &cobra.Command{
	Use: "absicht",
	Short: "Compose emails",
	Long: `Absicht is a mail composer.

Edit the text body with your $EDITOR. Manage attachments.
Save emails and send them through msmtp.
`,
	Run: func(cmd *cobra.Command, args []string) {
		zone.NewGlobal()
		defer zone.Close()

		model := internal.InitialModel()

		err := readEmail(&model)

		p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

		if err != nil {
			go p.Send(err)
		}

		if goEditNoTimeToWaste {
			go p.Send(internal.EditEmail{})
		}

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running the program: %v", err)
		}
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetEnvPrefix("ABSICHT")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
		viper.AutomaticEnv()
		return viper.BindPFlags(cmd.Flags())
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readEmail(model *internal.Model) error {
	if emailPath == "-" {
		stat, _ := os.Stdin.Stat()
		hasPipedInput := (stat.Mode() & os.ModeCharDevice) == 0
		if hasPipedInput {
			return model.ReadEmail(bufio.NewReader(os.Stdin))
		}
		return nil
	}

	file, err := os.Open(emailPath)
	if err == nil {
		err = model.ReadEmail(file)
		file.Close()
	}

	return err
}
