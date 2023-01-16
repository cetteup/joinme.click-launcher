# joinme.click-launcher

[![ci](https://img.shields.io/github/actions/workflow/status/cetteup/joinme.click-launcher/ci.yaml?label=ci)](https://github.com/cetteup/joinme.click-launcher/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/cetteup/joinme.click-launcher)](https://goreportcard.com/report/github.com/cetteup/joinme.click-launcher)
[![License](https://img.shields.io/github/license/cetteup/joinme.click-launcher)](/LICENSE)
[![Last commit](https://img.shields.io/github/last-commit/cetteup/joinme.click-launcher)](https://github.com/cetteup/joinme.click-launcher/commits/main)
[![Discord](https://img.shields.io/discord/1001891950544306259?label=Discord)](https://discord.gg/wwsuMk9g4E)

Launcher utility to handle custom game URL protocols supported on [joinme.click](https://joinme.click/).

![joinme click-launcher](https://user-images.githubusercontent.com/17167062/182002183-38b134c4-6749-4273-b08b-9147f80b7b81.gif)

## Supported games

| Game                             | URL protocol            | Minimum launcher version¹ | Supported mods                                                                                                              |
|----------------------------------|-------------------------|---------------------------|-----------------------------------------------------------------------------------------------------------------------------|
| Battlefield 1942                 | bf1942://{ip}:{port}    | v0.1.7-alpha              | `The Road to Rome`², `Secret Weapons of WWII`², `Battlefield 1918`, `Desert Combat (0.7)`, `Desert Combat Final`, `Pirates` |
| Battlefield Vietnam              | bfvietnam://{ip}:{port} | v0.1.7-alpha              | `Battlegroup 42`                                                                                                            |
| Battlefield 2                    | bf2://{ip}:{port}       | v0.1.11                   | `Special Forces`², `Allied Intent Xtended`, `Pirates (Yarr2)`, `Point of Existence 2`, `Arctic Warfare`                     |
| Battlefield 4                    | bf4://{gameid}          | v0.1.5-alpha              |
| Battlefield 1                    | bf1://{gameid}          | v0.1.5-alpha              |
| Call of Duty                     | cod://{ip}:{port}       | v0.1.3-alpha              |
| Call of Duty: United Offensive   | coduo://{ip}:{port}     | v0.1.3-alpha              |
| Call of Duty 2                   | cod2://{ip}:{port}      | v0.1.3-alpha              |
| Call of Duty 4: Modern Warfare   | cod4://{ip}:{port}      | v0.1.3-alpha              |
| Call of Duty: World at War       | codwaw://{ip}:{port}    | v0.1.3-alpha              |
| F.E.A.R. Combat (SEC2)           | fearsec2://{ip}:{port}  | v0.1.3-alpha              |
| ParaWorld                        | paraworld://{ip}:{port} | v0.1.7-alpha              |
| SWAT 4                           | swat4://{ip}:{port}     | v0.1.3-alpha              |
| SWAT 4: The Stetchkov Syndicate³ | swat4x://{ip}:{port}    | v0.1.3-alpha              |
| Unreal Tournament                | ut://{ip}:{port}        | v0.1.12                   |
| Unreal Tournament 2003           | ut2003://{ip}:{port}    | v0.1.12                   |
| Unreal Tournament 2004           | ut2004://{ip}:{port}    | v0.1.12                   |
| Vietcong                         | vietcong://{ip}:{port}  | v0.1.3-alpha              |

¹ refers to the minimum launcher version supporting all features relevant to the game

² these addons are considered mods for technical reasons, since they use the same game executable which is launched with different parameters

³ while technically an addon, it uses a separate game executable and is thus considered a different game

## Usage

### Registering URL handlers

Before you can launch games based on URLs, the launcher need to register as a URL handler for the supported URL protocols.
In order to do this, simply run the launcher once after download. It will check which of the supported games are
installed and register itself as a URL handler for each one it finds. After registering, the launcher shows the result
for each supported game.

```text
10: 37AM INF Checked status for game="Battlefield 1" result="launcher registered successfully"
10: 37AM INF Checked status for game="Battlefield 1942" result="not installed"
10: 37AM INF Checked status for game="Battlefield 2" result="launcher registered successfully"
10: 37AM INF Checked status for game="Battlefield 4" result="launcher registered successfully"
10: 37AM INF Checked status for game="Battlefield Vietnam" result="launcher registered successfully"
10: 37AM INF Checked status for game="Call of Duty" result="launcher registered successfully"
10: 37AM INF Checked status for game="Call of Duty 2" result="launcher registered successfully"
10: 37AM INF Checked status for game="Call of Duty 4: Modern Warfare" result="launcher registered successfully"
10: 37AM INF Checked status for game="Call of Duty: United Offensive" result="launcher registered successfully"
10: 37AM INF Checked status for game="Call of Duty: World at War" result="launcher registered successfully"
10: 37AM INF Checked status for game="F.E.A.R. Combat (SEC2)" result="not installed"
10: 37AM INF Checked status for game=ParaWorld result="launcher registered successfully"
10: 37AM INF Checked status for game="SWAT 4" result="launcher registered successfully"
10: 37AM INF Checked status for game="SWAT 4: The Stetchkov Syndicate" result="launcher registered successfully"
10: 37AM INF Checked status for game=Vietcong result="not installed"
10: 37AM INF Window will close in 15 seconds
```

### Launching a game based on a URL

No extra steps are required to launch a game based on one of the supported URL protocols. If you click a link
to [bf2://95.172.92.116:16567](bf2://95.172.92.116:16567) for example, the launcher will start Battlefield 2 and
join [2F4Y.com - Best Maps No Rules!](https://bf2.tv/servers/95.172.92.116:16567). If the game is already running and
cannot join a server from a running state, the launcher will close any existing game instance automatically before
launching a new one _(only supported by launcher v0.1.8-alpha and newer)_.

```text
10: 40AM INF Killing existing game process executable=BF2.exe pid=3916
10: 40AM INF Successfully launched game url=bf2://95.172.92.116:16567/
10: 40AM INF Window will close in 15 seconds
```

Depending on your browser and settings, you may need to confirm that you want to allow the launcher to start after
clicking the link.

![Browser URL protocol launch confirmation prompt](https://user-images.githubusercontent.com/17167062/179347704-8187a42a-9487-469e-b49c-fd56d8925136.png)

### Advanced configuration

You can customize some elements of how the launcher starts your games. For example, you can provide additional command line arguments on a per-game basis. The config needs to be placed in the same folder as the launcher executable as `config.yaml`.

#### General configuration options

| Option name     | Type    | Description                                                                  | Default value |
|-----------------|---------|------------------------------------------------------------------------------|---------------|
| `quiet_launch`  | boolean | do not leave the window open any longer than required                        | `false`       |
| `debug_logging` | boolean | show lots of information relevant for debugging any issues with the launcher | `false`       |

#### Per-game configuration options

These options can be configured (differently) on a per-game basis. They need to be placed in the `config.yaml` under `games` and then keyed by the game URL protocol (e.g. `bf2`). These options do not have any default values. Instead, they override dynamic values the launcher usually determines on its own (e.g. the game's install path). Some options also pass additional details to the launcher.  

| Option name       | Type     | Description                                                                                                               | 
|-------------------|----------|---------------------------------------------------------------------------------------------------------------------------|
| `executable_name` | string   | name of the game executable (usually statically defined per game)                                                         |
| `executable_path` | string   | relative path from the game's install path to folder containing the game executable (usually statically defined per game) |
| `install_path`    | string   | path where the game is installed (usually determined via the Windows registry)                                            |
| `args`            | string[] | array of additional arguments to pass the game when launching                                                             |

#### Example configuration

This example configuration would cause the launcher to not leave the launcher window open after performing any actions (meaning you will not see any output it printed). Also, Battlefield 2 would be launched in windowed mode with `C:\Games\Battlefield 2\bin\BF2.playbf2.exe` being started in `C:\Games\Battlefield 2`. Debug logging is disabled by default, meaning that option does not change any default behaviour.  

```yaml
quiet_launch: true
debug_logging: false
games:
    bf2:
        executable_name: BF2.playbf2.exe
        executable_path: bin
        install_path: C:\Games\Battlefield 2
        args: ["+fullscreen", "0", "+szx", "1600", "+szy", "900"]
```

You can also find the example configuration as a file: [config.example.yaml](config.example.yaml).

## Downloads

* https://joinme.click/download

License
-------

This is free software under the terms of the MIT license.
