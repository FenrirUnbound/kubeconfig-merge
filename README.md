# kubeconfig-merge
Merging of multiple kube config files

## Usage

```
merge multiple kube config files together

Usage:
  kubeconfig-merge [flags]

Flags:
  -o, --dest string     Destination kube config to write values to (default "$KUBECONFIG")
  -h, --help            help for kubeconfig-merge
  -s, --source string   Source kube config file to inherit from
  -y, --yes             actually write the merged result to the destination file
```

Example usage:

```
# will print out what the combined configuration would be
$ kubeconfig-merge -s ~/Downloads/new_config.yml -o ~/.kube/config 

# will actually write the combined configuration to the "-o" file
$ kubeconfig-merge -s ~/Downloads/new_config.yml -o ~/.kube/config --yes
```
