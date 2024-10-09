# ctRestClient

The ctRestClient allows exporting dynamic groups of ChurchTools as CSV files. This is required if one wants to send a serial letter using tools like Microsoft Word or Adobe InDesign to the group members.

## Configure it
The ctRestClient is configured via a YAML file. Here is an example configuration:

```
instances:
  - hostname: host1.yourchurch.de
    groups:
    - name: school kids
      fields: [id, firstName, lastName, street, zip]
    - name: new parents
      fields: [id, firstName, lastName, sexId]
    - name: newlyweds
      fields: [id, firstName, lastName, street, zip, sexId]
  - hostname: host2.yourchurch.de
    groups:
    - name: teenagers
      fields: [id, firstName, lastName, street, zip]
```

## Run it
```
Usage:
  ctRestClient.exe -c <path to config.yml> -t <access token> 
  
  NOTE: on MacOs the executable would be ctRestClient-arm64 or ctRestClient-amd64

Options:
  -c            By default, ctRestClient tries to load the configuration from a 
                config.yml file that is searched beside the executable.
                Using -c allows you to configure the path of the configuration YAML.

  -o            By default, ctRestClient creates an exports directory beside the
                executable. This exports directory is then filled with the CSV files.
                Using -o allows you to configure the path where the exports directory 
                should be created.

  -t            ctRestClient requires a token to access a ChurchTools instance.
```








## Develop it
The ctRestClient can be easily built from source using:
```
# Compile it for Windows
GOOS=windows GOARCH=amd64 go build

# Compile it for MacOS with Apple silicon
GOOS=darwin GOARCH=arm64 go build -o ctRestClient-arm64

# Compile it for MacOS with Intel silicon
GOOS=darwin GOARCH=amd64 go build -o ctRestClient-amd64
```

Tests can be execute via:
```
go test ./...
```