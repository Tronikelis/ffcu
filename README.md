# ffcu

Update your firefox with remote themes and remote user.js config + custom local overrides with one command

```
NAME:
   ffcu - Helper CLI tool to auto update your firefox config

USAGE:
   ffcu [global options] command [command options]

COMMANDS:
   config   Commands related to ffcu configuration
   update   Kills firefox and updates it with the latest downloaded files
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

![ffcu](https://github.com/Tronikelis/ffcu/assets/56039679/8bbf9241-4776-4e96-b7ec-9c1571d8be21)

## Installing

```
go install github.com/Tronikelis/ffcu
```

## Usage

> [!WARNING]  
> This tool will override your user.js and chrome directory with the downloaded ones

1. `ffcu config set-chrome-zip-url "the url that has the latest zip which contains the chrome dir"`
2. `ffcu config set-userjs-url "the url that has the latest user.js file"`
3. `ffcu config set-profile-dir "absolute path to the firefox profile directory"`
4. `ffcu update`

## Config

The ffcu config is stored in `$HOME/.ffcu/config.json`

## Overriding user.js

Create a file in the root of your firefox profile directory named `user.overrides.js`, once the new `user.js` is written,
that file will be appended at the end of the new `user.js`

## Example config

```
$ ffcu config print

{
    "ProfileDir": "C:\\Users\\tronikel\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles\\1coaibnj.default-release",
    "UserJsUrl": "https://raw.githubusercontent.com/yokoffing/Betterfox/main/user.js",
    "ZippedChromeUrl": "https://codeload.github.com/bmFtZQ/edge-frfox/zip/refs/heads/main"
}
```
