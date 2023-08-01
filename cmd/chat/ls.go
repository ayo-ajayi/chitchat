package chat

import (
	"fmt"
	"github.com/ayo-ajayi/chitchat/gpt"
	rdb "github.com/ayo-ajayi/chitchat/redis"
	"github.com/spf13/cobra"
)

var refresh bool
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list of available gpt models",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		var ls []string
		var err error
		if refresh {
			ls, err = gpt.LS()
			if err != nil {
				panic(err)
			}
			c := rdb.DefaultClient()
			c.SetList("chitchat:model", ls)
			defer c.C.Close()
		} else {
			ls = gpt.ListOfModels
		}
		for i, l := range ls {
			fmt.Printf("[%v]: %v\n", i, l)
		}
		fmt.Println("select the model of choice with the `-m` flag and the index of the model from the list")
		fmt.Printf("default model is %v\n", ls[len(ls)-1])
	},
}

func init() {
	lsCmd.Flags().BoolVarP(&refresh, "refresh", "r", false, "refresh the models list in the cache")
	ChatCmd.AddCommand(lsCmd)
}
