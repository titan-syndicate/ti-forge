package scaffold

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/titan-syndicate/ti-scaffold/internal/logger"
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
		logger.Log.Info("Starting scaffold command...")
		defer logger.Sync() // Ensure logs are flushed when the command exits

		// Initialize Viper
		viper.SetConfigName("scaffold")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		// If config file is specified, use it
		if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
			logger.Log.Debugw("Using config file", "path", configFile)
			viper.SetConfigFile(configFile)
		}

		// Bind flags
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("package", cmd.Flags().Lookup("package"))

		// Read config
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				logger.Log.Errorw("Error reading config file", "error", err)
				return fmt.Errorf("error reading config file: %v", err)
			}
			logger.Log.Warn("No config file found, using default values")
		}

		// Unmarshal config
		if err := viper.Unmarshal(&config); err != nil {
			logger.Log.Errorw("Error unmarshaling config", "error", err)
			return fmt.Errorf("error unmarshaling config: %v", err)
		}

		// Validate config
		if config.Name == "" {
			logger.Log.Error("Plugin name is required")
			return fmt.Errorf("plugin name is required")
		}
		if config.Package == "" {
			logger.Log.Error("Package name is required")
			return fmt.Errorf("package name is required")
		}

		logger.Log.Infow("Scaffolding plugin",
			"name", config.Name,
			"package", config.Package,
		)

		// Initialize Jet
		view := jet.NewSet(
			jet.NewOSFileSystemLoader("cmd/scaffold/templates"),
		)

		// Create output directory
		outputDir := "output"
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			logger.Log.Errorw("Error creating output directory",
				"path", outputDir,
				"error", err,
			)
			return fmt.Errorf("error creating output directory: %v", err)
		}

		// Process templates
		templateFiles := []string{"base.jet", "plugin.jet"}
		for _, tmpl := range templateFiles {
			logger.Log.Debugw("Processing template", "file", tmpl)

			// Read template from embedded FS
			templateContent, err := templates.ReadFile("templates/" + tmpl)
			if err != nil {
				logger.Log.Errorw("Error reading template",
					"file", tmpl,
					"error", err,
				)
				return fmt.Errorf("error reading template %s: %v", tmpl, err)
			}

			// Parse template
			t, err := view.Parse(tmpl, string(templateContent))
			if err != nil {
				logger.Log.Errorw("Error parsing template",
					"file", tmpl,
					"error", err,
				)
				return fmt.Errorf("error parsing template %s: %v", tmpl, err)
			}

			var buf bytes.Buffer
			if err := t.Execute(&buf, nil, config); err != nil {
				logger.Log.Errorw("Error executing template",
					"file", tmpl,
					"error", err,
				)
				return fmt.Errorf("error executing template %s: %v", tmpl, err)
			}

			outputFile := filepath.Join(outputDir, tmpl[:len(tmpl)-4]+".go")
			if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
				logger.Log.Errorw("Error writing output file",
					"file", outputFile,
					"error", err,
				)
				return fmt.Errorf("error writing output file %s: %v", outputFile, err)
			}

			logger.Log.Infow("Generated file", "path", outputFile)
		}

		logger.Log.Info("✅ Scaffolding completed successfully!")
		logger.Sync() // Explicitly flush logs before returning
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&devMode, "dev", false, "Run in development mode")
	rootCmd.Flags().String("name", "", "Plugin name")
	rootCmd.Flags().String("package", "", "Package name")
	rootCmd.Flags().String("config", "", "Path to config file")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Get log level from flag
	logLevel, _ := rootCmd.PersistentFlags().GetString("log-level")

	// Initialize logger if not already initialized
	if err := logger.Init(logLevel); err != nil {
		return fmt.Errorf("error initializing logger: %v", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting scaffold command execution...")
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Errorw("Error executing command", "error", err)
		// Ensure logs are flushed before returning error
		logger.Sync()
		return err
	}
	logger.Log.Info("Scaffold command execution completed")
	// Ensure logs are flushed before returning success
	logger.Sync()
	return nil
}
