package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/brunorwx/joni/internal/model"
	"github.com/brunorwx/joni/internal/store"
	"github.com/spf13/cobra"
)

var addLang string
var addTags string
var addDesc string

var addCmd = &cobra.Command{
	Use:   "add [content]",
	Short: "Add a new snippet",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content := strings.Join(args, " ")
		dir, _ := os.UserHomeDir()
		s, err := store.OpenStore(dir)
		if err != nil {
			fmt.Println("failed to open store:", err)
			return
		}
		defer s.Close()
		sn := &model.Snippet{
			Content:     content,
			Language:    addLang,
			Tags:        parseTags(addTags),
			Description: addDesc,
		}
		id, err := s.Add(sn)
		if err != nil {
			fmt.Println("error adding snippet:", err)
			return
		}
		fmt.Println("added snippet id", id)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&addLang, "lang", "", "language of snippet")
	addCmd.Flags().StringVar(&addTags, "tags", "", "comma separated tags")
	addCmd.Flags().StringVar(&addDesc, "desc", "", "description")
}

func parseTags(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
