package commands

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/leocov-dev/go-scad-render/config"
	"github.com/leocov-dev/go-scad-render/internal/scad"
	"github.com/leocov-dev/go-scad-render/internal/utils"
	"github.com/leocov-dev/go-scad-render/internal/utils/spinners"
	"github.com/spf13/cobra"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	imgCmd = &cobra.Command{
		Use:   "img <directory>",
		Short: "Render images for scad files",
		Run:   run,
		Args:  imgArgs(),
	}

	precision uint8
)

func init() {
	imgCmd.Flags().Uint8Var(&precision, "Precision", 24, "Set $fn value")
	rootCmd.AddCommand(imgCmd)
}

func imgArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("<directory> argument missing")
		}
		return nil
	}
}

func run(cmd *cobra.Command, args []string) {

	scanDirectory := filepath.Clean(args[0])

	s := spinners.StartNewSpinner(fmt.Sprintf("Scanning for OpenScad files in: %s...", config.ExecDir))
	files, err := scad.GatherFiles(scanDirectory)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	// Only render when no existing png file
	filtered, err := files.Filter(func(f *scad.File) (bool, error) {
		imagePath := fmt.Sprintf("%s%s.png", f.Dir, f.Filename)
		return !utils.FileExists(imagePath), nil
	})

	// Skip files in "lib_" prefix
	filtered, err = filtered.Filter(func(f *scad.File) (bool, error) {
		return !strings.HasPrefix(f.Filename, "lib_"), nil
	})

	count := len(filtered)

	utils.WriteMain("Rendering SCAD files to images", strconv.Itoa(count))

	if count > 0 {
		uiprogress.Start()
		pb := uiprogress.
			AddBar(count).
			AppendCompleted().
			PrependElapsed().
			PrependFunc(func(b *uiprogress.Bar) string {
				return fmt.Sprintf("Rendering (%d/%d)", b.Current(), count)
			})

		errors, err := scad.RenderImages(filtered, 512, 24, pb)

		uiprogress.Stop()

		if err != nil {
			utils.CmdFailed(cmd, err)
		}

		utils.PrintErrorList(errors)
	}

	fmt.Println("Done")
}
