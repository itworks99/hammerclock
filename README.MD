# Hammerclock

A simple terminal-based [chess clock](https://en.wikipedia.org/wiki/Chess_clock) and phase tracker for tabletop games. It is made with Warhammer40k/Age of Sigmar type of games in mind, but it can be used for any game that requires a timer and phase tracking.

# Running

By default the application searches for the options file in the current directory. You can specify a different configuration file by using the `-c` flag.

```bash
./hammerclock -o /path/to/config.json
```

# Configuration

The configuration file is a JSON file that contains the following fields:
```json
{
"name": "W40K 10th Edition",
"default": true,
"playerCount": 2,
"playerNames": [
"Player 1",
"Player 2"
],
"phases": [
"Command Phase",
"Movement Phase",
"Shooting Phase",
"Charge Phase",
"Fight Phase",
"End Phase"
],
"oneTurnForAllPlayers": false,
"colorPalette": "k9s",
"timeFormat": "AMPM"
}
```
Where:
- `name`: The name of the game. This is used to identify the configuration file.
- `default`: If true, this configuration file will be used as the default configuration file.
- `playerCount`: The number of players in the game. This is used to determine the number of players in the UI.
- `playerNames`: The names of the players. This is used to display the player names in the UI. The game supports infinite amount of players and will divide the screen space proportionally. The amount of player names must match the `playerCount` field.
- `phases`: The phases of the game. This is used to display the phases in the UI.
- `oneTurnForAllPlayers`: If true, all players will have just one turn skipping the phases.
- `colorPalette`: The color palette to use. This is used to display the colors in the UI. The default is `k9s`. Supported values are `k9s`, `dracula`, `monokai`