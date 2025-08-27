# Absicht

`absicht` is a modern mail composer TUI. Write mails in your `$EDITOR` and send
them through `mstmp`.

![A demo showing how to view and edit a saved file](demo.gif)

## Usage

    absicht [flags]

    Flags:
      -e, --edit              Start editing the email right away
      -f, --file string       Read initial email from this path; `-' means stdin (default "-")
      -h, --help              help for absicht
      -s, --sendmail string   Command to send mail; mail will be piped to stdin (default "msmtp -t --read-envelope-from")

## Build

    go build    # if you have go installed
    nix build   # if you have nix installed
