package commands

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

// ArtInfo contains information about an ASCII art piece
type ArtInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// artCollection holds all available ASCII art
var artCollection = map[string]string{
	"fist": `        .+.
        .+?:
         .+??.
           ??? .
           +???.
      +?????????=.
      .???????????.
      .????????????.

     ########### ########
     ############.#######.
     ####### ####  .......
     ######## #### #######                50 41 4E 54 48 45 4F 4E
     #########.####.######        _____________  __  ________  ____  ______
     ######  ...                 /_  __/ __/ _ \/  |/  /  _/ |/ / / / / __/
     #######.??.##########        / / / _// , _/ /|_/ // //    / /_/ /\ \
     #######~+??.#########       /_/ /___/_/|_/_/  /_/___/_/|_/\____/___/
     ########.??..
     #########.??.#######.
     #########.+?? ######.
               .+?.
         .????????????.
           +??????????,
            .????++++++.
              ????.
              .???,
               .~??.
                 .??
                  .?,`,
	"hello": `Hello World!`,
	"rocket": `       !
       !
       ^
      / \
     /___\
    |=   =|
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
    |     |
   /|##!##|\
  / |##!##| \
 /  |##!##|  \
|  / ^ | ^ \  |
| /  ( | )  \ |
|/   ( | )   \|
    ((   ))
   ((  :  ))
   ((  :  ))
    ((   ))
     (( ))
      ( )
       .
       .
       .`,
	"unicorn": "       \\\n        \\\n         \\\\\n          \\\\\n           >\\/7\n       _.-(6'  \\\n      (=___._/` \\\n           )  \\ |\n          /   / |\n         /    > /\n        j    < _\\\n    _.-' :      ``.\n    \\ r=._\\        `.\n   <`\\\\_  \\         .`-.\n    \\ r-7  `-. ._  ' .  `\\\n     \\`,      `-.`7  7)   )\n      \\/         \\|  \\'  / `-._\n                 ||    .'\n                  \\\\  (\n                   >\\  >\n               ,.-' >.'\n              <.'_.''\n                <'",
	"wordpress": `                                 ..............
                           ..........................
                      ..'...                        ...'..
                  ..'..          ..............          ..'..
               ..'..       ...''''''''''''''''''''...        .'..
             ....      ..'''''''''''''''''''''''''''''''..     ..'.
           .'.     ..''''''''''''''''''''''''''''''''''''''..     .'.
         .'.     .'''''''''''''''''''''''''''''''''''''''''''.      .'.
        ..     .''''''''''''''''''''''''''''''''''''''''''.           ...
      .'.    .'''''''''''''''''''''''''''''''''''''''''''.             .'.
     .'.    ............'''''.....................'''''''                '.
    .'                  .'''.                    .'''''''                 '.
   .'              ..'''''''''''..           ..''''''''''.            .    '.
  .'.   .           ''''''''''''''.          .''''''''''''.           ..   .'.
  '.    '.          .''''''''''''''.          .'''''''''''''.         .'.   .'
 .'    .''.          .''''''''''''''           .'''''''''''''.        '''    '.
 '.   .''''           .'''''''''''''.           ''''''''''''''        '''.   .'
.'.   .''''.           ''''''''''''''.          .'''''''''''''.      .'''.    '.
.'    ''''''.          .''''''''''''''.          .'''''''''''''     .'''''    '.
.'    '''''''.          .'''''''''''''.           '''''''''''''     ''''''    ''
.'    ''''''''           '''''''''''''            .'''''''''''.    .''''''    ''
.'    ''''''''.          .'''''''''''.             .''''''''''.   .'''''''    '.
.'.   .''''''''.          .'''''''''.   '.          '''''''''.    '''''''.    '.
 '.   .'''''''''.          '''''''''   .''          .''''''''    .'''''''.   .'
 .'    .'''''''''          .'''''''.  .'''.          .''''''.   .''''''''    '.
  '.    ''''''''''          .'''''.   '''''.          '''''.   .''''''''.   .'
  .'.   .'''''''''.          '''''   .''''''.         .''''.   ''''''''.   .'.
   .'.   .'''''''''.         .'''.  .''''''''          .''.   .'''''''.    '.
    .'    .'''''''''          .'.  .'''''''''.          .'   .'''''''.    '.
     .'.    .'''''''.          .   '''''''''''.         ..   ''''''.     '.
      .'.    .'''''''.            .''''''''''''.            .'''''.    .'.
        .'.    .''''''.          .''''''''''''''           .''''.     ..
         .'.     .'''''          '''''''''''''''.         .'''.     .'.
           .'.     ..''.        .''''''''''''''''.        '..     .'.
             ....      ..      .''''''''''''''''''.            ....
                .'..           '''''''''''''''''''.         ..'..
                   ....           ............           ..'.
                      ......                        ...'..
                           ..........................
                                 ......''......`,
	"druplicon": `                                ..
                                 ld'.
                                .XK;,.
                               'KMX:;;,.
                             .xWWXo;;;;;,'.
                          .l0WXKOc;;;;;;;;;,'..
                      .;dKMMMMXOOl;;;;;;;;;;;;;,,'.
                   ,d0WMMMMMMNOOl;;;;;;;;;;;;;;;;;;;,'.
                ;kNMMMMMMMMN0Od:;;;;;;;;;;;;;;;;;;;;;;;,'.
             .xNMMMMMMMMMN0Ox:;;;;;;;;;;;;;;;;;;;;;;;;;;;;,.
           ;0MMMMMMMMMNKOOd:;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,.
         ,0MMMMMMMMNK0Oxl;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,.
       .dMMMMMMWXKOOxl:;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;.
      'xNWWNXK0Okdl;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,'.
     'dOOOOOkdl:;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;'''.
    ';:lllc;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;''''.
   .;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,'''''.
  .,;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,'''''''.
  ,;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;''''''''.
 .;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,''''''''''.
 ';;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,'''''''''''.
 ,;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,''''''''''''''
 ,;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;,''''''''''''''''
 ,;;;;;;;;;;;;;;;;;;;;;lk0XNWWMWNKOdc;;;;;;;;;;;;;;;;,'''''''';lol:'''''
 ,;;;;;;;;;;;;;;;;;;l0WMMMMMMMMMMMMMMNOo;;;;;;;;;;,,'''''''lONMMMMMNo'''
 ';;;;;;;;;;;;;;;;cXMMMMMMMMMMMMMMMMMMMMW0d;;;;;,'''''''lOWMMMMMMMMMMl''
 .,;;;;;;;;;;;;;;oWMMMMMMMMMMMMMMMMMMMMMMMMW0o,''''',lOWMMMMMMMMMMMMMO,.
  ';;;;;;;;;;;;;cWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMXOk0XMMMMMMMMMMMMMMMMMO'
  .,;;;;;;;;;;;;kMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWXKNMMMMMMMMMMMMMMMMMMo.
   .,;;;;;;;;;;;kMMMMMMMMMMMMMMMMMMMMMMMMMMWOo;'''';oXMMMMMMMMMMMMMMX.
    .,;;;;;;;;;;cWMMMMMMMMMMMMMMMMMMMMMMXxc'''''''''''c0MMMMMMMMMMMN,
     .,;;;;;;;;;;oWMMMMMMMMMMMMMMMMMWOo;''',ldkkkkd;'''':0MMMMMMMMX,
       .',,,,;;;;,;d0WMMMMMMMMMMNOd:'''''l0WX0kxkONMx''''':OWMMMNd.
        .''''''''''''':loooool:,''''''''OWO:'''''''ON:'''''',clc.
          .''''''''''''''''''''''''''''''''''''''''''''''''',,'
            .'''''''''''''''''''';c'''''''''''''''''':o''',,.
              ..'''''''''''''''';WMNOdl:;;,,,,;:cox0NNk,,'.
                 ..'''''''''''''',oOXWMMMMMMMMMMNKko:,'.
                     ...'''''''''''''',;;:::;;,,''..
                          .....''''''''''''....`,
	"metal": `

        _                             _
       | |                           | |
       | |                           | |
       | |                           | |
       | |_                         _| |
       |   |_______________________|   |
        \                             /
         \      ___         ___      /
          |    |   |       |   |    |
          |    |   |       |   |    |
          |    |   |       |   |    |
          |    |___|       |___|    |
          |                         |
           \_______________________/

                    \m/ \m/                    `,
}

// artDescriptions provides descriptions for each art piece
var artDescriptions = map[string]string{
	"fist":      "The Pantheon fist logo",
	"hello":     "A friendly greeting",
	"rocket":    "A rocket ship blasting off",
	"unicorn":   "A magical unicorn",
	"wordpress": "The WordPress logo",
	"druplicon": "The Drupal druplicon mascot",
	"metal":     "The sign of the horns",
}

var artCmd = &cobra.Command{
	Use:   "art [name]",
	Short: "Display Pantheon ASCII art",
	Long: `Display one of several ASCII art images.

If no name is specified, a random art piece will be displayed.
Use 'terminus art:list' to see all available artwork.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runArt,
}

var artListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available ASCII art",
	Long:  "Display a list of all available ASCII artwork that can be displayed with the 'art' command.",
	RunE:  runArtList,
}

func init() {
	artCmd.AddCommand(artListCmd)
}

func runArt(_ *cobra.Command, args []string) error {
	var name string

	if len(args) == 0 {
		// Select random art
		name = getRandomArtName()
	} else {
		name = args[0]
	}

	art, exists := artCollection[name]
	if !exists {
		return fmt.Errorf("art '%s' not found. Use 'terminus art:list' to see available artwork", name)
	}

	// Print the art directly (not through printOutput since it's raw text)
	fmt.Println(art)
	return nil
}

func runArtList(_ *cobra.Command, _ []string) error {
	// Build list of art info
	artList := make([]ArtInfo, 0, len(artCollection))

	for name := range artCollection {
		artList = append(artList, ArtInfo{
			Name:        name,
			Description: artDescriptions[name],
		})
	}

	// Sort by name for consistent output
	sort.Slice(artList, func(i, j int) bool {
		return artList[i].Name < artList[j].Name
	})

	return printOutput(artList)
}

// getRandomArtName returns a random art name from the collection
func getRandomArtName() string {
	names := make([]string, 0, len(artCollection))
	for name := range artCollection {
		names = append(names, name)
	}

	//nolint:gosec // This is not security-sensitive - just selecting random art
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return names[r.Intn(len(names))]
}
