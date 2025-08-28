package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/brunorwx/joni/internal/store"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search snippets",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		q := strings.Join(args, " ")
		dir, _ := os.UserHomeDir()
		s, err := store.OpenStore(dir)
		if err != nil {
			fmt.Println("failed to open store:", err)
			return
		}
		defer s.Close()
		res, err := s.Search(q)
		if err != nil {
			fmt.Println("search error:", err)
			return
		}
		for _, r := range res {
			fmt.Printf("%d\t%s\t%s\t%v\n", r.ID, r.Language, r.Description, r.Tags)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
