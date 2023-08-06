package history

import (
	"fmt"
	"log"

	"github.com/ayo-ajayi/chitchat/history"
	"github.com/ayo-ajayi/chitchat/redis"
	"github.com/spf13/cobra"
)


type HistoryService struct {
	history *history.History
}
var defaultHistoryLimit int64 = 10
var limit int64 = defaultHistoryLimit
	
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Chat History",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		historyS := HistoryService{
			history: history.NewHistory(redis.DefaultClient()),
		}
		if limit <=0 {
			limit=defaultHistoryLimit
		}
		chat, err := historyS.history.GetChat(limit)
		if err != nil {
			log.Fatalln(err)
		}
		for _, c:=range chat{
			date, err:=historyS.history.GetDate(c.ID)
			if err != nil {
				date=""
			}
			fmt.Print(date, " ")
			for key, value := range c.Values {
				fmt.Println("Prompt:", key, "Answer:", value)
			}
			fmt.Println()
		}
	},
}
func init() {
	HistoryCmd.Flags().Int64VarP(&limit, "limit", "l", 10, "chat history limit")
}
