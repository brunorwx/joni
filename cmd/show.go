package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/brunorwx/joni/internal/store"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show a snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println("invalid id")
			return
		}
		dir, _ := os.UserHomeDir()
		s, err := store.OpenStore(dir)
		if err != nil {
			fmt.Println("failed to open store:", err)
			return
		}
		defer s.Close()
		sn, err := s.Get(id)
		if err != nil {
			fmt.Println("not found")
			return
		}
		fmt.Printf("ID: %d\nLang: %s\nDesc: %s\nTags: %v\n\n%s\n", sn.ID, sn.Language, sn.Description, sn.Tags, sn.Content)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
