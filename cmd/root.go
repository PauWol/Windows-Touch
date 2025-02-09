package cmd

import (
	"fmt"
	"os"

	"github.com/pauwol/touch/cmd/util"
)

type CMD struct {
	flags []Flag
}

func (c *CMD) AddFlag(Name string, Shorthand string, Usage string, Value string) {
	f := Flag{Name: Name, Shorthand: Shorthand, Usage: Usage, Value: Value}
	c.flags = append(c.flags, f)
}

type Flag struct {
	Name      string
	Shorthand string
	Usage     string
	Value     string
}

func (c *CMD) Execute() error {
	args := os.Args

	// Handle the case where no arguments are passed
	if len(args) == 1 {
		fmt.Println(util.Intro())
		return nil
	}

	// Extract flags and non-flag arguments
	rest, flags := c.extract(args[1:])

	// Handle help flag (-h, --help)
	if len(flags) > 0 {
		for _, v := range flags {
			if v.Name == "--help" || v.Shorthand == "-h" {
				fmt.Println(util.Intro())
				return nil
			}
		}
	}

	// Check if at least one file or directory is specified
	if len(rest) == 0 {
		fmt.Println("No file or directory specified.")
		return nil
	}

	// Process flags and arguments
	if err := c.process(rest, flags); err != nil {
		return err
	}

	return nil
}

// Check if the argument is a valid flag (by name or shorthand)
func (c *CMD) isFlag(arg string) (bool, Flag) {
	for _, v := range c.flags {
		if v.Name == arg || v.Shorthand == arg {
			return true, v
		}
	}
	return false, Flag{}
}

// Extract flags and non-flag arguments
func (c *CMD) extract(args []string) ([]string, []Flag) {
	var flags []Flag
	var rest []string

	for i := 0; i < len(args); i++ {
		v := args[i]

		// Handle timestamp (-t, --timestamp)
		if v == "-t" || v == "--timestamp" {
			if i+1 < len(args) { // Prevent out-of-bounds access
				flags = append(flags, Flag{Name: "--timestamp", Shorthand: "-t", Usage: "Set the creation timestamp for the file (YYYY-MM-DD HH:MM:SS)", Value: args[i+1]})
				i++ // Skip next value since it's part of this flag
			} else {
				fmt.Println("Error: Missing value for -t / --timestamp flag")
			}
			continue
		}

		// Handle permissions (-p, --permissions)
		if v == "-p" || v == "--permissions" {
			if i+1 < len(args) { // Prevent out-of-bounds access
				flags = append(flags, Flag{Name: "--permissions", Shorthand: "-p", Usage: "Set file permissions (USER or ADMIN)", Value: args[i+1]})
				i++ // Skip next value since it's part of this flag
			} else {
				fmt.Println("Error: Missing value for -p / --permissions flag")
			}
			continue
		}

		// Handle other flags using isFlag function
		isFlag, f := c.isFlag(v)
		if isFlag {
			flags = append(flags, f)
			continue
		}

		// If it's not a flag, add it to the rest of the arguments
		rest = append(rest, v)
	}

	return rest, flags
}

// Handle flag processing logic
func (c *CMD) process(rest []string, flags []Flag) error {
	for _, path := range rest {
		err := c.processPath(path, flags)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CMD) processPath(path string, flags []Flag) error {
	p := util.Path{Path: path}

	// Step 1: Parse flags into options
	options := map[string]struct {
		Index  int
		Active bool
	}{}

	for i, v := range flags {
		switch v.Name {
		case "--directory", "-d":
			options["directory"] = struct {
				Index  int
				Active bool
			}{Index: i, Active: true}
		case "--force", "-f":
			options["force"] = struct {
				Index  int
				Active bool
			}{Index: i, Active: true}
		case "--timestamp", "-t":
			options["timestamp"] = struct {
				Index  int
				Active bool
			}{Index: i, Active: true}
		case "--permissions", "-p":
			options["permissions"] = struct {
				Index  int
				Active bool
			}{Index: i, Active: true}
		case "--update", "-u":
			options["update"] = struct {
				Index  int
				Active bool
			}{Index: i, Active: true}
		case "--help", "-h":
			options["help"] = struct {
				Index  int
				Active bool
			}{Index: i, Active: true}
		default:
			fmt.Printf("Unknown flag: %s\n", v.Name)
		}
	}

	// Step 2: Apply logic in a structured way
	if options["directory"].Active {
		// If -d (directory) is set, -f (force) is ignored
		if err := p.CreateDir(); err != nil {
			return err
		}
	} else {
		if !options["update"].Active {

			// If -f (force) is set, replace existing file
			if options["force"].Active {
				if err := p.ForceCreate(); err != nil {
					return err
				}
			} else {
				if err := p.Create(); err != nil {
					return err
				}
			}

		}
		// Handle timestamp if -t is set
		if options["timestamp"].Active {
			// Example timestamp modification function (you need to implement this)
			if err := p.ModifyTimestamps(flags[options["timestamp"].Index].Value); err != nil {
				return err
			}
		}

		// Handle permissions if -p is set
		if options["permissions"].Active {
			// Example permission modification function (you need to implement this)
			if err := p.ModifyPermissions(flags[options["permissions"].Index].Value); err != nil {
				return err
			}
		}
	}

	return nil
}
