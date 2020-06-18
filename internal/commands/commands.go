package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"github.com/fenrirunbound/kubeconfig-merge/internal/source"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func Root(name string) *cobra.Command {
	var destFile string
	var commit bool
	var sourceFile string
	var overwrite bool

	cmd := &cobra.Command{
		Use:   name,
		Short: "merge multiple kube config files together",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf(`Loading source "%s"...`, sourceFile)
			src, err := source.New(sourceFile)
			if err != nil {
				return err
			}
			cmd.Println("done")

			cmd.Printf(`Loading destination "%s"...`, destFile)
			dest, err := source.New(destFile)
			if err != nil {
				return err
			}
			cmd.Println("done")

			err = dest.Combine(src, overwrite)
			if err != nil {
				return err
			}
			cmd.Println("Combined configs")

			output, err := yaml.Marshal(dest.RawConfig())
			if err != nil {
				return err
			}

			if commit {
				cmd.Printf("Committing combine file to %s\n", destFile)

				return ioutil.WriteFile(destFile, output, 0644)
			}

			cmd.Println("Combined kube configs:")
			cmd.Println(string(output))
			return nil
		},
	}

	defaultConfig := ""
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		defaultConfig = kubeconfig
	} else if home := homeDir(); home != "" {
		defaultConfig = filepath.Join(home, ".kube", "config")
	}

	cmd.Flags().StringVarP(&destFile, "dest", "o", defaultConfig, "Destination kube config to write values to")
	cmd.Flags().StringVarP(&sourceFile, "source", "s", "", "Source kube config file to inherit from")
	cmd.Flags().BoolVarP(&commit, "yes", "y", false, "actually write the merged result to the destination file")
	cmd.Flags().BoolVarP(&overwrite, "overwrite", "f", false, "overwrite the values the destination file already contains")

	return cmd
}
