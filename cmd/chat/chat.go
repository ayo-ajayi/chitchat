package chat

import (
	"bufio"
	"fmt"
	"github.com/ayo-ajayi/chitchat/gpt"
	"github.com/ayo-ajayi/chitchat/history"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var gptModel int16
var nosave bool
var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "interactive conversation with ai",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to My Chitchat CLI App! Enter Ctrl+C to exit.")
		fmt.Println("Type 'exit' to quit.")
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		var model string
		if gptModel < 0 {
			model = string(gpt.DefaultModel())
			fmt.Println("using default model: ", model)
			fmt.Println()
		} else {
			model = gpt.ListOfModels[gptModel]
			fmt.Println("using model: ", model)
			fmt.Println()
		}
		scanner := bufio.NewScanner(os.Stdin)
		go func() {
			for {
				fmt.Print(">> ")
				if !scanner.Scan() {
					break
				}
				input := scanner.Text()
				switch input {
				case "exit":
					fmt.Println("Exiting...")
					os.Exit(0)
				case "":

				case "cls":
					clearScreen()
				default:
					var res string
					res = gpt.Chat(input, &model)
					if !nosave {
						if err := history.SaveChat(input, res); err != nil {
							fmt.Println(err)
						}
					}
					fmt.Println(res)
					fmt.Println()
				}
			}
		}()
		<-interrupt
		fmt.Println("\nExiting...")
		os.Exit(0)
	},
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func init() {
	ChatCmd.Flags().Int16VarP(&gptModel, "model", "m", -1, "gpt model to use")
	ChatCmd.Flags().BoolVarP(&nosave, "nosave", "n", false, "do not save chat")
}
