# Metal Detector

One of your duo partners too heavy to carry? Figure out who it is with Metal Detector! Leverages the Riot API through [Golio](github.com/KnutZuidema/golio) to pull your match history and analyze your team's players.

## Installation

*Prerequisite - [Install Golang](https://go.dev/)*

1. Run `git clone https://github.com/Chickaloo/MetalDetector.git`
2. `cd` into the repo and run `go build`
3. Run the resulting executable with flags provided below

**Flags**

- `-s` - Summoner name
- `-n` - Number of games played
- `-t` - Threshold under which people won't be counted as your duo
- `-a` - API Key from developer.riotgames.com

**Example**
`./MetalDetector.exe -s Chickaloo -n 561 -t 3`

No warranty provided. Read code at your own peril.
