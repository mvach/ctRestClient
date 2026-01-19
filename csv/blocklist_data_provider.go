package csv

import (
	"bytes"
	"ctRestClient/config"
	"ctRestClient/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//counterfeiter:generate . BlockListDataProvider
type BlockListDataProvider interface {
	IsBlocked(personJson map[string]json.RawMessage, group config.Group) (bool, error)

	BlockListExists(group config.Group) bool
}

type addressData map[string]interface{}

type cacheEntry struct {
	data []addressData
	err  error
}

type blockListDataProvider struct {
	dataDir   string
	dataCache map[string]cacheEntry
	logger    logger.Logger
}

func NewBlockListDataProvider(dataDir string, logger logger.Logger) BlockListDataProvider {
	return &blockListDataProvider{
		dataDir:   dataDir,
		dataCache: make(map[string]cacheEntry),
		logger:    logger,
	}
}

func (bp *blockListDataProvider) IsBlocked(personJson map[string]json.RawMessage, group config.Group) (bool, error) {
	var entry cacheEntry
	entry, err := bp.loadBlocklist(group.BlocklistFileName())
	if err != nil {
		return false, fmt.Errorf("failed to load blocklist %s: %w", group.BlocklistFileName(), err)
	}

	for _, blockedAddress := range entry.data {
		matched := true

		for fieldName, blockedFieldValue := range blockedAddress {
			if personJsonFieldValue, exists := personJson[fieldName]; !exists {
				bp.logger.Warn(fmt.Sprintf("      Ignoring blocklist element %v since field '%s' is not available in the person data", blockedAddress, fieldName))
				matched = false
				break
			} else {

				// Unescape the json data that contains (\u00df and \u00fc instead of ß and ü) for later string comparison with data from the blocklist.
				personJsonValue, err := unescapeUnicodeCharacters(personJsonFieldValue)
				if err != nil {
					return false, fmt.Errorf("failed to unescape unicode characters for field %s: %w", fieldName, err)
				}

				if !bytes.Equal([]byte(personJsonValue), blockedFieldValue.(json.RawMessage)) {
					matched = false
					break
				}
			}
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

func (bp *blockListDataProvider) loadBlocklist(name string) (cacheEntry, error) {

	if entry, ok := bp.dataCache[name]; ok {
		// ---- CACHE HIT ----
		return entry, entry.err
	}

	path := filepath.Join(bp.dataDir, name)

	yamlData, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		entry := cacheEntry{}
		bp.dataCache[name] = entry
		return entry, nil
	}

	if err != nil {
		entry := cacheEntry{err: err}
		bp.dataCache[name] = entry
		return entry, err
	}

	if len(yamlData) == 0 {
		entry := cacheEntry{data: []addressData{}}
		bp.dataCache[name] = entry
		return entry, nil
	}

	var blocked []addressData
	if err := yaml.Unmarshal(yamlData, &blocked); err != nil {
		entry := cacheEntry{err: err}
		bp.dataCache[name] = entry
		return entry, err
	}

	for i, addr := range blocked {
		for key, val := range addr {
			jsonBytes, err := json.Marshal(val)
			if err != nil {
				return cacheEntry{}, fmt.Errorf("failed to marshal blocklist value for key %s: %w", key, err)
			}
			addr[key] = json.RawMessage(jsonBytes)
		}
		blocked[i] = addr
	}

	entry := cacheEntry{data: blocked}
	bp.dataCache[name] = entry

	return entry, nil
}

func (bp *blockListDataProvider) BlockListExists(group config.Group) bool {
	blocklistFilePath := filepath.Join(bp.dataDir, group.BlocklistFileName())
	_, err := os.Stat(blocklistFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			bp.logger.Error(fmt.Sprintf("      failed to evaluate blocklist existence: %v", err))
		}
	}

	return true
}

func unescapeUnicodeCharacters(jsonRaw json.RawMessage) (json.RawMessage, error) {
	var temp interface{}
	if err := json.Unmarshal(jsonRaw, &temp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal person json field")
	}

	result, err := json.Marshal(temp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal person json field")
	}
	return result, nil
}
