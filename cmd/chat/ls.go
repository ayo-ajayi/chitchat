package chat

import (
	"fmt"
	"github.com/ayo-ajayi/chitchat/gpt"
	"github.com/ayo-ajayi/chitchat/redis"
	"github.com/spf13/cobra"
	"log"
)

type DBClient interface {
	SetList(key string, value []string) error
	GetList(key string) ([]string, error)
	Close() error
}
type List struct {
	gpt *gpt.GPT
	db DBClient
}

func Newlist(db DBClient) *List {
	g, err := gpt.NewGPT()
	if err != nil {
		log.Fatalln(err)
	}
	return &List{gpt: g,
		db: db,}
}

var refresh bool
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list of available gpt models",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		list := Newlist(redis.DefaultClient())
		fmt.Println()
		var ls []string
		var err error
		if refresh {
			ls, err = list.gpt.LS()
			if err != nil {
				log.Fatalln(err)
			}
			c := list.db
			c.SetList("chitchat:model", ls)
			defer c.Close()
		} else {
			ls = list.gpt.ListOfModels()
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
