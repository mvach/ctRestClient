# ctRestClient

The ctRestClient allows exporting dynamic groups of [ChurchTools](https://church.tools/en/home/) as CSV files. This is required if one wants to send a serial letter using tools like Microsoft Word or Adobe InDesign to the group members.

## Configure it
The ctRestClient is mostly configured via a YAML file. Here is an example configuration:

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

This configuration contains multiple `token_name` properties. The value of each `token_name` is simply the name of an environment variable that is containing the relevant token. Therefore, such variables (here `TOKEN_HOST1` and `TOKEN_HOST2`) must be set in the environment before running the client.

```
# Set environment variables on Windows
set TOKEN_HOST1=6GN68NH4WM5OHRXDDEU1FHMMV...
set TOKEN_HOST2=QQO5N324JOWCO3RPOUTSSKXPY...

# Set environment variables on Mac
export TOKEN_HOST1=6GN68NH4WM5OHRXDDEU1FHMMV...
export TOKEN_HOST2=QQO5N324JOWCO3RPOUTSSKXPY...
```

## Run it
```
Usage:
  ctRestClient.exe -c <path to config.yml> -t <access token> 
  
  NOTE: on MacOs the executable would be ctRestClient-arm64 or ctRestClient-amd64

Options:
  -c            By default, ctRestClient tries to load the configuration from a 
                config.yml file that is searched beside the executable.
                Using -c allows to configure the path of the configuration YAML.

  -o            By default, ctRestClient creates an exports directory beside the
                executable. This exports directory is then filled with the CSV files.
                Using -o allows to configure the path where the exports directory 
                should be created.

  -k            By default, ctRestClient reads credentials from a Keepass passwords.kdbx file beside 
                the executable.
                Using -k allows to configure the path of the Keepass Database file.
```








## Develop it
The ctRestClient can be easily built from source using:
```
# Compile it for Windows
GOOS=windows GOARCH=amd64 go build -o ctRestClient.exe

# Compile it for MacOS with Apple silicon
GOOS=darwin GOARCH=arm64 go build -o ctRestClient-arm64

# Compile it for MacOS with Intel silicon
GOOS=darwin GOARCH=amd64 go build -o ctRestClient-amd64

# Compile it for Linux
GOOS=linux GOARCH=amd64 go build
```

Tests can be execute via:
```
go test ./...
```
