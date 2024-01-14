# ffcu

Update your firefox with remote themes and user.js config with one command

## Usage

> [!WARNING]  
> This tool will override your user.js and chrome directory with the downloaded ones

1. `ffcu config set-chrome-zip "the url that has the latest zip which contains the chrome dir"`
2. `ffcu config set-userjs-url "the url that has the latest user.js file"`
3. `ffcu config set-profile-dir "absolute path to the firefox profile directory"`
4. `ffcu update`

## Installing

```
go install github.com/Tronikelis/ffcu
```

## Config

The ffcu config is stored in `$HOME/.ffcu/config.json`
