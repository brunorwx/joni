package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/brunorwx/joni/internal/store"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a snippet",
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
		if err := s.Delete(id); err != nil {
			fmt.Println("delete failed:", err)
			return
		}
		fmt.Println("deleted")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
