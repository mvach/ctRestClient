var userHome = getUserHome();

if (typeof EXPORTS_DIR === 'undefined') {
    var EXPORTS_DIR = userHome + "/exports";
}
if (typeof TEMPLATES_DIR === 'undefined') {
    var TEMPLATES_DIR = userHome + "/templates";
}
if (typeof OUTPUT_DIR === 'undefined') {
    var OUTPUT_DIR = userHome + "/serial_letters";
}

var LOG_FILE;

/**
 * Opens an InDesign template file without user interaction
 * @param {File} templateFile - The template file to open
 * @returns {Document} The opened document
 */
function loadTemplate(templateFile) {
    app.scriptPreferences.userInteractionLevel = UserInteractionLevels.NEVER_INTERACT;
    var doc = app.open(templateFile);
    app.scriptPreferences.userInteractionLevel = UserInteractionLevels.INTERACT_WITH_ALL;
    return doc;
}

/**
 * Checks if a folder exists and raises an exception if not
 * @param {Folder} folder - The folder to check
 */
function assertFolder(folder) {
    if (!folder.exists) {
        throw new Error("The folder \"" + folder.fsName + "\" does not exist.");
    }
}

/**
 * Ensures a folder exists by creating it if it doesn't exist
 * @param {Folder} folder - The folder to ensure exists
 */
function ensureFolder(folder) {
    if (!folder.exists) {
        folder.create();
    }
}

/**
 * Append a line to the persistent log file.
 * @param {string} message
 */
function writeLogFile(message) {
    try {
        LOG_FILE.encoding = 'UTF-8';
        if (LOG_FILE.open('a')) {
            var ts = new Date();
            // ExtendScript Date may not implement toISOString(), so format manually
            var pad = function (n) { return (n < 10 ? '0' : '') + n; };
            var iso = ts.getFullYear() + '-' + pad(ts.getMonth() + 1) + '-' + pad(ts.getDate()) + 'T' + pad(ts.getHours()) + ':' + pad(ts.getMinutes()) + ':' + pad(ts.getSeconds());
            LOG_FILE.writeln('[' + iso + '] ' + String(message));
            LOG_FILE.close();
        }
    } catch (e) {
        $.writeln("[ERROR] Failed to write log file: " + e.message);
    }
}

/**
 * Logs a message and optionally shows an alert
 * @param {string} type - The type of log message (e.g., "ERROR", "WARN", "INFO")
 * @param {string} message - The message
 * @param {boolean} showAlert - Whether to show an alert dialog
 * @param {boolean} addNewline - Whether to add a newline after logging
 */
function log(type, message, showAlert, addNewline) {
    var logMessage;
    if (type === 'WARN' || type === 'warn') {
        logMessage = "[WARN] " + message;
    } else if (type === 'ERROR' || type === 'error') {
        logMessage = "[ERROR] " + message;
        if (showAlert) {
            alert("❌ " + message);
        }
    } else {
        logMessage = "[INFO] " + message;
    }
    $.writeln(logMessage);
    try { writeLogFile(logMessage); } catch (e) { }

    if (addNewline) {
        $.writeln("");
        try { writeLogFile(''); } catch (e) { }
    }
}

/**
 * Initializes the log file
 * @param {string} filename - Path to the log file
 */
function initLogFile(filename) {
    LOG_FILE = File(filename);
    try {
        if (LOG_FILE.exists) {
            LOG_FILE.remove();
        }
        LOG_FILE.encoding = 'UTF-8';
        if (LOG_FILE.open('w')) {
            var ts = new Date();
            var pad = function (n) { return (n < 10 ? '0' : '') + n; };
            var iso = ts.getFullYear() + '-' + pad(ts.getMonth() + 1) + '-' + pad(ts.getDate()) + 'T' + pad(ts.getHours()) + ':' + pad(ts.getMinutes()) + ':' + pad(ts.getSeconds());
            LOG_FILE.close();
        }
    } catch (e) {
        $.writeln("[ERROR] Could not initialize log file: " + e.message);
        alert("❌ Could not initialize log file: " + e.message);
    }
}

/**
 * Generates file paths for a community's CSV, INDD template, and output INDD
 * @param {File} csvFile - The CSV file
 * @param {Folder} templateFolder - The folder containing the INDD template
 * @param {Folder} outputFolder - The folder for output INDD files
 * @returns {Object} Object with csv, indd, and outputIndd File objects
 */
function getCommunityFile(csvFile, templateFolder, outputFolder) {
    var cleanName = csvFile.displayName.replace(/\.(csv|indd|pdf)$/i, "");
    var communityName = csvFile.parent.name.split('.')[0];

    return {
        csv: File(csvFile.fsName),
        indd: File(templateFolder + "/" + cleanName + ".indd"),
        outputIndd: File(outputFolder + "/" + communityName + "_" + cleanName + "_merged.indd")
    };
}

/**
 * Checks whether a document has overset text
 * @param {Document} doc - The document to check
 * @returns {Object} Object with details (array) and firstPage (number)
 */
function getOverflows(doc) {
    var details = [];
    var firstPage = -1;

    if (!doc || !doc.isValid) {
        return { details: [], firstPage: -1 };
    }

    // Check all pages for overset text
    for (var p = 0; p < doc.pages.length; p++) {
        var page = doc.pages[p];
        for (var f = 0; f < page.textFrames.length; f++) {
            if (page.textFrames[f].overflows) {
                if (firstPage === -1) {
                    firstPage = p;
                }
                details.push("- page " + (p + 1) + ", textbox " + (f + 1) + ": " + page.textFrames[f].contents.length + " chars");
            }
        }
    }
    if (details.length > 0) {
        log("WARN", "Overset text found: " + details.join("\n"), false, false);
    }

    return {
        details: details,
        firstPage: firstPage
    };
}

/**
 * Get user home directory (cross-platform)
 * @returns {string} The user's home directory path
 */
function getUserHome() {
    var home = $.getenv('USERPROFILE') || $.getenv('HOME');

    if (!home) {
        throw new Error('Could not determine user home directory. Neither USERPROFILE nor HOME environment variable is set.');
    }

    // Normalize path separators to forward slashes (works everywhere in ExtendScript)
    return home.replace(/\\/g, '/');
}

/**
 * Returns community folders that contain CSV files
 * @param {Folder} dataRoot
 * @returns {Array<Folder>}
 */
function getCommunityFolders(dataRoot) {
    return dataRoot.getFiles(function (f) {
        if (!(f instanceof Folder)) return false;
        var csvs = f.getFiles("*.csv");
        return csvs && csvs.length > 0;
    });
}

/**
 * Resolves community folders from a user's folder selection
 * @param {Folder} selectedFolder
 * @returns {Array<Folder>}
 */
function resolveCommunityFoldersFromFolder(selectedFolder) {
    if (!(selectedFolder instanceof Folder)) {
        return [];
    }
    var ownCSVs = selectedFolder.getFiles("*.csv");
    if (ownCSVs && ownCSVs.length > 0) {
        return [selectedFolder];
    }
    return getCommunityFolders(selectedFolder) || [];
}

/**
 * Checks whether the CSV file has at least header + one data row
 * @param {File} csvFile - The CSV file to check
 * @returns {boolean}
 */
function csvHasData(csvFile) {
    if (csvFile.length === 0) return false;

    try {
        csvFile.encoding = "UTF-16";
        if (!csvFile.open("r")) return false;

        var lineCount = 0;
        while (!csvFile.eof && lineCount < 2) {
            var line = csvFile.readln();
            if (line && String(line).replace(/^\s+|\s+$/g, "") !== "") {
                lineCount++;
            }
        }
        csvFile.close();
        return lineCount >= 2;
    } catch (e) {
        try { csvFile.close(); } catch (ee) { }
        return false;
    }
}

(function () {
    // === Configuration ===
    var max_data_merge_retries = 20;
    var data_merge_retry_delay_ms = 200;
    var log_file_name = "_merge_log.txt";

    // === Path Configuration ===
    var userSelection = Folder.selectDialog("Select a folder", Folder(EXPORTS_DIR));
    if (!userSelection) {
        log("ERROR", "Script aborted: No base folder selected.", true, true);
        return; // Exit the script
    }

    var templateRoot = Folder(TEMPLATES_DIR);
    var outputRoot = Folder(OUTPUT_DIR);

    assertFolder(userSelection);
    assertFolder(templateRoot);

    var communityFolders = resolveCommunityFoldersFromFolder(userSelection);

    if (communityFolders.length === 0) {
        log("ERROR", "Could not find any CSV files in the selected folder or its subfolders: '" + userSelection.fsName + "'", true, true);
        return;
    }

    var timestampName = communityFolders[0].parent.name;
    var outputSubFolder = communityFolders[0].parent.parent.name;
    if (outputSubFolder === "exports") {
        outputSubFolder = ""
    } else {
        outputSubFolder = "/" + outputSubFolder;
    }

    var outputFolder = Folder(outputRoot + outputSubFolder + "/" + timestampName);
    ensureFolder(outputFolder);
    initLogFile(outputFolder + "/" + log_file_name);

    // === Process each community folder ===
    for (var i = 0; i < communityFolders.length; i++) {
        var communityName = communityFolders[i].name;
        var templateFolder = Folder(templateRoot + "/" + communityName);
        if (!templateFolder.exists) {
            log("ERROR", "Skipping pdf generation for '" + communityName + "' since there is no templates folder: '" + templateFolder.fsName + "'", true, true);
            continue;
        }

        var csvFiles = communityFolders[i].getFiles("*.csv");

        // === Loop over all CSV files ===
        for (var j = 0; j < csvFiles.length; j++) {
            var csvFile = csvFiles[j];
            var communityFile = getCommunityFile(csvFile, templateFolder, outputFolder);

            if (!csvHasData(csvFile)) {
                log("ERROR", "Skipping pdf generation for '" + communityFile.csv.fsName + "' since it does not contain data rows", true, true);
                continue;
            }

            var templateFile = communityFile.indd;
            if (!templateFile.exists) {
                log("ERROR", "Skipping pdf generation for '" + communityFile.csv.fsName + "' since there is no template: '" + communityFile.indd.fsName + "'", true, true);
                continue;
            }

            var template = loadTemplate(templateFile);
            log("INFO", "Loaded template: " + communityFile.indd.fsName, false, false);
            try {
                // === Remove old data source (if present) ===
                try {
                    template.dataMergeProperties.removeDataSource();
                    log("INFO", "Removed old data source", false, false);
                } catch (e) {
                    // No old data source present
                }

                template.dataMergeProperties.selectDataSource(csvFile);

                // 🟢 Safely update data source (prevents "empty" first merges)
                // InDesign loads data fields asynchronously - there are no callbacks/promises in ExtendScript!
                // Therefore: Polling loop is the only option
                template.dataMergeProperties.updateDataSource();

                var attempt;
                for (attempt = 1; attempt <= max_data_merge_retries; attempt++) {
                    if (template.dataMergeProperties.dataMergeFields.length > 0) {
                        log("INFO", "Loaded CSV: " + csvFile.fsName + " on attempt " + attempt, false, false);
                        break;
                    }
                    if (attempt < max_data_merge_retries) {
                        $.sleep(data_merge_retry_delay_ms);
                    }
                }

                if (template.dataMergeProperties.dataMergeFields.length === 0) {
                    log("ERROR", "Could not load CSV data from: " + csvFile.name, true, true);
                    continue;
                }

                // Merge records - creates a new document window with merged data
                template.dataMergeProperties.mergeRecords(File(outputFolder));

                // The merged document is now the active document
                var mergedDoc = app.activeDocument;
                if (!mergedDoc || !mergedDoc.isValid) {
                    log("ERROR", "Could not create merged document.", true, true);
                    continue;
                }

                // Set the correct document title
                mergedDoc.name = communityFile.outputIndd.name + "_merged";

                if (mergedDoc && mergedDoc.isValid) {
                    log("INFO", "Merged document created with " + mergedDoc.pages.length + " page(s)", false, false);

                    // Check for overset text
                    var overflows = getOverflows(mergedDoc);

                    if (overflows.details.length > 0) {
                        // Scroll to first page with overset
                        if (overflows.firstPage >= 0 && mergedDoc.windows.length > 0) {
                            try {
                                var window = mergedDoc.windows[0];
                                window.activeSpread = mergedDoc.pages[overflows.firstPage].parent;
                                log("INFO", "Scrolled to page " + (overflows.firstPage + 1), false, false);
                            } catch (e) {
                                log("ERROR", "Could not scroll to page: " + e.message, false, true);
                            }
                        }

                        log("ERROR", "Found text overflows on:\n" + overflows.details.join("\n"), false, true);
                    } else {
                        log("INFO", "No overset text detected", false, false);

                        // Save indd file since no oversets were found
                        try {
                            ensureFolder(outputFolder);
                            mergedDoc.save(File(communityFile.outputIndd));
                            log("INFO", "Saved merged document: " + communityFile.outputIndd.fsName, false, true);

                            // Close the merged document after successful export
                            mergedDoc.close(SaveOptions.NO);
                        } catch (saveErr) {
                            log("ERROR", "Error saving merged document: " + saveErr.message, true, true);
                        }
                    }
                } else {
                    log("WARN", "Could not get merged document", false, true);
                }
            } catch (err) {
                log("ERROR", "Error while generating '" + communityFile.outputIndd.fsName + "': " + err.message, true, true);
            } finally {
                // Close the template document
                try {
                    if (template && template.isValid) {
                        template.close(SaveOptions.NO);
                    }
                } catch (e) {
                    log("ERROR", "Error while closing the template " + e.message, true, true);
                }
            }
        }
    }
})();
