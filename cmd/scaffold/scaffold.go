package scaffold

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	devMode bool
)

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold a new Titanium plugin",
	Long:  `Scaffold creates a new Titanium plugin project with the necessary boilerplate code.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Starting scaffold command...")

		// Prompt for plugin name using stderr
		fmt.Fprintf(os.Stderr, "Enter plugin name: ")

		reader := bufio.NewReader(os.Stdin)
		pluginName, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading plugin name: %v", err)
		}

		// Clean up the input
		pluginName = strings.TrimSpace(pluginName)
		log.Printf("Received plugin name: %q", pluginName)

		if pluginName == "" {
			return fmt.Errorf("plugin name is required")
		}

		log.Printf("Creating plugin: %s", pluginName)
		// TODO: Implement plugin creation logic
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&devMode, "dev", false, "Run in development mode")
	// Ensure output goes to stdout
	rootCmd.SetOut(os.Stdout)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	log.Println("Executing scaffold command...")
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error executing command: %v", err)
		os.Exit(1)
	}
}
