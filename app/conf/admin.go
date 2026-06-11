package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	activeConfigPathMu sync.RWMutex
	activeConfigPath   string
)

var ErrAdminConfigUnavailable = errors.New("admin config file is unavailable")

type AdminSecret struct {
	Index  int    `json:"index"`
	Set    bool   `json:"set"`
	Masked string `json:"masked,omitempty"`
	Value  string `json:"value,omitempty"`
}

type AdminChatGPT struct {
	Index        int         `json:"index"`
	IdToken      AdminSecret `json:"id_token"`
	AccessToken  AdminSecret `json:"access_token"`
	RefreshToken AdminSecret `json:"refresh_token"`
	AccountID    string      `json:"account_id"`
	LastRefresh  string      `json:"last_refresh"`
	Email        string      `json:"email"`
	Type         string      `json:"type"`
	Expired      string      `json:"expired"`
	Proxy        string      `json:"proxy"`
}

type AdminConfig struct {
	ConfigPath          string         `json:"config_path"`
	Proxy               string         `json:"proxy"`
	ChatGPTBaseURL      string         `json:"chatgpt_base_url"`
	AuthTokens          []AdminSecret  `json:"auth_tokens"`
	AccessTokenPrefixes []AdminSecret  `json:"access_token_prefixes"`
	ChatGPTs            []AdminChatGPT `json:"chatgpts"`
}

func setActiveConfigPath(path string) {
	activeConfigPathMu.Lock()
	defer activeConfigPathMu.Unlock()
	activeConfigPath = path
}

func getActiveConfigPath() string {
	activeConfigPathMu.RLock()
	defer activeConfigPathMu.RUnlock()
	return activeConfigPath
}

func LoadAdminConfig() (AdminConfig, error) {
	path := getActiveConfigPath()
	if strings.TrimSpace(path) == "" {
		return AdminConfig{}, ErrAdminConfigUnavailable
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return AdminConfig{}, err
	}
	cfg := defaultApp()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return AdminConfig{}, err
	}
	return adminConfigFromApp(path, cfg), nil
}

func SaveAdminConfig(req AdminConfig) (AdminConfig, error) {
	path := getActiveConfigPath()
	if strings.TrimSpace(path) == "" {
		return AdminConfig{}, ErrAdminConfigUnavailable
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return AdminConfig{}, err
	}
	current := defaultApp()
	if err := yaml.Unmarshal(data, &current); err != nil {
		return AdminConfig{}, err
	}

	next := current
	next.Proxy = strings.TrimSpace(req.Proxy)
	next.ChatGPTBaseUrl = strings.TrimSpace(req.ChatGPTBaseURL)
	next.Auth.AccessTokens = resolveAdminSecretList(req.AuthTokens, current.Auth.AccessTokens, normalizeAuthToken)
	next.Auth.AccessTokenPrefix = nonEmptyAccessTokenPrefixes(resolveAdminSecretList(req.AccessTokenPrefixes, current.Auth.AccessTokenPrefix, strings.TrimSpace))
	if len(nonEmptyAuthTokens(next.Auth.AccessTokens)) == 0 {
		return AdminConfig{}, fmt.Errorf("at least one local API key is required")
	}
	next.ChatGPTs = resolveAdminChatGPTs(req.ChatGPTs, current.ChatGPTs)

	if err := writeAdminConfig(path, data, next); err != nil {
		return AdminConfig{}, err
	}
	normalized := next
	normalizeConfig(&normalized)
	setApp(normalized)
	return adminConfigFromApp(path, next), nil
}

func adminConfigFromApp(path string, cfg app) AdminConfig {
	out := AdminConfig{
		ConfigPath:          path,
		Proxy:               cfg.Proxy,
		ChatGPTBaseURL:      cfg.ChatGPTBaseUrl,
		AuthTokens:          maskedAdminSecrets(cfg.Auth.AccessTokens, normalizeAuthToken),
		AccessTokenPrefixes: maskedAdminSecrets(cfg.Auth.AccessTokenPrefix, strings.TrimSpace),
		ChatGPTs:            make([]AdminChatGPT, 0, len(cfg.ChatGPTs)),
	}
	for i, account := range cfg.ChatGPTs {
		out.ChatGPTs = append(out.ChatGPTs, AdminChatGPT{
			Index:        i,
			IdToken:      maskedAdminSecret(i, account.IdToken, strings.TrimSpace),
			AccessToken:  maskedAdminSecret(i, account.AccessToken, normalizeAuthToken),
			RefreshToken: maskedAdminSecret(i, account.RefreshToken, strings.TrimSpace),
			AccountID:    account.AccountId,
			LastRefresh:  account.LastRefresh,
			Email:        account.Email,
			Type:         account.Type,
			Expired:      account.Expired,
			Proxy:        account.Proxy,
		})
	}
	return out
}

func maskedAdminSecrets(values []string, normalize func(string) string) []AdminSecret {
	out := make([]AdminSecret, 0, len(values))
	for i, value := range values {
		out = append(out, maskedAdminSecret(i, value, normalize))
	}
	return out
}

func maskedAdminSecret(index int, value string, normalize func(string) string) AdminSecret {
	value = normalize(value)
	return AdminSecret{
		Index:  index,
		Set:    value != "",
		Masked: maskToken(value),
	}
}

func resolveAdminSecretList(in []AdminSecret, current []string, normalize func(string) string) []string {
	out := make([]string, 0, len(in))
	for _, item := range in {
		value := resolveAdminSecret(item, current, normalize)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}

func resolveAdminSecret(item AdminSecret, current []string, normalize func(string) string) string {
	if value := normalize(item.Value); value != "" {
		return value
	}
	if item.Set && item.Index >= 0 && item.Index < len(current) {
		return normalize(current[item.Index])
	}
	return ""
}

func resolveAdminChatGPTs(in []AdminChatGPT, current []chatgpt) []chatgpt {
	out := make([]chatgpt, 0, len(in))
	for _, item := range in {
		var currentAccount chatgpt
		if item.Index >= 0 && item.Index < len(current) {
			currentAccount = current[item.Index]
		}
		account := chatgpt{
			IdToken:      resolveAdminSecretValue(item.IdToken, currentAccount.IdToken, strings.TrimSpace),
			AccessToken:  resolveAdminSecretValue(item.AccessToken, currentAccount.AccessToken, normalizeAuthToken),
			RefreshToken: resolveAdminSecretValue(item.RefreshToken, currentAccount.RefreshToken, strings.TrimSpace),
			AccountId:    strings.TrimSpace(item.AccountID),
			LastRefresh:  strings.TrimSpace(item.LastRefresh),
			Email:        strings.TrimSpace(item.Email),
			Type:         strings.TrimSpace(item.Type),
			Expired:      strings.TrimSpace(item.Expired),
			Proxy:        strings.TrimSpace(item.Proxy),
		}
		if isEmptyChatGPT(account) {
			continue
		}
		out = append(out, account)
	}
	return out
}

func resolveAdminSecretValue(item AdminSecret, current string, normalize func(string) string) string {
	if value := normalize(item.Value); value != "" {
		return value
	}
	if item.Set {
		return normalize(current)
	}
	return ""
}

func isEmptyChatGPT(account chatgpt) bool {
	return account.IdToken == "" &&
		account.AccessToken == "" &&
		account.RefreshToken == "" &&
		account.AccountId == "" &&
		account.LastRefresh == "" &&
		account.Email == "" &&
		account.Type == "" &&
		account.Expired == "" &&
		account.Proxy == ""
}

func writeAdminConfig(path string, original []byte, cfg app) error {
	var doc yaml.Node
	if err := yaml.Unmarshal(original, &doc); err != nil {
		return err
	}
	root := mappingRoot(&doc)
	setStringChild(root, "proxy", cfg.Proxy)
	setStringChild(root, "chatgpt_base_url", cfg.ChatGPTBaseUrl)

	authNode := ensureMappingChild(root, "auth")
	setStringSequenceChild(authNode, "access_tokens", cfg.Auth.AccessTokens, normalizeAuthToken)
	setStringSequenceChild(authNode, "access_token_prefix", cfg.Auth.AccessTokenPrefix, strings.TrimSpace)
	setMappingChild(root, "chatgpts", chatGPTSequenceNode(cfg.ChatGPTs))

	data, err := yaml.Marshal(&doc)
	if err != nil {
		return err
	}
	return writeFileWithBackup(path, original, data)
}

func setStringChild(root *yaml.Node, key string, value string) {
	setMappingChild(root, key, &yaml.Node{Kind: yaml.ScalarNode, Value: strings.TrimSpace(value)})
}

func setStringSequenceChild(root *yaml.Node, key string, values []string, normalize func(string) string) {
	seq := &yaml.Node{Kind: yaml.SequenceNode}
	for _, value := range values {
		value = normalize(value)
		if value == "" {
			continue
		}
		seq.Content = append(seq.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: value})
	}
	setMappingChild(root, key, seq)
}

func chatGPTSequenceNode(accounts []chatgpt) *yaml.Node {
	seq := &yaml.Node{Kind: yaml.SequenceNode}
	for _, account := range accounts {
		item := &yaml.Node{Kind: yaml.MappingNode}
		setStringChild(item, "id_token", strings.TrimSpace(account.IdToken))
		setStringChild(item, "access_token", normalizeAuthToken(account.AccessToken))
		setStringChild(item, "refresh_token", strings.TrimSpace(account.RefreshToken))
		setStringChild(item, "account_id", strings.TrimSpace(account.AccountId))
		setStringChild(item, "last_refresh", strings.TrimSpace(account.LastRefresh))
		setStringChild(item, "email", strings.TrimSpace(account.Email))
		setStringChild(item, "type", strings.TrimSpace(account.Type))
		setStringChild(item, "expired", strings.TrimSpace(account.Expired))
		setStringChild(item, "proxy", strings.TrimSpace(account.Proxy))
		seq.Content = append(seq.Content, item)
	}
	return seq
}

func writeFileWithBackup(path string, original []byte, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	backupPath := path + ".bak"
	if err := os.WriteFile(backupPath, original, 0600); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err == nil {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Rename(tmpPath, path)
}
