package history

import (
	"fmt"

	"github.com/ayo-ajayi/chitchat/history"
	"github.com/spf13/cobra"
)
var defaultHistoryLimit int64 = 10
var limit int64 = defaultHistoryLimit
	
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Chat History",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if limit <=0 {
			limit=defaultHistoryLimit
		}
		chat, err := history.GetChat(limit)
		if err != nil {
			fmt.Println(err)
		}
		for _, c:=range chat{
			date, err:=history.GetDate(c.ID)
			if err != nil {
				date=""
			}
			fmt.Print(date, " ")
			for key, value := range c.Values {
				fmt.Println("Prompt:", key, "Answer:", value)
			}
		}
	},
}
func init() {
	HistoryCmd.Flags().Int64VarP(&limit, "limit", "l", 10, "chat history limit")
}
