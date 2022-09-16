package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/vikpe/automapshot/internal/pkg/mapsettings"
	"github.com/vikpe/automapshot/internal/pkg/mapshot"
	"github.com/vikpe/go-ezquake"
	"github.com/vikpe/prettyfmt"
)

func NewApp() *cli.App {
	cli.AppHelpTemplate = `{{.Name}} [{{.Version}}]
{{.Description}}

  Usage:   {{.UsageText}}
Example:   {{.Name}} dm2 dm6
`

	return &cli.App{
		Name:        "automapshot",
		Description: "Automate screenshots of QuakeWorld maps.",
		UsageText:   "automapshot [<maps> ...]",
		Version:     "__VERSION__", // updated during build workflow
		Action: func(c *cli.Context) error {
			mapSettings, err := mapsettings.FromJsonFile("map_settings.json")

			if err != nil {
				return err
			}

			maps := c.Args().Slice()

			if len(maps) == 1 && "all" == maps[0] {
				maps = mapSettings.MapNames()
			}

			pfmt := prettyfmt.New("mapshot", color.FgHiCyan, "15:04:05", color.FgWhite)
			client := mapshot.NewClient(getEzquakeController())

			for _, mapName := range maps {
				if !mapSettings.HasMap(mapName) {
					pfmt.Printfln(`%s (skip, no settings defined)`, mapName)
					continue
				}

				client.Mapshot(mapName, mapSettings[mapName])
				pfmt.Printfln(`%s (success)`, mapName)
			}

			return nil
		},
		Before: func(context *cli.Context) error {
			return validateSetup()
		},
	}
}

func validateSetup() error {
	err := godotenv.Load()

	if err != nil {
		return errors.New("unable to load environment variables. create .env (see .env.example)")
	}

	ctrl := getEzquakeController()

	if !ctrl.Process.IsStarted() {
		return errors.New(fmt.Sprintf("ezQuake is not started (%s)", ctrl.Process.Path))
	}

	return nil
}

func getEzquakeController() *ezquake.ClientController {
	return ezquake.NewClientController(
		os.Getenv("EZQUAKE_PROCESS_USERNAME"),
		os.Getenv("EZQUAKE_BIN_PATH"),
	)
}
