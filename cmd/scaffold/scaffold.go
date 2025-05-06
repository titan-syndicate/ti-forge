package scaffold

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed templates
var templates embed.FS

type ScaffoldConfig struct {
	Name    string
	Package string
}

var (
	devMode bool
	config  ScaffoldConfig
)

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold a new Titanium plugin",
	Long:  `Scaffold creates a new Titanium plugin project with the necessary boilerplate code.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Starting scaffold command...")

		// Initialize Viper
		viper.SetConfigName("scaffold")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		// If config file is specified, use it
		if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
			viper.SetConfigFile(configFile)
		}

		// Bind flags
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("package", cmd.Flags().Lookup("package"))

		// Read config
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return fmt.Errorf("error reading config file: %v", err)
			}
		}

		// Unmarshal config
		if err := viper.Unmarshal(&config); err != nil {
			return fmt.Errorf("error unmarshaling config: %v", err)
		}

		// Validate config
		if config.Name == "" {
			return fmt.Errorf("plugin name is required")
		}
		if config.Package == "" {
			return fmt.Errorf("package name is required")
		}

		// Initialize Jet
		view := jet.NewSet(
			jet.NewOSFileSystemLoader("cmd/scaffold/templates"),
		)

		// Create output directory
		outputDir := "output"
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %v", err)
		}

		// Process templates
		templateFiles := []string{"base.jet", "plugin.jet"}
		for _, tmpl := range templateFiles {
			// Read template from embedded FS
			templateContent, err := templates.ReadFile("templates/" + tmpl)
			if err != nil {
				return fmt.Errorf("error reading template %s: %v", tmpl, err)
			}

			// Parse template
			t, err := view.Parse(tmpl, string(templateContent))
			if err != nil {
				return fmt.Errorf("error parsing template %s: %v", tmpl, err)
			}

			var buf bytes.Buffer
			if err := t.Execute(&buf, nil, config); err != nil {
				return fmt.Errorf("error executing template %s: %v", tmpl, err)
			}

			outputFile := filepath.Join(outputDir, tmpl[:len(tmpl)-4]+".go")
			if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
				return fmt.Errorf("error writing output file %s: %v", outputFile, err)
			}

			log.Printf("Generated %s", outputFile)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&devMode, "dev", false, "Run in development mode")
	rootCmd.Flags().String("name", "", "Plugin name")
	rootCmd.Flags().String("package", "", "Package name")
	rootCmd.Flags().String("config", "", "Path to config file")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	log.Println("Executing scaffold command...")
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error executing command: %v", err)
		return err
	}
	return nil
}
