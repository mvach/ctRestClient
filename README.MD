# ctRestClient

`ctRestClient` allows exporting dynamic groups from [ChurchTools](https://church.tools/en/home/) as CSV files. This is useful when you want to send serial letters to group members using tools like Microsoft Word or Adobe InDesign.

## Configuration
`ctRestClient` is primarily configured via a YAML file. Below is an example configuration:

```
instances:
  - hostname: host1.yourchurch.de
    token_name: TOKEN_HOST1
    groups:
    - name: school kids
      fields: [id, firstName, lastName, street, zip]
    - name: new parents
      fields: [id, firstName, lastName, sexId]
    - name: newlyweds
      fields: [id, firstName, lastName, street, zip, sexId]
  - hostname: host2.yourchurch.de
    token_name: TOKEN_HOST2
    groups:
    - name: teenagers
      fields: [id, firstName, lastName, street, zip]
```

This configuration includes multiple `token_name` entries. The corresponding tokens must be stored in a KeePass database. The tokens are retrieved using the provided names (e.g.,`TOKEN_HOST1`, `TOKEN_HOST2`)

## Execution

`ctRestClient` uses **KeePassXC** as user token store. Therefore, `keepassxc-cli` must be available in your system's `PATH`. Without it, the token store cannot be accessed.

```
Usage:
  ctRestClient.exe -c <path to config.yml> -t <access token> 
  
  NOTE: on MacOs the executable would be ctRestClient-arm64 or ctRestClient-amd64

Options:
  -c      Specify the path to the configuration YAML file.
          By default, ctRestClient looks for config.yml in the same directory as the executable.

  -o      Specify the output path for the CSV exports directory.
          By default, the exports directory is created next to the executable.

  -k      Specify the path to the KeePass database file.
          By default, it looks for passwords.kdbx in the same directory as the executable.
```

## Development
You can easily build `ctRestClient` from source:
```
./scripts/build_binaries
```

Tests can be execute via:
```
go test ./...
```
