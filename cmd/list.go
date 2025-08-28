package cmd

import (
	"fmt"
	"os"

	"github.com/brunorwx/joni/internal/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snippets",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.UserHomeDir()
		s, err := store.OpenStore(dir)
		if err != nil {
			fmt.Println("failed to open store:", err)
			return
		}
		defer s.Close()
		items, err := s.List()
		if err != nil {
			fmt.Println("error listing:", err)
			return
		}
		for _, it := range items {
			fmt.Printf("ID\tLanguage\tDescription\tTags\n")
			fmt.Printf("%d\t%s\t%s\t%v\n", it.ID, it.Language, it.Description, it.Tags)
			fmt.Printf("Content\n")
			fmt.Printf("%s\n", it.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
