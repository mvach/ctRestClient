# ctRestClient - Benutzerhandbuch

## Überblick

`ctRestClient` ist ein Kommandozeilen-Tool, das den Export von Gruppen aus [ChurchTools](https://church.tools/de/home/) als CSV-Dateien ermöglicht. Dies ist besonders nützlich für die Erstellung von Serienbriefen mit Tools wie Microsoft Word oder Adobe InDesign.

## Hauptfunktionen

- **Export von ChurchTools-Gruppen**: Exportiert Mitgliederdaten aus dynamischen Gruppen
- **CSV-Format**: Ausgabe in standardisierten CSV-Dateien
- **Sichere Token-Verwaltung**: Verwendung von KeePass-Datenbanken für API-Token
- **Multi-Instanz-Unterstützung**: Gleichzeitiger Export von mehreren ChurchTools-Instanzen
- **Konfigurierbare Datenfelder**: Flexible Auswahl der zu exportierenden Felder
- **Plattformübergreifend**: Verfügbar für Windows, macOS und Linux

## Voraussetzungen

### Software-Anforderungen

1. **KeePassXC**: KeePassXC muss installiert sein und das Kommandozeilen-Tool `keepassxc-cli` muss im System-PATH verfügbar sein
   - **Windows**: `keepassxc-cli.exe` befindet sich meist unter `Program Files/KeePassXC`
   - **macOS**: `keepassxc-cli` befindet sich meist unter `/Applications/KeePassXC.app/Contents/MacOS`
   - **Linux**: `keepassxc-cli` befindet sich meist unter `/usr/bin` oder `/usr/local/bin`
   
   **Wichtig**: Ohne `keepassxc-cli` im PATH kann ctRestClient nicht auf die Token-Datenbank zugreifen!

2. **ChurchTools-Zugang**: 
   - Gültiger API-Token für jede ChurchTools-Instanz
   - Berechtigung zum Lesen der entsprechenden Gruppen

### Hardware-Anforderungen

- Minimal: Beliebiges System mit ausreichend Speicher für die zu exportierenden Daten

## Installation

### 1. Executable herunterladen

Laden Sie die Software von [GitHub](https://github.com/mvach/ctRestClientReleaseTest/releases/latest) herunter.

Das herunterzuladende Archiv trägt den Namen ctRestClient_<version\>.tar.gz. In diesem Archiv befinden sich die ausführbaren Binaries für Windows, Linux, und MacOs.


### 2. KeePass-Datenbank einrichten

Erstellen Sie eine KeePass-Datenbank (`.kdbx`-Datei) und speichern Sie Ihre ChurchTools-API-Token:

1. Öffnen Sie KeePassXC
2. Erstellen Sie eine neue Datenbank
3. Für jede ChurchTools-Instanz erstellen Sie einen Eintrag:
   - **Titel**: Ein eindeutiger Name (z.B. `meineKirche`)
   - **Passwort**: Ihr ChurchTools-API-Token
4. Speichern Sie die Datenbank als `churchtools-tokens.kdbx` (empfohlener Name)

### 3. Ordnerstruktur

Erstellen Sie folgende Ordnerstruktur:
```
mein-projekt/
├── ctRestClient-[platform]        # Die Executable
├── config.yml                     # Konfigurationsdatei
├── churchtools-tokens.kdbx        # KeePass-Datenbank
├── data/                          # ggf. Daten für Wertumwandlungen
└── exports/                       # Ausgabeverzeichnis (wird automatisch erstellt)
```

## Konfiguration

### Grundlegende Konfiguration

Die Konfiguration erfolgt über eine YAML-Datei. Hier ein Beispiel (`config.yml`):

```yaml
instances:
  - hostname: meineKirche.de
    token_name: meineKirche
    groups:
    - name: Konfirmanden
      fields: [id, firstName, lastName, sexId, street, zip, city]
    - name: Eltern von Konfirmanden
      fields: [id, firstName, lastName, sexId, email]
```

### Konfigurationsparameter

#### Instanzen (`instances`)
- **hostname**: Die Domäne Ihrer ChurchTools-Instanz (ohne https://)
- **token_name**: Name des Token-Eintrags in der KeePass-Datenbank
- **groups**: Liste der zu exportierenden Gruppen

#### Gruppen (`groups`)
- **name**: Exakter Name der Gruppe in ChurchTools
- **fields**: Liste der zu exportierenden Datenfelder

### Erweiterte Feldkonfiguration

#### Wertumwandlung mit benutzerdefinierten Spaltennamen

`ctRestClient` unterstützt die automatische Umwandlung von ChurchTools-Werten in benutzerfreundlichere Ausgaben. Gleichzeitig können Sie dabei auch die Spaltennamen in der CSV-Ausgabe anpassen.

**Funktionsweise:**
- Werte werden aus YAML-Dateien im `data/`-Verzeichnis gelesen
- Der Pfad folgt dem Schema: `data/<modul>/<feldname>.yml`
- Für das Feld `sexId` im Modul `persons`: `data/persons/sexId.yml`
- Spaltennamen können mit `columnname` angepasst werden

**Beispiel für die Umwandlung der sexId-Werte:**

Erstellen Sie die Datei `data/persons/sexId.yml`:
```yaml
1: männlich
2: weiblich
3: divers
```

**Konfiguration:**
```yaml
fields: 
  - id
  - firstName
  - lastName
  - {fieldname: sexId, columnname: "geschlecht"}
  - birthday
```

### Beispielkonfigurationen

#### Geburtstagslisten

```yaml
instances:
  - hostname: meineKirche.de
    token_name: meineKirche
    groups:
    - name: Konfirmanden
      fields: [id, firstName, lastName, sexId, street, zip, city]
    - name: Eltern von Konfirmanden
      fields: [id, firstName, lastName, sexId, email]
```

#### Multi-Instanz-Setup

```yaml
instances:
  - hostname: kumulus.meineKirche.de
    token_name: kumulus
    groups:
    - name: Kindergottesdienst
      fields: [id, firstName, lastName]
  - hostname: stratus.meineKirche.de
    token_name: stratus
    groups:
    - name: Jugendgruppe
      fields: [id, firstName, lastName, email]
```

## Verwendung

### Kommandozeilen-Parameter

```bash
ctRestClient [OPTIONS]
```

#### Verfügbare Optionen:

- **`-c <pfad>`**: Pfad zur Konfigurationsdatei
  - Standard: `config.yml` im Verzeichnis der Executable
- **`-k <pfad>`**: Pfad zur KeePass-Datenbank
  - Standard: `passwords.kdbx` im Verzeichnis der Executable
- **`-o <pfad>`**: Ausgabeverzeichnis für CSV-Dateien
  - Standard: `exports/` im Verzeichnis der Executable
- **`-d <pfad>`**: Pfad zum Datenverzeichnis für Wertumwandlungen
  - Standard: `data/` im Verzeichnis der Executable

### Grundlegende Ausführung

#### Windows
```cmd
ctRestClient-windows-amd64.exe
```

#### macOS
```bash
./ctRestClient-darwin-arm64  # Für Apple Silicon
./ctRestClient-darwin-amd64  # Für Intel-Macs
```

#### Linux
```bash
./ctRestClient-linux-amd64
```

### Ausführung mit benutzerdefinierten Pfaden

```bash
# Mit spezifischer Konfiguration
./ctRestClient-linux-amd64 -c /pfad/zu/meiner/config.yml

# Mit spezifischer KeePass-Datenbank
./ctRestClient-linux-amd64 -k /pfad/zu/tokens.kdbx

# Mit spezifischem Ausgabeverzeichnis
./ctRestClient-linux-amd64 -o /pfad/zu/exports/

# Mit spezifischem Datenverzeichnis
./ctRestClient-linux-amd64 -d /pfad/zu/data/

# Kombination aller Parameter
./ctRestClient-linux-amd64 -c config.yml -k tokens.kdbx -o exports/ -d data/
```

### Automatisierung mit Skripten

#### Bash-Skript für Linux/macOS (`export.sh`)

```bash
#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Plattform erkennen
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
    "Darwin")
        export PATH=$PATH:/Applications/KeePassXC.app/Contents/MacOS
        case "$ARCH" in
            "arm64")
                EXECUTABLE="$SCRIPT_DIR/ctRestClient-darwin-arm64"
                ;;
            "x86_64")
                EXECUTABLE="$SCRIPT_DIR/ctRestClient-darwin-amd64"
                ;;
        esac
        ;;
    "Linux")
        export PATH=$PATH:/usr/bin:/usr/local/bin
        EXECUTABLE="$SCRIPT_DIR/ctRestClient-linux-amd64"
        ;;
esac

# Ausführung
chmod +x "$EXECUTABLE"
"$EXECUTABLE" -k "$HOME/churchtools-tokens.kdbx" -c "$SCRIPT_DIR/config.yml"
```

#### Batch-Skript für Windows (`export.bat`)

```batch
@echo off
set SCRIPT_DIR=%~dp0
set EXECUTABLE=%SCRIPT_DIR%ctRestClient-windows-amd64.exe

"%EXECUTABLE%" -k "%USERPROFILE%\churchtools-tokens.kdbx" -c "%SCRIPT_DIR%config.yml"
pause
```

## Ausgabe

### Dateistruktur

Der Export erstellt folgende Struktur:

```
exports/
└── [DATUM]_[ZEIT]/
    ├── ctRestClient.log
    └── [HOSTNAME]/
        ├── [Gruppenname_1].csv
        ├── [Gruppenname_2].csv
        └── ...
```

Beispiel:
```
exports/
└── 2025.08.06_14-30-15/
    ├── ctRestClient.log
    └── ihre-kirche.krz.tools/
        ├── Konfirmanden.csv
        └── Eltern_von_Konfirmanden.csv
```

### CSV-Format

Die CSV-Dateien verwenden:
- **Trennzeichen**: Semikolon (`;`)
- **Kodierung**: UTF-16 Little Endian mit BOM
- **Erste Zeile**: Spaltenüberschriften

Beispiel-Inhalt:
```csv
id;firstName;lastName;street;zip;city
123;Max;Mustermann;Musterstraße 1;12345;Musterstadt
124;Maria;Musterfrau;Beispielweg 2;54321;Beispielort
```

## Logging

Detaillierte Informationen über die Ausführung finden Sie in der `ctRestClient.log`-Datei im Ausgabeverzeichnis. Diese enthält:
- Zeitstempel aller Aktionen
- Erfolgreiche Exports
- Fehlermeldungen mit Details
- Performance-Informationen

## Best Practices

### Sicherheit

1. **KeePass-Datenbank schützen**: Verwenden Sie ein starkes Master-Passwort
2. **API-Token rotieren**: Erneuern Sie regelmäßig Ihre ChurchTools-Token
3. **KeePassXC aktuell halten**: Installieren Sie regelmäßig Updates für KeePassXC, um Sicherheitslücken zu schließen
4. **Dateiberechtigungen**: Beschränken Sie den Zugriff auf Konfigurationsdateien
5. **Sichere Übertragung**: Stellen Sie sicher, dass ChurchTools über HTTPS erreichbar ist

### Performance

1. **Feldauswahl optimieren**: Exportieren Sie nur benötigte Felder
2. **Gruppengröße beachten**: Bei sehr großen Gruppen kann der Export länger dauern
3. **Netzwerkverbindung**: Verwenden Sie eine stabile Internetverbindung

### Organisation

1. **Sinnvolle Dateinamen**: Verwenden Sie aussagekräftige Konfigurationsdateinamen
2. **Konfigurations-Backups**: Erstellen Sie regelmäßig Backups Ihrer Konfigurationsdateien
3. **Dokumentation**: Dokumentieren Sie spezielle Konfigurationen für Ihr Team
4. **Automatisierung**: Nutzen Sie Skripte für wiederkehrende Exporte

---

### Lizenz

Dieses Projekt steht unter der MIT-Lizenz. Siehe [LICENSE](LICENSE.md)-Datei für Details.
