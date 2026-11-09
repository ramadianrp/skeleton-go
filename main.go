package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ramadianrp/skeleton-go/cmd"
)

const banner = `
        __          .__          __                 
  _____|  | __ ____ |  |   _____/  |_  ____   ____  
 /  ___/  |/ // __ \|  | _/ __ \   __\/  _ \ /    \ 
 \___ \|    <\  ___/|  |_\  ___/|  | (  <_> )   |  \
/____  >__|_ \\___  >____/\___  >__|  \____/|___|  /
     \/     \/    \/          \/                 \/ 
`

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			fmt.Printf("error loading location '%s': %v\n", tz, err)
		} else {
			fmt.Printf("location loaded '%s'\n", tz)
		}
	}

	fmt.Print(banner)
	if err := cmd.Run(); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}
