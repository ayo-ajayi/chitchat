package chat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ayo-ajayi/chitchat/history"
	"github.com/ayo-ajayi/chitchat/redis"
	"github.com/spf13/cobra"
)

type ChatService struct {
	history *history.History
	list *List
}

var gptModel int16
var nosave bool
var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "interactive conversation with ai",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		chat:=ChatService{
			history: history.NewHistory(redis.DefaultClient()),
			list: Newlist(redis.DefaultClient()),
		}
		fmt.Println("Welcome to My Chitchat CLI App! Enter Ctrl+C to exit.")
		fmt.Println("Type 'exit' to quit.")
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		var model string
		if gptModel < 0 {
			model = string(chat.list.gpt.DefaultModel()) 
			fmt.Println("using default model: ", model)
			fmt.Println()
		} else {
			model = chat.list.gpt.ListOfModels()[gptModel]
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
					res = chat.list.gpt.Chat(input, &model)
					if !nosave {
						if err := chat.history.SaveChat(input, res); err != nil {
							log.Fatalln(err)
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
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		return
	}
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func init() {
	ChatCmd.Flags().Int16VarP(&gptModel, "model", "m", -1, "gpt model to use")
	ChatCmd.Flags().BoolVarP(&nosave, "nosave", "n", false, "do not save chat")
}
