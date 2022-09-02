package runner

import "github.com/projectdiscovery/gologger"

const banner = `
                __                        
 _      _____  / /_  ______________ _____ 
| | /| / / _ \/ __ \/ ___/ ___/ __  / __ \
| |/ |/ /  __/ /_/ (__  ) /__/ /_/ / / / /
|__/|__/\___/_.___/____/\___/\__,_/_/ /_/

	` + Version + ` by zp857
`

const Version = `v2.0`

func showBanner() {
	gologger.Print().Msgf("%v\n", banner)
}
